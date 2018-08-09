package repl

import (
	"bufio"
	"fmt"
	"io"

	"monkey/lexer"
	"monkey/token"
)

// PROMPT REPLのプロンプト文字列
const PROMPT = ">> "

// Start REPL
func Start(in io.Reader, out io.Writer) {
	// bufio Reader/Writerを引数に取りバッファリング用機能を追加したReader/Writerを返す
	// bufio.Scanner 文字列を特定の区切り文字で区切るようにバッファリングを行う デフォルトは改行区切り
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
