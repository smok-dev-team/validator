package validator

import (
	"fmt"
	"reflect"
	"errors"
)

const (
	k_VALIDATOR_FUNC_SUFFIX = "Validator"
)

////////////////////////////////////////////////////////////////////////////////
type Validator interface {
	ErrorList() []error
	ErrorMap() map[string][]error
	ErrorListWithField(field string) []error
	Error() error
	OK() bool
}

////////////////////////////////////////////////////////////////////////////////
type validator struct {
	ErrMap    map[string][]error `json:"error_map"`
	errList   []error            `json:"-"`
	fieldList []string           `json:"-"`
	lazy      bool               `json:"-"`
}

func (this *validator) String() string {
	return fmt.Sprintf("[validator]: Valid:%t, Error:%s", this.OK(), this.ErrMap)
}

func (this *validator) ErrorList() []error {
	if this.errList == nil {
		if len(this.ErrMap) > 0 {
			this.errList = make([]error, 0, len(this.fieldList))
			for _, field := range this.fieldList {
				this.errList = append(this.errList, this.ErrMap[field]...)
			}
		}
	}
	return this.errList
}

func (this *validator) ErrorMap() map[string][]error {
	return this.ErrMap
}

func (this *validator) ErrorListWithField(field string) []error {
	return this.ErrMap[field]
}

func (this *validator) Error() error {
	if len(this.ErrorList()) > 0 {
		return this.ErrorList()[0]
	}
	return nil
}

func (this *validator) OK() bool {
	return (this.ErrMap != nil && len(this.ErrMap) == 0)
}

////////////////////////////////////////////////////////////////////////////////
func Validate(obj interface{}) Validator {
	return _validate(obj, false)
}

func LazyValidate(obj interface{}) Validator {
	return _validate(obj, true)
}

func _validate(obj interface{}, lazy bool) Validator {
	var objType = reflect.TypeOf(obj)
	var objValue = reflect.ValueOf(obj)
	var objValueKind = objValue.Kind()

	var val = &validator{}
	val.ErrMap = make(map[string][]error)
	val.lazy = lazy

	for {
		if objValueKind == reflect.Ptr && objValue.IsNil() {
			var errList = make([]error, 1)
			errList[0] = errors.New("object passed for validation is nil")
			val.fieldList = make([]string, 1)
			val.fieldList[0] = "Object"
			val.ErrMap["Object"] = errList
			return val
		}
		if objValueKind == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			objValueKind = objValue.Kind()
			continue
		}
		break
	}

	val.fieldList = make([]string, 0, objType.NumField())

	validate(objType, objValue, objValue, val)
	return val
}

func validate(objType reflect.Type, currentObjValue, objValue reflect.Value, val *validator) {
	var numField = objType.NumField()
	for i := 0; i < numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			validate(fieldValue.Type(), currentObjValue, fieldValue, val)
			if val.lazy && len(val.ErrMap) > 0 {
				return
			}
			continue
		}

		var funcName = fieldStruct.Name + k_VALIDATOR_FUNC_SUFFIX
		var funcValue = getFuncWithName(funcName, currentObjValue, objValue)

		if funcValue.IsValid() {
			var eList = funcValue.Call([]reflect.Value{fieldValue})

			if !eList[0].IsNil() {
				val.fieldList = append(val.fieldList, fieldStruct.Name)
				if eList[0].Kind() == reflect.Slice {
					val.ErrMap[fieldStruct.Name] = eList[0].Interface().([]error)
				} else {
					val.ErrMap[fieldStruct.Name] = []error{eList[0].Interface().(error)}
				}
				if val.lazy {
					return
				}
			}
		}
	}
}

func getFuncWithName(funcName string, currentObjValue, objValue reflect.Value) reflect.Value {
	var funcValue = currentObjValue.MethodByName(funcName)
	if funcValue.IsValid() == false {
		if currentObjValue.CanAddr() {
			funcValue = currentObjValue.Addr().MethodByName(funcName)
		}
	}
	if funcValue.IsValid() == false && currentObjValue != objValue {
		return getFuncWithName(funcName, objValue, objValue)
	}
	return funcValue
}
