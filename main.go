package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bradford-hamilton/godax/pkg/godax"
)

// NOTE:
// This is a dumping ground for testing godax along the way until it's complete.
// The code lives in "/pkg/godax"
// TODO: consider moving everything under a "v1" dir

func main() {
	client, err := godax.NewSandboxClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	accounts, err := client.ListAccounts()
	fmt.Printf("\naccounts: %+v\n", accounts)

	account, err := client.GetAccount(accounts[0].ID)
	fmt.Printf("\naccount: %+v\n", account)

	his, err := client.GetAccountHistory(accounts[1].ID)
	if err != nil {
		fmt.Println("err getting account history", err)
		os.Exit(1)
	}
	fmt.Printf("\nhistory: %+v\n", his)

	holds, err := client.GetAccountHolds(accounts[0].ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("\nholds: %+v\n", holds)

	o, err := client.PlaceOrder(godax.OrderParams{
		CommonOrderParams: godax.CommonOrderParams{
			Side:      "buy",
			ProductID: "BTC-USD",
			Type:      "limit",
			Size:      "0.01",
			Price:     "0.100",
		},
	})
	if err != nil {
		fmt.Println("err ordering:", err)
	}

	fmt.Printf("\norder: %+v\n", o)

	orderID1, err := client.CancelOrderByID(o.ID, godax.QueryParams{godax.ProductIDParam: o.ProductID})
	if err != nil {
		fmt.Println("err canceling:", err)
	}

	fmt.Printf("\norderID1: %s\n", orderID1)

	// o, err = client.PlaceOrder(godax.OrderParams{
	// 	CommonOrderParams: godax.CommonOrderParams{
	// 		Side:      "buy",
	// 		ProductID: "BTC-USD",
	// 		Type:      "market",
	// 		Size:      "0.01",
	// 		ClientOID: "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3",
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // _, err = client.CancelOrderByClientOID("c6dfb02e-7f65-4e02-8fa3-866d46ed15b3", &o.ProductID)
	// // if err != nil {
	// // 	fmt.Println("err canceling:", err)
	// // }

	// // fmt.Printf("orderID2: %s\n", orderID2)

	for i := 0; i < 5; i++ {
		_, err = client.PlaceOrder(godax.OrderParams{
			CommonOrderParams: godax.CommonOrderParams{
				Side:      "buy",
				ProductID: "BTC-USD",
				Type:      "limit",
				Size:      "0.01",
				Price:     "0.100",
				ClientOID: "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3",
			},
		})
		if err != nil {
			fmt.Println("err placing an order: ", err)
			os.Exit(1)
		}
	}

	ords, err := client.ListOrders(godax.QueryParams{godax.ProductIDParam: "BTC-USD"})
	if err != nil {
		fmt.Println("err listing orders: ", err)
		os.Exit(1)
	}
	fmt.Printf("\norders: %+v\n", ords)

	// gotOrder, err := client.GetOrderByID(ords[0].ID)
	// fmt.Printf("gotOrder: %+v\n", gotOrder)

	orderIDs, err := client.CancelAllOrders(nil)
	if err != nil {
		fmt.Println("err canceling:", err)
		os.Exit(1)
	}

	fmt.Printf("\norderIDs: %+v\n", orderIDs)

	// fmt.Printf("Orders should be empty now: %+v\n", ords)
	qp := godax.QueryParams{godax.OrderIDParam: o.ID}
	fills, err := client.ListFills(qp)
	if err != nil {
		fmt.Printf("err getting fills: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nfills: %+v\n", fills)

	limits, err := client.GetCurrentExchangeLimits()
	if err != nil {
		fmt.Printf("err getting exchange limits: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nlimits: %+v\n", limits)

	conv, err := client.StableCoinConversion("USD", "USDC", "10")
	if err != nil {
		fmt.Printf("\nerr converting coins: %+v\n", err)
	}
	fmt.Printf("\nconv: %+v \n", conv)

	pms, err := client.ListPaymentMethods()
	if err != nil {
		fmt.Printf("\nerr listing payment methods: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\npms: %+v\n", pms)

	cbActs, err := client.ListCoinbaseAccounts()
	if err != nil {
		fmt.Printf("\nerr listing payment methods: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\ncbActs: %+v\n", cbActs)

	fees, err := client.GetCurrentFees()
	if err != nil {
		fmt.Printf("\nerr getting current fees: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nfees: %+v\n", fees)

	v, err := client.GetTrailingVolume()
	if err != nil {
		fmt.Printf("\nerr getting trailing volume: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nuser accounts: %+v\n", v)

	profiles, err := client.ListProfiles()
	if err != nil {
		fmt.Printf("\nerr listing profiles: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nprofiles: %+v\n", profiles)

	profile, err := client.GetProfile(profiles[0].ID)
	if err != nil {
		fmt.Printf("\nerr getting profile: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nprofile: %+v\n", profile)

	// if err := client.ProfileTransfer(godax.TransferParams{
	// 	From:     profiles[0].ID,
	// 	To:       profiles[1].ID,
	// 	Currency: "BTC-USD",
	// 	Amount:   "0.05",
	// }); err != nil {
	// 	fmt.Println("err transferring profile:", err)
	// 	os.Exit(1)
	// }

	products, err := client.ListProducts()
	if err != nil {
		fmt.Printf("\nerr getting profile: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nproducts: %+v\n", products)

	product, err := client.GetProductByID("BTC-USD")
	if err != nil {
		fmt.Printf("\nerr getting product: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nproduct: %+v\n", product)

	ob, err := client.GetProductOrderBook("BTC-USD", godax.QueryParams{godax.LevelParam: "1"})
	if err != nil {
		fmt.Printf("\nerr getting order book: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nob: %+v\n", ob)

	ticker, err := client.GetProductTicker("BTC-USD")
	if err != nil {
		fmt.Printf("\nerr getting ticker: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nticker: %+v\n", ticker)

	trades, err := client.ListTradesByProduct("BTC-USD")
	if err != nil {
		fmt.Printf("\nerr getting trades: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\ntrades: %+v\n", trades)

	rates, err := client.GetHistoricRatesForProduct("BTC-USD", godax.QueryParams{godax.StartParam: "", godax.EndParam: "", godax.GranularityParam: "60"})
	if err != nil {
		fmt.Printf("\nerr getting rates: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nrates: %+v\n", rates)

	stats, err := client.Get24HourStatsForProduct("BTC-USD")
	if err != nil {
		fmt.Printf("\nerr getting stats: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nstats: %+v\n", stats)

	currencies, err := client.ListCurrencies()
	if err != nil {
		fmt.Printf("\nerr getting currencies: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\ncurrencies: %+v\n", currencies)

	srvTime, err := client.GetServerTime()
	if err != nil {
		fmt.Printf("\nerr getting srvTime: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nsrvTime: %+v\n", srvTime)

	report, err := client.CreateReport(godax.ReportParams{
		Type:      "account",
		StartDate: "2020-03-01T00:00:00.000Z",
		EndDate:   "2020-07-26T00:00:00.000Z",
		AccountID: "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
		Email:     "brad.lamson@gmail.com",
	})
	if err != nil {
		fmt.Printf("\nerr getting report: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nreport: %+v\n", report)

	time.Sleep(2 * time.Second)

	status, err := client.GetReportStatus(report.ID)
	if err != nil {
		fmt.Printf("\nerr getting status: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\nstatus: %+v\n", status)

	oracle, err := client.GetOracle()
	if err != nil {
		fmt.Printf("\nerr getting oracle: %+v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\noracle: %+v\n", oracle)

	// TODO: message coinbase about sandbox margin so I can test these
	bp, err := client.GetBuyingPower(godax.QueryParams{godax.ProductIDParam: "BTC-USD"})
	if err != nil {
		fmt.Printf("\nerr getting bp: %+v\n", err)
	}
	fmt.Printf("\nbuying power: %+v\n", bp)

	mp, err := client.GetMarginProfile(godax.QueryParams{godax.ProductIDParam: "BTC-USD"})
	if err != nil {
		fmt.Printf("\nerr getting mp: %+v\n", err)
	}
	fmt.Printf("\nmargin profile: %+v\n", mp)

	wp, err := client.GetWithdrawalPowerForCurrency(godax.QueryParams{godax.CurrencyParam: "BTC-USD"})
	if err != nil {
		fmt.Printf("\nerr getting wp: %+v\n", err)
	}
	fmt.Printf("\nwithdrawal power: %+v\n", wp)

	allWp, err := client.GetAllWithdrawalPower()
	if err != nil {
		fmt.Printf("\nerr getting allWp: %+v\n", err)
	}
	fmt.Printf("\nallWp: %+v\n", allWp)
}
