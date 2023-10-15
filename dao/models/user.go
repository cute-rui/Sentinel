package models

import "time"

type User struct {
	ID       int    `gorm:"column:id;AUTO_INCREMENT;NOT NULL;PRIMARY KEY;"`
	Username string `gorm:"uniqueIndex;column:username;NOT NULL;type:VARCHAR(255);"`
	Realname string `gorm:"column:realname;NOT NULL"`
	Hash     string `gorm:"column:hash;NOT NULL;"`
	Email    string `gorm:"uniqueIndex;column:email;NOT NULL;type:VARCHAR(255);"`
	//Leave blank for now
	QQ int64 `gorm:"column:qq;"`
	//Leave blank for now
	Feishu string `gorm:"column:feishu;"`
	//Leave blank for now, permission should be a JSON Object or something could parse.
	Permission string `gorm:"column:permission;"`
	//System wide admin
	IsAdmin   bool      `gorm:"column:is_admin;NOT NULL;"`
	CreatedAt time.Time `gorm:"autoCreateTime:true"`
}

func (User) TableName() string {
	return "user"
}
