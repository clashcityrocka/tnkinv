package portfolio

import (
	"fmt"
	"log"
	"math"
	"sort"
	"time"

	"../client"
	"../schema"
)

var beginning = time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

type Portfolio struct {
	client *client.MyClient

	data struct {
		ops schema.OperationsResponse
	}

	tickers   map[string]string
	positions map[string]*schema.PositionInfo
	mcandles  map[string]*schema.CandlesResponse

	figisSorted []string

	totals *schema.Balance
}

func NewPortfolio(c *client.MyClient) *Portfolio {
	return &Portfolio{
		client:    c,
		tickers:   make(map[string]string),
		positions: make(map[string]*schema.PositionInfo),
		mcandles:  make(map[string]*schema.CandlesResponse),
		totals:    schema.NewBalance(),
	}
}

// =============================================================================

func (p *Portfolio) getTicker(figi string) string {
	ticker, ok := p.tickers[figi]
	if !ok {
		ticker = p.client.RequestTicker(figi)
		p.tickers[figi] = ticker
	}
	return ticker
}

// =============================================================================

func (p *Portfolio) processPortfolio() {
	pfResp := p.client.RequestPortfolio()

	for _, pos := range pfResp.Payload.Positions {
		p.tickers[pos.Figi] = pos.Ticker
	}
}

// =============================================================================

func (p *Portfolio) preprocessOperations(start time.Time) {
	ops := p.client.RequestOperations(start)
	for i := range ops.Payload.Operations {
		var err error
		op := &ops.Payload.Operations[i]
		op.DateParsed, err = time.Parse(time.RFC3339, op.Date)
		if err != nil {
			log.Fatal("Failed to parse time: %v", err)
		}
	}

	sort.Slice(ops.Payload.Operations, func(i, j int) bool {
		return ops.Payload.Operations[i].DateParsed.Before(ops.Payload.Operations[j].DateParsed)
	})

	p.data.ops = ops
}

// =============================================================================

func (p *Portfolio) processOperation(op schema.Operation) (deal *schema.Deal) {
	if op.Figi == "" {
		// payins, service commissions
		return
	}

	pinfo := p.positions[op.Figi]
	if pinfo == nil {
		pinfo = &schema.PositionInfo{
			Figi:              op.Figi,
			Ticker:            p.getTicker(op.Figi),
			AccumulatedIncome: schema.NewCValue(0, op.Currency),
		}
		p.positions[op.Figi] = pinfo
	}

	if op.IsTrading() {
		deal = &schema.Deal{
			Date:       op.DateParsed,
			Price:      schema.NewCValue(op.Price, op.Currency),
			Quantity:   int(op.Quantity),
			Commission: op.Commission.Value,
		}
		if op.OperationType == "Sell" {
			deal.Quantity = -deal.Quantity
		}

	} else if op.OperationType == "BrokerCommission" {
		// negative
		pinfo.AccumulatedIncome.Value += op.Payment
	} else if op.OperationType == "Dividend" || op.OperationType == "TaxDividend" {
		// positive, negative
		pinfo.AccumulatedIncome.Value += op.Payment
		pinfo.Dividends = append(pinfo.Dividends,
			&schema.Dividend{
				Date:  op.DateParsed,
				Value: op.Payment,
			})
	} else {
		log.Printf("Unprocessed transaction %v", op)
	}

	return
}

// =============================================================================

/*

 Balance consists of:

 USD
 + Assets:
    1. Cash balance
        1.1 Direct payins
        1.2 Exchanges
        1.3 Sold stocks
        1.4 - Bought stocks
        1.5 - Service commissions
    2. Open USD positions
 - Payins
    3. Directs payins
    4. Exchanges

RUB
 + Assets:
    1. Cash balance
        1.1 Direct payins
        1.3 Sold stocks & dollars
        1.4 - Bought stocks & dollars
        1.5 - Service commissions
    2. Open RUB positions
 - Payins:
    3. Direct payins
    5. - Exchanged money

*/

func (p *Portfolio) addOpToBalance(bal *schema.Balance, op schema.Operation) {
	if op.OperationType == "Dividend" || op.OperationType == "TaxDividend" {
		bal.Assets[op.Currency].Value += op.Payment
	}

	if op.Figi != "" {
		return
	}

	if op.OperationType == "PayIn" {
		// 1.1
		bal.Assets[op.Currency].Value += op.Payment
		// 3
		bal.Payins[op.Currency].Value += op.Payment
	} else if op.OperationType == "ServiceCommission" {
		bal.Commissions[op.Currency].Value += op.Payment
		// 1.5
		bal.Assets[op.Currency].Value -= -op.Payment
	} else {
		log.Printf("Unprocessed transaction 2 %v", op)
	}
}

