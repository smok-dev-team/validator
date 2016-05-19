package form

import (
	"fmt"
	"reflect"
)

const (
	K_VALIDATOR_FUNC_SUFFIX = "Validator"
)

////////////////////////////////////////////////////////////////////////////////
type ValidatorError struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
}

func NewValidatorError(code int, msg string) (*ValidatorError) {
	var err = &ValidatorError{}
	err.Code = code
	err.Message = msg
	return err
}

func (this *ValidatorError) Error() string {
	return fmt.Sprintf("[%d]%s", this.Code, this.Message)
}

////////////////////////////////////////////////////////////////////////////////
type IValidator interface {
	ErrorList()                      []error
	ErrorMap()                       map[string][]error
	ErrorListWithField(field string) []error
	Error()                          error
	OK()                             bool
}

////////////////////////////////////////////////////////////////////////////////
type validator struct {
	errorMap  map[string][]error  `json:"error_map"`
	errorList []error             `json:"-"`
	fieldList []string            `json:"-"`
}

func (this *validator) String() string {
	return fmt.Sprintf("[validator]: Valid:%t, Error:%s", this.OK(), this.errorMap)
}

func (this *validator) ErrorList() []error {
	if this.errorList == nil {
		if len(this.errorMap) > 0 {
			this.errorList = make([]error, 0, len(this.fieldList))
			for _, field := range this.fieldList {
				this.errorList = append(this.errorList, this.errorMap[field]...)
			}
		}
	}
	return this.errorList
}

func (this *validator) ErrorMap() map[string][]error {
	return this.errorMap
}

func (this *validator) ErrorListWithField(field string) []error {
	return this.errorMap[field]
}

func (this *validator) Error() error {
	if len(this.ErrorList()) > 0 {
		return this.ErrorList()[0]
	}
	return nil
}

func (this *validator) OK() bool {
	return (len(this.errorMap) == 0)
}

////////////////////////////////////////////////////////////////////////////////
func Validate(obj interface{}) (IValidator) {
	var objType = reflect.TypeOf(obj)
	var objValue = reflect.ValueOf(obj)

	for {
		if objValue.Kind() == reflect.Ptr {
			objValue = objValue.Elem()
			objType = objType.Elem()
			continue
		}
		break
	}

	if !objValue.IsValid() {
		return nil
	}


	var val = &validator{}
	val.errorMap = make(map[string][]error)
	val.fieldList = make([]string, 0, objType.NumField())

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
			continue
		}

		var mName  = fieldStruct.Name + K_VALIDATOR_FUNC_SUFFIX
		var mValue = objValue.MethodByName(mName)
		if mValue.IsValid() {
			var eList = mValue.Call([]reflect.Value{fieldValue})

			if !eList[0].IsNil() {
				val.fieldList = append(val.fieldList, fieldStruct.Name)
				if eList[0].Kind() == reflect.Slice {
					val.errorMap[fieldStruct.Name] = eList[0].Interface().([]error)
				} else {
					val.errorMap[fieldStruct.Name] = []error{eList[0].Interface().(error)}
				}
			}
		}
	}
}