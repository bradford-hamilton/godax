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

	accounts, err := client.ListAccounts()
	fmt.Println(accounts)

	account, err := client.GetAccount(accounts[0].ID)
	fmt.Println(account)

	his, err := client.GetAccountHistory(accounts[1].ID)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("history!!!!!!!!!!!", his)
}
