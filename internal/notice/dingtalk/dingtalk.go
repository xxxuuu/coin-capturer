package dingtalk

import (
	"bytes"
	"coin-capturer/internal/capturer"
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

const URL = "https://oapi.dingtalk.com/robot/send?access_token=%s"

//go:embed post.json
var POST_CONTENT string

//go:embed msg.md
var MSG_CONTENT string

//go:embed coins-msg.md
var COINS_CONTENT string

type DingTalk struct {
	token string
}

func New(token string) *DingTalk {
	return &DingTalk{token: token}
}

func (d *DingTalk) OnTransfer(event *capturer.TransferEvent) {
	url := fmt.Sprintf(URL, d.token)

	var coinsMsg string
	for _, c := range event.FromBalance {
		coinsMsg += fmt.Sprintf(COINS_CONTENT, c.Name, c.Balance, c.Address)
	}

	post := fmt.Sprintf(POST_CONTENT,
		event.Tokens,
		fmt.Sprintf(MSG_CONTENT,
			event.TxHash,
			event.TxHash,
			event.TxHash,
			event.From,
			event.From,
			event.To,
			event.To,
			event.Tokens,
			coinsMsg,
		),
	)

	fmt.Println(post)
	var jsonStr = []byte(post)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
