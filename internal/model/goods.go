package model

type Goods struct {
	GID    int     `gorm:"column:gid;primaryKey;autoIncrement" json:"gid"`
	GName  string  `gorm:"column:gName" json:"gName"`
	GShelf string  `gorm:"column:gShelf" json:"gShelf"`
	GCnt   int     `gorm:"column:gCnt" json:"gCnt"`
	GPrice float64 `gorm:"column:gPrice" json:"gPrice"`
}

func (g Goods) Clone() *Goods {
	copy := g
	return &copy
}

func (Goods) TableName() string {
	return "goods"
}
