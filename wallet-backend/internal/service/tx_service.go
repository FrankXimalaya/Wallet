package service

import (
	"context"
	"math/big"
	"time"
	"wallet-backend/internal/eth"
	"wallet-backend/internal/model"
	"wallet-backend/internal/repository"

	"github.com/ethereum/go-ethereum/common"
)

type TxService struct {
	txRepo     *repository.TxRepository
	walletRepo *repository.WalletRepository
	ethClient  *eth.Client
}

func NewTxService(txRepo *repository.TxRepository, walletRepo *repository.WalletRepository, ethClient *eth.Client) *TxService {
	return &TxService{
		txRepo:     txRepo,
		walletRepo: walletRepo,
		ethClient:  ethClient,
	}
}

// Transfer 发起转账
func (s *TxService) Transfer(userID uint, to string, amountStr string) (string, error) {
	ctx := context.Background()

	// 获取用户主钱包
	wallet, err := s.walletRepo.FindMainWalletByUserID(userID)
	if err != nil {
		return "", err
	}

	// 解析金额
	amount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok {
		return "", err
	}

	// 获取 nonce
	fromAddress := common.HexToAddress(wallet.Address)
	nonce, err := s.ethClient.GetNonce(ctx, fromAddress)
	if err != nil {
		return "", err
	}

	// 获取 gas price
	gasPrice, err := s.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	// 获取 chain ID
	chainID, err := s.ethClient.GetChainID(ctx)
	if err != nil {
		return "", err
	}

	// 签名交易
	toAddress := common.HexToAddress(to)
	signedTx, err := eth.SignTransaction(
		wallet.PrivateKey,
		toAddress,
		amount,
		nonce,
		21000, // 标准 ETH 转账的 gas limit
		gasPrice,
		chainID,
	)
	if err != nil {
		return "", err
	}

	// 发送交易
	if err := s.ethClient.SendTransaction(ctx, signedTx); err != nil {
		return "", err
	}

	// 记录交易
	tx := &model.Transaction{
		UserID:    userID,
		TxHash:    signedTx.Hash().Hex(),
		From:      wallet.Address,
		To:        to,
		Amount:    amountStr,
		GasPrice:  gasPrice.String(),
		Status:    "pending",
		TxType:    "transfer",
		Timestamp: time.Now().Unix(),
	}

	if err := s.txRepo.Create(tx); err != nil {
		return "", err
	}

	return signedTx.Hash().Hex(), nil
}

// GetTransactionByHash 根据交易哈希查询交易
func (s *TxService) GetTransactionByHash(hash string) (*model.Transaction, error) {
	return s.txRepo.FindByHash(hash)
}

// GetTransactionsByAddress 根据地址查询交易记录
func (s *TxService) GetTransactionsByAddress(address string) ([]*model.Transaction, error) {
	return s.txRepo.FindByAddress(address)
}

// GetUserTransactions 获取用户的交易记录
func (s *TxService) GetUserTransactions(userID uint) ([]*model.Transaction, error) {
	return s.txRepo.FindByUserID(userID)
}
