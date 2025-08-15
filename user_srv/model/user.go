package model

import (
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primary_kay"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	IsDeleted bool      `gorm:"column:is_deleted;default:false"`
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"type:varchar(11);not null;index:idx_mobile;unique"`
	Password string     `gorm:"type:varchar(255);not null"`
	Nickname string     `gorm:"type:varchar(255)"`
	Birthday *time.Time `gorm:"type:datetime;default:null"`
	Gender   string     `gorm:"type:varchar(10);default:'male'';column:gender"`
	Role     int        `gorm:"type:int;default:1"`
}
