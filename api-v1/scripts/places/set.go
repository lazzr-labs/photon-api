package places

import "googlemaps.github.io/maps"

var MapsClient *MapsClientService

type MapsClientService struct {
	client *maps.Client
}

func SetClient(apiKey string) {
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		panic("failed to create Google Maps client")
	}

	MapsClient = &MapsClientService{
		client: client,
	}
}
