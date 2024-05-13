package assistant

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jjkirkpatrick/clara-ollama/internal/chatui"
	"github.com/jjkirkpatrick/clara-ollama/internal/config"
	"github.com/jjkirkpatrick/clara-ollama/internal/plugins"
	"github.com/ollama/ollama/api"
	"github.com/sashabaranov/go-openai"
)

type assistant struct {
	cfg                 config.Cfg
	Client              *api.Client
	functionDefinitions []openai.FunctionDefinition
	chat                *chatui.ChatUI
}

var conversation []api.Message

var systemPrompt = `
You have access to the following plugins: ${functionDefinitionsJSON}
Follow these instructions when responding to user queries:
1. Evaluate the user's request to determine if it matches the capabilities of the available plugins.
2. If a suitable plugin is identified, respond strictly with a JSON object that specifies the selected plugin and its required input parameters, following the JSON schema specified at the plugin definition.
`

func (assistant assistant) getSystemPrompt() string {
	// Get the function definitions
	functionDefinitions := plugins.GenerateOpenAIFunctionsDefinition()

	// Serialize definitions to JSON
	b, err := json.Marshal(functionDefinitions)
	if err != nil {
		return "{}" // return empty JSON object on error
	}

	// Replace the placeholder with the actual function definitions
	prompt := strings.Replace(systemPrompt, "${functionDefinitionsJSON}", string(b), 1)

	return prompt
}

func appendMessage(role string, message string) {
	conversation = append(conversation, api.Message{
		Role:    role,
		Content: message,
	})
}

func (assistant assistant) restartConversation() {
	resetConversation()
	// append the system prompt to the conversation
	appendMessage("system", assistant.getSystemPrompt())

	// send the system prompt to openai
	response, err := assistant.sendMessage()

	if err != nil {
		assistant.cfg.AppLogger.Fatalf("Error sending system prompt to OpenAI: %v", err)
	}

	// append the assistant message to the conversation
	appendMessage("assistant", response)

}

func resetConversation() {
	conversation = []api.Message{}
}

func (assistant assistant) Message(message string) (string, error) {

	assistant.chat.DisableInput()
	assistant.cfg.AppLogger.Info("Message input disabled")
	//check to see if the message is a command
	//if it is, handle the command and return
	if assistant.paraseCommandsFromInput(message) {
		return "", nil
	}

	// append the user message to the conversation
	appendMessage(openai.ChatMessageRoleUser, message)

	response, err := assistant.sendMessage()

	if err != nil {
		return "", err
	}

	// append the assistant message to the conversation
	appendMessage(openai.ChatMessageRoleAssistant, response)
	// print the conversation
	assistant.chat.AddMessage("Clara", response)

	assistant.chat.EnableInput()
	assistant.cfg.AppLogger.Info("Message input enabled")

	return response, nil
}

func (assistant assistant) sendMessage() (string, error) {
	resp, err := assistant.sendRequestToOllama()

	if err != nil {
		return "", err
	}

	if strings.Contains(resp.Message.Content, "{\"plugin\":") {
		assistant.cfg.AppLogger.Info("Handling function call")
		responseContent, err := assistant.handleFunctionCall(&resp.Message)
		if err != nil {
			return "", err
		}
		return responseContent, nil
	}
	return resp.Message.Content, nil
}

func (assistant assistant) handleFunctionCall(funcCall *api.Message) (string, error) {
	var functionCall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
	assistant.cfg.AppLogger.Info(funcCall.Content)
	err := json.Unmarshal([]byte(funcCall.Content), &functionCall)
	if err != nil {
		return "", fmt.Errorf("failed to parse function call: %v", err)
	}

	// check to see if a plugin is loaded with the same name as the function call
	ok := plugins.IsPluginLoaded(functionCall.Name)

	if !ok {
		return "", fmt.Errorf("no plugin loaded with name %v", functionCall.Name)
	}

	// convert arguments to string
	jsonInput, err := json.Marshal(functionCall.Arguments)
	if err != nil {
		return "", fmt.Errorf("failed to marshal function arguments: %v", err)
	}

	assistant.cfg.AppLogger.Info(functionCall.Name)
	// call the plugin with the arguments
	jsonResponse, err := plugins.CallPlugin(functionCall.Name, string(jsonInput))

	if err != nil {
		return "", err
	}
	appendMessage("assistant", string(jsonInput))
	appendMessage("assistant", jsonResponse)

	resp, err := assistant.sendRequestToOllama()
	if err != nil {
		return "", err
	}

	// Check if the response is another function call
	if strings.Contains(resp.Message.Content, "{\"plugin\":") {
		return assistant.handleFunctionCall(&api.Message{Content: resp.Message.Content})
	}

	return resp.Message.Content, nil
}

func (assistant assistant) sendRequestToOllama() (*api.ChatResponse, error) {
	ctx := context.Background()
	assistant.cfg.AppLogger.Info(conversation)
	var stream bool = false
	stream = false
	req := &api.ChatRequest{
		Model:    "llama3",
		Messages: conversation,
		Stream:   &stream,
	}

	var chatResponse *api.ChatResponse
	respFunc := func(resp api.ChatResponse) error {
		chatResponse = &resp
		assistant.cfg.AppLogger.Info(chatResponse)
		return nil
	}

	err := assistant.Client.Chat(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}
	return chatResponse, nil
}

func Start(cfg config.Cfg, ollamaClient *api.Client, chat *chatui.ChatUI) assistant {
	if err := plugins.LoadPlugins(cfg, ollamaClient, chat); err != nil {
		cfg.AppLogger.Fatalf("Error loading plugins: %v", err)
	}

	assistant := assistant{
		cfg:                 cfg,
		Client:              ollamaClient,
		functionDefinitions: plugins.GenerateOpenAIFunctionsDefinition(),
		chat:                chat,
	}

	assistant.chat.ClearHistory()

	assistant.restartConversation()

	return assistant

}
