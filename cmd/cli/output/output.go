package output

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	yellow = "\033[33m"
)

// Print to the output with added newline.
func Println(data any) {
	var msg []byte
	switch v := data.(type) {
	case []byte:
		msg = v
	case string:
		msg = []byte(v)
	default:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			msg = []byte(fmt.Sprintf("%v", data))
		} else {
			msg = b
		}
	}
	os.Stdout.Write(append(msg, '\n'))
}

// PrintlnErr prints to the output in red with added newline.
func PrintlnErr(data any) {
	var msg []byte
	switch v := data.(type) {
	case []byte:
		msg = v
	case error:
		msg = []byte(v.Error())
	}

	msg = append([]byte(red), msg...)
	msg = append(msg, []byte(reset)...)
	msg = append(msg, '\n')
	os.Stderr.Write(msg)
}
