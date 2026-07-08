package places

import (
	"context"
	"errors"

	"googlemaps.github.io/maps"
)

type AddressPlace struct {
	PlaceID          string            `json:"place_id"`
	Description      string            `json:"description"`
	FormattedAddress string            `json:"formatted_address,omitempty"`
	Components       map[string]string `json:"components,omitempty"`
	Latitude         float64           `json:"latitude,omitempty"`
	Longitude        float64           `json:"longitude,omitempty"`
	Country          string            `json:"country,omitempty"`
	Alpha            string            `json:"alpha,omitempty"`
}

func (client *MapsClientService) PlaceGet(placeID string) (*AddressPlace, error) {
	if placeID == "" {
		return nil, errors.New("place ID cannot be empty")
	}

	ctx := context.Background()

	detailsRequest := &maps.PlaceDetailsRequest{
		PlaceID: placeID,
		Fields: []maps.PlaceDetailsFieldMask{
			maps.PlaceDetailsFieldMaskFormattedAddress,
			maps.PlaceDetailsFieldMaskAddressComponent,
			maps.PlaceDetailsFieldMaskPlaceID,
			maps.PlaceDetailsFieldMaskGeometry,
		},
	}

	details, err := client.client.PlaceDetails(ctx, detailsRequest)
	if err != nil {
		return nil, errors.New("failed to get place details")
	}

	components := make(map[string]string)
	var countryLongName, countryShortName string
	for _, component := range details.AddressComponents {
		for _, addressType := range component.Types {
			if addressType == "street_number" || addressType == "route" || addressType == "locality" || addressType == "administrative_area_level_1" || addressType == "country" || addressType == "postal_code" {
				components[addressType] = component.LongName
			}
			if addressType == "country" {
				countryLongName = component.LongName
				countryShortName = component.ShortName
			}
		}
	}

	var latitude, longitude float64
	latitude = details.Geometry.Location.Lat
	longitude = details.Geometry.Location.Lng

	return &AddressPlace{
		PlaceID:          details.PlaceID,
		Description:      details.FormattedAddress,
		FormattedAddress: details.FormattedAddress,
		Components:       components,
		Latitude:         latitude,
		Longitude:        longitude,
		Country:          countryLongName,
		Alpha:            countryShortName,
	}, nil
}
