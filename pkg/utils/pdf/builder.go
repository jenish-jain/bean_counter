package pdf

import (
	"fmt"
	"log/slog"

	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func (p *pdfBuilderImpl) Build() error {
	document, err := p.builder.Generate()
	if err != nil {
		slog.Error("error generating pdf document: %+v", err)
	}
	err = document.Save("docs/assets/pdf/billingv2.pdf")
	if err != nil {
		slog.Error("error saving pdf document: %+v", err)
	}

	// err = document.GetReport().Save("docs/assets/text/billingv2.txt")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	return err

}

func (p *pdfBuilderImpl) WithHeader() PDFBuilder {
	p.builder.RegisterHeader(getPageHeader())
	return p
}

func (p *pdfBuilderImpl) WithFooter() PDFBuilder {
	p.builder.RegisterFooter(getPageFooter())
	return p
}

func (p *pdfBuilderImpl) AddRow(value string) PDFBuilder {
	// TODO: take these styling param as a config in builder input later
	p.builder.AddRows(text.NewRow(ROW_HEIGHT, value, BOLD_GREY_TEXT, CENTER_ALIGN))
	return p
}

func (p *pdfBuilderImpl) AddTable(headers []string, rows [][]interface{}, indexed bool) PDFBuilder {
	tableRows := []core.Row{}
	columns := []core.Col{}

	if len(headers) > 0 {
		if indexed {
			headers = append([]string{"S.No."}, headers...)
		}
		for header := range headers {
			column := text.NewCol(COLUMN_SIZE, headers[header], BOLD_WHITE_CENTERED_TEXT).
				WithStyle(TABLE_HEADER_STYLE)
			columns = append(columns, column)
		}

		p.builder.AddRows(row.New(TABLE_ROW_HEIGHT).Add(columns...))
	}

	for rowNo, rowElement := range rows {
		rowColumns := []core.Col{}
		if indexed {
			rowElement = append([]interface{}{rowNo + 1}, rowElement...)
		}
		if len(headers) != 0 && (len(rowElement) != len(headers)) {
			slog.Error("length of row %d does not match headers %d", len(rowElement), len(headers))
			panic("length of row does not match headers")
		}
		for _, rowContent := range rowElement {
			rowElement := text.NewCol(COLUMN_SIZE, fmt.Sprint(rowContent), NORMAL_GREY_CENTERED_TEXT).
				WithStyle(TABLE_ROW_STYLE)
			rowColumns = append(rowColumns, rowElement)
		}
		tableRows = append(tableRows, row.New(TABLE_ROW_HEIGHT).Add(rowColumns...))
	}
	p.builder.AddRows(tableRows...)

	return p
}

func getPageFooter() core.Row {
	return row.New(15).Add(
		col.New(12).Add(
			text.New("provided to you by munshi ji", props.Text{
				Top:   20,
				Style: fontstyle.BoldItalic,
				Size:  6,
				Align: align.Center,
				Color: DARK_GREY_COLOR,
			}),
		),
	)
}

func getPageHeader() core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(3, "docs/assets/images/biplane.jpg", props.Rect{
			Center:  true,
			Percent: 80,
		}),
		col.New(6),
		col.New(3).Add(
			text.New("Jainco Textiles, 101 sarvodaya textile market, ring road. gujarat. surat", props.Text{
				Size:  8,
				Align: align.Right,
				Color: RED_COLOR,
			}),
			text.New("Tel: 9825120344", props.Text{
				Top:   10,
				Style: fontstyle.BoldItalic,
				Size:  8,
				Align: align.Right,
				Color: BLUE_COLOR,
			}),
		),
	)
}
