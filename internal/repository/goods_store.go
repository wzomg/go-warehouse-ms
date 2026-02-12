package repository

import "go-warehouse-ms/internal/model"

type GoodsStore interface {
	FindAll() ([]model.Goods, error)
	FindByID(gid int) (model.Goods, error)
	CreateMany(items []model.Goods) error
	DeleteByID(gid int) error
	UpdateStock(gid int, det int) error
	Save(item *model.Goods) error
}
