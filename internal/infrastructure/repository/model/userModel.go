package model

type UserEntity struct {
	id    string `gorm:"primaryKey"`
	name  string
	email string
	token string `gorm:"not null"`
}
