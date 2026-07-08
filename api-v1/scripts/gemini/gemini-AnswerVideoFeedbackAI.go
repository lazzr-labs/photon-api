package gemini

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var answerVideoFeedbackInstructions = "Analyze this student's answer video. Describe what is happening in the video and give short, kind, specific feedback to help the student improve their video answer. Use simple language appropriate for elementary school students. Mention one thing the student did well and one clear way to make the video answer better. Do not grade the answer. Do not use HTML markup or markdown formatting."

var answerVideoFeedbackToolDescription = "This tool analyzes a student's answer video, describes what is happening, and gives short, kind feedback."
var answerVideoDescription = "This property, description, explains what is happening in the student's video in simple language."
var answerVideoFeedbackDescription = "This property, feedback, contains short, kind, specific feedback for the student's video answer. It should mention one thing the student did well and one clear way to make the video answer better."

var answerVideoFeedbackTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "answerVideoFeedbackTool",
		Description: answerVideoFeedbackToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"description": {
					Type:        genai.TypeString,
					Description: answerVideoDescription,
				},
				"feedback": {
					Type:        genai.TypeString,
					Description: answerVideoFeedbackDescription,
				},
			},
			Required: []string{"description", "feedback"},
		},
	}},
}

func AnswerVideoFeedbackAI(video []byte) (string, string, error) {
	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return "", "", errors.New("cannot create ai client")
	}
	defer client.Close()

	mimeType := http.DetectContentType(video)
	if mimeType == "application/octet-stream" {
		mimeType = "video/mp4"
	}

	file, err := client.UploadFile(ctx, "", bytes.NewReader(video), &genai.UploadFileOptions{MIMEType: mimeType})
	if err != nil {
		return "", "", err
	}
	defer client.DeleteFile(ctx, file.Name)

	for file.State == genai.FileStateProcessing {
		time.Sleep(5 * time.Second)
		file, err = client.GetFile(ctx, file.Name)
		if err != nil {
			return "", "", err
		}
	}
	if file.State != genai.FileStateActive {
		return "", "", errors.New("file not active")
	}

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{answerVideoFeedbackTool}
	model.ToolConfig = &genai.ToolConfig{
		FunctionCallingConfig: &genai.FunctionCallingConfig{
			Mode: genai.FunctionCallingAny,
		},
	}

	session := model.StartChat()
	response, err := session.SendMessage(ctx, genai.Text(answerVideoFeedbackInstructions), genai.FileData{URI: file.URI})
	if err != nil {
		return "", "", err
	}

	if len(response.Candidates) == 0 || response.Candidates[0].Content == nil || len(response.Candidates[0].Content.Parts) == 0 {
		return "", "", errors.New("no ai response")
	}

	functionCall, ok := response.Candidates[0].Content.Parts[0].(genai.FunctionCall)
	if !ok {
		return "", "", errors.New("no function call in response")
	}

	description, ok := functionCall.Args["description"].(string)
	if !ok {
		return "", "", errors.New("description value error")
	}

	feedback, ok := functionCall.Args["feedback"].(string)
	if !ok {
		return "", "", errors.New("feedback value error")
	}

	return description, feedback, nil
}
