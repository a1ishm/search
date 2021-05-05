package search

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	// "time"
)

// Result описывает один результат поиска
type Result struct {
	// Фраза, которую искали
	Phrase string
	// Целиком вся строка, в которой нашли вхождение (без \n или \r\n в конце)
	Line string
	// Номер строки (начиная с 1), на которой нашли вхождение
	LineNum int64
	// Номер позиции (начиная с 1), на которой нашли вхождение
	ColNum int64
}

// All ищет все вхождения phrase в текстовых файлах files
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	fileExist := true

	for _, file := range files {
		wg.Add(1)
		go func(ctx context.Context, path string, substr string, ch chan<- []Result) {
			defer wg.Done()
			filePath, err := filepath.Abs(path)
			if err != nil {
				log.Print(err)
				return
			}

			content, err := os.Open(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					fileExist = false
				} else {
					log.Print(err)
					return
				}
			}

			lineNum := 0
			results := []Result{}
			var scanner *bufio.Scanner

			if fileExist {
				scanner = bufio.NewScanner(content)
			}

			for scanner.Scan() {
				if !fileExist {
					break
				}
				line := scanner.Text()

				lineNum++
				if strings.Contains(line, substr) {
					res := Result{
						Phrase:  substr,
						Line:    strings.Trim(line, "\r\n"),
						LineNum: int64(lineNum),
						ColNum:  int64(strings.Index(line, substr) + 1),
					}
					results = append(results, res)
				}
			}

			if len(results) == 0 {
				return
			}
			ch <- results
		}(ctx, file, phrase, ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()
	cancel()

	return ch
}

// Any ищет любое одно вхождение phrase в текстовых files
func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan Result)
	fileExist := true
	ctx, cancel := context.WithCancel(ctx)
	wg := sync.WaitGroup{}

	for _, file := range files {
		wg.Add(1)
		go func(ctx context.Context, path string, substr string, ch chan<- Result) {
			defer wg.Done()
			filePath, err := filepath.Abs(path)
			if err != nil {
				log.Print(err)
				cancel()
			}

			content, err := os.Open(filePath)
			if err != nil {
				if os.IsNotExist(err) {
					fileExist = false
				} else {
					log.Print(err)
					cancel()
				}
			}

			lineNum := 0
			res := Result{}
			var scanner *bufio.Scanner

			if fileExist {
				scanner = bufio.NewScanner(content)
			}

			for scanner.Scan() {
				if !fileExist {
					break
				}
				line := scanner.Text()

				lineNum++
				select {
				case <-ctx.Done():
					return
				default:
					if strings.Contains(line, substr) {
						res.Phrase = substr
						res.Line = strings.Trim(line, "\r\n")
						res.LineNum = int64(lineNum)
						res.ColNum = int64(strings.Index(line, substr) + 1)

						ch <- res
						cancel()
					}
				}
			}
		}(ctx, file, phrase, ch)
	}

	go func() {
		defer close(ch)
		wg.Wait()
		cancel()
	}()

	return ch
}
