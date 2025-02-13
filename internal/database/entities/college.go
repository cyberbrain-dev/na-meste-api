// Contains all the representations of database tables
package entities

// Represents a college record in db
type College struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:200; not null; unique"`
	Users       []User
	Attendances []Attendance
}
