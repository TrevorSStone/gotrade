package gotrade

import (
	"fmt"
	"log"
	"time"
)

const (
	marketURLPath   = "market"
	getQuoteURLPath = "quote"
)

type Quote struct {
	AdjNonAdjFlag    bool
	AnnualDividend   IntDollar
	Ask              IntDollar
	AskExchange      string
	AskSize          int64
	AskTime          time.Time
	Beta             float64
	Bid              IntDollar
	BidExchange      string
	BidSize          int64
	BidTime          time.Time
	ChgClose         IntDollar
	ChgClosePrcn     float64
	CompanyName      string
	DaysToExpiration int64
	DirLast          string
	Dividend         IntDollar
	Eps              IntDollar
	EstEarnings      IntDollar
	ExDivDate        string
	ExchgLastTrade   string
	Fsi              string
	High             IntDollar
	High52           IntDollar
	HighAsk          IntDollar
	HighBid          IntDollar
	LastTrade        IntDollar
	Low              IntDollar
	Low52            IntDollar
	LowAsk           IntDollar
	LowBid           IntDollar
	NumTrades        int64
	Open             IntDollar
	OpenInterest     int64
	OptionStyle      string
	OptionUnderlier  string
	PrevClose        IntDollar
	PrevDayVolume    int64
	PrimaryExchange  string
	SymbolDesc       string
	TodayClosed      IntDollar
	TotalVolume      int64
	Upc              int64
	Volume10Day      int64
	DateTime         time.Time
	Symbol           string
	ProductType      string
	Exchange         string
}

type quoteRaw struct {
	AdjNonAdjFlag    bool      `json:"adjNonAdjFlag"`
	AnnualDividend   IntDollar `json:"annualDividend"`
	Ask              IntDollar `json:"ask"`
	AskExchange      string    `json:"askExchange"`
	AskSize          int64     `json:"askSize"`
	AskTime          string    `json:"askTime"`
	Beta             float64   `json:"beta"`
	Bid              IntDollar `json:"bid"`
	BidExchange      string    `json:"bidExchange"`
	BidSize          int64     `json:"bidSize"`
	BidTime          string    `json:"bidTime"`
	ChgClose         IntDollar `json:"chgClose"`
	ChgClosePrcn     float64   `json:"chgClosePrcn"`
	CompanyName      string    `json:"companyName"`
	DaysToExpiration int64     `json:"daysToExpiration"`
	DirLast          string    `json:"dirLast"`
	Dividend         IntDollar `json:"dividend"`
	Eps              IntDollar `json:"eps"`
	EstEarnings      IntDollar `json:"estEarnings"`
	ExDivDate        string    `json:"exDivDate"`
	ExchgLastTrade   string    `json:"exchgLastTrade"`
	Fsi              string    `json:"fsi"`
	High             IntDollar `json:"high"`
	High52           IntDollar `json:"high52"`
	HighAsk          IntDollar `json:"highAsk"`
	HighBid          IntDollar `json:"highBid"`
	LastTrade        IntDollar `json:"lastTrade"`
	Low              IntDollar `json:"low"`
	Low52            IntDollar `json:"low52"`
	LowAsk           IntDollar `json:"lowAsk"`
	LowBid           IntDollar `json:"lowBid"`
	NumTrades        int64     `json:"numTrades"`
	Open             IntDollar `json:"open"`
	OpenInterest     int64     `json:"openInterest"`
	OptionStyle      string    `json:"optionStyle"`
	OptionUnderlier  string    `json:"optionUnderlier"`
	PrevClose        IntDollar `json:"prevClose"`
	PrevDayVolume    int64     `json:"prevDayVolume"`
	PrimaryExchange  string    `json:"primaryExchange"`
	SymbolDesc       string    `json:"symbolDesc"`
	TodayClosed      IntDollar `json:"todayClosed"`
	TotalVolume      int64     `json:"totalVolume"`
	Upc              int64     `json:"upc"`
	Volume10Day      int64     `json:"volume10Day"`
	DateTime         time.Time `json:"dateTime"`
	Product          struct {
		Symbol      string `json:"symbol"`
		ProductType string `json:"type"`
		Exchange    string `json:"exchange"`
	} `json:"product"`
}

func (client ETradeClient) GetQuote(symbol ...string) (quotes []Quote, raw string, err error) {
	url := fmt.Sprintf(client.url, marketURLPath, getQuoteURLPath)
	url = url + fmt.Sprintf("/%s%s", strURLParameter(symbol), jsonURL)
	log.Println(url)
	var response quoteRaw
	raw, err = client.requestAndUnmarshal(url, &response)
	if err != nil {
		return quotes, raw, err
	}
	//quotes, err = response.convert()
	return
}
