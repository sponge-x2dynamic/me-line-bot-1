package model

import "time"

type Shop struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	Name        string  `gorm:"size:100;not null"`
	Lat         float64 `gorm:"type:decimal(10,8);not null"`
	Lng         float64 `gorm:"type:decimal(11,8);not null"`
	RadiusM     int     `gorm:"default:200"`
	LineGroupID string  `gorm:"size:64;not null"`
	CreatedAt   time.Time
}

// Role เช่น customer, staff, admin
type Role struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:50;not null;uniqueIndex"`
}

// User คือทุกคนที่ add bot
type User struct {
	ID          uint       `gorm:"primaryKey;autoIncrement"`
	LineUserID  string     `gorm:"size:64;not null;uniqueIndex"`
	Name        string     `gorm:"size:100;not null"`
	DisplayName string     `gorm:"size:100"`      // ชื่อที่แสดงใน LINE
	PictureURL  string     `gorm:"size:500"`      // รูป profile
	IsActive    bool       `gorm:"default:false"` // admin approve ก่อนใช้งานได้
	UserRoles   []UserRole `gorm:"foreignKey:UserID"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole เชื่อม User กับ Role (many-to-many)
type UserRole struct {
	UserID uint  `gorm:"primaryKey"`
	RoleID uint  `gorm:"primaryKey"`
	ShopID *uint // nil = customer ทั่วไป, มีค่า = staff ของร้านนั้น
	Role   Role  `gorm:"foreignKey:RoleID"`
	Shop   *Shop `gorm:"foreignKey:ShopID"`
}

type Attendance struct {
	ID              uint `gorm:"primaryKey;autoIncrement"`
	UserID          uint `gorm:"not null;index:idx_user_date"`
	User            User `gorm:"foreignKey:UserID"`
	ShopID          uint `gorm:"not null"`
	Shop            Shop `gorm:"foreignKey:ShopID"`
	CheckInTime     *time.Time
	CheckInLat      float64 `gorm:"type:decimal(10,8)"`
	CheckInLng      float64 `gorm:"type:decimal(11,8)"`
	CheckOutTime    *time.Time
	CheckOutLat     float64 `gorm:"type:decimal(10,8)"`
	CheckOutLng     float64 `gorm:"type:decimal(11,8)"`
	WorkDurationMin int
	WorkDate        string `gorm:"size:10;not null;index:idx_user_date"`
	CreatedAt       time.Time
}
