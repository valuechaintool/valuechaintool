package models

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Session struct {
	ID         string    `json:"id" sql:"primary key"`
	UserID     uuid.UUID `json:"user_id"`
	Creation   time.Time `json:"creation"`
	Expiration time.Time `json:"expiration"`
}

func NewSession(userID uuid.UUID) (*Session, error) {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	s := Session{
		ID:         base64.URLEncoding.EncodeToString(b),
		UserID:     userID,
		Creation:   time.Now(),
		Expiration: time.Now().Add(24 * time.Hour),
	}
	if err := session.Create(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func GetSession(id string) (*Session, error) {
	var s Session
	err := session.Where("id = ?", id).First(&s).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if time.Now().After(s.Expiration) {
		if err = session.Delete(s).Error; err != nil {
			return nil, err
		}
		return nil, errors.New("provided token is expired")
	}
	return &s, nil
}
