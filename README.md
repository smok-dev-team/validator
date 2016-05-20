## Validator
这是一个用 Golang 写的数据验证组件,主要用于验证 Struct 的属性是否满足特定的规范。

Golang 已经存在很多类似的第三方组件,但是总感觉使用起来不够灵活,尤其是在开发 API 接口的应用程序中。[Validator](https://github.com/smartwalle/going/tree/master/validator) 这是我之前写的一个,使用起来总觉得不方便。

新版本的 Validator 将数据验证的权限完全将由者处理，开发者为每一个需要验证的属性写相关的验证方法，例：

``` 
import (
	"github.com/smartwalle/validator"
	"github.com/smartwalle/errors"
)

type Human struct {
	Name string
	Age  int
}

// 为 Human 需要验证的属性添加方法
func (this Human) NameValidator(n string) error {
	if n == "" {
		// 可以只返回一个 error
		return errors.NewWithCode(1001, "请提供你的名字哦")
	}
	return nil
}

func (this Human) AgeValidator(a int) []error {
	// 也可以返回一个 error slice
	var errList = make([]error, 0, 0)
	if a <= 0 {
		errList = append(errList, errors.NewWithCode(1002, "你确定这是你的年龄？"))
	}

	if a > 100 {
		errList = append(errList, errors.NewWithCode(1003, "你也太长命了吧"))
	}
	return errList
}


// 验证
var h Human
var v = validator.Validate(h)
if !v.OK() {
	fmt.Println("抱歉，验证没有通过")
}


```

如上所示，验证方法命名规则为：属性名+Validator（如: NameValidator），方法有唯一的参数，即对应属性的值，方法需要返回一个满足 error 接口的对象，如果返回 nil，则表示该验证通过。