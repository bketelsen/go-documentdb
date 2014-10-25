package docdb

import (
	"fmt"
	"io"
	"io/ioutil"
)

func logBody(reader io.Reader) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
