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
	fmt.Printf("%+v", accounts)
}
