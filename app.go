package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
)

func main() {
	metrics := metrics.NewUserMetrics("3zcurdia")
	_ = metrics.GetLanguagesCount(true)
	fmt.Printf("%+v\n", metrics)
}
