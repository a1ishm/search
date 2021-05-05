package main

import (
	"context"
	"log"
	// "sync"
	"github.com/a1ishm/search/pkg/search"
)

func main() {
	context := context.Background()
	ch := search.Any(context, "ayo", []string{"files/text1.txt", "files/text2.txt", "files/text3.txt"})

	for v := range ch {
		log.Print(v)
	}
}
