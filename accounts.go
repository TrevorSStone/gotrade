package gotrade

import "fmt"

type Account struct {
	Description      string
	ID               int
	MarginLevel      string
	NetAccountValue  IntDollar
	RegistrationType string
}

type AccountBalance struct {
	AccountID         int
	AccountType       string
	OptionLevel       string
	BasicBalance      AccountBasicBalance
	MarginBalance     AccountMarginBalance
	DayTradingBalance AccountDayTradingBalance
	CashBalance       AccountCashBalance
}

type AccountBasicBalance struct {
	CashAvailableForWithdrawal  IntDollar
	CashCall                    IntDollar
	FundsWithheldFromWithdrawal IntDollar
	NetAccountValue             IntDollar
	NetCash                     IntDollar
	SweepDepositAmount          IntDollar
	TotalLongValue              IntDollar
	TotalSecuritiesMktValue     IntDollar
	TotalCash                   IntDollar
}
type AccountMarginBalance struct {
	FedCall                           IntDollar
	MarginBalance                     IntDollar
	MarginBalanceWithdrawal           IntDollar
	MarginEquity                      IntDollar
	MarginEquityPercent               float32
	MarginableSecurities              IntDollar
	MaxAvailableForWithdrawal         IntDollar
	MinEquityCall                     IntDollar
	NonMarginableSecuritiesAndOptions IntDollar
	TotalShortValue                   IntDollar
	ShortReserve                      IntDollar
}

type AccountDayTradingBalance struct {
}

type AccountCashBalance struct {
}

const (
	accountsURLPath       = "accounts"
	accountListURLPath    = "accountlist"
	accountBalanceURLPath = "accountbalance"
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
	url := fmt.Sprintf(client.url, accountsURLPath, accountListURLPath) + jsonURL
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
		netValue, err := convertToIntDollar(r.NetAccountValue)
		if err != nil {
			return accounts, err
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
	url := fmt.Sprintf(client.url, accountsURLPath, accountBalanceURLPath)
	url = url + fmt.Sprintf("/%d%s", accountID, jsonURL)
	fmt.Println(url)
	raw, err = client.requestAndUnmarshal(url, nil)
	fmt.Println(raw)
	if err != nil {
		return raw, err
	}
	return
}
