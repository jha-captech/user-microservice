package main

type User struct {
	ID        uint   `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Role      string `json:"role,omitempty"`
	UserID    uint   `json:"user_id,omitempty"`
}