func (p *Portfolio) addDealToBalance(bal *schema.Balance, figi string, deal *schema.Deal) {
	pinfo := p.positions[figi]
	if pinfo.Figi == schema.FigiUSD {
		// Exchanges
		// 1.2
		bal.Assets["USD"].Value += float64(deal.Quantity)
		// 4
		bal.Payins["USD"].Value += float64(deal.Quantity)
		// 5
		bal.Payins["RUB"].Value -= deal.Value()
	}
	// 1.3, 1.4, 2
	bal.Assets[deal.Price.Currency].Value -= deal.Value() - deal.Commission
}

// =============================================================================

func (p *Portfolio) getOpenPortion(pinfo *schema.PositionInfo) *schema.Portion {
	if len(pinfo.Portions) == 0 {
		return nil
	}

	po := pinfo.Portions[len(pinfo.Portions)-1]
	if po.IsClosed {
		return nil
	}

	return po
}

func (p *Portfolio) addToPortions(pinfo *schema.PositionInfo, deal *schema.Deal) {
	po := p.getOpenPortion(pinfo)
	if po == nil {
		po = &schema.Portion{
			Balance: schema.NewCValue(0, deal.Price.Currency),
			AvgDate: deal.Date,
		}
		pinfo.Portions = append(pinfo.Portions, po)
	}

	pinfo.OpenQuantity += deal.Quantity
	pinfo.OpenSpent += deal.Value()

	if deal.Quantity > 0 { // buy
		// TODO think again is this correct?
		mult := deal.Value() / pinfo.OpenSpent

		biasDays := int(math.Round(deal.Date.Sub(po.AvgDate).Hours() * mult / 24))
		po.AvgDate.AddDate(0, 0, biasDays)

		po.AvgPrice.Value = deal.Price.Value*mult + po.AvgPrice.Value*(1-mult)
		po.Buys = append(po.Buys, deal)

	} else { // sell
		if pinfo.OpenQuantity > 0 {
			log.Printf("Partial sells are not handled nicely yet")
			return
		}
		if pinfo.OpenQuantity < 0 {
			log.Fatal("negative balance? %v", pinfo)
		}

		pinfo.OpenSpent = 0

		// complete sell
		po.Close = deal
		po.IsClosed = true
	}
}

func (p *Portfolio) makeOpenDeal(pinfo *schema.PositionInfo, date time.Time, price float64, setClose bool) *schema.Deal {
	po := p.getOpenPortion(pinfo)
	if po == nil {
		return nil
	}

	deal := &schema.Deal{
		Date:     date,
		Price:    schema.NewCValue(price, po.Balance.Currency),
		Quantity: -pinfo.OpenQuantity,
	}

	if setClose {
		po.Close = deal
		pinfo.OpenDeal = deal
	}

	return deal
}

func (p *Portfolio) makePortionYields(pinfo *schema.PositionInfo) {
	for _, po := range pinfo.Portions {
		var expense float64

		profit := po.Close.Price.Mult(float64(-po.Close.Quantity))

		for _, div := range pinfo.Dividends {
			if div.Date.Before(po.Buys[0].Date) {
				continue
			}
			if div.Date.After(po.Close.Date) {
				// TODO not quite right. Dividends come with delay
				continue
			}
			profit.Value += div.Value
		}

		for _, deal := range po.Buys {
			expense += deal.Value()
			expense += -deal.Commission
		}
		expense += -po.Close.Commission

		po.Yield = profit.Div(expense / 100)
		po.Yield.Value -= 100

		po.YieldAnnual = po.Yield.Value * 365 / (po.Close.Date.Sub(po.AvgDate).Hours() / 24)

		po.Balance = profit
		po.Balance.Value -= expense
	}
}

// =============================================================================

func (p *Portfolio) processOperations(cb func(*schema.Balance, time.Time) bool) *schema.Balance {
	p.preprocessOperations(beginning)

	//log.Print("== Transaction log ==")

	bal := schema.NewBalance()

	for _, op := range p.data.ops.Payload.Operations {
		/* log.Printf("at %s %s some %s",
		op.DateParsed.String(), op.OperationType+"-ed",
		p.tickers[op.Figi])*/

		if op.Status != "Done" {
			// cancelled declined etc
			// noone is interested in that
			continue
		}

		if !cb(bal, op.DateParsed) {
			break
		}

		deal := p.processOperation(op)
		if deal != nil {
			pinfo := p.positions[op.Figi]
			pinfo.Deals = append(pinfo.Deals, deal)

			p.addToPortions(pinfo, deal)
			p.addDealToBalance(bal, pinfo.Figi, deal)
		}

		p.addOpToBalance(bal, op)
	}

	return bal
}

