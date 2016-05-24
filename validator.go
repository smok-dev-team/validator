package validator

import (
	"fmt"
	"reflect"
)

const (
	k_VALIDATOR_FUNC_SUFFIX = "Validator"
)

////////////////////////////////////////////////////////////////////////////////
type Validator interface {
	ErrorList()                      []error
	ErrorMap()                       map[string][]error
	ErrorListWithField(field string) []error
	Error()                          error
	OK()                             bool
}

////////////////////////////////////////////////////////////////////////////////
type validator struct {
	ErrMap    map[string][]error  `json:"error_map"`
	errList   []error             `json:"-"`
	fieldList []string            `json:"-"`
	lazy      bool                `json:"-"`
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
func Validate(obj interface{}) (Validator) {
	return _validate(obj, false)
}

func LazyValidate(obj interface{}) (Validator) {
	return _validate(obj, true)
}

func _validate(obj interface{}, lazy bool) (Validator) {
	var objType = reflect.TypeOf(obj)
	var objValue = reflect.ValueOf(obj)
	var objValueKind = objValue.Kind()

	for {
		if objValueKind == reflect.Ptr && objValue.IsNil() {
			panic("object passed for validation is nil")
		}
		if objValueKind == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			objValueKind = objValue.Kind()
			continue
		}
		break
	}

	var val = &validator{}
	val.ErrMap = make(map[string][]error)
	val.fieldList = make([]string, 0, objType.NumField())
	val.lazy = lazy

	validate(objType, objValue, val)
	return val
}

func validate(objType reflect.Type, objValue reflect.Value, val *validator) {
	var numField = objType.NumField()
	for i:=0; i<numField; i++ {
		var fieldStruct = objType.Field(i)
		var fieldValue = objValue.Field(i)

		if fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Struct {
			validate(fieldValue.Type(), fieldValue, val)
			if val.lazy && len(val.ErrMap) > 0 {
				return
			}
			continue
		}

		var mName  = fieldStruct.Name + k_VALIDATOR_FUNC_SUFFIX
		var mValue = objValue.MethodByName(mName)

		if mValue.IsValid() == false {
			if objValue.CanAddr() {
				mValue = objValue.Addr().MethodByName(mName)
			}
		}

		if mValue.IsValid() {
			var eList = mValue.Call([]reflect.Value{fieldValue})

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