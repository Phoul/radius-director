package resources

import _ "embed"

// ExampleConfig contains the canonical example RADIUS Director configuration.
//
//go:embed example.yaml
var ExampleConfig []byte
