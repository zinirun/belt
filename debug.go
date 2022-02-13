package belt

import (
	"fmt"
	"reflect"
)

func ListItems(items []I) {
	for idx, item := range items {
		fmt.Println(idx, ": ", item)
	}
}

func Type(item I) reflect.Type {
	return reflect.TypeOf(item)
}
