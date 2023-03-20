package main

import (
	"coin-capturer/internal/capturer"
	"coin-capturer/internal/config"
	"coin-capturer/internal/erc20"
	"coin-capturer/internal/listener"
	"coin-capturer/internal/notice"
	"coin-capturer/internal/notice/dingtalk"
	"coin-capturer/internal/notice/terminal"
	"coin-capturer/internal/util"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config: %+v", cfg)

	l := listener.New(cfg)
	logs := make(chan types.Log)
	sub, err := l.Run(logs)
	if err != nil {
		log.Fatal(err)
	}

	var notices []notice.Notice
	notices = append(notices, terminal.New())
	if cfg.DingtalkToken != "" {
		notices = append(notices, dingtalk.New(cfg.DingtalkToken))
	}
	log.Printf("notices: %+v", cfg)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			event, err := erc20.ParseTransfer(vLog)
			if err != nil {
				log.Println(err)
			}

			if util.ToDecimal(event.Tokens, 18).LessThan(util.ToDecimal(cfg.LowerLimitValue, 0)) {
				continue
			}

			coins, err := capturer.GetBalance(cfg, event.From)
			if err != nil {
				log.Println(err)
			}
			if coins == nil {
				coins = []capturer.Coin{}
			}
			transactionList, err := capturer.GetTransactionList(cfg, event.From)
			if err != nil {
				log.Println(err)
				transactionList = []capturer.Transaction{}
			}

			transferEvent := &capturer.TransferEvent{
				TxHash:              vLog.TxHash,
				From:                event.From,
				FromBalance:         coins,
				FromTransactionList: transactionList,
				To:                  event.To,
				Tokens:              util.ToDecimal(event.Tokens, 18),
			}

			for _, n := range notices {
				n.OnTransfer(transferEvent)
			}
		}
	}
}
