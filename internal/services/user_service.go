package services

import (
	"context"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/musorok/server/internal/domain"
)

type UserService struct{ db *gorm.DB }

func NewUserService(db *gorm.DB) *UserService { return &UserService{db: db} }

func (s *UserService) Create(ctx context.Context, phone, email, name, password string, role domain.Role) (*domain.User, error) {
	var hash *string
	if password != "" {
		h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil { return nil, err }
		hs := string(h); hash = &hs
	}
	u := &domain.User{Phone: strings.TrimSpace(phone), Name: name, Role: role, PasswordHash: hash}
	if email != "" { u.Email = &email }
	if err := s.db.WithContext(ctx).Create(u).Error; err != nil { return nil, err }
	return u, nil
}

func (s *UserService) Authenticate(ctx context.Context, login, password string) (*domain.User, error) {
	var u domain.User
	q := s.db.WithContext(ctx).Where("phone = ? OR email = ?", login, login).First(&u)
	if q.Error != nil { return nil, q.Error }
	if u.PasswordHash == nil { return nil, gorm.ErrRecordNotFound }
	if bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(password)) != nil { return nil, gorm.ErrRecordNotFound }
	return &u, nil
}
