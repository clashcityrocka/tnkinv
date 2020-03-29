package client

import (
	"encoding/json"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"

	swagger "../go-client"
	"../schema"
)

type MyClient struct {
	swc    *swagger.APIClient
	tokenf string
}

func NewClient(tokenf string) *MyClient {
	return &MyClient{
		tokenf: tokenf,
	}
}

type optional struct {
	value string
}

func (o optional) IsSet() bool {
	return o.value != ""
}

func (o optional) Value() interface{} {
	return o.value
}

func (c *MyClient) getToken(fname string) string {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

func (c *MyClient) TrySandbox() error {
	conf := swagger.NewConfiguration()
	conf.BasePath = "https://api-invest.tinkoff.ru/openapi/sandbox/"
	conf.AddDefaultHeader("Authorization", "Bearer "+c.getToken(c.tokenf))

	swc := swagger.NewAPIClient(conf)

	sand := swc.SandboxApi

	_, err := sand.SandboxRegisterPost(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("sandbox register complete")
	return nil
}

func (c *MyClient) getAPI() *swagger.APIClient {
	if c.swc == nil {
		conf := swagger.NewConfiguration()
		conf.BasePath = "https://api-invest.tinkoff.ru/openapi/"
		conf.AddDefaultHeader("Authorization", "Bearer "+c.getToken(c.tokenf))

		c.swc = swagger.NewAPIClient(conf)
	}
	return c.swc
}

func (c *MyClient) RequestCurrentPrice(figi string) float64 {
	mktApi := c.getAPI().MarketApi
	mktResp := schema.OrderbookResponse{}

	body, err := mktApi.MarketOrderbookGet(nil, figi, 1)
	if err != nil {
		log.Fatalf("price(%s): %s", figi, err)
	}

	err = json.Unmarshal(body, &mktResp)
	if err != nil {
		log.Fatalf("price(%s): %s", figi, err)
	}

	log.Trace(string(body))

	return mktResp.Payload.LastPrice
}

func (c *MyClient) RequestTicker(figi string) string {
	mktApi := c.getAPI().MarketApi
	resp := schema.SearchByFigiResponse{}

	body, err := mktApi.MarketSearchByFigiGet(nil, figi)
	if err != nil {
		log.Fatalf("by figi(%s): %s", figi, err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatalf("by figi(%s): %s", figi, err)
	}

	return resp.Payload.Ticker
}

func (c *MyClient) RequestFigi(ticker string) string {
	mktApi := c.getAPI().MarketApi
	resp := schema.SearchByTickerResponse{}

	body, err := mktApi.MarketSearchByTickerGet(nil, ticker)
	if err != nil {
		log.Fatalf("by ticker(%s): %s", ticker, err)
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		log.Fatalf("by ticker(%s): %s", ticker, err)
	}

	log.Trace(string(body))

	return resp.Payload.Instruments[0].Figi
}

func (c *MyClient) RequestPortfolio(acc string) schema.PortfolioResponse {
	pfApi := c.getAPI().PortfolioApi
	pfResp := schema.PortfolioResponse{}
	opts := &swagger.PortfolioGetOpts{
		BrokerAccountId: optional{acc},
	}

	body, err := pfApi.PortfolioGet(nil, opts)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &pfResp)
	if err != nil {
		log.Fatal(err)
	}

	log.Trace(string(body))

	return pfResp
}

func (c *MyClient) RequestOperations(start time.Time, acc string) schema.OperationsResponse {
	timeStartStr := start.Format(time.RFC3339)
	timeNow := time.Now()

	opsApi := c.getAPI().OperationsApi
	opsResp := schema.OperationsResponse{}
	opts := &swagger.OperationsGetOpts{
		Figi:            optional{},
		BrokerAccountId: optional{acc},
	}

	body, err := opsApi.OperationsGet(nil, timeStartStr, timeNow.Format(time.RFC3339), opts)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &opsResp)
	if err != nil {
		log.Fatal(err)
	}

	log.Trace(string(body))

	return opsResp
}

func (c *MyClient) RequestCandles(figi string, t1, t2 time.Time, interval string) schema.CandlesResponse {
	t1Str := t1.Format(time.RFC3339)
	t2Str := t2.Format(time.RFC3339)

	mktApi := c.getAPI().MarketApi
	mktResp := schema.CandlesResponse{}

	body, err := mktApi.MarketCandlesGet(nil, figi, t1Str, t2Str, interval)
	if err != nil {
		log.Fatalf("candles(%s, %s): %s", figi, t1, err)
	}

	err = json.Unmarshal(body, &mktResp)
	if err != nil {
		log.Fatalf("candles(%s, %s): %s", figi, t1, err)
	}

	log.Trace(string(body))

	return mktResp
}

func (c *MyClient) RequestAccounts() schema.AccountsResponse {
	userApi := c.getAPI().UserApi
	accResp := schema.AccountsResponse{}

	body, err := userApi.UserAccountsGet(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &accResp)
	if err != nil {
		log.Fatal(err)
	}

	log.Trace(string(body))

	return accResp
}

func (c *MyClient) Stop() {
}
