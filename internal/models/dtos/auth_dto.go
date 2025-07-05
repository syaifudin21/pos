package dtos

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	OutletID *uint  `json:"outlet_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       uint      `json:"id"`
	Uuid     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	OutletID *uint     `json:"outlet_id,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}