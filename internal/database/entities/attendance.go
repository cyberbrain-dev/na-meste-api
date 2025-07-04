package entities

import "time"

// Represents an attendance record in th db
type Attendance struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"<-:create;not null;constraint:OnDelete:CASCADE;"`
	CollegeID uint      `gorm:"<-:create;not null;constraint:OnDelete:CASCADE;"`
	Date      time.Time `gorm:"not null;index"`
}