func (p *Portfolio) addOpenDealsToBalance(bal *schema.Balance, time time.Time, pricef func(string) float64) {
	for _, pinfo := range p.positions {
		od := p.makeOpenDeal(pinfo, time, pricef(pinfo.Figi), true)

		if od != nil && pinfo.Figi != schema.FigiUSD {
			p.addDealToBalance(bal, pinfo.Figi, od)
		}
	}
}

func (p *Portfolio) Collect(at time.Time) {
	c := p.client

	p.processPortfolio()

	p.totals = p.processOperations(func(bal *schema.Balance, opTime time.Time) bool {
		return opTime.Before(at)
	})

	p.addOpenDealsToBalance(p.totals, time.Now(), c.RequestCurrentPrice) // TODO only true with at==now

	for _, pinfo := range p.positions {
		p.makePortionYields(pinfo)
	}
}

// =============================================================================

func (p *Portfolio) ListDeals(start time.Time) {
	empty := true

	p.processPortfolio()

	p.preprocessOperations(start)

	bal := schema.NewBalance()
	for _, op := range p.data.ops.Payload.Operations {
		if op.Status != "Done" {
			// cancelled declined etc
			// noone is interested in that
			continue
		}

		if op.Figi != "" {
			op.Ticker = p.getTicker(op.Figi)
		}
		fmt.Printf("%s\n", op)

		// exploit those balance maps for totals
		if op.IsTrading() {
			bal.Assets[op.Currency].Value += math.Abs(op.Payment)
			empty = false
		} else if op.OperationType == "ServiceCommission" || op.OperationType == "BrokerCommission" {
			bal.Payins[op.Currency].Value += math.Abs(op.Payment)
			empty = false
		}
	}

	if empty {
		return
	}

	fmt.Printf(" - Total deals:\n")
	for c := range schema.Currencies {
		if bal.Assets[c].Value != 0 {
			fmt.Printf("\t %s\n", bal.Assets[c])
		}
	}
	fmt.Printf("   commissions:\n")
	for c := range schema.Currencies {
		if bal.Payins[c].Value != 0 {
			fmt.Printf("\t %s\n", bal.Payins[c])
		}
	}

	comms, deals, _ := bal.GetTotal(p.client.RequestCurrentPrice(schema.FigiUSD), 0)
	fmt.Printf("   percentage: %.2f%%\n", comms/deals*100)
}

// =============================================================================

func (p *Portfolio) getCandles(figi string, start time.Time, period string) *schema.CandlesResponse {
	pcandles := p.mcandles[figi]
	if pcandles != nil {
		return pcandles
	}

	resp := p.client.RequestCandles(figi, start, time.Now(), period)
	pcandles = &resp

	for i := range pcandles.Payload.Candles {
		var err error
		c := &pcandles.Payload.Candles[i]

		c.TimeParsed, err = time.Parse(time.RFC3339, c.Time)
		if err != nil {
			log.Fatal("failed to parse time %v", err)
		}
	}

	sort.Slice(pcandles.Payload.Candles, func(i, j int) bool {
		return pcandles.Payload.Candles[i].TimeParsed.Before(pcandles.Payload.Candles[j].TimeParsed)
	})

	p.mcandles[figi] = pcandles
	return pcandles
}

func (p *Portfolio) ListBalances(start time.Time, period string) {
	// just for time reference, can be any figi
	candles := p.getCandles(schema.FigiUSD, start, period).Payload.Candles

	cidx := 0
	num := len(candles)

	bal := p.processOperations(func(bal *schema.Balance, opTime time.Time) (cont bool) {
		cont = true

		if cidx == num {
			return
		}

		t := candles[cidx].TimeParsed

		if opTime.Before(t) {
			// proceed with current candle
			return
		}

		pricef := func(figi string) float64 {
			return p.getCandles(figi, start, period).Payload.Candles[cidx].O
		}

		localBal := bal.Copy()
		p.addOpenDealsToBalance(localBal, t, pricef)

		pbal(t, localBal, pricef(schema.FigiUSD))

		cidx += 1
		return
	})

	pricef := func(figi string) float64 {
		return p.getCandles(figi, start, period).Payload.Candles[num-1].C
	}

	p.addOpenDealsToBalance(bal, time.Now(), pricef)
	pbal(time.Now(), bal, pricef(schema.FigiUSD))
}

func pbal(t time.Time, b *schema.Balance, usdprice float64) {
	for _, cur := range []string{"USD", "RUB"} {
		fmt.Printf("%s, %s, %f, %f, %f\n",
			t.Format("2006/01/02"), cur,
			b.Payins[cur].Value, b.Assets[cur].Value, b.Get(cur).Value)
	}
	p, a, d := b.GetTotal(usdprice, 0)
	fmt.Printf("%s, %s, %f, %f, %f\n",
		t.Format("2006/01/02"), "---",
		p, a, d)
}
