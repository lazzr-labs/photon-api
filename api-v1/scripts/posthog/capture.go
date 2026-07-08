package posthog

import "github.com/posthog/posthog-go"

func Capture(userID string, eventName string) {
	client := posthog.New(posthogApiKey)
	defer client.Close()

	client.Enqueue(posthog.Capture{
		DistinctId: userID,
		Event:      eventName,
	})
}

func CaptureWithProperties(userID string, eventName string, properties posthog.Properties) {
	client := posthog.New(posthogApiKey)
	defer client.Close()

	client.Enqueue(posthog.Capture{
		DistinctId: userID,
		Event:      eventName,
		Properties: properties,
	})
}
