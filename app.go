package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
)

func main() {
	um := metrics.NewUserMetrics("3zcurdia")
	_, err := um.FetchLanguagesCount(true)
	if err != nil {
		panic(err)
	}
	fmt.Println(um.Languages)
}
