package utils

import (
	"errors"
	"reflect"
)

func ToM(value interface{}) map[string]interface{} {
	if value == nil{
		return nil
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		value = v.Interface()
		v = reflect.ValueOf(value)
	}
	t := reflect.TypeOf(value)
	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//
func ToMs(values []interface{}) []interface{} {
	list:= make([]interface{},len(values))
	l:= len(values)
	for i:=0;i<l;i++{
		value:=values[i]
		if value == nil{
			return nil
		}
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
			value = v.Interface()
			v = reflect.ValueOf(value)
		}
		t := reflect.TypeOf(value)
		data := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			data[t.Field(i).Name] = v.Field(i).Interface()
		}
		list[i] = data
	}
	return list
}

func setField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New("No such field: %s in obj" + name)
	}

	if !structFieldValue.CanSet() {
		return errors.New("Cannot set %s field value" + name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

type RPCConfig struct {
	ip   string
	port string
}

func (s *RPCConfig) FillStruct(m map[string]string) error {
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
