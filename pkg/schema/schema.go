package schema

import (
	"time"
)

type InsType string

const (
	InsTypeEtf      InsType = "Etf"
	InsTypeBond             = "Bond"
	InsTypeStock            = "Stock"
	InsTypeCurrency         = "Currency"
)

type Section string

const (
	BondRub  Section = "Bond.RUB"
	BondUsd          = "Bond.USD"
	StockRub         = "Stock.RUB"
	StockUsd         = "Stock.USD"
	CashRub          = "Cash.RUB"
	CashUsd          = "Cash.USD"
)

type Instrument struct {
	Figi      string `json:"figi"`
	Ticker    string `json:"ticker"`
	Name      string `json:"name"`
	Currency  string `json:"currency"`
	FaceValue int    `json:"faceValue"`

	Type    InsType
	Section Section
}

type PortfolioResponse struct {
	Payload struct {
		Positions []struct {
			AveragePositionPrice      CValue
			AveragePositionPriceNoNkd CValue
			Balance                   float64 `json:"balance"`
			Blocked                   float64 `json:"blocked"`
			ExpectedYield             CValue
			Figi                      string  `json:"figi"`
			InstrumentType            string  `json:"instrumentType"`
			Isin                      string  `json:"isin"`
			Lots                      float64 `json:"lots"`
			Ticker                    string  `json:"ticker"`
		} `json:"positions"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type Trade struct {
	Date     string  `json:"date"`
	Price    float64 `json:"price"`
	Quantity uint    `json:"quantity"`
	TradeID  string  `json:"tradeId"`
}

/* Payment:
 * buy: negative
 * dividends, coupons: positive
 * taxes: negative
 * service commission: negative
 */
type Operation struct {
	Commission     CValue
	Currency       string  `json:"currency"`
	Date           string  `json:"date"`
	Figi           string  `json:"figi"`
	ID             string  `json:"id"`
	InstrumentType string  `json:"instrumentType"`
	IsMarginCall   bool    `json:"isMarginCall"`
	OperationType  string  `json:"operationType"`
	Payment        float64 `json:"payment"`
	Price          float64 `json:"price"`
	Quantity       uint    `json:"quantity"`
	Status         string  `json:"status"`
	Trades         []Trade `json:"trades"`
	// Added fields below
	DateParsed time.Time `json:"-"`
	Ticker     string    `json:"-"`
}

type OperationsResponse struct {
	Payload struct {
		Operations []Operation `json:"operations"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type OrderbookResponse struct {
	Payload struct {
		Asks []struct {
			Price    float64 `json:"price"`
			Quantity uint    `json:"quantity"`
		} `json:"asks"`
		Bids []struct {
			Price    float64 `json:"price"`
			Quantity uint    `json:"quantity"`
		} `json:"bids"`
		ClosePrice        float64 `json:"closePrice"`
		Depth             uint    `json:"depth"`
		Figi              string  `json:"figi"`
		LastPrice         float64 `json:"lastPrice"`
		LimitDown         float64 `json:"limitDown"`
		LimitUp           float64 `json:"limitUp"`
		MinPriceIncrement float64 `json:"minPriceIncrement"`
		TradeStatus       string  `json:"tradeStatus"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type SearchByFigiResponse struct {
	Payload struct {
		Currency          string  `json:"currency"`
		Figi              string  `json:"figi"`
		Isin              string  `json:"isin"`
		Lot               int     `json:"lot"`
		MinPriceIncrement float64 `json:"minPriceIncrement"`
		Name              string  `json:"name"`
		Ticker            string  `json:"ticker"`
		FaceValue         float64 `json:"faceValue"`
		Type              string  `json:"type"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type SearchByTickerResponse struct {
	Payload struct {
		Instruments []struct {
			Currency          string  `json:"currency"`
			Figi              string  `json:"figi"`
			Isin              string  `json:"isin"`
			Lot               int     `json:"lot"`
			MinPriceIncrement float64 `json:"minPriceIncrement"`
			Name              string  `json:"name"`
			Ticker            string  `json:"ticker"`
			FaceValue         float64 `json:"faceValue"`
			Type              string  `json:"type"`
		} `json:"instruments"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type Candle struct {
	C        float64 `json:"c"`
	Figi     string  `json:"figi"`
	H        float64 `json:"h"`
	Interval string  `json:"interval"`
	L        float64 `json:"l"`
	O        float64 `json:"o"`
	Time     string  `json:"time"`
	V        float64 `json:"v"`
}

type CandlesResponse struct {
	Payload struct {
		Candles  []Candle `json:"candles"`
		Figi     string   `json:"figi"`
		Interval string   `json:"interval"`
	} `json:"payload"`
	Status     string `json:"status"`
	TrackingID string `json:"trackingId"`
}

type AccountsResponse struct {
	TrackingID string `json:"trackingId"`
	Status     string `json:"status"`
	Payload    struct {
		Accounts []struct {
			BrokerAccountType string `json:"brokerAccountType"`
			BrokerAccountID   string `json:"brokerAccountId"`
		} `json:"accounts"`
	} `json:"payload"`
}

type AutoGenerated struct {
	Empty struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Payload struct {
				Type string `json:"type"`
			} `json:"payload"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
		} `json:"properties"`
	} `json:"Empty"`
	Error struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Type       string `json:"type"`
				Properties struct {
					Message struct {
						Type string `json:"type"`
					} `json:"message"`
					Code struct {
						Type string `json:"type"`
					} `json:"code"`
				} `json:"properties"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"Error"`
	PortfolioResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"PortfolioResponse"`
	Portfolio struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Positions struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"positions"`
		} `json:"properties"`
	} `json:"Portfolio"`
	PortfolioCurrenciesResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"PortfolioCurrenciesResponse"`
	Currencies struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Currencies struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"currencies"`
		} `json:"properties"`
	} `json:"Currencies"`
	CurrencyPosition struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Currency struct {
				Ref string `json:"$ref"`
			} `json:"currency"`
			Balance struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"balance"`
			Blocked struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"blocked"`
		} `json:"properties"`
	} `json:"CurrencyPosition"`
	PortfolioPosition struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Ticker struct {
				Type string `json:"type"`
			} `json:"ticker"`
			Isin struct {
				Type string `json:"type"`
			} `json:"isin"`
			InstrumentType struct {
				Ref string `json:"$ref"`
			} `json:"instrumentType"`
			Balance struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"balance"`
			Blocked struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"blocked"`
			Lots struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"lots"`
			ExpectedYield struct {
				Ref string `json:"$ref"`
			} `json:"expectedYield"`
			AveragePositionPrice struct {
				Ref string `json:"$ref"`
			} `json:"averagePositionPrice"`
			AveragePositionPriceNoNkd struct {
				Ref string `json:"$ref"`
			} `json:"averagePositionPriceNoNkd"`
		} `json:"properties"`
	} `json:"PortfolioPosition"`
	MoneyAmount struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Currency struct {
				Ref string `json:"$ref"`
			} `json:"currency"`
			Value struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"value"`
		} `json:"properties"`
	} `json:"MoneyAmount"`
	OrderbookResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"OrderbookResponse"`
	Orderbook struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Depth struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"depth"`
			Bids struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"bids"`
			Asks struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"asks"`
			TradeStatus struct {
				Ref string `json:"$ref"`
			} `json:"tradeStatus"`
			MinPriceIncrement struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"minPriceIncrement"`
			LastPrice struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"lastPrice"`
			ClosePrice struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"closePrice"`
			LimitUp struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"limitUp"`
			LimitDown struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"limitDown"`
		} `json:"properties"`
	} `json:"Orderbook"`
	OrderResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Price struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"price"`
			Quantity struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"quantity"`
		} `json:"properties"`
	} `json:"OrderResponse"`
	CandlesResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"CandlesResponse"`
	Candles struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Interval struct {
				Ref string `json:"$ref"`
			} `json:"interval"`
			Candles struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"candles"`
		} `json:"properties"`
	} `json:"Candles"`
	Candle struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Interval struct {
				Ref string `json:"$ref"`
			} `json:"interval"`
			O struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"o"`
			C struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"c"`
			H struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"h"`
			L struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"l"`
			V struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"v"`
			Time struct {
				Type        string    `json:"type"`
				Format      string    `json:"format"`
				Description string    `json:"description"`
				Example     time.Time `json:"example"`
			} `json:"time"`
		} `json:"properties"`
	} `json:"Candle"`
	CandleResolution struct {
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Enum        []string `json:"enum"`
	} `json:"CandleResolution"`
	OperationsResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"OperationsResponse"`
	Operations struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Operations struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"operations"`
		} `json:"properties"`
	} `json:"Operations"`
	OperationTrade struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TradeID struct {
				Type string `json:"type"`
			} `json:"tradeId"`
			Date struct {
				Type        string    `json:"type"`
				Format      string    `json:"format"`
				Description string    `json:"description"`
				Example     time.Time `json:"example"`
			} `json:"date"`
			Price struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"price"`
			Quantity struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"quantity"`
		} `json:"properties"`
	} `json:"OperationTrade"`
	Operation struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			ID struct {
				Type string `json:"type"`
			} `json:"id"`
			Status struct {
				Ref string `json:"$ref"`
			} `json:"status"`
			Trades struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"trades"`
			Commission struct {
				Ref string `json:"$ref"`
			} `json:"commission"`
			Currency struct {
				Ref string `json:"$ref"`
			} `json:"currency"`
			Payment struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"payment"`
			Price struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"price"`
			Quantity struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"quantity"`
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			InstrumentType struct {
				Ref string `json:"$ref"`
			} `json:"instrumentType"`
			IsMarginCall struct {
				Type string `json:"type"`
			} `json:"isMarginCall"`
			Date struct {
				Type        string    `json:"type"`
				Format      string    `json:"format"`
				Description string    `json:"description"`
				Example     time.Time `json:"example"`
			} `json:"date"`
			OperationType struct {
				Ref string `json:"$ref"`
			} `json:"operationType"`
		} `json:"properties"`
	} `json:"Operation"`
	OrdersResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"OrdersResponse"`
	Order struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			OrderID struct {
				Type string `json:"type"`
			} `json:"orderId"`
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Operation struct {
				Ref string `json:"$ref"`
			} `json:"operation"`
			Status struct {
				Ref string `json:"$ref"`
			} `json:"status"`
			RequestedLots struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"requestedLots"`
			ExecutedLots struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"executedLots"`
			Type struct {
				Ref string `json:"$ref"`
			} `json:"type"`
			Price struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"price"`
		} `json:"properties"`
	} `json:"Order"`
	LimitOrderRequest struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Lots struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"lots"`
			Operation struct {
				Ref string `json:"$ref"`
			} `json:"operation"`
			Price struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"price"`
		} `json:"properties"`
	} `json:"LimitOrderRequest"`
	LimitOrderResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"LimitOrderResponse"`
	PlacedLimitOrder struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			OrderID struct {
				Type string `json:"type"`
			} `json:"orderId"`
			Operation struct {
				Ref string `json:"$ref"`
			} `json:"operation"`
			Status struct {
				Ref string `json:"$ref"`
			} `json:"status"`
			RejectReason struct {
				Type string `json:"type"`
			} `json:"rejectReason"`
			RequestedLots struct {
				Type string `json:"type"`
			} `json:"requestedLots"`
			ExecutedLots struct {
				Type string `json:"type"`
			} `json:"executedLots"`
			Commission struct {
				Ref string `json:"$ref"`
			} `json:"commission"`
		} `json:"properties"`
	} `json:"PlacedLimitOrder"`
	TradeStatus struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"TradeStatus"`
	OperationType struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"OperationType"`
	OperationTypeWithCommission struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"OperationTypeWithCommission"`
	OperationStatus struct {
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Enum        []string `json:"enum"`
	} `json:"OperationStatus"`
	OrderStatus struct {
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Enum        []string `json:"enum"`
	} `json:"OrderStatus"`
	OrderType struct {
		Description string   `json:"description"`
		Type        string   `json:"type"`
		Enum        []string `json:"enum"`
	} `json:"OrderType"`
	SandboxSetCurrencyBalanceRequest struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Currency struct {
				Ref string `json:"$ref"`
			} `json:"currency"`
			Balance struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"balance"`
		} `json:"properties"`
	} `json:"SandboxSetCurrencyBalanceRequest"`
	SandboxSetPositionBalanceRequest struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Balance struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"balance"`
		} `json:"properties"`
	} `json:"SandboxSetPositionBalanceRequest"`
	MarketInstrumentListResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"MarketInstrumentListResponse"`
	MarketInstrumentList struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Total struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"total"`
			Instruments struct {
				Type  string `json:"type"`
				Items struct {
					Ref string `json:"$ref"`
				} `json:"items"`
			} `json:"instruments"`
		} `json:"properties"`
	} `json:"MarketInstrumentList"`
	MarketInstrumentResponse struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			TrackingID struct {
				Type string `json:"type"`
			} `json:"trackingId"`
			Status struct {
				Type    string `json:"type"`
				Default string `json:"default"`
			} `json:"status"`
			Payload struct {
				Ref string `json:"$ref"`
			} `json:"payload"`
		} `json:"properties"`
	} `json:"MarketInstrumentResponse"`
	MarketInstrument struct {
		Type       string   `json:"type"`
		Required   []string `json:"required"`
		Properties struct {
			Figi struct {
				Type string `json:"type"`
			} `json:"figi"`
			Ticker struct {
				Type string `json:"type"`
			} `json:"ticker"`
			Isin struct {
				Type string `json:"type"`
			} `json:"isin"`
			MinPriceIncrement struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"minPriceIncrement"`
			Lot struct {
				Type   string `json:"type"`
				Format string `json:"format"`
			} `json:"lot"`
			Currency struct {
				Ref string `json:"$ref"`
			} `json:"currency"`
			Name struct {
				Type string `json:"type"`
			} `json:"name"`
		} `json:"properties"`
	} `json:"MarketInstrument"`
	SandboxCurrency struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"SandboxCurrency"`
	Currency struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"Currency"`
	InstrumentType struct {
		Type string   `json:"type"`
		Enum []string `json:"enum"`
	} `json:"InstrumentType"`
}
