package dtos

import "github.com/google/uuid"

type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,passwordstrength"`
	Role     string `json:"role" validate:"required"`
	OutletID *uint  `json:"outlet_id,omitempty" validate:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"required"`
	Password *string `json:"password,omitempty" validate:"required,passwordstrength"`
	Role     *string `json:"role,omitempty" validate:"required,oneof=admin owner manager cashier"`
	OutletID *uint   `json:"outlet_id,omitempty" validate:"required"`
}

type RegisterAdminRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,passwordstrength"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
