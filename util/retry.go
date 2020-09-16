package util

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func Retry(retry int, duration time.Duration, fn interface{}, args ...interface{}) error {
	if err := checkFn(fn); err != nil {
		return err
	}

	callBack := reflect.ValueOf(fn)
	passedArgs := setUpArgs(args ...)
	if err := doFunc(callBack, passedArgs);err == nil {
		return nil
	}

	var err error
	for i := 1; i < retry; i++ {
		time.Sleep(computeDuration(i, duration))
		if err = doFunc(callBack, passedArgs);err == nil {
			return nil
		}
	}
	return err
}

func RetryAsync(retry int, duration time.Duration, fn interface{}, args ...interface{}) {
	go func() {
		if err := Retry(retry, duration, fn, args ...); err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
			fmt.Println(fmt.Sprintf("%s retry %d times failed, args:%s, err:%v", funcName, retry, fmt.Sprint(args), err))
		}
	}()
}

func doFunc(callBack reflect.Value, passedArgs []reflect.Value) error {
	out := callBack.Call(passedArgs)[0]
	err, _ := out.Interface().(error)
	return err
}


//fibonacci
func computeDuration(cur int, duration time.Duration) time.Duration{
	if cur <= 1 {
		return duration
	}
	a, b := duration, duration
	for i := 1; i < cur; i++ {
		a = a + b
		a, b = b, a
	}
	return b
}



func setUpArgs(args ...interface{}) []reflect.Value {
	passedArgs := make([]reflect.Value, 0)
	for _, arg := range args {
		passedArgs = append(passedArgs, reflect.ValueOf(arg))
	}
	return passedArgs
}

func checkFn(fn interface{}) error {
	tp := reflect.TypeOf(fn)
	if tp.Kind() != reflect.Func {
		return notFuncError()
	}

	if tp.NumOut() != 1 || !tp.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return signatureError()
	}

	return nil
}

func notFuncError() error {
	return errors.New("fn is not func")
}

func signatureError() error {
	return errors.New("only [func(args ...interface{}) error] support")
}