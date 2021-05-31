package notion

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

type JSON map[string]interface{}

func (j JSON) String() string {
	b, _ := json.MarshalIndent(j, "", "  ")

	return string(b)
}

func (j JSON) Set(name string, v interface{}) {
	t := reflect.TypeOf(v).Kind()

	switch t {
	case reflect.Struct, reflect.Ptr, reflect.Slice, reflect.Map:
		jj := JSON{}
		jj.Marshal(v)
		j[name] = jj
	default:
		j[name] = v
	}
}

func (j JSON) Append(name string, v interface{}) {
	vv, ok := j.GetJSONList(name)
	if ok {
		jj := JSON{}
		jj.Marshal(v)
		j[name] = append(vv, jj)
	}
}

func (j JSON) Get(name string) interface{} {
	return j[name]
}

func (j JSON) GetString(name string) string {
	return fmt.Sprint(j.Get(name))
}

func (j JSON) GetInt(name string) int {
	i, _ := strconv.Atoi(j.GetString(name))

	return i
}

func (j JSON) GetBool(name string) bool {
	b, _ := strconv.ParseBool(j.GetString(name))

	return b
}

func (j JSON) GetFloat(name string) float64 {
	f, _ := strconv.ParseFloat(j.GetString(name), 64)

	return f
}

func (j JSON) GetJSON(name string) (JSON, bool) {
	if v := j.Get(name); v != nil {
		if m, ok := v.(JSON); ok {
			return m, ok
		} else if m, ok := v.(map[string]interface{}); ok {
			return m, ok
		}
	}
	return JSON{}, false
}

func (j JSON) GetJSONList(name string) ([]JSON, bool) {
	l := []JSON{}

	if v := j.Get(name); v != nil {
		if m, ok := v.([]JSON); ok {
			for _, mm := range m {
				jj := JSON{}
				jj.Marshal(mm)
				l = append(l, jj)
			}

			return l, ok
		} else if m, ok := v.([]interface{}); ok {
			for _, mm := range m {
				jj := JSON{}
				jj.Marshal(mm)
				l = append(l, jj)
			}

			return l, ok
		}
	}

	return []JSON{}, false
}

func (j JSON) Marshal(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &j)
}

func (j JSON) Unmarshal(v interface{}) error {
	b, err := json.Marshal(j)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
