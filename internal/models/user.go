package models

type User struct {
	UUID     string `json:"UUID"`
	Mail     string `json:"mail"`
	Password string `json:"password,omitempty"`
	Type     string `json:"type,omitempty"`
}
