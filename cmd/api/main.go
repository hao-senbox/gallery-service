package main

import (
	"fmt"
	"gallery-service/cmd/api/cli"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Some error occurred during execute app. Error: %v\n", err)

		os.Exit(2)
	}
}
