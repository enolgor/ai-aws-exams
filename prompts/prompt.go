package prompts

import "embed"

//go:embed data/certificates/*
var certDefs embed.FS

//go:embed data/question-preamble.md
var definitionPreamble string

//go:embed data/question-postamble.md
var definitionPostamble string

//go:embed data/question-preamble.md
var questionPreamble string

//go:embed data/question-postamble.md
var questionPrompt string

func QuestionsPrompt(certification string) (string, error) {
	file, err := certDefs.ReadFile("data/certificates/" + certification + ".md")
	if err != nil {
		return "", err
	}
	return questionPreamble + "\n" + string(file) + "\n" + questionPrompt, nil
}

func DefinitionPrompt(certification string) (string, error) {
	file, err := certDefs.ReadFile("data/certificates/" + certification + ".md")
	if err != nil {
		return "", err
	}
	return definitionPreamble + "\n" + string(file) + "\n" + definitionPostamble, nil
}
