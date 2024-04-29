package uuidv7

import "github.com/google/uuid"

func New() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}

func NewString() string {
	return New().String()
}
