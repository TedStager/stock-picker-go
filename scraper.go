package main

import (
	"fmt"
	"errors"
	"time"
	"github.com/gocolly/colly"
	"strconv"
	"os"
	"encoding/csv"
)

type chart struct {
	sym Symbol
	date []string
	open, close, high, low []float32
}

func scraper(ch chan chart) {
	// read list of symbols to check
	file, err := os.Open("stocks.csv")
	checkError(err)
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	checkError(err)

	// set up collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)
	c.AllowURLRevisit = true
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       3 * time.Second,
		RandomDelay: 2 * time.Second,
	})

	// set up callbacks and pointers
	symc := new(comb_sym) // symc for symbol (combined)
	symb := new(sym) // symb to diff from sym (type)
	dat := new(chart)
	
	c.OnHTML("tbody", func(el *colly.HTMLElement) {
		if symc.symbol != "" { 
			*dat = parseTable(el, *symc)
		} else {
			*dat = parseTable(el, *symb)
		}
	})

	// loop through symbols
	for _, line := range lines {
		var err error

		if line[1] != "" { // combined symbol
			*symc = comb_sym{line[0], line[1]}
			*symb = sym{""} // placeholder so we know to look at other ptr

			err = getDatFromSymbol(c, *symc) 
		} else {
			*symb = sym{line[0]}
			*symc = comb_sym{"",""}

			err = getDatFromSymbol(c, *symb) 
		}

		checkError(err)
		ch <- *dat
	}

	close(ch)
}

func getDatFromSymbol(c *colly.Collector, stock Symbol, days ...int64) (error) {
	now := time.Now().Unix()
	var day int64 = 60 * 60 * 24
	var span timeframe

	// days is 0, 1 or 2 ints
	switch num_args := len(days); num_args {

	case 0: // default is 1 year data from today
		span = timeframe{(now - 365 * day), now}
	case 1: // next is 1 year from an arbitrary day ago
		span = timeframe{now - (365 + days[0]) * day, now - days[0] * day}
	case 2: // last is any 2 days ago
		span = timeframe{now - days[0] * day, now - days[1] * day}
	default: // invalid input
		return errors.New("Only 0, 1, or 2 timestamps are accepted as the timeframe.")
	}

	fmt.Printf("Searching for %s...\n", stock.GetString())
	c.Visit(stock.GetURL(span))
	return nil // no error
}

func parseTable(t *colly.HTMLElement, stock Symbol) (chart) {
	dat := chart{
		sym: stock,
	}

	// loop through table rows
	t.ForEach("tr", func (idx int, row *colly.HTMLElement) {
		if row.ChildText("td:nth-child(5)") == "" { return }

		dat.date = append(dat.date, row.ChildText("td:nth-child(1)"))

		var price float64
		var err error

		price, err = strconv.ParseFloat(row.ChildText("td:nth-child(3)"), 32)
		checkError(err)
		dat.high = append(dat.high, float32(price))

		price, err = strconv.ParseFloat(row.ChildText("td:nth-child(4)"), 32)
		checkError(err)
		dat.low = append(dat.low, float32(price))

		price, err = strconv.ParseFloat(row.ChildText("td:nth-child(2)"), 32)
		checkError(err)
		dat.open = append(dat.open, float32(price))

		price, err = strconv.ParseFloat(row.ChildText("td:nth-child(5)"), 32)
		checkError(err)
		dat.close = append(dat.close, float32(price))

	})

	return dat
}
