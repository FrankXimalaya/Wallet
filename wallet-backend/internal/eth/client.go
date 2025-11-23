package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	client *ethclient.Client
}

func NewClient(nodeURL string) (*Client, error) {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Close() {
	c.client.Close()
}

// GetBalance 获取地址余额
func (c *Client) GetBalance(ctx context.Context, address common.Address) (*big.Int, error) {
	return c.client.BalanceAt(ctx, address, nil)
}

// GetLatestBlockNumber 获取最新区块号
func (c *Client) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	return c.client.BlockNumber(ctx)
}

// GetClient 返回原始客户端（用于高级操作）
func (c *Client) GetClient() *ethclient.Client {
	return c.client
}
