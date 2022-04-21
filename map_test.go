package linkedMap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
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
	m := New[User, User]()
	m.Set(User{
		Name: "g",
	}, User{Name: "g"})

	fmt.Println(m.Get(User{
		Name: "g",
	}))
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

func BenchmarkName(b *testing.B) {
	l := sync.Mutex{}
	b.SetParallelism(10)
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Lock()
			l.Unlock()
		}
	})
}

func BenchmarkChan(b *testing.B) {
	c := make(chan bool, 1)
	c <- true
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			go func() {

			}()
		}
	})
}

func BenchmarkGO(b *testing.B) {
	c := make(chan bool, 1)
	c <- true

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			runtime.Gosched()
		}
	})
}

func BenchmarkMaap(b *testing.B) {
	m := map[string]int{}
	b.ReportAllocs()
	m["1"] = 5
	for i := 0; i < b.N; i++ {
		_ = m["1"]
	}
}

func BenchmarkArr(b *testing.B) {
	b.ReportAllocs()
	a := make([]string, 30)
	for i := 0; i < b.N; i++ {
		for j := 0; j < len(a); j++ {

		}
	}
}
func BenchmarkMap2(b *testing.B) {
	b.ReportAllocs()
	a := New[string, int]()
	a.Set("ee", 1)
	for i := 0; i < b.N; i++ {
		a.Get("ee")
	}
}
