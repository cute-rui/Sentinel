package models

import "time"

type Verify struct {
	ID        int       `gorm:"column:id;AUTO_INCREMENT;NOT NULL;PRIMARY KEY;"`
	Email     string    `gorm:"column:email;NOT NULL;"`
	Verify    string    `gorm:"column:verify;NOT NULL"`
	CreatedAt time.Time `gorm:"autoCreateTime:true"`
}

func (Verify) TableName() string {
	return "verify"
}
