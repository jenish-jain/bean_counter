package main

import (
	"bean_counter/internal/reporter"
	"bean_counter/internal/types/invoice"
	"bean_counter/pkg/gsheets"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	http.HandleFunc("/gstReport/monthly", generateMonthlyGstReport)
	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	// Start HTTP server.
	log.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func generateMonthlyGstReport(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Month time.Month `json:"month"`
		Year  int        `json:"year"`
	}

	var req request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading body err : %+v \n", err)
		return
	}
	err = json.Unmarshal(body, &req)
	if err != nil {
		fmt.Printf("Error unmarhsaling request err : %+v \n", err)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "invalid request type")
		return
	}
	ctx := context.Background()
	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(os.Getenv("KEY_JSON_BASE64"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// authenticate and get configuration
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
	invoiceService := invoice.NewService(sheetsRepo)
	reporter := reporter.NewReporter(invoiceService)

	month := req.Month
	year := req.Year
	sheetName := fmt.Sprintf("%s %d", month, year)
	taxReport := reporter.GetTaxReportOfMonth(month, year)
	spreadsheetID := os.Getenv("REPORTER_SPREADSHEET_ID")
	sheetsRepo.AddNewWorksheet(spreadsheetID, sheetName)
	values := reporter.GetSheetValuesToPublishReport(taxReport, month, year)

	sheetsRepo.WriteToSheet(spreadsheetID, sheetName, values)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "report generated successfully")
	return
}
