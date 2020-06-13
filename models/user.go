package models

type User struct {
	Username  string `json:"username"`
	AuthToken string `json:"auth_token"`
}
