# In-depth-GO-reflection
In-depth GO reflection 深入Go反射

深入Go反射
### reflect.Type中的方法：
`func (t *rtype) String() string` // 获取 t 类型的字符串描述，不要通过 String 来判断两种类型是否一致。

`func (t *rtype) Name() string` // 获取 t 类型在其包中定义的名称，未命名类型则返回空字符串。

`func (t *rtype) PkgPath() string` // 获取 t 类型所在包的名称，未命名类型则返回空字符串。

`func (t *rtype) Kind() reflect.Kind` // 获取 t 类型的类别。

`func (t *rtype) Size() uintptr` // 获取 t 类型的值在分配内存时的大小，功能和 unsafe.SizeOf 一样。

`func (t *rtype) Align() int`  // 获取 t 类型的值在分配内存时的字节对齐值。

`func (t *rtype) FieldAlign() int`  // 获取 t 类型的值作为结构体字段时的字节对齐值。

`func (t *rtype) NumMethod() int`  // 获取 t 类型的方法数量。

`func (t *rtype) NumField() int` //返回一个struct 类型 的属性个数，如果非struct类型会抛异常

`func (t *rtype) Method() reflect.Method`  // 根据索引获取 t 类型的方法，如果方法不存在，则 panic。
// 如果 t 是一个实际的类型，则返回值的 Type 和 Func 字段会列出接收者。
// 如果 t 只是一个接口，则返回值的 Type 不列出接收者，Func 为空值。

`func (t *rtype) MethodByName(string) (reflect.Method, bool)` // 根据名称获取 t 类型的方法。

`func (t *rtype) Implements(u reflect.Type) bool` // 判断 t 类型是否实现了 u 接口。

`func (t *rtype) ConvertibleTo(u reflect.Type) bool` // 判断 t 类型的值可否转换为 u 类型。

`func (t *rtype) AssignableTo(u reflect.Type) bool` // 判断 t 类型的值可否赋值给 u 类型。

`func (t *rtype) Comparable() bool` // 判断 t 类型的值可否进行比较操作
//注意对于：数组、切片、映射、通道、指针、接口 
`func (t *rtype) Elem() reflect.Type` // 获取元素类型、获取指针所指对象类型，获取接口的动态类型
有个方法是Elem()，获取元素类型、获取指针所指对象类型，获取接口的动态类型。对指针类型进行反射的时候，可以通过reflect.Elem()获取这个指针指向元素的类型。

### 获取指针的内容
```go
func main() {
    c := data{Nmae:"1",Code:110}
    t1()
}
func t1(obj interface{}) {
    // 这里是值类型
    reflect.TypeOf(obj)  // 获取类型
    reflect.ValueOf(obj) // 获取值
}

func t2(obj interface{}) {
    // 这里是引用类型
    reflect.TypeOf(obj).Elem()
    reflect.ValueOf(obj).Elem()    
}
```

### 遍历获取结构体名称和value
```go
func main() {
	c := data{Name:"1",Code:110}
	t1()
}
func t1(obj interface{}) {
	// 这里以值类型为例  引用类型也如此
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)
	for i := 0;i<typ.NumField();i++ {
		fmt.Println(typ.Field(i).Tag)  // 打印当前行的Tag
		fmt.Println(typ.Field(i).Name) // 当前tag的名称
		fmt.Println(val.Field(i).Interface()) // 获取当前value
		
		// 设置 哈哈你值传递是改不了值的
		val.Field(i).SetString()
	}
}
```

### Reflection Function Map 实现
```
package reflect

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type MapFunc struct {
	db map[string]reflect.Value
	mu sync.Mutex
}

func New() *MapFunc {
	return &MapFunc{
		db: make(map[string]reflect.Value),
	}
}

func (m *MapFunc) Add(key string,fn interface{}) (err error) {
	defer func() {
		if erc := recover();erc != nil {
			err = fmt.Errorf("%v",erc)
			return
		}
	}()
	m.mu.Lock()
	defer m.mu.Unlock()
	of := reflect.ValueOf(fn)
	m.db[key] = of
	return
}

func (m *MapFunc) Call(key string,params ...interface{}) (result []reflect.Value,err error) {
	defer func() {
		if erc := recover();erc != nil {
			err = fmt.Errorf("%v",erc)
			return
		}
	}()
	m.mu.Lock()
	defer m.mu.Unlock()
	value,ex := m.db[key]
	if !ex {
		return nil,errors.New("NOT EX")
	}
	if len(params) != value.Type().NumIn() {
		err = fmt.Errorf("%d parameters are required, but %d parameters are entered",value.Type().NumIn(),len(params))
		return
	}
	req := make([]reflect.Value,len(params))
	for k,v := range params {
		req[k] = reflect.ValueOf(v)
	}

	return value.Call(req),nil
}
```
测试
```
func TestReflectMap(t *testing.T) {
	sp := New()

	sp.Add("add", func(a,b int) int {return a + b})

	call, err := sp.Call("add", 12, 23)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(call)
	i,b := call[0].Interface().(int)
	if b {
		fmt.Println(i)
	}
}
```
