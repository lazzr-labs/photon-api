package posthog

import "github.com/posthog/posthog-go"

func UserCaptureWithProperty(userID string, propertyName string, propertyValue interface{}) {
	client := posthog.New(posthogApiKey)
	defer client.Close()

	client.Enqueue(posthog.Capture{
		DistinctId: userID,
		Event:      "$identify",
		Properties: map[string]interface{}{
			"$set": map[string]interface{}{
				propertyName: propertyValue,
			},
		},
	})
}

func UserCaptureWithProperties(userID string, properties map[string]interface{}) {
	client := posthog.New(posthogApiKey)
	defer client.Close()

	client.Enqueue(posthog.Capture{
		DistinctId: userID,
		Event:      "$identify",
		Properties: map[string]interface{}{
			"$set": properties,
		},
	})
}
