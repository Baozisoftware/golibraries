package json

import "encoding/json"

type jObject map[string]interface{}

func PruseJSON(data string) *jObject {
	var o interface{}
	if json.Unmarshal([]byte(data), &o) == nil {
		if m, ok := o.(map[string]interface{}); ok {
			j := jObject(m)
			return &j
		}
	}
	return nil
}

func (v *jObject) JToken(key string) *jObject {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.(map[string]interface{}); ok {
			j := jObject(m)
			return &j
		}
	}
	return nil
}

func (v *jObject) jArray(key string) []interface{} {
	t, ok := (*v)[key]
	if ok {
		if m, ok := t.([]interface{}); ok {
			return m
		}
	}
	return nil
}

func (v *jObject) JTokens(key string) []*jObject {
	t := v.jArray(key)
	if t != nil {
		r := make([]*jObject, 0)
		for _, v := range t {
			if m, ok := v.(map[string]interface{}); ok {
				j := jObject(m)
				r = append(r, &j)
			}
		}
		return r
	}
	return nil
}
