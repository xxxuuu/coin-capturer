package listener

import (
	config2 "coin-capturer/internal/config"
	"coin-capturer/internal/erc20"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const USDT = "0x55d398326f99059fF775485246999027B3197955"

type Listener struct {
	config *config2.Config
}

func New(config *config2.Config) *Listener {
	return &Listener{
		config: config,
	}
}

func (l *Listener) Run(logs chan types.Log) (ethereum.Subscription, error) {
	var monitoredWallet []common.Hash
	for _, addr := range l.config.MonitoredWallet {
		monitoredWallet = append(monitoredWallet, common.HexToHash(addr))
	}

	// 连接到BSC节点
	client, err := ethclient.Dial(l.config.NodeAddress)
	if err != nil {
		return nil, err
	}

	// 监听指定钱包地址的转账事件
	address := common.HexToAddress(USDT)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{address},
		Topics: [][]common.Hash{
			{
				common.HexToHash(erc20.TRANSFER_EVENT),
			},
			{},
			monitoredWallet,
		},
	}
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, err
	}

	return sub, err
}
