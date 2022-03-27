package log

import "fmt"

func Proxy(res interface{}, err error) interface{} {
	if err != nil {
		fmt.Errorf("error %s", err)
		return nil
	}
	return res
}
