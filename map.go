package linkedMap

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)
import (
	"github.com/tidwall/gjson"
)

type Map[K comparable, V any] struct {
	elem  map[K]*Elem[K, V]
	elist *list[K, V]
	testK any
}

func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		elem:  make(map[K]*Elem[K, V]),
		elist: newList[K, V](),
		testK: *new(K),
	}
}

func (m *Map[K, V]) Set(key K, v V) {
	e, ok := m.elem[key]
	if ok {
		e.Val = v
		return
	}
	e = &Elem[K, V]{}
	e.Val = v
	e.Key = key
	m.elist.pushBack(e)
	m.elem[key] = e
}
func (m *Map[K, V]) Of(key K, v V) *Map[K, V] {
	m.Set(key, v)
	return m
}

func (m *Map[K, V]) Get(key K) (v V, ok bool) {
	e, ok := m.elem[key]
	if ok {
		v = e.Val
		ok = true
		return
	}
	ok = false
	return
}

func (m *Map[K, V]) Delete(key K) {
	e, ok := m.elem[key]
	if !ok {
		return
	}
	delete(m.elem, key)
	m.elist.remove(e)
	return
}

func (m *Map[K, V]) Range(f func(key K, val V) bool) {
	m.elist.foreach(func(e *Elem[K, V]) bool {
		return f(e.Key, e.Val)
	})
}

func (m *Map[K, V]) MarshalJSON() (res []byte, err error) {
	res = make([]byte, 0, 128)
	lenght := 0
	res = append(res, '{')
	m.elist.foreach(func(e *Elem[K, V]) bool {
		var vb []byte
		vb, err = json.Marshal(e.Val)
		if err != nil {
			return false
		}
		res = append(res, '"')
		res = append(res, stringOf(e.Key)...)
		res = append(res, '"')
		res = append(res, ':')
		res = append(res, vb...)
		res = append(res, ',')
		lenght++
		return true
	})
	if err != nil {
		return nil, err
	}
	if lenght > 0 {
		res = res[:len(res)-1]
	}
	res = append(res, '}')
	return res, nil
}

func (m *Map[K, V]) UnmarshalJSON(b []byte) (err error) {
	_, ok := m.testK.(string)
	if !ok {
		panic(fmt.Sprintf("json unmarshal key type must be string: but now is %s", reflect.TypeOf(m.testK)))
	}
	if m.elem == nil {
		m.elem = make(map[K]*Elem[K, V])
	}
	if m.elist == nil {
		m.elist = newList[K, V]()
	}

	p := gjson.Parse(tostring(b))
	if !p.IsObject() {
		return fmt.Errorf("value is not object:%v", p.Type)
	}
	p.ForEach(func(key, value gjson.Result) bool {
		e := new(V)
		err = unmarshalObject2Struct(key.Str, &value, reflect.ValueOf(e).Elem())
		if err != nil {
			return false
		}
		m.Set(any(key.Str).(K), *e)
		return true
	})
	return err
}

func (m *Map[K, V]) String() string {
	sb := &strings.Builder{}
	m.Range(func(k K, v V) bool {
		sb.WriteString(fmt.Sprintf("%v:%v,", k, v))
		return true
	})
	return sb.String()
}

func toValue(r *gjson.Result) any {
	switch r.Type {
	case gjson.String:
		return r.Str
	case gjson.Number:
		return r.Num
	case gjson.True:
		return true
	case gjson.False:
		return false
	case gjson.Null:
		return nil
	case gjson.JSON:
		if r.IsObject() {
			m := New[string, any]()
			r.ForEach(func(key, value gjson.Result) bool {
				m.Set(key.Str, toValue(&value))
				return true
			})
			return m
		}
		if r.IsArray() {
			m := make([]any, 0, 2)
			r.ForEach(func(key, value gjson.Result) bool {
				m = append(m, toValue(&value))
				return true
			})
			return m
		}
	}
	return nil
}

func stringOf(k interface{}) string {
	switch k := k.(type) {
	case string:
		return k
	default:
		panic("k is not string")
	}
}

