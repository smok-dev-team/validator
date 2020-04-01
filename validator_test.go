package validator

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

type Human struct {
	Name string
	Age  int
}

func (this *Human) NameValidator(n string) error {
	if n == "" {
		return errors.New("请输入名字")
	}
	return nil
}

func (this Human) AgeValidator(a int) error {
	if a <= 0 {
		return errors.New("请输入年龄")
	}

	if a > 100 {
		return errors.New("你也太长命了吧")
	}
	return nil
}

func TestValidate(t *testing.T) {
	var h = &Human{}
	var r = Validate(&h)
	if !r.Passed() {
		fmt.Println(r.ErrorList())
	}
}

type Student struct {
	Human  *Human
	Number int
	Time   *time.Time
}

func (this Student) NameValidator(n string) error {
	return errors.New("不管输入什么，都不会通过的")
}

func (this Student) NumberValidator(n int) error {
	if n <= 0 {
		return errors.New("应该有一个学号")
	}
	return nil
}

func (this Student) TimeValidator(p *time.Time) error {
	if p == nil {
		return errors.New("oh no")
	}
	return nil
}

func TestLazyValidate(t *testing.T) {
	var s = &Student{Human: &Human{}}

	var r = LazyValidate(&s)
	if !r.Passed() {
		fmt.Println(r.ErrorList())
	}
}
