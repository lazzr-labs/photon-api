package auth_api

import (
	"context"
	"strconv"

	"github.com/danielgtaylor/huma/v2"

	"api-go/ent/user"
	"api-go/scripts/emails"
	"api-go/scripts/generator"
	"api-go/utils/cache"
	"api-go/utils/db"
)

type PasswordForgotPostInput struct {
	Body struct {
		Email string `json:"email"`
	}
}

type PasswordForgotPostOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func PasswordForgotPostAPI(ctx context.Context, input *PasswordForgotPostInput) (*PasswordForgotPostOutput, error) {
	userObj, err := db.EntDB.User.Query().
		Where(user.EmailEqualFold(input.Body.Email)).
		Only(ctx)
	if err != nil {
		return nil, huma.Error404NotFound("User not found.")
	}

	code := generator.RandomLetters(3) + strconv.Itoa(userObj.ID)
	cache.SetKey(code, strconv.Itoa(userObj.ID), 21600)
	emails.ForgotPasswordEmail(userObj.Email, userObj.Name, code)

	response := &PasswordForgotPostOutput{}
	response.Body.Message = "Code sent to email."
	return response, nil
}
