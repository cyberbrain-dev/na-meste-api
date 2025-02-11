package entities

// Represents a user record in db
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:100; not null"`
	Email        string `gorm:"size:200; not null; unique"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"check:role IN ('teacher', 'scanner', 'student')"`

	CollegeID uint `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Attendances []Attendance
}
