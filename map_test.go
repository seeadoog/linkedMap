package linkedMap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m := New[string, any]()
	m.Set("1", 5)
	m.Set("2", 10)
	m.Set("c", []interface{}{
		1, 2, 3,
	})
	m.Of("1", 4).Of("adf", New[string, any]().Of("name", "bon").Of("haha", 3))
	b, _ := json.Marshal(m)
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
}

func TestMap(t *testing.T) {

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
		json.Unmarshal(js, m)
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
