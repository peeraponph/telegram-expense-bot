package sheet

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func WriteRow(row []interface{}) error {
	ctx := context.Background()

	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	if spreadsheetID == "" {
		return fmt.Errorf("ไม่พบค่า SPREADSHEET_ID")
	}

	// โหลด credentials.json
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		return fmt.Errorf("อ่าน credentials.json ไม่ได้: %v", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(b))
	if err != nil {
		return fmt.Errorf("เชื่อม Google Sheets ไม่ได้: %v", err)
	}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	_, err = srv.Spreadsheets.Values.Append(spreadsheetID, "Expenses!A:F", vr).
		ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		return fmt.Errorf("เขียนข้อมูลลงชีทไม่สำเร็จ: %v", err)
	}

	return nil
}



// Expenses!A:F