package domain

import (
	"time"
)

type Preference struct {
	ID                string
	ExternalReference string
	InitPoint         string
	SandboxInitPoint  string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewPreference(preference, externalReference, initPoint, sandboxInitPoint string) (*Preference, error) {
	return &Preference{
		ID:                preference,
		ExternalReference: externalReference,
		InitPoint:         initPoint,
		SandboxInitPoint:  sandboxInitPoint,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}
