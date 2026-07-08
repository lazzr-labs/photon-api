package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var answerOutlineInstructions = "Create short, clear bullet points that a student can use as a small teleprompter for a video answer. Use simple language appropriate for elementary school students. Base the outline only on the student's answer text. Do not grade the answer. Do not add unrelated ideas. Do not use HTML markup. Keep the outline concise with 3-5 talking points. Answer text: "

var answerOutlineToolDescription = "This tool creates short, clear bullet points that a student can use as teleprompter talking points for a video answer."
var answerOutlineDescription = "This property, outline, contains short, clear bullet points based only on the student's answer text. It should use simple language appropriate for elementary school students, include 3-5 concise talking points, avoid grading, avoid unrelated ideas, and avoid HTML markup."

var answerOutlineTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "answerOutlineTool",
		Description: answerOutlineToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"outline": {
					Type:        genai.TypeString,
					Description: answerOutlineDescription,
				},
			},
			Required: []string{"outline"},
		},
	}},
}

func AnswerOutlineAI(text string) (string, error) {
	ctx := context.Background()

	prompt := answerOutlineInstructions + text

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return "", errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{answerOutlineTool}
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

	outline, ok := functionCall.Args["outline"].(string)
	if !ok {
		return "", errors.New("outline value error")
	}

	return outline, nil
}
