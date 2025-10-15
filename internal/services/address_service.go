package services

import (
	"context"
	"gorm.io/gorm"
	"github.com/musorok/server/internal/domain"
)

type AddressService struct{ db *gorm.DB }

func NewAddressService(db *gorm.DB) *AddressService { return &AddressService{db: db} }

func (s *AddressService) Create(ctx context.Context, a *domain.Address, polygon *domain.Polygon) error {
	if polygon != nil {
		a.PolygonID = &polygon.ID
		a.PolygonName = &polygon.Name
	}
	return s.db.WithContext(ctx).Create(a).Error
}

func (s *AddressService) List(ctx context.Context, userID string, out *[]domain.Address) error {
	return s.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(out).Error
}
