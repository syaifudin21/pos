package models

type UserPayment struct {
	UserID          uint          `gorm:"not null"`
	User            User          `gorm:"foreignKey:UserID"`
	PaymentMethodID uint          `gorm:"not null"`
	PaymentMethod   PaymentMethod `gorm:"foreignKey:PaymentMethodID"`
	IsActive        bool          `gorm:"default:true"`
}
