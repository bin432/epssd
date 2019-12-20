package epss

import (
	"fmt"
)

var d = def{}

type def struct{}

func (d def) Debug(v ...interface{}) {
	fmt.Println(v...)
}
func (d def) Error(v ...interface{}) {
	fmt.Println(v...)
}
