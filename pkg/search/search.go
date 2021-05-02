package search

import (
	"context"
)

// Result описывает один результат поиска
type Result struct {
	// Фраза, которую искали
	Phrase  string
	// Целиком вся строка, в которой нашли вхождение (без \n или \r\n в конце)
	Line    string
	// Номер строки (начиная с 1), на которой нашли вхождение
	LineNum int64
	// Номер позиции (начиная с 1), на которой нашли вхождение
	ColNum  int64
}

// All ищет все вхождения phrase в текстовых файлах files
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)

	// TODO: some code

	return ch
}