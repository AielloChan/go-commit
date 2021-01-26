package pipline

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"yuanling.com/go-commit/config"
	"yuanling.com/go-commit/model"
	"yuanling.com/go-commit/tools"
)

const END_ROUTE = "$END"

// StringList a type for store string list
type StringList []string

// FindIndex a method for find index from a string list
func (sa *StringList) FindIndex(s string) int {
	for i, str := range *sa {
		if str == s {
			return i
		}
	}
	return -1
}

// GoOriginalNext Goto next stage
func GoOriginalNext(next string, stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	if next == END_ROUTE {
		return nil
	}
	if idx, ok := config.FindStageIndexByName(stages, next); ok {
		// process configured next stage
		return RunJob(stages, idx, store)
	}
	if curStage.Next == END_ROUTE {
		return nil
	}
	if idx, ok := config.FindStageIndexByName(stages, curStage.Next); ok {
		// process configured next stage
		return RunJob(stages, idx, store)
	}
	if index < len(*stages)-1 {
		// process rest stage
		return RunJob(stages, index+1, store)
	}
	return nil
}

// ProcessString string input handler
func ProcessString(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	value := ""
	if curStage.Config.Default != nil {
		value, err := tools.ExecCommand(curStage.Config.Default.(string), store)
		if err != nil {
			return err
		}
		fmt.Println(value)
	}
	res, err := tools.ExecCommand(curStage.Label, store)
	if err != nil {
		return err
	}

	prompt := &survey.Input{
		Message: res,
	}
	err = survey.AskOne(prompt, &value,
		survey.WithValidator(survey.MinLength(curStage.Config.Min)),
		survey.WithValidator(survey.MaxLength(curStage.Config.Max)),
	)
	if err != nil {
		return err
	}
	(*store)[curStage.Name] = value
	return GoOriginalNext("", stages, index, store)
}

// ProcessMultiLine multiLine input handler
func ProcessMultiLine(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]

	value := ""
	if curStage.Config.Default != nil {
		value, err := tools.ExecCommand(curStage.Config.Default.(string), store)
		if err != nil {
			return err
		}
		fmt.Println(value)
	}
	res, err := tools.ExecCommand(curStage.Label, store)
	if err != nil {
		return err
	}

	prompt := &survey.Multiline{
		Message: res,
	}
	err = survey.AskOne(prompt, &value,
		survey.WithValidator(survey.MinLength(curStage.Config.Min)),
		survey.WithValidator(survey.MaxLength(curStage.Config.Max)),
	)
	if err != nil {
		return err
	}
	(*store)[curStage.Name] = value
	return GoOriginalNext("", stages, index, store)
}

// ProcessSelect select hander
func ProcessSelect(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	selectedLabel := ""
	options := make(StringList, len(curStage.Config.Options))
	for i := 0; i < len(curStage.Config.Options); i++ {
		res, err := tools.ExecCommand(curStage.Config.Options[i].Label, store)
		if err != nil {
			return err
		}
		options[i] = res
	}
	res, err := tools.ExecCommand(curStage.Label, store)
	if err != nil {
		return err
	}

	prompt := &survey.Select{
		Message:  res,
		Options:  options,
		PageSize: curStage.Config.Size,
	}
	err = survey.AskOne(prompt, &selectedLabel)
	if err != nil {
		return err
	}
	selectOption := curStage.Config.Options[options.FindIndex(selectedLabel)]
	(*store)[curStage.Name] = selectOption.Value
	if selectOption.Next == END_ROUTE {
		return nil
	}
	return GoOriginalNext(selectOption.Next, stages, index, store)
}

// ProcessMultiSelect multi select handler
func ProcessMultiSelect(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	selectedLabels := []string{}
	options := make(StringList, len(curStage.Config.Options))
	for i := 0; i < len(curStage.Config.Options); i++ {
		res, err := tools.ExecCommand(curStage.Config.Options[i].Label, store)
		if err != nil {
			return err
		}
		options[i] = res
	}
	res, err := tools.ExecCommand(curStage.Label, store)
	if err != nil {
		return err
	}
	prompt := &survey.MultiSelect{
		Message:  res,
		Options:  options,
		PageSize: curStage.Config.Size,
	}
	err = survey.AskOne(prompt, &selectedLabels)
	if err != nil {
		return err
	}
	// set values to store
	if len(selectedLabels) > 0 {
		selectedValues := []string{}
		for _, selectedLabel := range selectedLabels {
			selectedOption := curStage.Config.Options[options.FindIndex(selectedLabel)]
			selectedValues = append(selectedValues, selectedOption.Value)
		}
		(*store)[curStage.Name] = strings.Join(selectedValues[:], ",")
	} else {
		// FIXME: has not choose
	}
	return GoOriginalNext("", stages, index, store)
}

// ProcessConfirm confirm question handler
func ProcessConfirm(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	answer := false
	res, err := tools.ExecCommand(curStage.Label, store)
	if err != nil {
		return err
	}
	defaultValue := true
	if curStage.Config.Default != nil {
		defaultValue = curStage.Config.Default.(bool)
	}
	prompt := &survey.Confirm{
		Message: res,
		Default: defaultValue,
	}
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return err
	}
	(*store)[curStage.Name] = strconv.FormatBool(answer)
	return GoOriginalNext("", stages, index, store)
}

// ProcessCommand command handler
// just for route command by excute result
func ProcessCommand(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	_, err := tools.ExecCommand(curStage.Config.Cmd, store)
	var next string
	if err != nil {
		// Goto failed route
		next = curStage.Config.Failed
	} else {
		// Goto success route
		next = curStage.Config.Success
	}
	return GoOriginalNext(next, stages, index, store)
}

// RunJob run every stage jobs
func RunJob(stages *[]config.Stage, index int, store *model.Store) error {
	curStage := (*stages)[index]
	switch curStage.Type {
	case "string":
		return ProcessString(stages, index, store)
	case "multiline":
		return ProcessMultiLine(stages, index, store)
	case "select":
		return ProcessSelect(stages, index, store)
	case "multi-select":
		return ProcessMultiSelect(stages, index, store)
	case "confirm":
		return ProcessConfirm(stages, index, store)
	case "command":
		// Validate data, just for route, with any UI
		return ProcessCommand(stages, index, store)
	default:
		return errors.New("Unknown type: " + curStage.Type)
	}
}
