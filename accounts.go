package gotrade

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Account struct {
	Description      string
	ID               int
	MarginLevel      string
	NetAccountValue  IntDollar
	RegistrationType string
}

const (
	accountsUrlPath       = "accounts"
	accountListUrlPath    = "accountlist"
	accountBalanceUrlPath = "accountbalance"
)

type accountListResponse struct {
	AccountListResponse struct {
		Accounts []accountRaw `json:"response"`
	} `json:"json.accountListResponse"`
}

type accountRaw struct {
	Description string `json:"accountDesc"`
	AccountID   int    `json:"accountId"`
	MarginLevel string `json:"marginLevel"`
	//NetAccountValue can be either a string of a number...seriously
	NetAccountValue  interface{} `json:"netAccountValue"`
	RegistrationType string      `json:"registrationType"`
}

func (client ETradeClient) ListAccounts() (accounts []Account, raw string, err error) {
	url := fmt.Sprintf(client.url, accountsUrlPath, accountListUrlPath) + jsonURL
	var response accountListResponse
	raw, err = client.requestAndUnmarshal(url, &response)
	if err != nil {
		return accounts, raw, err
	}
	accounts, err = response.convert()
	return
}

func (rawAccountListResponse accountListResponse) convert() (accounts []Account, err error) {
	rawList := rawAccountListResponse.AccountListResponse.Accounts
	accounts = make([]Account, len(rawList), len(rawList))
	for i, r := range rawList {
		var netValue IntDollar
		switch v := r.NetAccountValue.(type) {
		default:
			return accounts, errors.New("Unexpected type in response for Value")
		case float32:
			netValue = IntDollar(math.Floor((float64(v) * 100) + .5))
		case string:
			v = strings.Replace(v, ".", "", -1)
			value, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return accounts, err
			}
			netValue = IntDollar(value)
		case float64:
			netValue = IntDollar(math.Floor((v * 100) + .5))
		}
		accounts[i].ID = r.AccountID
		accounts[i].Description = r.Description
		accounts[i].MarginLevel = r.MarginLevel
		accounts[i].NetAccountValue = netValue
		accounts[i].RegistrationType = r.RegistrationType
	}
	return
}

func (client ETradeClient) AccountBalance(accountID int) (raw string, err error) {
	url := fmt.Sprintf(client.url, accountsUrlPath, accountBalanceUrlPath)
	url = url + fmt.Sprintf("/%d%s", accountID, jsonURL)
	fmt.Println(url)
	raw, err = client.requestAndUnmarshal(url, nil)
	fmt.Println(raw)
	if err != nil {
		return raw, err
	}
	return
}
