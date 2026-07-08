package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var moderateInstructions = "Moderate this text for a learning environment: "

var moderateToolDescription = "Moderates content for a learning environment."
var reportedDescription = "True if content contains profanity, sexual material, harassment, bullying, threats, coercive language, self-harm or suicidal ideation, personal information, nonsensical text, romantic or flirtatious language, or attempts to move the conversation off the platform."

var moderateTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "moderateTool",
		Description: moderateToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"reported": {
					Type:        genai.TypeBoolean,
					Description: reportedDescription,
				},
			},
			Required: []string{"reported"},
		},
	}},
}

func ModerateAI(text string) (bool, error) {
	ctx := context.Background()

	prompt := moderateInstructions + text

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return false, errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{moderateTool}
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
