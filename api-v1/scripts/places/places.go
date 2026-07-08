package places

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

type AddressPlaceAutocomplete struct {
	PlaceID          string            `json:"place_id"`
	Description      string            `json:"description"`
	StructuredFormat *StructuredFormat `json:"structured_formatting,omitempty"`
}

type StructuredFormat struct {
	MainText      string `json:"main_text"`
	SecondaryText string `json:"secondary_text"`
}

func (client *MapsClientService) PlacesGet(search string) ([]*AddressPlaceAutocomplete, error) {
	if search == "" {
		return nil, errors.New("address cannot be empty")
	}

	ctx := context.Background()

	autocompleteRequest := &maps.PlaceAutocompleteRequest{
		Input: search,
		Types: maps.AutocompletePlaceTypeGeocode,
	}

	autocomplete, err := client.client.PlaceAutocomplete(ctx, autocompleteRequest)
	if err != nil {
		return nil, errors.New("failed to search address")
	}

	if len(autocomplete.Predictions) == 0 {
		return nil, errors.New("no address found")
	}

	var places []*AddressPlaceAutocomplete

	for _, place := range autocomplete.Predictions {
		address := &AddressPlaceAutocomplete{
			PlaceID:     place.PlaceID,
			Description: place.Description,
		}
		if place.StructuredFormatting.MainText != "" {
			address.StructuredFormat = &StructuredFormat{
				MainText:      place.StructuredFormatting.MainText,
				SecondaryText: place.StructuredFormatting.SecondaryText,
			}
		}
		places = append(places, address)
	}

	return places, nil
}
