package main

import (
	"fmt"
	"os"

	"github.com/bradford-hamilton/godax/pkg/godax"
)

// NOTE:
// This is a dumping ground for testing godax until it's more complete.
// Godax lives at "github.com/bradford-hamilton/godax/pkg/godax"

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

	orderID1, err := client.CancelOrderByID(o.ID, godax.QueryParams{godax.ProductID: o.ProductID})
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

	ords, err := client.ListOrders(godax.QueryParams{godax.ProductID: "BTC-USD"})
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
	qp := godax.QueryParams{godax.OrderID: o.ID}
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
}

func stringPtr(str string) *string {
	return &str
}
