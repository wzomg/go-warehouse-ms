package service

import (
	"time"

	"go-warehouse-ms/internal/model"
	"go-warehouse-ms/internal/repository"
)

type GoodsItem struct {
	model.Goods
	Cost float64 `json:"cost"`
}

type GoodsList struct {
	Items     []GoodsItem `json:"items"`
	TotalCost float64     `json:"totalCost"`
}

type GoodsService struct {
	goods     repository.GoodsStore
	bus       *EventBus
	caretaker *GoodsCaretaker
}

func NewGoodsService(goods repository.GoodsStore, bus *EventBus, caretaker *GoodsCaretaker) *GoodsService {
	return &GoodsService{goods: goods, bus: bus, caretaker: caretaker}
}

func (s *GoodsService) List() (GoodsList, error) {
	items, err := s.goods.FindAll()
	if err != nil {
		return GoodsList{}, err
	}
	result := make([]GoodsItem, 0, len(items))
	total := 0.0
	for _, item := range items {
		cost := item.GPrice * float64(item.GCnt)
		result = append(result, GoodsItem{Goods: item, Cost: cost})
		total += cost
	}
	return GoodsList{Items: result, TotalCost: total}, nil
}

func (s *GoodsService) Add(items []model.Goods) error {
	if len(items) == 0 {
		return ErrInvalidPayload
	}
	clones := make([]model.Goods, 0, len(items))
	for _, item := range items {
		clones = append(clones, *item.Clone())
	}
	if err := s.goods.CreateMany(clones); err != nil {
		return err
	}
	for _, item := range clones {
		s.bus.Publish(GoodsEvent{Type: EventGoodsAdded, Item: item, At: time.Now()})
	}
	return nil
}

func (s *GoodsService) Delete(gid int) error {
	item, err := s.goods.FindByID(gid)
	if err != nil {
		return err
	}
	s.caretaker.Push(GoodsMemento{Item: item, Action: ActionDelete, At: time.Now()})
	if err := s.goods.DeleteByID(gid); err != nil {
		return err
	}
	s.bus.Publish(GoodsEvent{Type: EventGoodsDeleted, Item: item, At: time.Now()})
	return nil
}

func (s *GoodsService) ReduceStock(gid int, det int) error {
	if det <= 0 {
		return ErrInvalidPayload
	}
	item, err := s.goods.FindByID(gid)
	if err != nil {
		return err
	}
	s.caretaker.Push(GoodsMemento{Item: item, Action: ActionReduce, Det: det, At: time.Now()})
	if err := s.goods.UpdateStock(gid, det); err != nil {
		return err
	}
	s.bus.Publish(GoodsEvent{Type: EventStockReduced, Item: item, Det: det, At: time.Now()})
	return nil
}

func (s *GoodsService) UndoLast() error {
	m, ok := s.caretaker.Pop()
	if !ok {
		return ErrInvalidPayload
	}
	item := m.Item
	if err := s.goods.Save(&item); err != nil {
		return err
	}
	eventType := EventStockUndo
	if m.Action == ActionDelete {
		eventType = EventDeleteRestored
	}
	s.bus.Publish(GoodsEvent{Type: eventType, Item: item, Det: m.Det, At: time.Now()})
	return nil
}
