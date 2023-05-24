package interactive

import "github.com/AlecAivazis/survey/v2"

type Input struct {
	survey.Input

	validate survey.Validator
}

func NewInput(message string) *Input {
	return &Input{
		Input: survey.Input{
			Message: message,
		},
	}
}

func (question *Input) AskString() (string, error) {
	var answer string

	return answer, askPrompt(&question.Input, &answer, question.validate)
}

func (question *Input) SetValidate(validate survey.Validator) *Input {
	question.validate = validate

	return question
}
