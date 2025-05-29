package main

import (
	"github.com/thomas-maurice/api/go-vue/pkg/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
