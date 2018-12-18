package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type JObject map[string]interface{}

func PruseJSON(data string) *JObject {
	var o interface{}
	if json.Unmarshal([]byte(data), &o) == nil {
		if m, ok := o.(map[string]interface{}); ok {
			j := JObject(m)
			return &j
		} else if list, ok := o.([]interface{}); ok {
			m := make(map[string]interface{})
			for i, v := range list {
				m[strconv.Itoa(i)] = v
			}
			j := JObject(m)
			return &j
		}
	}
	return nil
}

func (v *JObject) JToken(key string) *JObject {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.(map[string]interface{}); ok {
			j := JObject(m)
			return &j
		}
	}
	return nil
}

func (v *JObject) jArray(key string) []interface{} {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.([]interface{}); ok {
			if len(m) == 1 {
				if x, ok := m[0].([]interface{}); ok {
					return x
				}
			}
			return m
		} else if t, ok := t.(map[string]interface{}); ok {
			r := make([]interface{}, 0)
			for _, v := range t {
				if m, ok := v.(map[string]interface{}); ok {
					r = append(r, m)
				}
			}
			return r
		}
	}
	return nil
}

func (v *JObject) JTokens(key string) []*JObject {
	t := v.jArray(key)
	if t != nil {
		r := make([]*JObject, 0)
		for _, v := range t {
			if m, ok := v.(map[string]interface{}); ok {
				j := JObject(m)
				r = append(r, &j)
			}
		}
		return r
	}
	return nil
}

func (v *JObject) String() string {
	if s, err := json.Marshal(v); err == nil {
		return string(s)
	}
	return ""
}

func JTokensToString(jTokens []*JObject) string {
	arr := make([]string, 0)
	for _, v := range jTokens {
		arr = append(arr, v.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(arr, ","))
}

func (v *JObject) Unmarshal(p interface{}) error {
	return json.Unmarshal([]byte(v.String()), p)
}

func UnmarshalJTokens(jTokens []*JObject, p interface{}) error {
	return json.Unmarshal([]byte(JTokensToString(jTokens)), p)
}
