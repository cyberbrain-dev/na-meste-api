package models

type User struct {
	ID           uint
	Username     string
	Email        string
	PasswordHash string
	Role         string

	CollegeID uint
}
