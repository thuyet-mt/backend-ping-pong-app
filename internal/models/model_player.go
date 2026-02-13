package models

import "time"

type Player struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	BirthYear *int      `json:"birth_year,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	CCCD      *string   `json:"cccd,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type PlayerListResponse struct {
	FullName    string     `json:"full_name"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	//Phone     *string `json:"phone"`
	AvatarPath *string `json:"avatar_url"`
}
