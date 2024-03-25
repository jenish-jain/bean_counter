package pdf

import (
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

const (
	ROW_HEIGHT       = 4
	COLUMN_SIZE      = 2
	TABLE_ROW_HEIGHT = 7
)

// Colors
var (
	DARK_GREY_COLOR = &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
	GREY_COLOR = &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
	WHITE_COLOR = &props.Color{
		Red:   255,
		Green: 255,
		Blue:  255,
	}
	RED_COLOR = &props.Color{
		Red:   255,
		Green: 0,
		Blue:  0,
	}
	BLUE_COLOR = &props.Color{
		Red:   0,
		Green: 0,
		Blue:  255,
	}
)

// Text styles
var (
	BOLD_WHITE_CENTERED_TEXT = props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: WHITE_COLOR,
		Size:  7,
	}

	NORMAL_GREY_CENTERED_TEXT = props.Text{
		Top:   3,
		Style: fontstyle.Normal,
		Align: align.Center,
		Color: DARK_GREY_COLOR,
		Size:  7,
	}

	BOLD_GREY_CENTERED_TEXT = props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Align: align.Center,
		Color: DARK_GREY_COLOR,
		Size:  7,
	}

	BOLD_GREY_TEXT = props.Text{
		Top:   3,
		Style: fontstyle.Bold,
		Color: DARK_GREY_COLOR,
		Size:  7,
	}

	CENTER_ALIGN = props.Text{
		Align: align.Center,
	}
)

// Table Formattings
var (
	TABLE_HEADER_STYLE = &props.Cell{
		BackgroundColor: GREY_COLOR,
		BorderType:      border.Full,
		BorderThickness: 0.3,
	}

	TABLE_ROW_STYLE = &props.Cell{
		BackgroundColor: WHITE_COLOR,
		BorderType:      border.Full,
		BorderThickness: 0.1,
	}
)
