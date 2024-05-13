package assistant

import (
	"fmt"
	"os"

	"github.com/inancgumus/screen"
)

var commands = []string{
	"/restart",
	"/exit",
}

func (assistant assistant) paraseCommandsFromInput(message string) bool {

	// check to see if the message is a command
	for _, command := range commands {
		if message == command {

			// handle the command
			assistant.cfg.AppLogger.Info("Command received: " + message)
			assistant.handleCommand(message)
			return true
		}
	}
	return false
}

func (assistant assistant) handleCommand(command string) {
	switch command {
	case "/restart":
		assistant.cfg.AppLogger.Info("Restart command received")
		assistant.restartConversation()
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Println("Conversation restarted")
	case "/exit":
		assistant.cfg.AppLogger.Info("Exit command received")
		screen.Clear()
		screen.MoveTopLeft()
		fmt.Println("Exiting...")
		os.Exit(0)
	}
}
