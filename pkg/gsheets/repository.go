package gsheets

import (
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
)

type Repository interface {
	GetAllRecords(spreadSheetID string, sheetID int) [][]interface{}
	AddNewWorksheet(spreadSheetID string, sheetName string) bool
	WriteToSheet(spreadSheetID string, sheetName string, values [][]interface{})
}

type repositoryImpl struct {
	service sheets.Service
}

func (r repositoryImpl) GetAllRecords(spreadSheetID string, sheetID int) [][]interface{} {
	sheetName := r.getSheetName(spreadSheetID, sheetID)
	response, err := r.service.Spreadsheets.Values.
		Get(spreadSheetID, sheetName).Do()
	if err != nil || response.HTTPStatusCode != 200 {
		fmt.Printf("error getting all records for spreadsheetId : %s & sheetId : %d \n", spreadSheetID, sheetID)
		panic(err)
	}

	return response.Values
}

func (r repositoryImpl) AddNewWorksheet(spreadSheetID string, sheetName string) bool {
	sheetReq := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetName,
				//SheetId: sheetID,
			},
		},
	}
	batchRequest := sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: true,
		Requests:                     []*sheets.Request{&sheetReq},
	}
	resp, err := r.service.Spreadsheets.BatchUpdate(spreadSheetID, &batchRequest).Do()
	if err != nil {
		fmt.Errorf("error creating new worksheet for spreadSheetID : %s , error : %+v", spreadSheetID, err)
		return false
	}
	fmt.Printf("successfully created new worksheet for spreadSheetID : %s , response : %+v", spreadSheetID, resp)
	return true
}

func (r repositoryImpl) WriteToSheet(spreadSheetID string, sheetName string, values [][]interface{}) {
	row := &sheets.ValueRange{
		Values: values,
	}
	range_ := fmt.Sprintf("%s!A1", sheetName)

	_, err := r.service.Spreadsheets.Values.
		Update(spreadSheetID, range_, row).
		ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func (r repositoryImpl) getSheetName(spreadSheetID string, sheetID int) string {
	// Convert sheet ID to sheet name.
	response1, err := r.service.Spreadsheets.
		Get(spreadSheetID).
		Fields("sheets(properties(sheetId,title))").
		Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		fmt.Printf("error getting response from gsheets service for spreadsheetId %s \n", spreadSheetID)
		panic(err)
	}

	sheetName := ""
	for _, v := range response1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetID) {
			sheetName = prop.Title
			break
		}
	}
	return sheetName
}

func NewRepository(service sheets.Service) Repository {
	return &repositoryImpl{
		service: service,
	}
}
