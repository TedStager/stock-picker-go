package main

import "fmt"

func picker(ch chan chart) (hits []string) {
	var dat chart
	var acc float32 = 0.02

	// loop infinitely while channel is open
	for dat = range ch {
		fmt.Println("Got", dat.sym.GetString())

		avg_val := avg(dat, 200)
		margin := dat.close[0] * acc
		if (avg_val > dat.close[0] - margin) && (avg_val < dat.close[0] + margin) {
			hits = append(hits, dat.sym.GetString())
		}
	}

	return
}

func avg(ch chart, days int) float32 {
	var sum float32 = 0
	for i := 0; i < days; i++ {
		sum += ch.close[i]
	}

	return sum / float32(days)
}