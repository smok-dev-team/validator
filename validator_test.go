package validator

import (
	"testing"
	"fmt"
	"github.com/smartwalle/errors"
)

////////////////////////////////////////////////////////////////////////////////////////////////////
type Human struct {
	Name string
	Age  int
}

func (this *Human) NameValidator(n string) error {
	if n == "" {
		return errors.NewWithCode(1001, "请输入名字")
	}
	return nil
}

func (this Human) AgeValidator(a int) error {
	if a <= 0 {
		return errors.NewWithCode(1002, "请输入年龄")
	}

	if a > 100 {
		return errors.NewWithCode(1003, "你也太长命了吧")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func TestValidator(t *testing.T) {
	var h *Human = &Human{}
	var r = Validate(&h)
	if !r.OK() {
		var e = r.Error()
		fmt.Println(errors.ErrorCode(e), errors.ErrorMessage(e))
		fmt.Println(r.ErrorList())
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////

type Student struct {
	Human
	Number int
}

func (this Student) NameValidator(n string) error {
	return nil
}

func (this Student) NumberValidator(n int) error {
	if n <= 0 {
		return errors.NewWithCode(1004, "应该有一个学号")
	}
	return nil
}

func TestLazyValidator(t *testing.T) {
	var s Student

	var r = Validate(&s)
	if !r.OK() {
		var e = r.Error()
		fmt.Println(errors.ErrorCode(e), errors.ErrorMessage(e))
		fmt.Println(r.ErrorList())
	}
}