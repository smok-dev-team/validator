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

func (this Human) NameValidator(n string) error {
	if n == "" {
		return errors.NewWithCode(1001, "请提供你的名字哦")
	}
	return nil
}

func (this Human) AgeValidator(a int) error {
	if a <= 0 {
		return errors.NewWithCode(1002, "你确定这是你的年龄？")
	}

	if a > 100 {
		return errors.NewWithCode(1003, "你也太长命了吧")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func TestValidator(t *testing.T) {
	var h Human
	h.Name = "这是我的名字"

	var r = Validate(h)
	if !r.OK() {
		var e = r.Error()
		fmt.Println(errors.Code(e), errors.Message(e))
	}
}