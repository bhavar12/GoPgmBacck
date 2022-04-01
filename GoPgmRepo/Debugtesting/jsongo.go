package main

import (
	"fmt"
	"time"
)

var input = `
{
"created_at" : "Thu May 31 00:00:01 +0000 2012"
}`

// Timestamp for custom
type Timestamp time.Time

//UnmarshalJSON custom implementation
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	v, err := time.Parse(time.RubyDate, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = Timestamp(v)
	fmt.Println(v)
	return nil
}

func main() {
	var t Timestamp
	tt := &t
	tt.UnmarshalJSON([]byte(input))
}
