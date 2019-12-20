package json

import (
	"encoding/json"
	"errors"
	"io"
)

// Kind json kind
type Kind uint8

const (
	// Null 空
	Null Kind = iota
	// Bool b
	Bool
	// Number n
	Number // int float 他们可以 强转
	// String s
	String
	// Object o
	Object
	// Array a
	Array
)

// JSON s 简单 json 的封装
// JSON 只能是 jobject 和 jarray
type JSON struct {
	v interface{}
}

// LoadBytes is new
func LoadBytes(data []byte) (j *JSON, err error) {
	var d interface{}
	err = json.Unmarshal(data, &d)
	if err != nil {
		return
	}
	switch d.(type) {
	case map[string]interface{}, []interface{}:
		j = &JSON{v: d}
	default:
		err = errors.New("type is err")
	}
	return
}

// LoadString s
func LoadString(data string) (j *JSON, err error) {
	return LoadBytes([]byte(data))
}

// LoadReader like
func LoadReader(r io.Reader) (j *JSON, err error) {
	var d interface{}
	err = json.NewDecoder(r).Decode(&d)
	if err != nil {
		return
	}
	switch d.(type) {
	case map[string]interface{}, []interface{}:
		j = &JSON{v: d}
	default:
		err = errors.New("type is err")
	}
	return
}

// Kind get json kind
func (j *JSON) Kind() Kind {
	switch j.v.(type) {
	case bool:
		return Bool
	case float64:
		return Number
	case string:
		return String
	case map[string]interface{}:
		return Object
	case []interface{}:
		return Array
	default:

	}
	return Null
}

// ToFloat64  作为 value 的 方法	ok 表示 类型 是否 正确
func (j *JSON) ToFloat64() (ret float64, ok bool) {
	// json.Decode 是用 float64 来表示 数字
	ret, ok = j.v.(float64)
	return
}

// ToFloat32 float32
func (j *JSON) ToFloat32() (ret float32, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = float32(f)
	}
	return
}

// ToInt8 to int
func (j *JSON) ToInt8() (ret int8, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = int8(f)
	}
	return
}

// ToUint8 to
func (j *JSON) ToUint8() (ret uint8, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = uint8(f)
	}
	return
}

// ToInt16 to int
func (j *JSON) ToInt16() (ret int16, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = int16(f)
	}
	return
}

// ToUint16 to
func (j *JSON) ToUint16() (ret uint16, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = uint16(f)
	}
	return
}

// ToInt32 to int
func (j *JSON) ToInt32() (ret int32, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = int32(f)
	}
	return
}

// ToUint32 to
func (j *JSON) ToUint32() (ret uint32, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = uint32(f)
	}
	return
}

// ToInt to int
func (j *JSON) ToInt() (ret int, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = int(f)
	}
	return
}

// ToUint to
func (j *JSON) ToUint() (ret uint, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = uint(f)
	}
	return
}

// ToInt64 t
func (j *JSON) ToInt64() (ret int64, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = int64(f)
	}
	return
}

// ToUint64 ut
func (j *JSON) ToUint64() (ret uint64, ok bool) {
	var f float64
	f, ok = j.ToFloat64()
	if ok {
		ret = uint64(f)
	}
	return
}

// ToBool b
func (j *JSON) ToBool() (ret bool, ok bool) {
	ret, ok = j.v.(bool)
	return
}

// ToString ss
func (j *JSON) ToString() (ret string, ok bool) {
	ret, ok = j.v.(string)
	return
}

