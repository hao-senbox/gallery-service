package decorator

import (
	"fmt"
	"reflect"
	"time"
)

type LoggingDecorator struct {
	wrapped Decorator
}

// Logging là một hàm decorator để thêm chức năng ghi log
func (l *LoggingDecorator) Handle() error {
	start := time.Now()
	err := l.wrapped.Handle()
	duration := time.Since(start)

	structType := reflect.TypeOf(l.wrapped)
	funcName := getFunctionName(l.wrapped)

	// Ghi log thông tin
	fmt.Printf("Handle function(%s) of %s struct in %v\n", funcName, structType.Name(), duration)

	return err
}

// Hàm gốc muốn trang trí
func getFunctionName(i interface{}) string {
	return fmt.Sprintf("%T", i)
}
