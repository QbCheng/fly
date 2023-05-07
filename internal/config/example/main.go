package main

import (
	"fly/internal/config"
	"fmt"
)

func main() {
	//file := config.FilePath()
	file := "D:\\meta\\luckyshake\\luckyshake-server\\internal\\config\\test.yaml"
	_, c, err := config.NewViper(file)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v", c)
}
