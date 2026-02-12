package api

import (
	"net/http"
	"strconv"

	"go-warehouse-ms/internal/model"
	"go-warehouse-ms/internal/service"

	"github.com/gin-gonic/gin"
)

type GoodsHandler struct {
	goods *service.GoodsService
}

func NewGoodsHandler(goods *service.GoodsService) *GoodsHandler {
	return &GoodsHandler{goods: goods}
}

type AddGoodsRequest struct {
	Items []model.Goods `json:"items"`
}

type UpdateStockRequest struct {
	GID int `json:"gid"`
	Det int `json:"det"`
}

func (h *GoodsHandler) List(c *gin.Context) {
	list, err := h.goods.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *GoodsHandler) Add(c *gin.Context) {
	var req AddGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	if err := h.goods.Add(req.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *GoodsHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	gid, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	if err := h.goods.Delete(gid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *GoodsHandler) UpdateStock(c *gin.Context) {
	var req UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return
	}
	if err := h.goods.ReduceStock(req.GID, req.Det); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *GoodsHandler) Undo(c *gin.Context) {
	if err := h.goods.UndoLast(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "撤销失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
