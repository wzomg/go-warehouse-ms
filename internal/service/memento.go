package service

import (
	"sync"
	"time"

	"go-warehouse-ms/internal/model"
)

type GoodsAction string

const (
	ActionReduce GoodsAction = "reduce"
	ActionDelete GoodsAction = "delete"
)

type GoodsMemento struct {
	Item   model.Goods
	Action GoodsAction
	Det    int
	At     time.Time
}

type GoodsCaretaker struct {
	mu    sync.Mutex
	stack []GoodsMemento
	limit int
}

func NewGoodsCaretaker(limit int) *GoodsCaretaker {
	return &GoodsCaretaker{limit: limit, stack: []GoodsMemento{}}
}

func (c *GoodsCaretaker) Push(m GoodsMemento) {
	c.mu.Lock()
	if len(c.stack) >= c.limit {
		c.stack = c.stack[1:]
	}
	c.stack = append(c.stack, m)
	c.mu.Unlock()
}

func (c *GoodsCaretaker) Pop() (GoodsMemento, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.stack) == 0 {
		return GoodsMemento{}, false
	}
	idx := len(c.stack) - 1
	m := c.stack[idx]
	c.stack = c.stack[:idx]
	return m, true
}
