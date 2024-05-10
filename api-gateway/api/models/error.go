package models

// Error ...
type Error struct {
	Message string `json:"message"`
}

const (
	InternalMessage = "Error happened during process"
)
