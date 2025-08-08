package dto

import "time"

type UserEntityResponse struct {
	ID           string   `json:"id"`
	Username     string   `json:"username"`
	Nickname     string   `json:"nickname"`
	Fullname     string   `json:"fullname"`
	Phone        string   `json:"phone"`
	Email        string   `json:"email"`
	Dob          string   `json:"dob"`
	QRLogin      string   `json:"qr_login"`
	Avatar       string   `json:"avatar"`
	AvatarURL    string   `json:"avatar_url"`
	IsBlocked    bool     `json:"is_blocked"`
	BlockedAt    string   `json:"blocked_at"`
	Organization []string `json:"organizations"`
	CreatedAt    string   `json:"created_at"`

	Roles   *[]RoleResponse `json:"roles"`
	Devices *[]string       `json:"devices"`

	OrganizationAdmin *OrganizationAdmin `json:"organization_admin"`
}

type RoleResponse struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role"`
}

type OrganizationAdmin struct {
	ID               string    `json:"id"`
	OrganizationName string    `json:"organization_name"`
	Avatar           string    `json:"avatar"`
	AvatarURL        string    `json:"avatar_url"`
	Address          string    `json:"address"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
