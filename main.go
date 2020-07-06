package main

import (
	"fmt"
	"os"

	"github.com/bradford-hamilton/godax/pkg/godax"
)

func main() {
	client, err := godax.NewSandboxClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// accounts, err := client.ListAccounts()
	// fmt.Println(accounts)

	// account, err := client.GetAccount(accounts[0].ID)
	// fmt.Println(account)

	// his, err := client.GetAccountHistory(accounts[1].ID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println(his)

	// holds, err := client.GetAccountHolds(accounts[0].ID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println("HODLS", holds)

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
		fmt.Println(err)
	}

	_, err = client.CancelOrderByID(o.ID, &o.ProductID)
	if err != nil {
		fmt.Println("err canceling:", err)
	}

	// fmt.Printf("orderID1: %s\n", orderID1)

	o, err = client.PlaceOrder(godax.OrderParams{
		CommonOrderParams: godax.CommonOrderParams{
			Side:      "buy",
			ProductID: "BTC-USD",
			Type:      "market",
			Size:      "0.01",
			ClientOID: "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3",
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	// _, err = client.CancelOrderByClientOID("c6dfb02e-7f65-4e02-8fa3-866d46ed15b3", &o.ProductID)
	// if err != nil {
	// 	fmt.Println("err canceling:", err)
	// }

	// fmt.Printf("orderID2: %s\n", orderID2)

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
			fmt.Println(err)
		}
	}

	ords, err := client.ListOrders(nil, nil)
	if err != nil {
		fmt.Println("err listing orders: ", err)
		os.Exit(1)
	}
	// fmt.Printf("Orders here: %+v\n", ords)

	gotOrder, err := client.GetOrderByID(ords[0].ID)
	fmt.Printf("gotOrder: %+v\n", gotOrder)

	_, err = client.CancelAllOrders(nil)
	if err != nil {
		fmt.Println("err canceling:", err)
		os.Exit(1)
	}

	// fmt.Printf("orderIDs: %+v", orderIDs)

	// fmt.Printf("Orders should be empty now: %+v\n", ords)

	fills, err := client.ListFills(&o.ID, nil)
	if err != nil {
		fmt.Printf("err getting fills: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("fills: %+v\n", fills)

	limits, err := client.GetCurrentExchangeLimits()
	if err != nil {
		fmt.Printf("err getting exchange limits: %+v\n", err)
		os.Exit(1)
	}

	fmt.Printf("limits: %+v", limits)
}

func stringPtr(str string) *string {
	return &str
}
