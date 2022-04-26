package models

import "time"

type Post struct {
	UUID         string `json:"uuid"`
	CreatorUUID  string `json:"creator_uuid"`
	Headline     string `json:"headline"`
	Summary      string `json:"summary"`
	PreviewImage []byte `json:"preview_img,omitempty" db:"preview_img"`
	Content      string `json:"content,omitempty"`
	IsProtected  bool   `json:"is_private,omitempty" db:"is_protected"`

	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}
