package capturer

import (
	"coin-capturer/internal/config"
	"coin-capturer/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"io/ioutil"
	"net/http"
)

const BALANCE_API = "https://www.oklink.com/api/v5/explorer/address/address-balance-fills?chainShortName=bsc&address=%s&protocolType=token_20"

type BalanceResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		Page           string `json:"page"`
		Limit          string `json:"limit"`
		TotalPage      string `json:"totalPage"`
		ChainFullName  string `json:"chainFullName"`
		ChainShortName string `json:"chainShortName"`
		TokenList      []struct {
			Token                string `json:"token"`
			HoldingAmount        string `json:"holdingAmount"`
			TotalTokenValue      string `json:"totalTokenValue"`
			Change24H            string `json:"change24h"`
			PriceUsd             string `json:"priceUsd"`
			ValueUsd             string `json:"valueUsd"`
			TokenContractAddress string `json:"tokenContractAddress"`
		} `json:"tokenList"`
	} `json:"data"`
}

func GetBalance(config *config.Config, addr common.Address) ([]Coin, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(BALANCE_API, addr.String()), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Ok-Access-Key", config.OklinkToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	balanceResp := &BalanceResponse{}
	if err := json.Unmarshal(body, balanceResp); err != nil {
		return nil, err
	}

	if balanceResp.Code != "0" {
		return nil, errors.New(fmt.Sprintf("oklink balance api response failed: %+v", balanceResp))
	}
	if len(balanceResp.Data) == 0 {
		return nil, errors.New(fmt.Sprintf("oklink balance api response data is empty: %+v", balanceResp))
	}

	var coins []Coin
	for _, token := range balanceResp.Data[0].TokenList {
		coins = append(coins, Coin{
			Name:    token.Token,
			Address: common.HexToAddress(token.TokenContractAddress),
			Balance: util.ToDecimal(token.HoldingAmount, 18),
		})
	}
	return coins, nil
}
