package main

import (
	"encoding/json"
	"fmt"

	// "io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/johnfercher/maroto/v2"
	// "github.com/johnfercher/maroto/v2/pkg/components/page"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"

	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/core"

	// "github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func main() {
	m := GetMaroto()

	document, err := m.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = document.Save("docs/pdf/platnova-invoice.pdf")
	if err != nil {
		log.Fatal(err.Error())
	}
}

var background = &props.Color{
	Red:   200,
	Green: 200,
	Blue:  200,
}

type Data struct {
	Date               string  `json:"date"`
	PaymentDescription string  `json:"paymentDescription"`
	TransactionType    string  `json:"transactionType"`
	Amount             float64 `json:"amount"`
	Balance            float64 `json:"balance"`
}

func GetMaroto() core.Maroto {
	cfg := config.NewBuilder().
		WithPageNumber("{current} / {total}", props.Bottom).
		Build()

	mrt := maroto.New(cfg)
	m := maroto.NewMetricsDecorator(mrt)

	m.AddRow(20,
		image.NewFromFileCol(2, "images/platnova2.jpg", props.Rect{
			// Center:  true,
			Percent: 80,
		}),
		col.New(6),
		image.NewFromFileCol(4, "images/platnova.png", props.Rect{
			Center:  true,
			Percent: 80,
		}),
	)

	m.AddRow(15,
		col.New(9),
		text.NewCol(3, "Contact tel 03457404404 Text phone 03451235674 used by handicap customers www.platnova.com", props.Text{
			Align: align.Right,
			Size:  10,
		}),
	)
	m.AddRow(2)
	m.AddRow(4,
		col.New(9),
		text.NewCol(3, "Your Statement", props.Text{
			Align: align.Center,
			Size:  14,
			Color: &props.RedColor,
		}),
	)
	m.AddRow(3)
	m.AddRow(30,
		text.NewCol(2, "Mrs I.P Jones April Cottage, 18 East Street, Turners Hill Crawley West Sussex. RH10 4PU", props.Text{
			Align: align.Left,
			Size:  12,
		}),
		col.New(10),
	)
	m.AddRow(5)
	m.AddRow(5,
		text.NewCol(3, getFormattedDateRange(), props.Text{
			Align: align.Left,
			Size:  12,
			Style: fontstyle.Bold,
		}),
		col.New(4),
		text.NewCol(5, "International Bank Account Number", props.Text{
			Align: align.Left,
			Size:  12,
			Style: fontstyle.Bold,
		}),
	)
	m.AddRow(5,
		col.New(7),
		text.NewCol(5, "HGFGHG75434JB5532HHGFG", props.Text{
			Align: align.Left,
			Size:  12,
		}),
	)
	m.AddRow(3)
	m.AddRow(5,
		col.New(7),
		text.NewCol(5, "Bank Identification Code", props.Text{
			Align: align.Left,
			Size:  12,
			Style: fontstyle.Bold,
		}),
	)
	m.AddRow(5,
		col.New(7),
		text.NewCol(5, "JB5532HHGFG", props.Text{
			Align: align.Left,
			Size:  12,
		}),
	)
	m.AddRow(6)
	m.AddRow(5,
		text.NewCol(5, "Account Name", props.Text{
			Align: align.Left,
			Size:  10,
			Style: fontstyle.Bold,
		}),
		col.New(1),
		text.NewCol(2, "Sortcode", props.Text{
			Align: align.Left,
			Size:  10,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, "Act Number", props.Text{
			Align: align.Left,
			Size:  10,
			Style: fontstyle.Bold,
		}),
		text.NewCol(2, "Sheet Number", props.Text{
			Align: align.Left,
			Size:  10,
			Style: fontstyle.Bold,
		}),
	)
	m.AddRow(5,
		text.NewCol(5, "Mrs Lisa Patricia Jones", props.Text{
			Align: align.Left,
			Size:  10,
		}),
		col.New(1),
		text.NewCol(2, "40-18-22", props.Text{
			Align: align.Left,
			Size:  10,
		}),
		text.NewCol(2, "7274792", props.Text{
			Align: align.Left,
			Size:  10,
		}),
		text.NewCol(2, "314", props.Text{
			Align: align.Left,
			Size:  10,
		}),
	)
	m.AddRow(5)
	m.AddRow(2,
		line.NewCol(12),
	)
	m.AddRow(2,
		text.NewCol(12, "Your Platnova Advance details", props.Text{
			Color: &props.BlueColor,
			Size:  15,
			Align: align.Center,
			Style: fontstyle.Bold,
		}),
	)
	m.AddRow(8)

	datas, err := getDatasFromJSON("json/account_statement.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	rows, err := list.Build[Data](datas)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.AddRows(rows...)

	return m
}

func getFormattedDateRange() string {
	today := time.Now()
	oneWeekAgo := today.AddDate(0, 0, -7)

	todayFormatted := today.Format("2nd January 2006")
	oneWeekAgoFormatted := oneWeekAgo.Format("2th January 2006")

	return fmt.Sprintf("%s - %s", oneWeekAgoFormatted, todayFormatted)
}

func (d Data) GetHeader() core.Row {
	return row.New(12).Add(
		text.NewCol(2, "Date", props.Text{Style: fontstyle.Bold, Size: 12}),
		text.NewCol(3, "Payment Description", props.Text{Style: fontstyle.Bold, Size: 12}),
		text.NewCol(3, "TransactionType", props.Text{Style: fontstyle.Bold, Size: 12}),
		text.NewCol(2, "Amount", props.Text{Style: fontstyle.Bold, Size: 12}),
		text.NewCol(2, "Balance", props.Text{Style: fontstyle.Bold, Size: 12}),
	)
}

func (d Data) GetContent(i int) core.Row {
	r := row.New(10).Add(
		text.NewCol(2, d.Date, props.Text{Align: align.Center, Top: 3}),
		text.NewCol(3, d.PaymentDescription, props.Text{Align: align.Center, Top: 3}),
		text.NewCol(3, d.TransactionType, props.Text{Align: align.Center, Top: 3}),
		text.NewCol(2, strconv.FormatFloat(d.Amount, 'f', -1, 64), props.Text{Align: align.Center, Top: 3}),
		text.NewCol(2, strconv.FormatFloat(d.Balance, 'f', -1, 64), props.Text{Align: align.Center, Top: 3}),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: background,
		})
	}

	return r
}

func getDatasFromJSON(filename string) ([]Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var datas []Data
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&datas)
	if err != nil {
		return nil, err
	}

	return datas, nil
}
