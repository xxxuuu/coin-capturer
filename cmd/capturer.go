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
	config, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config: %+v", config)

	l := listener.New(config)
	logs := make(chan types.Log)
	sub, err := l.Run(logs)
	if err != nil {
		log.Fatal(err)
	}

	var notices []notice.Notice
	notices = append(notices, terminal.New())
	if config.DingtalkToken != "" {
		notices = append(notices, dingtalk.New(config.DingtalkToken))
	}
	log.Printf("notices: %+v", config)

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			event, err := erc20.ParseTransfer(vLog)
			if err != nil {
				log.Println(err)
			}
			coins, err := capturer.GetBalance(config, event.From)
			if err != nil {
				log.Println(err)
			}
			if coins == nil {
				coins = []capturer.Coin{}
			}

			transferEvent := &capturer.TransferEvent{
				TxHash:      vLog.TxHash,
				From:        event.From,
				FromBalance: coins,
				To:          event.To,
				Tokens:      util.ToDecimal(event.Tokens, 18),
			}

			for _, n := range notices {
				n.OnTransfer(transferEvent)
			}
		}
	}
}
