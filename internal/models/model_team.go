package models

// Team represents a team in a season

type Team struct {
	ID        string  `json:"id"`
	SeasonID  string  `json:"season_id"`
	Name      string  `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}

// TeamMember represents a player membership in a team
type TeamMember struct {
	TeamID       string `json:"team_id"`
	PlayerID     string `json:"player_id"`
	JoinedRound  int    `json:"joined_round"`
	LeftRound    *int   `json:"left_round,omitempty"`
	TransferType string `json:"transfer_type"` // INITIAL, TRANSFER
}
