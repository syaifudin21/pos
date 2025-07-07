package dtos

import "github.com/google/uuid"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required"`
	OutletID *uint  `json:"outlet_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       uint      `json:"id"`
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	OutletID *uint     `json:"outlet_id,omitempty"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"required"`
	Email    *string `json:"email,omitempty" validate:"required,email"`
	Password *string `json:"password,omitempty" validate:"required,passwordstrength"`
	Role     *string `json:"role,omitempty" validate:"required,oneof=admin owner manager cashier"`
	OutletID *uint   `json:"outlet_id,omitempty"`
}

type RegisterAdminRequest struct {
	Name        string `json:"name" validate:"required"`
	Password    string `json:"password" validate:"required,passwordstrength"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,passwordstrength"`
}

type UpdateEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required,len=6"`
}

type SendOTPRequest struct {
	Email string `json:"email" validate:"required,email"`
}