// GetString s
func (j *JSON) GetString(k string, def ...string) (ret string) {
	v := j.get(k)
	if v != nil {
		if s, ok := v.(string); ok {
			ret = s
			return
		}
	}
	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetBytes 将 对象 序列化 字符串 包括 双引号
func (j *JSON) GetBytes(k string) (bs []byte) {
	v := j.get(k)
	if v != nil {
		bs, _ = json.Marshal(v)
	}
	return
}

// GetBool bool
func (j *JSON) GetBool(k string, def ...bool) (ret bool) {
	v := j.get(k)
	if v != nil {
		if b, ok := v.(bool); ok {
			ret = b
			return
		}
	}
	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetFloat64 float64
func (j *JSON) GetFloat64(k string, def ...float64) (ret float64) {
	num, ok := j.getNumber(k)
	if ok {
		ret = num
		return
	}

	if len(def) > 0 {
		return def[0]
	}
	return 0.0
}

// GetFloat32 float32
func (j *JSON) GetFloat32(k string, def ...float32) (ret float32) {
	num, ok := j.getNumber(k)
	if ok {
		ret = float32(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetInt int
func (j *JSON) GetInt(k string, def ...int) (ret int) {
	num, ok := j.getNumber(k)
	if ok {
		ret = int(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetUint uint
func (j *JSON) GetUint(k string, def ...uint) (ret uint) {
	num, ok := j.getNumber(k)
	if ok {
		ret = uint(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetInt64 int64
func (j *JSON) GetInt64(k string, def ...int64) (ret int64) {
	num, ok := j.getNumber(k)
	if ok {
		ret = int64(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetUint64 uint64
func (j *JSON) GetUint64(k string, def ...uint64) (ret uint64) {
	num, ok := j.getNumber(k)
	if ok {
		ret = uint64(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetInt8 int8
func (j *JSON) GetInt8(k string, def ...int8) (ret int8) {
	num, ok := j.getNumber(k)
	if ok {
		ret = int8(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetUint8 uint8
func (j *JSON) GetUint8(k string, def ...uint8) (ret uint8) {
	num, ok := j.getNumber(k)
	if ok {
		ret = uint8(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetInt16 int16
func (j *JSON) GetInt16(k string, def ...int16) (ret int16) {
	num, ok := j.getNumber(k)
	if ok {
		ret = int16(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetUint16 uint16
func (j *JSON) GetUint16(k string, def ...uint16) (ret uint16) {
	num, ok := j.getNumber(k)
	if ok {
		ret = uint16(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetInt32 int32
func (j *JSON) GetInt32(k string, def ...int32) (ret int32) {
	num, ok := j.getNumber(k)
	if ok {
		ret = int32(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}

// GetUint32 uint32
func (j *JSON) GetUint32(k string, def ...uint32) (ret uint32) {
	num, ok := j.getNumber(k)
	if ok {
		ret = uint32(num)
		return
	}

	if len(def) > 0 {
		ret = def[0]
	}
	return
}
func (j *JSON) get(k string) interface{} {
	if j.v == nil {
		return nil
	}
	if mm, ok := j.v.(map[string]interface{}); ok {
		if v, ok := mm[k]; ok {
			return v
		}
	}

	return nil
}

func (j *JSON) getNumber(k string) (float64, bool) {
	v := j.get(k)
	if v != nil {
		// json.Decode 是用 float64 来表示 数字
		if ret, ok := v.(float64); ok {
			return ret, true
		}

		if ret, ok := v.(json.Number); ok {
			n, err := ret.Float64()
			return n, err == nil
		}
		// 数字 解析 就 这两个 类型了
		// 就算是 有值 类型不对 也 当成 没有 number 是 一个类型 这里会 强转
	}
	return 0.0, false
}

// GetJSON 嵌套 对象
func (j *JSON) GetJSON(k string) *JSON {
	v := j.get(k)
	if v == nil {
		return nil
	}

	return &JSON{v: v}
}

// IsArray is
func (j *JSON) IsArray() (b bool) {
	_, b = j.v.([]interface{})
	return
}

// GetSize len
func (j *JSON) GetSize() int {
	if a, ok := j.v.([]interface{}); ok {
		return len(a)
	}
	return 0
}

// GetFloat64At at
func (j *JSON) GetFloat64At(index int) (ret float64) {
	v := j.index(index)
	if v != nil {
		ret, _ = v.(float64)
		return
	}
	return
}

// GetFloat32At at
func (j *JSON) GetFloat32At(index int) (ret float32) {
	f := j.GetFloat64At(index)
	ret = float32(f)
	return
}

// GetIntAt at
func (j *JSON) GetIntAt(index int) (ret int) {
	f := j.GetFloat64At(index)
	ret = int(f)
	return
}

//GetUintAt at
func (j *JSON) GetUintAt(index int) (ret uint) {
	f := j.GetFloat64At(index)
	ret = uint(f)
	return
}

// GetInt8At at
func (j *JSON) GetInt8At(index int) (ret int8) {
	f := j.GetFloat64At(index)
	ret = int8(f)
	return
}

// GetUint8At at
func (j *JSON) GetUint8At(index int) (ret uint8) {
	f := j.GetFloat64At(index)
	ret = uint8(f)
	return
}

// GetInt16At at
func (j *JSON) GetInt16At(index int) (ret int16) {
	f := j.GetFloat64At(index)
	ret = int16(f)
	return
}

// GetUint16At at
func (j *JSON) GetUint16At(index int) (ret uint16) {
	f := j.GetFloat64At(index)
	ret = uint16(f)
	return
}

// GetInt32At at
func (j *JSON) GetInt32At(index int) (ret int32) {
	f := j.GetFloat64At(index)
	ret = int32(f)
	return
}

// GetUint32At at
func (j *JSON) GetUint32At(index int) (ret uint32) {
	f := j.GetFloat64At(index)
	ret = uint32(f)
	return
}

// GetInt64At at
func (j *JSON) GetInt64At(index int) (ret int64) {
	f := j.GetFloat64At(index)
	ret = int64(f)
	return
}

// GetUint64At at
func (j *JSON) GetUint64At(index int) (ret uint64) {
	f := j.GetFloat64At(index)
	ret = uint64(f)
	return
}

// GetBoolAt at
func (j *JSON) GetBoolAt(index int) (ret bool) {
	v := j.index(index)
	if v != nil {
		ret, _ = v.(bool)
		return
	}
	return
}

// GetStringAt at
func (j *JSON) GetStringAt(index int) (ret string) {
	v := j.index(index)
	if v != nil {
		ret, _ = v.(string)
		return
	}
	return
}

func (j *JSON) index(index int) interface{} {
	if a, ok := j.v.([]interface{}); ok {
		return a[index]
	}
	return nil
}

// ToArray arr
func (j *JSON) ToArray() (arr []*JSON) {
	if a, ok := j.v.([]interface{}); ok {
		for i := 0; i < len(a); i++ {
			arr = append(arr, &JSON{a[i]})
		}
	}
	return
}

// ToStrings ss
func (j *JSON) ToStrings() (arr []string) {
	if a, ok := j.v.([]interface{}); ok {
		for i := 0; i < len(a); i++ {
			if s, ok := a[i].(string); ok {
				arr = append(arr, s)
			}
		}
	}
	return
}

// ToInts ins
func (j *JSON) ToInts() (arr []int) {
	if a, ok := j.v.([]interface{}); ok {
		for i := 0; i < len(a); i++ {
			if n, ok := a[i].(float64); ok {
				arr = append(arr, int(n))
			}
		}
	}
	return
}

// ToInt8s ins
func (j *JSON) ToInt8s() (arr []int8) {
	if a, ok := j.v.([]interface{}); ok {
		for i := 0; i < len(a); i++ {
			if n, ok := a[i].(float64); ok {
				arr = append(arr, int8(n))
			}
		}
	}
	return
}

// ToUint8s ins
func (j *JSON) ToUint8s() (arr []uint8) {
	if a, ok := j.v.([]interface{}); ok {
		for i := 0; i < len(a); i++ {
			if n, ok := a[i].(float64); ok {
				arr = append(arr, uint8(n))
			}
		}
	}
	return
}

// Marshal to
func (j *JSON) Marshal() (p []byte) {
	p, _ = json.Marshal(j.v)
	return
}

// Add add
func (j *JSON) Add(k string, v interface{}) {
	if a, ok := j.v.(map[string]interface{}); ok {
		a[k] = v
	}
}

// JObject -------------------------------------------
// 简单 生成 bytes
type JObject struct {
	v map[string]interface{}
}

// JArray js
type JArray struct {
	v []interface{}
}

// NewJObject new
func NewJObject() *JObject {
	return &JObject{
		v: make(map[string]interface{}),
	}
}

// NewJArray new
func NewJArray() *JArray {
	return &JArray{
		v: make([]interface{}, 0),
	}
}

// Marshal to
func (j *JObject) Marshal() (p []byte) {
	p, _ = json.Marshal(j.v)
	return
}

// Add a
func (j *JObject) Add(k string, v interface{}) {
	j.v[k] = v
}

// AddObject aa
func (j *JObject) AddObject(k string, o *JObject) {
	j.v[k] = o.v
}
