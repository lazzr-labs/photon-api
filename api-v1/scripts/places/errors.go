package places

import "errors"

var ErrClientNotConfigured = errors.New("places client is not configured")
var ErrPlaceDetailsFailed = errors.New("failed to get place details")
var ErrPlaceIDRequired = errors.New("place ID cannot be empty")
var ErrSearchRequired = errors.New("address cannot be empty")
var ErrPlacesNotFound = errors.New("no address found")
