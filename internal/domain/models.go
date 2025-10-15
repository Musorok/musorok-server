package domain

import (
	"time"
	"github.com/google/uuid"
)

type Role string
const (
	RoleUser Role = "USER"
	RoleCourier Role = "COURIER"
	RoleAdmin Role = "ADMIN"
)

type User struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Phone string `gorm:"uniqueIndex"`
	Email *string `gorm:"uniqueIndex"`
	Name string
	PasswordHash *string
	Role Role `gorm:"type:role_enum;default:'USER'"`
	IsDeleted bool `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Address struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;index"`
	Label *string
	Lat float64
	Lng float64
	City string
	Street string
	House string
	Entrance string
	Floor string
	Apartment string
	Intercom *string
	IsDefault bool
	PolygonID *uuid.UUID `gorm:"type:uuid;index"`
	PolygonName *string
	CreatedAt time.Time
}

type Polygon struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name string
	City string
	GeoJSON string `gorm:"type:jsonb"`
	IsActive bool `gorm:"default:true"`
	CreatedAt time.Time
}

type SubscriptionPlan string
const (
	PlanP7 SubscriptionPlan = "P7"
	PlanP15 SubscriptionPlan = "P15"
	PlanP30 SubscriptionPlan = "P30"
)

type SubscriptionStatus string
const (
	SubActive SubscriptionStatus = "ACTIVE"
	SubPaused SubscriptionStatus = "PAUSED"
	SubCanceled SubscriptionStatus = "CANCELED"
	SubExpired SubscriptionStatus = "EXPIRED"
)

type Subscription struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;index"`
	Plan SubscriptionPlan `gorm:"type:plan_enum"`
	TotalBags int
	RemainingBags int
	PriceKZT int
	Status SubscriptionStatus `gorm:"type:sub_status_enum;default:'ACTIVE'"`
	StartedAt time.Time
	ExpiresAt *time.Time
}

type OrderType string
const (
	OrderOneTime OrderType = "ONE_TIME"
	OrderSubscription OrderType = "SUBSCRIPTION"
)

type TimeOption string
const (
	ASAP TimeOption = "ASAP"
	SCHEDULED TimeOption = "SCHEDULED"
)

type OrderStatus string
const (
	StatusNew OrderStatus = "NEW"
	StatusPaid OrderStatus = "PAID"
	StatusAssigned OrderStatus = "ASSIGNED"
	StatusPickingUp OrderStatus = "PICKING_UP"
	StatusDone OrderStatus = "DONE"
	StatusCanceled OrderStatus = "CANCELED"
	StatusRefunded OrderStatus = "REFUNDED"
)

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;index"`
	AddressID uuid.UUID `gorm:"type:uuid;index"`
	PolygonID uuid.UUID `gorm:"type:uuid;index"`
	Type OrderType `gorm:"type:order_type_enum"`
	BagsCount int
	PriceKZT int
	Comment string
	TimeOption TimeOption `gorm:"type:time_option_enum"`
	ScheduledAt *time.Time
	CourierID *uuid.UUID `gorm:"type:uuid;index"`
	Status OrderStatus `gorm:"type:order_status_enum;default:'NEW'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentProvider string
const (
	ProviderPaynetworks PaymentProvider = "PAYNETWORKS"
)

type PaymentStatus string
const (
	PayInit PaymentStatus = "INIT"
	PayRequiresAction PaymentStatus = "REQUIRES_ACTION"
	PaySucceeded PaymentStatus = "SUCCEEDED"
	PayFailed PaymentStatus = "FAILED"
	PayCanceled PaymentStatus = "CANCELED"
)

type Payment struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;index"`
	OrderID *uuid.UUID `gorm:"type:uuid;index"`
	SubscriptionID *uuid.UUID `gorm:"type:uuid;index"`
	AmountKZT int
	Provider PaymentProvider `gorm:"type:payment_provider_enum"`
	Status PaymentStatus `gorm:"type:payment_status_enum;default:'INIT'"`
	ProviderIntentID string
	ProviderPayload string `gorm:"type:jsonb"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
