package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jjkirkpatrick/clara-ollama/internal/chatui"
	"github.com/jjkirkpatrick/clara-ollama/internal/config"
	"github.com/jjkirkpatrick/clara-ollama/internal/plugins"
	"github.com/ollama/ollama/api"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

var Plugin plugins.Plugin = &DateTime{}

type DateTime struct {
	cfg config.Cfg
}

func (c *DateTime) Init(cfg config.Cfg, ollamaClient *api.Client, chat *chatui.ChatUI) error {
	c.cfg = cfg

	c.cfg.AppLogger.Info("DateTime plugin initialized successfully")
	return nil
}

func (c DateTime) ID() string {
	return "datetime"
}

func (c DateTime) Description() string {
	return "parse date and time from natural language input."
}

func (c DateTime) FunctionDefinition() openai.FunctionDefinition {
	return openai.FunctionDefinition{
		Name:        "datetime",
		Description: "Parse date and time from natural language input such as 'now', 'tomorrow', 'yesterday', '1 month from now', etc. and return the exact date.",
		Parameters: jsonschema.Definition{
			Type: jsonschema.Object,
			Properties: map[string]jsonschema.Definition{
				"input": {
					Type:        jsonschema.String,
					Description: "The natural language date/time input to parse.",
				},
			},
			Required: []string{"input"},
		},
	}
}

func (c DateTime) Execute(jsonInput string) (string, error) {
	// marshal jsonInput to inputDefinition
	var args map[string]string
	err := json.Unmarshal([]byte(jsonInput), &args)
	if err != nil {
		c.cfg.AppLogger.Info("Error unmarshalling JSON input: ", err)
		return "", err
	}

	input, ok := args["input"]
	if !ok {
		return fmt.Sprintf(`%v`, "input is required but was not provided"), nil
	}

	// Parse the date/time from the input
	t, err := dateparse.ParseAny(input)
	if err != nil {
		c.cfg.AppLogger.Info("Error parsing date/time: ", err)
		return fmt.Sprintf(`%v`, err), err
	}

	// Format the date/time as a string
	dateStr := t.Format(time.RFC3339)

	return dateStr, nil
}
