package log

import "fmt"

func Proxy(res interface{}, err error) interface{} {
	if err != nil {
		fmt.Println("error", err)
		return nil
	}
	return res
}