func tostring(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

var (
	bytesType = reflect.TypeOf([]byte(nil))
)

func unmarshalObject2Struct(path string, in *gjson.Result, v reflect.Value) (err error) {
	if in == nil {
		return nil
	}
	// 是非导出的变量
	if v.Kind() != reflect.Ptr && !v.CanSet() {
		return nil
	}

	switch {
	// 目标是字节数组
	case bytesType == v.Type():
		if in.Type != gjson.String {
			return fmt.Errorf("%s value type is byte , but field in json is not base64", path)
		}

		bytes, err := base64.StdEncoding.DecodeString(in.Str)
		if err != nil {
			return fmt.Errorf("%s value type is byte , but field in json is not valid base64", path)
		}

		v.Set(reflect.ValueOf(bytes))
		return nil
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			vt := v.Type()
			elemType := vt.Elem()
			var nv reflect.Value
			switch elemType.Kind() {
			default:
				nv = reflect.New(elemType)
			}
			err := unmarshalObject2Struct(path, in, nv.Elem())
			if err != nil {
				return err
			}
			v.Set(nv)
			return nil
		}
		return unmarshalObject2Struct(path, in, v.Elem())
	case reflect.Slice:
		if !in.IsArray() {
			return fmt.Errorf("type of %s should be slice", path)
		}
		t := v.Type()
		elemType := t.Elem()
		slice := reflect.MakeSlice(t, 0, 0)
		in.ForEach(func(key, value gjson.Result) bool {
			elemVal := reflect.New(elemType)
			err = unmarshalObject2Struct(path, &value, elemVal)
			if err != nil {
				return false
			}
			slice = reflect.Append(slice, elemVal.Elem())
			return true
		})
		if err != nil {
			return err
		}
		v.Set(slice)
		return nil
	case reflect.String:
		if in.Type != gjson.String {
			return fmt.Errorf("type of %s should be string", path)
		}
		v.SetString(in.Str)
	case reflect.Map:
		if !in.IsObject() {
			return fmt.Errorf("type of %s should be object", path)
		}
		t := v.Type()
		elemT := t.Elem()
		newV := v
		if v.IsNil() {
			newV = reflect.MakeMap(v.Type())
		}
		in.ForEach(func(key, value gjson.Result) bool {
			elemV := reflect.New(elemT)
			err = unmarshalObject2Struct(key.Str, &value, elemV)
			if err != nil {
				return false
			}
			newV.SetMapIndex(reflect.ValueOf(key), elemV.Elem())
			return true
		})

		v.Set(newV)
		return nil
	case reflect.Struct:
		t := v.Type()

		if !in.IsObject() {
			return fmt.Errorf("type of %s should be object", path)
		}
		vmap := make(map[string]*gjson.Result)
		in.ForEach(func(key, value gjson.Result) bool {
			vmap[key.Str] = &value
			return true
		})
		for i := 0; i < t.NumField(); i++ {
			fieldT := t.Field(i)
			name := fieldT.Tag.Get("json")
			if name == "" {
				name = fieldT.Name
			}
			if fieldT.Anonymous {
				err := unmarshalObject2Struct(name, in, v.Field(i))
				if err != nil {
					return err
				}
				continue
			}

			elemV := vmap[name]
			if elemV == nil {
				continue
			}
			// 是包进

			err := unmarshalObject2Struct(name, elemV, v.Field(i))
			if err != nil {
				return err
			}

		}
		return nil
	case reflect.Interface:
		inVal := reflect.ValueOf(toValue(in))
		if inVal.Type().Implements(v.Type()) {
			v.Set(inVal)
		}
		return nil
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int16, reflect.Int8:
		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetInt(intV)
		return nil
	case reflect.Bool:
		boolV, err := boolValueOf(in)
		if err != nil {
			return fmt.Errorf("%s error:%w", path, err)
		}
		v.SetBool(boolV)
		return nil
	case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
		intV, err := intValueOf(in)
		if err != nil {
			return err
		}
		v.SetUint(uint64(intV))
		return nil
	case reflect.Float64, reflect.Float32:
		floatV, err := floatValueOf(in)
		if err != nil {
			return err
		}
		v.SetFloat(floatV)
		return nil
	case reflect.Array:
		if !in.IsArray() {
			return fmt.Errorf("type of %s should be slice", path)
		}

		arType := reflect.ArrayOf(v.Len(), v.Type().Elem())
		arrv := reflect.New(arType)
		pointer := arrv.Pointer()
		eleSize := v.Type().Elem().Size()
		//if v.Len() < len(arr) {
		//	return fmt.Errorf("length of %s is %d . but target value length is %d", path, v.Len(), len(arr))
		//}
		in.ForEach(func(key, value gjson.Result) bool {
			elemV := reflect.New(v.Type().Elem())
			if key.Index >= v.Len() {
				return false
			}
			err = unmarshalObject2Struct(path, &value, elemV)
			if err != nil {
				return false
			}
			memCopy(pointer+uintptr(key.Index)*eleSize, elemV.Pointer(), eleSize)
			return true
		})
		v.Set(arrv.Elem())
	default:
		panic("not support :" + v.Kind().String())
	}
	return nil
}

func intValueOf(v *gjson.Result) (int64, error) {
	switch v.Type {
	case gjson.Number:
		return int64(v.Int()), nil
	default:
		return 0, fmt.Errorf("type is %v ,not int ", reflect.TypeOf(v))
	}
}

func boolValueOf(v *gjson.Result) (bool, error) {
	switch v.Type {
	case gjson.True:
		return true, nil
	case gjson.False:
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool value:%v", v)
	}
}

func floatValueOf(v *gjson.Result) (float64, error) {
	switch v.Type {
	case gjson.Number:
		return v.Float(), nil
	default:
		return 0, fmt.Errorf("invalid float value:%v", v)
	}
}

func bytesOf(p uintptr, len uintptr) []byte {
	h := &reflect.SliceHeader{
		Data: p,
		Len:  int(len),
		Cap:  int(len),
	}
	return *(*[]byte)(unsafe.Pointer(h))
}

func memCopy(dst, src uintptr, len uintptr) {
	db := bytesOf(dst, len)
	sb := bytesOf(src, len)
	copy(db, sb)
}
