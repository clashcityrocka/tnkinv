package schema

import (
	log "github.com/sirupsen/logrus"

	"strings"

	"../aux"
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

// TODO why json tags?
type Instrument struct {
	Figi      string `json:"figi"`
	Ticker    string `json:"ticker"`
	Name      string `json:"name"`
	Currency  string `json:"currency"`
	Lot       int
	FaceValue int `json:"faceValue"`

	Type    InsType
	Section Section
}

func NewInstrument(figi, ticker, name, typ, currency string, faceValue, lot int) Instrument {
	if !Currencies.Has(currency) {
		log.Fatalf("unknown currency %s (%s)", currency, ticker)
	}

	ins := Instrument{
		Figi:      figi,
		Ticker:    ticker,
		Name:      name,
		Currency:  currency,
		FaceValue: faceValue,
		Lot:       lot,
	}
	ins.Type = getInstrumentType(typ, ticker)
	ins.Section = getSection(ins)
	return ins
}

func getInstrumentType(typ string, ticker string) InsType {
	if !map[InsType]bool{
		InsTypeEtf:      true,
		InsTypeStock:    true,
		InsTypeBond:     true,
		InsTypeCurrency: true,
	}[InsType(typ)] {
		log.Warnf("Unhandled type %s: %s", ticker, typ)
	}
	return InsType(typ)
}

func (ins Instrument) Benchmark() string {
	if bench, ok := map[Section]string{
		BondRub:  "VTBB",
		BondUsd:  "FXRU",
		StockRub: "FXRL",
		StockUsd: "FXUS", // see below
	}[ins.Section]; ok {
		if bench == "FXUS" && aux.IsIn(ins.Ticker,
			"AAPL", // 18.1%
			"MSFT", // 16.3%
			"GOOG", //  9.6%
			"FB",   //  6.2%
			"V",    //  3.6%
			"MA",   //  2.9%
			"INTC", //  2.8%
			"NVDA", //  2.6%
			"NFLX", //  2.3%
		) {
			return "FXIT"
		}
		if bench != ins.Ticker {
			return bench
		}
	}

	return ""
}

func (s Section) Currency() string {
	for cur := range Currencies {
		if strings.Contains(string(s), cur) {
			return cur
		}
	}
	return ""
}

func GetEtfSection(ticker string) (Section, bool) {
	s, ok := map[string]Section{
		"VTBB": BondRub,
		"FXRB": BondRub,

		// T* funds are (25x4 gold, stocks, long and short bonds)
		// TODO proper accounting
		// consider them bonds for now
		"TRUR": BondRub,
		"TUSD": BondUsd,

		"FXRU": BondUsd,

		"SBMX": StockRub,
		"FXRL": StockRub,
		"TMOS": StockRub,

		"AKNX": StockUsd,
		"FXIT": StockUsd,
		"FXUS": StockUsd,
		// it's actually StockEur, but leave it for now
		"FXDE": StockUsd,
		"TECH": StockUsd,
		"TSPX": StockUsd,
		"TIPO": StockUsd,
		"TBIO": StockUsd,
		"VTBE": StockUsd,

		"FXMM": CashRub,
		"FXTB": CashUsd,
	}[ticker]
	return s, ok
}

func getSection(ins Instrument) Section {
	if s, ok := map[string]Section{
		InsTypeBond + "RUB": BondRub,
		InsTypeBond + "USD": BondUsd,

		InsTypeStock + "RUB": StockRub,
		InsTypeStock + "USD": StockUsd,

		InsTypeCurrency + "RUB": CashRub,
		InsTypeCurrency + "USD": CashUsd,
	}[string(ins.Type)+ins.Currency]; ok {
		return s
	}

	if ins.Type == InsTypeEtf {
		s, ok := GetEtfSection(ins.Ticker)
		if !ok {
			log.Warnf("Uncatched ETF %s", ins.Ticker)
		}
		return s
	}

	return ""
}
