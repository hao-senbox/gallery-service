package dto

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
}

type RoleResponse struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role"`
}
