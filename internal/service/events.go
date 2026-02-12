package service

import (
	"sync"
	"time"

	"go-warehouse-ms/internal/model"

	"go.uber.org/zap"
)

type EventType string

const (
	EventGoodsAdded     EventType = "goods_added"
	EventGoodsDeleted   EventType = "goods_deleted"
	EventStockReduced   EventType = "stock_reduced"
	EventStockUndo      EventType = "stock_undo"
	EventDeleteRestored EventType = "delete_restored"
)

type GoodsEvent struct {
	Type EventType
	Item model.Goods
	Det  int
	At   time.Time
}

type Observer interface {
	Handle(GoodsEvent)
}

type EventBus struct {
	mu        sync.RWMutex
	observers []Observer
}

func NewEventBus() *EventBus {
	return &EventBus{observers: []Observer{}}
}

func (b *EventBus) Subscribe(o Observer) {
	b.mu.Lock()
	b.observers = append(b.observers, o)
	b.mu.Unlock()
}

func (b *EventBus) Publish(event GoodsEvent) {
	b.mu.RLock()
	list := append([]Observer{}, b.observers...)
	b.mu.RUnlock()
	for _, o := range list {
		o.Handle(event)
	}
}

type AuditObserver struct {
	logger *zap.Logger
}

func NewAuditObserver(logger *zap.Logger) *AuditObserver {
	return &AuditObserver{logger: logger}
}

func (o *AuditObserver) Handle(event GoodsEvent) {
	o.logger.Info("goods.event",
		zap.String("type", string(event.Type)),
		zap.Int("gid", event.Item.GID),
		zap.String("name", event.Item.GName),
		zap.Int("det", event.Det),
		zap.Time("at", event.At),
	)
}
