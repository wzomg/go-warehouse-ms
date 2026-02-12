package repository

import (
	"go-warehouse-ms/internal/model"

	"go.uber.org/zap"
)

type GoodsRepositoryProxy struct {
	target GoodsStore
	logger *zap.Logger
}

func NewGoodsRepositoryProxy(target GoodsStore, logger *zap.Logger) *GoodsRepositoryProxy {
	return &GoodsRepositoryProxy{target: target, logger: logger}
}

func (p *GoodsRepositoryProxy) FindAll() ([]model.Goods, error) {
	items, err := p.target.FindAll()
	p.logger.Info("goods.findAll", zap.Int("count", len(items)), zap.Error(err))
	return items, err
}

func (p *GoodsRepositoryProxy) FindByID(gid int) (model.Goods, error) {
	item, err := p.target.FindByID(gid)
	p.logger.Info("goods.findByID", zap.Int("gid", gid), zap.Error(err))
	return item, err
}

func (p *GoodsRepositoryProxy) CreateMany(items []model.Goods) error {
	err := p.target.CreateMany(items)
	p.logger.Info("goods.createMany", zap.Int("count", len(items)), zap.Error(err))
	return err
}

func (p *GoodsRepositoryProxy) DeleteByID(gid int) error {
	err := p.target.DeleteByID(gid)
	p.logger.Info("goods.deleteByID", zap.Int("gid", gid), zap.Error(err))
	return err
}

func (p *GoodsRepositoryProxy) UpdateStock(gid int, det int) error {
	err := p.target.UpdateStock(gid, det)
	p.logger.Info("goods.updateStock", zap.Int("gid", gid), zap.Int("det", det), zap.Error(err))
	return err
}

func (p *GoodsRepositoryProxy) Save(item *model.Goods) error {
	err := p.target.Save(item)
	p.logger.Info("goods.save", zap.Int("gid", item.GID), zap.Error(err))
	return err
}
