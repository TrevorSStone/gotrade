package gotrade

import (
	"errors"
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
	TodayClose       IntDollar
	TotalVolume      int64
	Upc              int64
	Volume10Day      int64
	DateTime         time.Time
	Symbol           string
	ProductType      string
	Exchange         string
}

type quoteContainer interface {
	Convert() ([]Quote, error)
}

type multiQuoteContainer struct {
	QuoteResponse struct {
		QuoteData []quoteRaw `json:"quoteData"`
	} `json:"quoteResponse"`
}

type singleQuoteContainer struct {
	QuoteResponse struct {
		QuoteData quoteRaw `json:"quoteData"`
	} `json:"quoteResponse"`
}

type quoteRaw struct {
	All struct {
		AdjNonAdjFlag    bool        `json:"adjNonAdjFlag"`
		AnnualDividend   interface{} `json:"annualDividend"`
		Ask              interface{} `json:"ask"`
		AskExchange      string      `json:"askExchange"`
		AskSize          int64       `json:"askSize"`
		AskTime          string      `json:"askTime"`
		Beta             float64     `json:"beta"`
		Bid              interface{} `json:"bid"`
		BidExchange      string      `json:"bidExchange"`
		BidSize          int64       `json:"bidSize"`
		BidTime          string      `json:"bidTime"`
		ChgClose         interface{} `json:"chgClose"`
		ChgClosePrcn     float64     `json:"chgClosePrcn"`
		CompanyName      string      `json:"companyName"`
		DaysToExpiration int64       `json:"daysToExpiration"`
		DirLast          string      `json:"dirLast"`
		Dividend         interface{} `json:"dividend"`
		Eps              interface{} `json:"eps"`
		EstEarnings      interface{} `json:"estEarnings"`
		ExDivDate        string      `json:"exDivDate"`
		ExchgLastTrade   string      `json:"exchgLastTrade"`
		Fsi              string      `json:"fsi"`
		High             interface{} `json:"high"`
		High52           interface{} `json:"high52"`
		HighAsk          interface{} `json:"highAsk"`
		HighBid          interface{} `json:"highBid"`
		LastTrade        interface{} `json:"lastTrade"`
		Low              interface{} `json:"low"`
		Low52            interface{} `json:"low52"`
		LowAsk           interface{} `json:"lowAsk"`
		LowBid           interface{} `json:"lowBid"`
		NumTrades        int64       `json:"numTrades"`
		Open             interface{} `json:"open"`
		OpenInterest     int64       `json:"openInterest"`
		OptionStyle      string      `json:"optionStyle"`
		OptionUnderlier  string      `json:"optionUnderlier"`
		PrevClose        interface{} `json:"prevClose"`
		PrevDayVolume    int64       `json:"prevDayVolume"`
		PrimaryExchange  string      `json:"primaryExchange"`
		SymbolDesc       string      `json:"symbolDesc"`
		TodayClose       interface{} `json:"todayClose"`
		TotalVolume      int64       `json:"totalVolume"`
		Upc              int64       `json:"upc"`
		Volume10Day      int64       `json:"volume10Day"`
	} `json:"all"`
	DateTime string `json:"dateTime"`
	Product  struct {
		Symbol      string `json:"symbol"`
		ProductType string `json:"type"`
		Exchange    string `json:"exchange"`
	} `json:"product"`
}

func (client ETradeClient) GetQuote(symbol ...string) (quotes []Quote, raw string, err error) {
	numSymbols := len(symbol)
	if numSymbols < 1 {
		return quotes, raw, errors.New("Need at least one symbol")
	}
	url := fmt.Sprintf(client.url, marketURLPath, getQuoteURLPath)
	url = url + fmt.Sprintf("/%s%s", strURLParameter(symbol), jsonURL)
	log.Println(url)

	var container quoteContainer

	if numSymbols == 1 {
		response := singleQuoteContainer{}
		raw, err = client.requestAndUnmarshal(url, &response)
		if err != nil {
			return quotes, raw, err
		}
		container = response
	} else {
		response := multiQuoteContainer{}
		raw, err = client.requestAndUnmarshal(url, &response)
		if err != nil {
			return quotes, raw, err
		}
		container = response
	}

	quotes, err = container.Convert()

	return
}

func (container multiQuoteContainer) Convert() (quotes []Quote, err error) {
	for _, v := range container.QuoteResponse.QuoteData {
		quote, err := v.Convert()
		if err != nil {
			return quotes, err
		}
		quotes = append(quotes, quote)
	}
	return
}

