package ci

import (
	"fmt"
	"slices"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

var ciTemplates = map[string]func(name string) error{
	"none": func(name string) error {
		return nil
	},
}

func SelectCITemplate() (string, error) {
	candidates := make([]string, 0, len(ciTemplates))
	for k := range ciTemplates {
		candidates = append(candidates, k)
	}
	slices.Sort(candidates)

	var template string
	if err := survey.AskOne(&survey.Select{
		Message: "Select the prefer CI tool:",
		Options: candidates,
	}, &template, survey.WithValidator(survey.Required)); err != nil {
		return "", err
	}

	template = strings.TrimSpace(template)

	return template, nil
}

func CreateCITemplate(template string) error {
	if fn, ok := ciTemplates[template]; ok {
		if err := fn("merge"); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid ci template: %s", template)
	}

	return nil
}
