package dtos

import "github.com/google/uuid"

type CreateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
}

type UpdateSupplierRequest struct {
	Name    string `json:"name"`
	Contact string `json:"contact,omitempty"`
	Address string `json:"address,omitempty"`
}

type SupplierResponse struct {
	ID      uint      `json:"id"`
	Uuid    uuid.UUID `json:"uuid"`
	Name    string    `json:"name"`
	Contact string    `json:"contact,omitempty"`
	Address string    `json:"address,omitempty"`
}
