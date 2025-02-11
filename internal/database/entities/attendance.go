package entities

import "time"

// Represents an attendance record in th db
type Attendance struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint
	CollegeID uint
	Date      time.Time `gorm:"index"`
}
