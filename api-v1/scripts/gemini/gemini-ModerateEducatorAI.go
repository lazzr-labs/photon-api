package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var moderateEducatorInstructions = "Moderate this text for an educator-focused learning environment: "

var moderateEducatorToolDescription = "Moderates content for educator-focused learning."
var reportedEducatorDescription = "True if content contains profanity, self-harm, sexual material, Harassment, bullying, threats, or coercive language, Self-harm or suicidal ideation."

var moderateEducatorTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "moderateEducatorTool",
		Description: moderateEducatorToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"reported": {
					Type:        genai.TypeBoolean,
					Description: reportedEducatorDescription,
				},
			},
			Required: []string{"reported"},
		},
	}},
}

func ModerateEducatorAI(text string) (bool, error) {
	ctx := context.Background()

	prompt := (moderateEducatorInstructions + text)

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return false, errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{moderateEducatorTool}
	model.ToolConfig = &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{
			Mode: genai.FunctionCallingAny,
		},
	}

	session := model.StartChat()
	response, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return false, err
	}

	functionCall, ok := response.Candidates[0].Content.Parts[0].(genai.FunctionCall)
	if !ok {
		return false, errors.New("no function call in response")
	}

	reported, ok := functionCall.Args["reported"].(bool)
	if !ok {
		return false, errors.New("reported value error")
	}

	return reported, nil
}
