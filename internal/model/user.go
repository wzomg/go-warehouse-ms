package model

type User struct {
	UserID  string `gorm:"column:userid;primaryKey" json:"userId"`
	UserPwd string `gorm:"column:userpwd" json:"userPwd"`
}

func (User) TableName() string {
	return "user"
}
