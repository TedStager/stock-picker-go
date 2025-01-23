// file for various interfaces and structs

package main

import "fmt"

type timeframe struct {
	start, end int64
}

// symbol stuff

type Symbol interface {
	GetURL(tf timeframe) string
	GetString() string
}

// combined symbol like XSP.TO

type comb_sym struct { // symbol and market
	symbol, market string
}

func (sym comb_sym) GetURL(tf timeframe) string {
	return fmt.Sprintf(
		"https://ca.finance.yahoo.com/quote/%s.%s/history/?period1=%d&period2=%d", 
		sym.symbol, sym.market,
		tf.start, tf.end,
	)
}

func (sym comb_sym) GetString() string {
	return fmt.Sprintf("%s.%s", sym.symbol, sym.market)
}

// lone symbol, like AMZN

type sym struct {
	symbol string
}

func (sym sym) GetURL(tf timeframe) string {
	return fmt.Sprintf(
		"https://ca.finance.yahoo.com/quote/%s/history/?period1=%d&period2=%d",
		sym.symbol,
		tf.start, tf.end,
	)
}

func (sym sym) GetString() string {
	return sym.symbol
}