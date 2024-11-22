package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model        // adds ID, created_at etc.
	ID         uint   `json:"id" gorm:"primaryKey"`
	Email      string `json:"email" gorm:"unique"`
	Username   string `json:"username"`
	Password   string `json:"author"`
	Phone      string `json:"phone" gorm:"unique"`
	Bearer     string `json:"bearer" gorm:"null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u *User) HashPassword() string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

// CheckPassword сравнивает пароль с хэшированным
func (u *User) CheckPassword(password string) error {
	fmt.Println(u.Password)
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
