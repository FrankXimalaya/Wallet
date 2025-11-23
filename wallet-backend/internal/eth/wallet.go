package eth

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// GenerateWallet 生成新钱包
func GenerateWallet() (address string, privateKey string, err error) {
	// 生成私钥
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// 获取地址
	addr := crypto.PubkeyToAddress(key.PublicKey)
	
	// 导出私钥
	privateKeyBytes := crypto.FromECDSA(key)
	
	return addr.Hex(), common.Bytes2Hex(privateKeyBytes), nil
}

// SignTransaction 签名交易
func SignTransaction(privateKeyHex string, to common.Address, amount *big.Int, nonce uint64, gasLimit uint64, gasPrice *big.Int, chainID *big.Int) (*types.Transaction, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	// 创建交易
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	// 签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// SendTransaction 发送交易
func (c *Client) SendTransaction(ctx context.Context, signedTx *types.Transaction) error {
	return c.client.SendTransaction(ctx, signedTx)
}

// GetNonce 获取地址的 nonce
func (c *Client) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	return c.client.PendingNonceAt(ctx, address)
}

// SuggestGasPrice 建议 gas 价格
func (c *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return c.client.SuggestGasPrice(ctx)
}

// GetChainID 获取链 ID
func (c *Client) GetChainID(ctx context.Context) (*big.Int, error) {
	return c.client.ChainID(ctx)
}

// GetPrivateKeyFromHex 从十六进制字符串获取私钥
func GetPrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hexKey)
}

// GetAddressFromPrivateKey 从私钥获取地址
func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) common.Address {
	return crypto.PubkeyToAddress(privateKey.PublicKey)
}
