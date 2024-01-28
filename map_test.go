package linkedMap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

func TestNew(t *testing.T) {
	m := New[string, any]()
	m.Set("1", 5)
	m.Set("2", 10)
	m.Set("\"66", `"dddd"
	`)
	m.Set("c", []interface{}{
		1, 2, 3,
	})
	m.Of("1", 4).Of("adf", New[string, any]().Of("name", "bon").Of("haha", 3))
	b, _ := m.MarshalJSON()
	fmt.Println(string(b))
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Chs  [2]int `json:"chs"`
}

func TestNew2(t *testing.T) {
	js := `{
"bob":[{
	"name":"blb1",
	"age":5
}],
"alice":[
{
	"name":"ace",
	"age":5
},
{
	"name":"ace2",
	"age":6
}
]

}`

	m := New[string, any]()
	err := json.Unmarshal([]byte(js), m)
	if err != nil {
		panic(err)
	}
	a, _ := m.Get("bog")
	fmt.Println(a)

	m.Range(func(k string, v interface{}) bool {
		fmt.Println(k, v, reflect.TypeOf(v))
		return true
	})
	ace, ok := m.Get("alice")
	if ok {
		fmt.Println(reflect.TypeOf(ace.([]any)[0]))
	}

}

func TestMap(t *testing.T) {
	a := New[string, *Map[string, any]]()
	js := `
	{
		"name":{
			"anc":5
		}
	}
	
	`
	err := json.Unmarshal([]byte(js), a)
	if err != nil {
		panic(err)
	}
	m, _ := a.Get("name")
	fmt.Println(m.Get("anc"))

}

func BenchmarkUnmarshal(b *testing.B) {
	js := []byte(`{
"bob":[{
	"name":"blb1",
	"age":5
}],
"alice":[
{
	"name":"ace",
	"age":5
},
{
	"name":"ace2",
	"age":6
}
]

}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := New[string, any]()
		m.UnmarshalJSON(js)
	}
}

func BenchmarkUnmarshal2(b *testing.B) {
	js := []byte(`{
"bob":[{
	"name":"blb1",
	"age":5
}],
"alice":[
{
	"name":"ace",
	"age":5
},
{
	"name":"ace2",
	"age":6
}
]

}`)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := map[string]any{}
		json.Unmarshal(js, &m)
	}
}

func TestUm(t *testing.T) {
	u := User{}

	js := `{
"chs":[1,2]
}`
	json.Unmarshal([]byte(js), &u)

	fmt.Println(u)
}

func BenchmarkGJSON(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var t Temp
		t.Decode(`{"name":"haha","age":5,"sum":false}`)
	}
}

func BenchmarkGJSON2(b *testing.B) {
	b.ReportAllocs()
	bs := []byte(`{"name":"haha","age":5,"sum":false}`)
	for i := 0; i < b.N; i++ {
		var t Temp
		json.Unmarshal(bs, &t)
	}
}

type Temp struct {
	Name string `json:"name"`
	Age  int64  `json:"age"`
	Sum  bool
	Big  string `json:"big"`
	Ss   sss    `json:"ss"`
}
type sss string

func (t *sss) UnmarshalJSON(b []byte) error {
	*t = sss(string(b))
	return nil
}

func (t *Temp) Decode(b string) {
	res := gjson.Parse(b)
	//t.Name = res.Get("name").String()
	//t.Age = res.Get("age").Int()
	//t.Sum = res.Get("sum").Bool()
	res.ForEach(func(key, value gjson.Result) bool {
		switch key.String() {
		case "name":
			t.Name = value.String()
		case "age":
			t.Age = value.Int()
		case "sum":
			t.Sum = value.Bool()
		case "big":
			t.Big = value.String()

		}
		return true
	})
}

func TestJSON2(t *testing.T) {
	m := New[string, Temp]()
	err := json.Unmarshal([]byte(`{"ss":{"ss":"dd"}}`), m)
	fmt.Println(err, m)
	//tv, _ := m.Get("ss")
	//fmt.Println(*tv.Ss)

	//var a json.Unmarshaler = new(sss)

}

func TestJSON(t *testing.T) {
	m := map[string]Temp{}
	err := json.Unmarshal([]byte(`{"ss":{"ss":"dd"}}`), &m)
	fmt.Println(err, m)

}

func Test_GOOO(t *testing.T) {

}

func Test_Unmarshal(t *testing.T) {
	js := `
	{
		"key":{
			"kls":{
				"age":0
			}
		},
		"vals":{
			"kls":{
				"age":5
			}
		}
	}
	
	`
	type A struct {
		Age int `json:"age,omitempty"`
	}

	res := New[string, *Map[string, *A]]()

	err := json.Unmarshal([]byte(js), res)
	if err != nil {
		panic(err)
	}

	ns, _ := res.MarshalJSON()

	fmt.Println(res, string(ns))
}

/*

 */
