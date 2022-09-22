package goconvrt

import (
	"reflect"
)

func ConvertNilSlice2Empty(data interface{}) interface{} {
	value := reflect.ValueOf(data)
	if !value.IsValid() {
		return data
	}
	return convertEmpty(value).Interface()
}

func convertEmpty(data reflect.Value) reflect.Value {
	switch data.Kind() {
	case reflect.Struct:
		if data.IsZero() {
			return newZeroStruct(data.Type())
		}
		return cloneStruct(data)
	case reflect.Slice:
		if data.IsNil() {
			return newEmptySlice(data.Type())
		}
		return cloneSlice(data)
	case reflect.Map:
		if data.IsNil() {
			return data
		}
		return cloneMap(data)
	case reflect.Ptr:
		if data.IsNil() {
			return data
		}
		return clonePtr(data)
	default:
		return data
	}
}

// type 是slice的type
func newEmptySlice(typ reflect.Type) reflect.Value {
	if typ.Kind() != reflect.Slice {
		panic("typ must be slice")
	}
	elemType := typ.Elem()
	switch elemType.Kind() {
	case reflect.Struct:
		field := newZeroStruct(elemType)
		return reflect.MakeSlice(reflect.SliceOf(field.Type()), 0, 0)
	default:
		return reflect.MakeSlice(reflect.SliceOf(elemType), 0, 0)
	}
}

func newZeroStruct(typ reflect.Type) reflect.Value {
	if typ.Kind() != reflect.Struct {
		panic("must be struct")
	}
	value := reflect.New(typ).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			if field.CanSet() {
				field.Set(newZeroStruct(field.Type()))
			}
		case reflect.Slice:
			if field.CanSet() {
				field.Set(newEmptySlice(field.Type()))
			}
		case reflect.Ptr:
			if field.CanSet() {
				field.Set(newZeroPtr(field))
			}
		}
	}
	return value
}

func newZeroPtr(val reflect.Value) reflect.Value {
	if val.Kind() != reflect.Ptr {
		panic("must be ptr")
	}
	elem := val.Elem()
	switch elem.Kind() {
	case reflect.Struct:
		return newZeroStruct(elem.Type()).Addr()
	case reflect.Slice:
		return newEmptySlice(elem.Type()).Addr()
	case reflect.Ptr:
		return newZeroPtr(elem).Addr()
	case reflect.Interface:
		return newZeroIface(elem).Addr()
	default:
		return val
	}
}

func newZeroIface(val reflect.Value) reflect.Value {
	if val.Kind() != reflect.Interface {
		panic("must be interface")
	}
	elem := val.Elem()
	switch elem.Kind() {
	case reflect.Struct:
		return newZeroStruct(elem.Type()).Addr()
	case reflect.Slice:
		return newEmptySlice(elem.Type()).Addr()
	case reflect.Ptr:
		return newZeroPtr(elem).Addr()
	case reflect.Interface:
		return newZeroIface(val)
	default:
		return val
	}
}

func clonePtr(oldPtr reflect.Value) reflect.Value {
	elem := oldPtr.Elem()
	if elem.CanAddr() {
		return convertEmpty(elem).Addr()
	}

	return oldPtr
}

func cloneSlice(oldSlice reflect.Value) reflect.Value {
	newSlice := reflect.MakeSlice(oldSlice.Type(), oldSlice.Len(), oldSlice.Cap())
	for i := 0; i < oldSlice.Len(); i++ {
		newField := newSlice.Index(i)
		if !newField.CanSet() {
			continue
		}
		cloneField(oldSlice.Index(i), newField)
	}
	return newSlice
}

func cloneMap(oldMap reflect.Value) reflect.Value {
	newMap := reflect.MakeMapWithSize(oldMap.Type(), oldMap.Len())
	iter := oldMap.MapRange()
	for iter.Next() {
		newMap.SetMapIndex(iter.Key(), getCloneField(iter.Value()))
	}
	return newMap
}

func cloneStruct(oldObj reflect.Value) reflect.Value {
	newObj := reflect.New(oldObj.Type())
	newVal := newObj.Elem()
	for i := 0; i < oldObj.NumField(); i++ {
		newValField := newVal.Field(i)
		oldField := oldObj.Field(i)
		if !newValField.CanSet() {
			continue
		}
		cloneField(oldField, newValField)
	}
	return newObj.Elem()
}

func cloneField(oldField reflect.Value, newValField reflect.Value) {
	if newValField.CanSet() {
		newValField.Set(getCloneField(oldField))
	}
}

func getCloneField(oldField reflect.Value) reflect.Value {
	switch oldField.Kind() {
	case reflect.Struct:
		return cloneStruct(oldField)
	case reflect.Slice:
		return cloneSlice(oldField)
	case reflect.Map:
		return cloneMap(oldField)
	case reflect.Ptr:
		return clonePtr(oldField)
	case reflect.Interface:
		return getCloneField(oldField.Elem())
	default:
		return oldField
	}
}
