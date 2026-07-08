package places_api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"

	"api-go/scripts/places"
	"api-go/utils/auth"
)

type PlaceGetInput struct {
	auth.AuthParam
	PlaceID string `path:"place_id"`
}

func (input *PlaceGetInput) Resolve(ctx huma.Context) []error {
	var err error
	input.User, err = auth.AuthUser(input.Auth)
	if err != nil {
		return []error{huma.Error401Unauthorized("Unable to authenticate.")}
	}
	return nil
}

type PlaceGetOutput struct {
	Body struct {
		Object places.AddressPlace `json:"object"`
	}
}

func PlaceGetAPI(ctx context.Context, input *PlaceGetInput) (*PlaceGetOutput, error) {
	placeObj, err := places.MapsClient.PlaceGet(input.PlaceID)
	if err != nil {
		return nil, huma.Error404NotFound("Place not found.")
	}

	response := &PlaceGetOutput{}
	object := &places.AddressPlace{}
	copier.Copy(&object, &placeObj)
	response.Body.Object = *object
	return response, nil
}
