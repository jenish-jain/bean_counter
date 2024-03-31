package reporter

import (
	"bean_counter/internal/types/invoice"
	"bean_counter/pkg/gsheets"
	"bean_counter/pkg/utils"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jenish-jain/logger"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Handler struct {
}

type monthlyReportRequest struct {
	Month time.Month `json:"month,omitempty"`
	Year  int        `json:"year,omitempty"`
}

func (h *Handler) GenerateMonthlyGSTReport(ctx *gin.Context) {
	var req monthlyReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.ErrorWithCtx(ctx, "error reading request body err : %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid request received"})
		return
	}

	// Use req.Month and req.Year as needed

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
	reporter := NewReporter(invoiceService)

	month := getMonth(req.Month)
	year := getYear(req.Year, month)
	sheetName := fmt.Sprintf("%s %d", month, year)
	taxReport := reporter.GetTaxReportOfMonth(ctx, month, year)
	spreadsheetID := os.Getenv("REPORTER_SPREADSHEET_ID")
	sheetsRepo.AddNewWorksheet(ctx, spreadsheetID, sheetName)
	values := reporter.GetSheetValuesToPublishReport(taxReport, month, year)

	sheetsRepo.WriteToSheet(spreadsheetID, sheetName, values)
	ctx.JSON(http.StatusCreated, gin.H{"message": "report generated successfully"})
	return

}

func getMonth(reqMonth time.Month) time.Month {
	if reqMonth != 0 {
		return reqMonth
	}
	return utils.GetPreviousMonth()
}

func getYear(reqYear int, month time.Month) int {
	if reqYear != 0 {
		return reqYear
	}
	year := utils.GetCurrentYear()
	if month == time.December { // for requests of january we need to generate reports of december last year
		return year - 1
	}
	return year
}

func NewHandler() *Handler {
	return &Handler{}
}
