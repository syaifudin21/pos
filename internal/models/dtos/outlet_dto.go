package dtos

import "github.com/google/uuid"

type OutletCreateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

type OutletUpdateRequest struct {
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

type OutletResponse struct {
	ID      uint      `json:"id"`
	Uuid    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Type    string    `json:"type"`
}
