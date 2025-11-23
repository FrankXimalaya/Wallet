package handlers

import (
	"net/http"
	"wallet-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TxHandler struct {
	txService *service.TxService
}

func NewTxHandler(txService *service.TxService) *TxHandler {
	return &TxHandler{txService: txService}
}

type TransferRequest struct {
	To     string `json:"to" binding:"required"`
	Amount string `json:"amount" binding:"required"`
}

// Transfer 发起转账
func (h *TxHandler) Transfer(c *gin.Context) {
	userID := c.GetUint("userID")

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txHash, err := h.txService.Transfer(userID, req.To, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer initiated",
		"tx_hash": txHash,
	})
}

// GetTransaction 查询交易详情
func (h *TxHandler) GetTransaction(c *gin.Context) {
	txHash := c.Param("hash")

	tx, err := h.txService.GetTransactionByHash(txHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": tx})
}

// GetTransactionList 查询地址交易记录
func (h *TxHandler) GetTransactionList(c *gin.Context) {
	address := c.Param("address")

	transactions, err := h.txService.GetTransactionsByAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// GetUserTransactions 获取用户的交易记录
func (h *TxHandler) GetUserTransactions(c *gin.Context) {
	userID := c.GetUint("userID")

	transactions, err := h.txService.GetUserTransactions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
