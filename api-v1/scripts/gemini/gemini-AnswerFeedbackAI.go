package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var answerFeedbackInstructions = "Give short, kind, specific feedback on this student's answer text. Use simple language appropriate for elementary school students. Mention one thing the student did well and one clear way to make the answer better. Do not grade the answer. Do not use HTML markup or markdown formatting. Keep the feedback to 2-3 sentences. Answer text: "

var answerFeedbackToolDescription = "This tool gives short, kind, specific feedback on a student's answer text. The feedback uses simple language appropriate for elementary school students, mentions one strength, and gives one clear improvement suggestion."
var answerFeedbackDescription = "This property, feedback, contains short, kind, specific feedback for a student's answer. It should be 2-3 sentences, use simple language appropriate for elementary school students, mention one thing the student did well, and give one clear way to make the answer better. Do not grade the answer. Do not use HTML markup or markdown formatting."

var answerFeedbackTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "answerFeedbackTool",
		Description: answerFeedbackToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"feedback": {
					Type:        genai.TypeString,
					Description: answerFeedbackDescription,
				},
			},
			Required: []string{"feedback"},
		},
	}},
}

func AnswerFeedbackAI(text string) (string, error) {
	ctx := context.Background()

	prompt := answerFeedbackInstructions + text

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return "", errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{answerFeedbackTool}
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

	feedback, ok := functionCall.Args["feedback"].(string)
	if !ok {
		return "", errors.New("feedback value error")
	}

	return feedback, nil
}
