package main

import (
	"fmt"
	"os"

	"transaction-service/internal/app"
	"transaction-service/internal/transaction"
)

const port = 8080

func main() {
	application := app.New(app.NewDependencies(transaction.NewUUIDGenerator(), app.NewHttpClient()))
	if err := application.Start(port); err != nil {
		fmt.Printf("An error occured: %v", err)
		os.Exit(1)
	}
}
