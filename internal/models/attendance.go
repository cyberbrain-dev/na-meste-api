package models

import "time"

type Attendance struct {
	ID        uint
	UserID    uint
	CollegeID uint
	Date      time.Time
}
