package model

import (
	"errors"

	argon2 "github.com/mdouchement/simple-argon2"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
    Name 		string `gorm:"not null"`
	Role_ID		uint   `gorm:"not null"`
    Email		string `gorm:"not null"`
    Phone_Num   string `gorm:"not null"`
    Password  	string `gorm:"not null"`
}

func (u *User) TableName() string {
	return "users"
}

func MigrateUser(db *gorm.DB) {
	db.AutoMigrate(&User{})
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) UpdateUser(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *User) DeleteUser(db *gorm.DB) error {
	return db.Delete(u).Error
}

func (u *User) GetUser(db *gorm.DB) error {
	return db.First(u).Error
}

func (u *User) HashPassword() error {

	if u.Password == "" {
        return errors.New("password cannot be empty")
    }
	
	hash, _ := argon2.GenerateFromPasswordString(u.Password, argon2.Default)
	u.Password = hash
	
	return nil
}

func CompareHashAndPassword(hashedPassword, password string) bool {
	if err := argon2.CompareHashAndPasswordString(hashedPassword, password); err != nil {
		return false
	}
	return true
}
