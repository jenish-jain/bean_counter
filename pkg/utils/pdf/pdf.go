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

func newBuilder() core.Maroto {
	cnf := config.NewBuilder().
		WithPageNumber("Page {current} of {total}", props.RightBottom).
		WithMargins(12, 15, 10).
		WithAuthor("bean counter", false).
		Build()

	mrt := maroto.New(cnf)
	m := maroto.NewMetricsDecorator(mrt)

	return m
}

func NewPDFGenerator() PDFBuilder {
	return &pdfBuilderImpl{
		builder: newBuilder(),
	}
}
