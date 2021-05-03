package main

import (
	"context"
	"log"
	"github.com/a1ishm/search/pkg/search"
)

func main() {
	context := context.Background()
	ch := search.All(context, "ayo", []string{"files/text1.txt", "files/text2.txt", "files/text3.txt"})
	for val := range ch {
		log.Print(val)
	}
}
