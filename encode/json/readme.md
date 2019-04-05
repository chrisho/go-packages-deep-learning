### encoding/json

```go
func Marshal(v interface{}) ([]byte, error) {
	e := newEncodeState()                               // 获取编码声明实例
 
	err := e.marshal(v, encOpts{escapeHTML: true})      // 设置编码参数并开始编码
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)            // 生成空字节切片，并插入编码结果

	e.Reset()                                           // 重置编码缓冲
	encodeStatePool.Put(e)                              // 回收编码缓冲

	return buf, nil
}
```

```go
// 编码缓冲
type encodeState struct {
	bytes.Buffer // accumulated output
	scratch      [64]byte
}

// 编码缓冲池
var encodeStatePool sync.Pool // 定义为同步池类型

// 生成编码状态实例
func newEncodeState() *encodeState {
	if v := encodeStatePool.Get(); v != nil {  // 如果编码状态池里已生成实例
		e := v.(*encodeState)                  // 则取出已有实例
		e.Reset()                              // 并重置字节缓存区
		return e
	}
	return new(encodeState)                    // 生成编码缓冲实例
}

```

```go
func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) {
	// 省略捕获恐慌
	... 
	
	e.reflectValue(reflect.ValueOf(v), opts) // 映射Value实体并再下一步映射
	return nil
}
```

```go
func (e *encodeState) reflectValue(v reflect.Value, opts encOpts) {
	valueEncoder(v)(e, v, opts)             // 获取编码方法并执行encoderFunc的方法
}
```

```go
func valueEncoder(v reflect.Value) encoderFunc {
	// 省略异常验证
	return typeEncoder(v.Type())
}
```

```go
var encoderCache sync.Map // map[reflect.Type]encoderFunc

func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {     // 如果存在对应类型的encoderFunc
		return fi.(encoderFunc)                 // 则直接调用
	}
	
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	// 如果存在键则返回方法，不存在则保存及返回方法
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(e *encodeState, v reflect.Value, opts encOpts) {
		wg.Wait()
		f(e, v, opts)                           // 这是闭包方法
	}))
	if loaded {                                 // 如果不是新储存的方法则直接返回调用
		return fi.(encoderFunc)             
	}

	f = newTypeEncoder(t, true)                 // 获取当前类型编码方法
	wg.Done()
	encoderCache.Store(t, f)                    // 最后保存
	return f
}
```

```go
func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	
	/*-----------
	检测是否实现了MarshalJSON或者MarshalText的接口
	*/
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(marshalerType) {
			return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
		}
	}

	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}
	if t.Kind() != reflect.Ptr && allowAddr {
		if reflect.PtrTo(t).Implements(textMarshalerType) {
			return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
		}
	}
	/*
	-------------
	*/

	/*
	检测类型并且返回对应编码方法
	*/
	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Ptr:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}
```