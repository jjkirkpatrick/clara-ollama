package main

import (
	"time"

	"github.com/jjkirkpatrick/clara-ollama/internal/assistant"
	"github.com/jjkirkpatrick/clara-ollama/internal/chatui"
	"github.com/jjkirkpatrick/clara-ollama/internal/config"
	"github.com/ollama/ollama/api"
)

var cfg = config.New()
var ollamaClient *api.Client

func main() {
	cfg.AppLogger.Info("Clara is starting up... Please wait a moment.")

	var err error
	ollamaClient, err = api.ClientFromEnvironment()
	if err != nil {
		cfg.AppLogger.Fatal(err)
	}

	chat, err := chatui.NewChatUI()
	clara := assistant.Start(cfg, ollamaClient, chat)

	if err != nil {
		cfg.AppLogger.Fatalf("Error initializing chat UI: %v", err)
	}

	go func() {
		if err := chat.Run(); err != nil {
			cfg.AppLogger.Fatalf("Error running chat UI: %v", err)
		}
	}()

	userMessagesChan := chat.GetUserMessagesChannel()
	for {
		select {
		case userMessage, ok := <-userMessagesChan: // userMessage is a string containing the user's message.
			if !ok {
				// If the channel is closed, exit the loop.
				cfg.AppLogger.Info("User message channel closed. Exiting.")
				return
			}

			clara.Message(userMessage)
		case <-time.After(10 * time.Minute):
			// Timeout: if there's no activity for 5 minutes, exit.
			cfg.AppLogger.Info("No activity for 10 minutes. Exiting.")
			return
		}
	}

}
