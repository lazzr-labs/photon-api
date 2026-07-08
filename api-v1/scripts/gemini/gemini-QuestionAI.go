package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var questionInstructions = "Analyze this question and return a diverse list of 5 to 10 relevant search tags (search). Question: "

var questionToolDescription = "This tool analyzes a question and extracts relevant keywords to help build effective searches."
var searchDescription = "The search tags property provides a list of extracted keywords and key phrases that are relevant to the question and can be used to optimize search queries. Aim for a *minimum* of 5 search tags and a *maximum* of 10 search tags per question. Include *general keywords*, *specific topics*, and *related concepts* relevant to the question. For example, for the question 'What is photosynthesis?', good search tags might be 'photosynthesis', 'plants', 'energy', 'sunlight', 'biology', 'science & nature'. Aim for a diverse range of relevant terms. Include common abbreviations. For example, if the question includes ´Martin Luther King´a good tag would be ´MLK´."

var questionTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "questionTool",
		Description: questionToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"search": {
					Type:        genai.TypeString,
					Description: searchDescription,
				},
			},
			Required: []string{"search"},
		},
	}},
}

func QuestionAI(text string) (string, error) {
	ctx := context.Background()

	prompt := (questionInstructions + text)

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return "", errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{questionTool}
	model.ToolConfig = &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{
			Mode: genai.FunctionCallingAny,
		},
	}

	session := model.StartChat()
	response, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	functionCall, ok := response.Candidates[0].Content.Parts[0].(genai.FunctionCall)
	if !ok {
		return "", errors.New("no function call in response")
	}

	search, ok := functionCall.Args["search"].(string)
	if !ok {
		return "", errors.New("search value error")
	}

	return search, nil
}
