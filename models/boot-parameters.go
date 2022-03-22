package models

// User model
type BootParameters struct {
	Image ImageObject `json:"image"`
}

type ImageObject struct {
	Version string `json:"version"`
	Path    string `json:"path"`
}
