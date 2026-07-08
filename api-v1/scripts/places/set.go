package places

import (
	"strings"

	"googlemaps.github.io/maps"
)

var MapsClient *MapsClientService

type MapsClientService struct {
	client *maps.Client
}

func SetClient(apiKey string) error {
	if strings.TrimSpace(apiKey) == "" {
		MapsClient = nil
		return nil
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return err
	}

	MapsClient = &MapsClientService{
		client: client,
	}

	return nil
}
