package main

import (
	"bean_counter/internal/reporter"
	"bean_counter/internal/types/invoice"
	"bean_counter/pkg/gsheets"
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"time"
)

func main() {

	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// authenticate and get configuration
	fmt.Println(string(credBytes))
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println(err)
		return
	}

	sheetsRepo := gsheets.NewRepository(*srv)
	invoiceService := invoice.NewService(sheetsRepo) // check on value and pointer once
	reporter := reporter.NewReporter(invoiceService)

	finalTxn := reporter.GetTaxReportOfMonth(time.March, 2023)
	fmt.Printf("\n %+v \n", finalTxn)
	
}
