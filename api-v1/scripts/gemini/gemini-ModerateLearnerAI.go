package gemini

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var moderateLearnerInstructions = "Moderate this text for a kid-focused learning environment: "

var moderateLearnerToolDescription = "Moderates content for kid-focused learning."
var reportedLearnerDescription = "True if content contains profanity, self-harm, sexual material, PII (surname, address, username, phone, SSN, email), nonsensical text (random characters, repetitive sequences, incoherent phrases), romantic/flirtatious language (including compliments implying attraction, dating interest, or admiration of appearance), Comments about another student's physical appearance, attractiveness, or body (even if phrased politely), Harassment, bullying, threats, or coercive language, Self-harm or suicidal ideation, Requests for or sharing of personal information (PII: full names, addresses, usernames, phone numbers, emails, social handles, etc.), or Attempts to move the conversation off the platform (e.g., 'message me,' 'add me,' 'follow me'). Note that typos are not the same as nonsensical text and should be flagged as true.\n\nExamples:\n\"jafsjf0iuwnfowof9232\" = true (nonsensical text)\n\"Nice v8deo!\" = false (typographical error, not nonsensical text)\n\"You are beautiful\" = true (flirtatious language, comment about another student's physical appearance)\n\"I like your glasses\" = false (a nice compliment, not too flirtatious)\n\"You have nice eyes\" = false (a nice compliment, not too flirtatious)\n\"You have nice legs\" = true (flirtatious language, comment about another student's physical appearance)\n\"Hi hi hi hi hi hi hi hi hi hi hi hi hi hi hi\" = true (nonsensical text, repetitive sequence)\n\"Whats your insta?\" = true (request for sharing personal information, attempt to move the conversation off the platform)"

var moderateLearnerTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{{
		Name:        "moderateLearnerTool",
		Description: moderateLearnerToolDescription,
		Parameters: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"reported": {
					Type:        genai.TypeBoolean,
					Description: reportedLearnerDescription,
				},
			},
			Required: []string{"reported"},
		},
	}},
}

func ModerateLearnerAI(text string) (bool, error) {
	ctx := context.Background()

	prompt := (moderateLearnerInstructions + text)

	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		return false, errors.New("cannot create ai client")
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash")
	model.Tools = []*genai.Tool{moderateLearnerTool}
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