func (container singleQuoteContainer) Convert() (quotes []Quote, err error) {
	quote, err := container.QuoteResponse.QuoteData.Convert()
	if err != nil {
		return
	}
	quotes = append(quotes, quote)
	return
}

const timeform = "15:04:05 MST 01-02-2006"

//Beta doesn't exist in sandbox
func (raw quoteRaw) Convert() (quote Quote, err error) {
	quote.AdjNonAdjFlag = raw.All.AdjNonAdjFlag
	quote.AnnualDividend, err = convertToIntDollar(raw.All.AnnualDividend)
	if err != nil {
		return
	}
	quote.Ask, err = convertToIntDollar(raw.All.Ask)
	if err != nil {
		return
	}
	quote.AskExchange = raw.All.AskExchange
	quote.AskSize = raw.All.AskSize

	quote.AskTime, err = time.Parse(timeform, raw.All.AskTime)
	if err != nil {
		return
	}
	quote.Bid, err = convertToIntDollar(raw.All.Bid)
	if err != nil {
		return
	}
	quote.BidExchange = raw.All.BidExchange
	quote.BidSize = raw.All.BidSize
	quote.BidTime, err = time.Parse(timeform, raw.All.BidTime)
	if err != nil {
		return
	}
	quote.ChgClose, err = convertToIntDollar(raw.All.ChgClose)
	if err != nil {
		return
	}
	quote.ChgClosePrcn = raw.All.ChgClosePrcn
	quote.CompanyName = raw.All.CompanyName
	quote.DaysToExpiration = raw.All.DaysToExpiration
	quote.DirLast = raw.All.DirLast
	quote.Dividend, err = convertToIntDollar(raw.All.Dividend)
	if err != nil {
		return
	}
	quote.Eps, err = convertToIntDollar(raw.All.Eps)
	if err != nil {
		return
	}
	quote.EstEarnings, err = convertToIntDollar(raw.All.EstEarnings)
	if err != nil {
		return
	}
	quote.ExDivDate = raw.All.ExDivDate
	quote.ExchgLastTrade = raw.All.ExchgLastTrade
	quote.Fsi = raw.All.Fsi

	quote.High, err = convertToIntDollar(raw.All.High)
	if err != nil {
		return
	}
	quote.High52, err = convertToIntDollar(raw.All.High52)
	if err != nil {
		return
	}
	quote.HighAsk, err = convertToIntDollar(raw.All.HighAsk)
	if err != nil {
		return
	}
	quote.HighBid, err = convertToIntDollar(raw.All.HighBid)
	if err != nil {
		return
	}
	quote.LastTrade, err = convertToIntDollar(raw.All.LastTrade)
	if err != nil {
		return
	}
	quote.Low, err = convertToIntDollar(raw.All.Low)
	if err != nil {
		return
	}
	quote.Low52, err = convertToIntDollar(raw.All.Low52)
	if err != nil {
		return
	}
	quote.LowAsk, err = convertToIntDollar(raw.All.LowAsk)
	if err != nil {
		return
	}
	quote.LowBid, err = convertToIntDollar(raw.All.LowBid)
	if err != nil {
		return
	}
	quote.NumTrades = raw.All.NumTrades
	quote.Open, err = convertToIntDollar(raw.All.Open)
	if err != nil {
		return
	}
	quote.OpenInterest = raw.All.OpenInterest

	quote.OptionStyle = raw.All.OptionStyle
	quote.OptionUnderlier = raw.All.OptionUnderlier
	quote.PrevClose, err = convertToIntDollar(raw.All.PrevClose)
	if err != nil {
		return
	}
	quote.PrevDayVolume = raw.All.PrevDayVolume
	quote.PrimaryExchange = raw.All.PrimaryExchange
	quote.SymbolDesc = raw.All.SymbolDesc
	quote.TodayClose, err = convertToIntDollar(raw.All.TodayClose)
	if err != nil {
		return
	}
	quote.TotalVolume = raw.All.TotalVolume
	quote.Upc = raw.All.Upc
	quote.Volume10Day = raw.All.Volume10Day
	quote.DateTime, err = time.Parse(timeform, raw.DateTime)
	quote.Symbol = raw.Product.Symbol
	quote.ProductType = raw.Product.ProductType
	quote.Exchange = raw.Product.Exchange

	return
}
