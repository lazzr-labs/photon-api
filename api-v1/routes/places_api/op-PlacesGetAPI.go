package places_api

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/jinzhu/copier"

	"api-go/scripts/places"
	"api-go/utils/auth"
)

type PlacesGetInput struct {
	auth.AuthParam
	Search string `query:"search"`
}

func (input *PlacesGetInput) Resolve(ctx huma.Context) []error {
	var err error
	input.User, err = auth.AuthUser(input.Auth)
	if err != nil {
		return []error{huma.Error401Unauthorized("Unable to authenticate.")}
	}
	return nil
}

type PlacesGetOutput struct {
	Body struct {
		List []places.AddressPlaceAutocomplete `json:"list"`
	}
}

func PlacesGetAPI(ctx context.Context, input *PlacesGetInput) (*PlacesGetOutput, error) {
	placesList, err := places.MapsClient.PlacesGet(input.Search)
	if err != nil {
		return nil, huma.Error404NotFound("Places not found.")
	}

	response := &PlacesGetOutput{}
	list := &[]places.AddressPlaceAutocomplete{}
	copier.Copy(&list, &placesList)
	response.Body.List = *list
	return response, nil
}
