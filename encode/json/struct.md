
#####encodeState
```go
// 编码声名
type encodeState struct {
	bytes.Buffer // 累计输出
	scratch      [64]byte
}

// 从Pool中取出实例，没有则new一个
func newEncodeState() *encodeState

// 在代码中有一个defer recover，如果是jsonError类型就返回，否则panic
// 调用e.reflectValue
func (e *encodeState) marshal(v interface{}, opts encOpts) (err error)

// 抛出带jsonError类型参数的panic
func (e *encodeState) error(err error)

// 调用valueEncoder，判断reflect.Value是否有效
// 无效则buffer写入null字符串
// 有效则调用typeEncoder
func (e *encodeState) reflectValue(v reflect.Value, opts encOpts)

// encoderCache是sync.Map类型
// 首先查找encoderCache是否已经绑定传入参数类型的typeEncoder，有则返回该类型的encoderFunc
// 无则创建一个encoderFunc并保存到encoderCache
// 这里有个技巧，先定义f类型为encdoerFunc的空变量
// 用sync.Map的LoadOrStore和sync.WaitGroup防止在没有赋值f的情况下被调用
// 最后f赋值后调用，这样做可以增加性能，不用产生更多的GC
// 如果只是用sync.Mutex代替sync.WaitGroup的话，可能会发race事件
func typeEncoder(t reflect.Type) encoderFunc 

func (e *encodeState) string(s string, escapeHTML bool) 
func (e *encodeState) stringBytes(s []byte, escapeHTML bool)
```