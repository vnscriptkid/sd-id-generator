package main

import (
	"fmt"

	"github.com/vnscriptkid/sd-id-generator/twitter-snowflake/lib"
)

func main() {
	flake, err := lib.NewSnowflake(1, 1)
	if err != nil {
		panic(err)
	}

	prev, err := flake.NextID()
	if err != nil {
		panic(err)
	}

	for i := 1; i < 10; i++ {
		cur, err := flake.NextID()
		if err != nil {
			panic(err)
		}

		if cur <= prev {
			panic("cur <= prev")
		}

		fmt.Println(cur)

		prev = cur
	}
}
