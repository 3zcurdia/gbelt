package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
)

func main() {
	metrics := metrics.NewUserMetrics("3zcurdia")
	_ = metrics.GetLanguagesCount(false)
	fmt.Printf("%+v\n", metrics)
}
