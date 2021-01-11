package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"yuanling.com/go-commit/configs"
	"yuanling.com/go-commit/model"
	"yuanling.com/go-commit/tools"
)

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 // not found.
}

type StringArray []string

func (sa *StringArray) FindIndex(s string) int {
	for i, str := range *sa {
		if str == s {
			return i
		}
	}
	return -1
}

func doSurvey(stages *[]configs.Stage, index int, store *model.Store) error {
	currentStage := (*stages)[index]

	switch currentStage.Type {
	case "select":
		selectedLabel := ""
		options := make(StringArray, len(currentStage.Config.Options))
		for i := 0; i < len(currentStage.Config.Options); i++ {
			options[i] = tools.GetString(currentStage.Config.Options[i].Label, store)
		}
		prompt := &survey.Select{
			Message:  tools.GetString(currentStage.Label, store),
			Options:  options,
			PageSize: currentStage.Config.Max,
		}
		survey.AskOne(prompt, &selectedLabel)

		selectOption := currentStage.Config.Options[options.FindIndex(selectedLabel)]
		(*store)[currentStage.Name] = selectOption.Value
		optionNext := selectOption.Next
		if optionNext != "" {
			for i, item := range *stages {
				if item.Name == optionNext {
					doSurvey(stages, i, store)
					return nil
				}
			}
		}
	case "multi-select":
		selectedLabels := []string{}
		options := make(StringArray, len(currentStage.Config.Options))
		for i := 0; i < len(currentStage.Config.Options); i++ {
			options[i] = tools.GetString(currentStage.Config.Options[i].Label, store)
		}
		prompt := &survey.MultiSelect{
			Message:  tools.GetString(currentStage.Label, store),
			Options:  options,
			PageSize: currentStage.Config.Max,
		}
		survey.AskOne(prompt, &selectedLabels)
		if len(selectedLabels) > 0 {
			selectedValues := []string{}
			for _, selectedLabel := range selectedLabels {
				selectedOption := currentStage.Config.Options[options.FindIndex(selectedLabel)]
				selectedValues = append(selectedValues, selectedOption.Value)
			}
			(*store)[currentStage.Name] = strings.Join(selectedValues[:], ",")
		} else {
			// 没有选任何项
		}
	case "string":
		value := ""
		if currentStage.Config.Default != nil {
			value = tools.GetString(currentStage.Config.Default.(string), store)
			fmt.Println(value)
		}
		prompt := &survey.Input{
			Message: tools.GetString(currentStage.Label, store),
		}
		survey.AskOne(prompt, &value,
			survey.WithValidator(survey.MinLength(currentStage.Config.Min)),
			survey.WithValidator(survey.MaxLength(currentStage.Config.Max)),
		)
		(*store)[currentStage.Name] = value
	case "multiline":
		value := ""
		if currentStage.Config.Default != nil {
			value = tools.GetString(currentStage.Config.Default.(string), store)
			fmt.Println(value)
		}
		prompt := &survey.Multiline{
			Message: tools.GetString(currentStage.Label, store),
		}
		survey.AskOne(prompt, &value,
			survey.WithValidator(survey.MinLength(currentStage.Config.Min)),
			survey.WithValidator(survey.MaxLength(currentStage.Config.Max)),
		)
		(*store)[currentStage.Name] = value
	case "confirm":
		answer := false
		prompt := &survey.Confirm{
			Message: tools.GetString(currentStage.Label, store),
			Default: currentStage.Config.Default.(bool),
		}
		survey.AskOne(prompt, &answer)
		(*store)[currentStage.Name] = strconv.FormatBool(answer)
	default:
		// error
		return errors.New("Unknown type: " + currentStage.Type)
	}

	stageNext := currentStage.Next
	if stageNext != "" {
		for i, item := range *stages {
			if item.Name == stageNext {
				doSurvey(stages, i, store)
				return nil
			}
		}
	}
	if index < len(*stages)-1 {
		doSurvey(stages, index+1, store)
	}
	return nil
}

func main() {
	cfg, err := configs.GetConfig("./commit.config.json")
	if err != nil {
		fmt.Printf("Error while read config: %s", err)
	}
	if len(cfg.Stages) == 0 {
		fmt.Println("无可用的 stages 配置")
		os.Exit(0)
	}

	store := model.GetInstance()
	doSurvey(&cfg.Stages, 0, store)

	confirmed := true
	if cfg.Preview {
		fmt.Println(tools.GetString(`#${{.type}}({{.scope}}{{.customScope}}): {{.title}}\n\n{{.body}}`, store))
		prompt := &survey.Confirm{
			Message: "确认提交信息?",
			Default: true,
		}
		survey.AskOne(prompt, &confirmed)
	}
	if confirmed {
		fmt.Println("Success")
	} else {
		fmt.Println("Cancel")
	}
}
