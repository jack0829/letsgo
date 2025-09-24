package tfidf

import "fmt"

var Debug bool

func debug(tpl string, v ...any) {
	if Debug {
		fmt.Printf(tpl, v...)
	}
}
