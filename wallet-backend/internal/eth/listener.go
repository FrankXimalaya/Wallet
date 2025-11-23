package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"math/big"
	"time"
	"wallet-backend/internal/model"
	"wallet-backend/internal/repository"
)

type Listener struct {
	client     *Client
	walletRepo *repository.WalletRepository
	txRepo     *repository.TxRepository
}

func NewListener(client *Client, walletRepo *repository.WalletRepository, txRepo *repository.TxRepository) *Listener {
	return &Listener{
		client:     client,
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

// Start 开始监听区块
func (l *Listener) Start(startBlock uint64) {
	ctx := context.Background()
	currentBlock := startBlock

	ticker := time.NewTicker(12 * time.Second) // 以太坊平均出块时间约 12 秒
	defer ticker.Stop()

	log.Println("Block listener started")

	for range ticker.C {
		latestBlock, err := l.client.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Printf("Failed to get latest block: %v", err)
			continue
		}

		// 处理新区块
		for currentBlock <= latestBlock {
			if err := l.processBlock(ctx, currentBlock); err != nil {
				log.Printf("Failed to process block %d: %v", currentBlock, err)
			}
			currentBlock++
		}
	}
}

// processBlock 处理单个区块
func (l *Listener) processBlock(ctx context.Context, blockNumber uint64) error {
	block, err := l.client.GetClient().BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return err
	}

	// 获取我们监控的所有钱包地址
	wallets, err := l.walletRepo.FindAll()
	if err != nil {
		return err
	}

	// 创建地址映射
	watchedAddresses := make(map[string]uint)
	for _, wallet := range wallets {
		watchedAddresses[wallet.Address] = wallet.UserID
	}

	// 检查区块中的交易
	for _, tx := range block.Transactions() {
		if to := tx.To(); to != nil {
			// 检查是否是发送到我们监控的地址
			if userID, exists := watchedAddresses[to.Hex()]; exists {
				// 记录充值交易
				if err := l.recordDeposit(ctx, tx, userID, block.Time()); err != nil {
					log.Printf("Failed to record deposit: %v", err)
				}
			}
		}
	}

	return nil
}

// recordDeposit 记录充值交易
func (l *Listener) recordDeposit(ctx context.Context, tx *types.Transaction, userID uint, timestamp uint64) error {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		return err
	}

	txRecord := &model.Transaction{
		UserID:    userID,
		TxHash:    tx.Hash().Hex(),
		From:      from.Hex(),
		To:        tx.To().Hex(),
		Amount:    tx.Value().String(),
		Status:    "confirmed",
		TxType:    "deposit",
		Timestamp: int64(timestamp),
	}

	return l.txRepo.Create(txRecord)
}
