package repository

import (
	"go-warehouse-ms/internal/model"

	"gorm.io/gorm"
)

type GoodsRepository struct {
	db *gorm.DB
}

func NewGoodsRepository(db *gorm.DB) *GoodsRepository {
	return &GoodsRepository{db: db}
}

func (r *GoodsRepository) FindAll() ([]model.Goods, error) {
	var items []model.Goods
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *GoodsRepository) FindByID(gid int) (model.Goods, error) {
	var item model.Goods
	if err := r.db.First(&item, gid).Error; err != nil {
		return model.Goods{}, err
	}
	return item, nil
}

func (r *GoodsRepository) CreateMany(items []model.Goods) error {
	return r.db.Create(&items).Error
}

func (r *GoodsRepository) DeleteByID(gid int) error {
	return r.db.Delete(&model.Goods{}, gid).Error
}

func (r *GoodsRepository) UpdateStock(gid int, det int) error {
	return r.db.Model(&model.Goods{}).
		Where("gid = ?", gid).
		UpdateColumn("gCnt", gorm.Expr("gCnt - ?", det)).Error
}

func (r *GoodsRepository) Save(item *model.Goods) error {
	return r.db.Save(item).Error
}
