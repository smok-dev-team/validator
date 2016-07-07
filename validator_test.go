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
		return errors.NewWithCode(1001, "name pls")
	}
	return nil
}

func (this Human) AgeValidator(a int) error {
	if a <= 0 {
		return errors.NewWithCode(1002, "age pls")
	}

	if a > 100 {
		return errors.NewWithCode(1003, "你也太长命了吧")
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////
func TestValidator(t *testing.T) {
	var h *Human = nil
	var r = Validate(&h)
	if !r.OK() {
		var e = r.Error()
		fmt.Println(errors.ErrorCode(e), errors.ErrorMessage(e))
		fmt.Println(r.ErrorList())
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////////////
//func TestLazyValidator(t *testing.T) {
//	var h Human
//
//	var r = LazyValidate(&h)
//	if !r.OK() {
//		var e = r.Error()
//		fmt.Println(errors.ErrorCode(e), errors.ErrorMessage(e))
//		fmt.Println(r.ErrorList())
//	}
//}