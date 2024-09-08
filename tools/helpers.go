package tools

import (
	"fmt"
	"runtime/debug"

	"gorm.io/gorm"
)

func IF(bool2 bool, a interface{}, b interface{}) interface{} {
	if bool2 {
		return a
	}
	return b
}

func WrapDbErr(err error) error {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		fmt.Printf("db err:%v\n", err)
		return fmt.Errorf("db err")
	}
	return err
}

func ErrToPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			fmt.Printf("[runSafe] err:%v\n", err)
		}
	}()
	fn()
}