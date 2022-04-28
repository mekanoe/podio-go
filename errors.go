package podio

import "fmt"

var (
	// ErrTrustLevel is returned when the API key trust level is not high enough.
	ErrTrustLevel = fmt.Errorf("podio-go: API key is underneath the trust floor of this call, must be at least 2. https://help.podio.com/hc/en-us/community/posts/200515598-Client-trust-level")

	// ErrNotFound is returned when the requested resource is not found.
	ErrNotFound = fmt.Errorf("podio-go: not found")
)
