package pdf

import (
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type pdfBuilderImpl struct {
	builder core.Maroto
}

type PDFBuilder interface {
	Build() error
	WithHeader() PDFBuilder
	WithFooter() PDFBuilder
	AddRow(value string) PDFBuilder
	AddTable(headers []string, rows [][]interface{}, indexed bool) PDFBuilder
}

type PDFConfig struct {
	LeftMargin  float64
	RightMargin float64
	TopMargin   float64
	Author      string
}

func newBuilder(pdfConfig PDFConfig) core.Maroto {
	cnf := config.NewBuilder().
		WithPageNumber("Page {current} of {total}", props.RightBottom).
		WithMargins(pdfConfig.LeftMargin, pdfConfig.TopMargin, pdfConfig.RightMargin).
		WithAuthor(pdfConfig.Author, false).
		Build()

	mrt := maroto.New(cnf)
	m := maroto.NewMetricsDecorator(mrt)

	return m
}

func NewPDFGenerator(config PDFConfig) PDFBuilder {
	return &pdfBuilderImpl{
		builder: newBuilder(config),
	}
}
