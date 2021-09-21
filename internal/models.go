package kelp

import (
	"time"
)

// file that contains all GORM models for
// kelp.ml

// fields used in all models
type CommonModelFields struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// user model
type KelpUser struct {
	CommonModelFields
	Username    string `json:"username"`
	Password    string `json:"password"`
	ApiKey      string `json:"api_key"`
	AccessLevel int    `json:"access_level"` // 0/nil for standard access, 3 for admin (iirc)
}

// uploaded file model
type KelpFile struct {
	CommonModelFields
	UserId  int    `json:"user_id"` // corresponding user ID
	Type    string `json:"type"`
	Name    string `json:"name"`          // name used when saving to upload folder
	OrgName string `json:"original_name"` // name given when uploading
}

// uploaded paste model
type KelpPaste struct {
	CommonModelFields
	UserId  int    `json:"user_id"` // corresponding user ID
	Name    string `json:"name"`    // name used in path (server-side generated)
	Title   string `json:"title"`   // actual title of paste
	Content string `json:"content"`
}

// small model for storing invite
type KelpInvite struct {
	CommonModelFields
	InviteKey string `json:"invite_key"`
}
