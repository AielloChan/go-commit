package main

import (
	"os"

	"github.com/AlecAivazis/survey/v2/terminal"
	"yuanling.com/go-commit/config"
	"yuanling.com/go-commit/model"
	"yuanling.com/go-commit/pipline"
	"yuanling.com/go-commit/tools"
)

func main() {
	// Create data store
	store := model.GetInstance()
	// Get config
	cfg, err := config.InitConfig()
	if err != nil {
		tools.PrintError(err.Error() + "\n")
		os.Exit(0)
	}
	// Run pipline
	err = pipline.RunJob(&cfg.Stages, 0, store)
	// Ctrl + c
	if err == terminal.InterruptErr {
		os.Exit(0)
	} else if err != nil {
		tools.PrintError(err.Error() + "\n")
		os.Exit(0)
	}
}
