package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	Id        uuid.UUID    `json:"id"`
	Email     string       `json:"email"`
	Password  PasswordType `json:"password"`
	IsActive  bool         `json:"is_active"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type PasswordType struct {
	passwordText string
	PasswordHash []byte
}

func (p *PasswordType) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.passwordText = text
	p.PasswordHash = hash
	return nil
}

func (p *PasswordType) Check(text string) error {
	return bcrypt.CompareHashAndPassword(p.PasswordHash, []byte(text))
}
