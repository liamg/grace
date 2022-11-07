package tracer

import (
	"errors"
	"fmt"
	"reflect"
)

func decodeStruct(memory []byte, target interface{}) error {

	sType := reflect.TypeOf(target)
	if sType.Kind() != reflect.Ptr {
		return errors.New("target must be a pointer")
	}

	var index uintptr

	sPtrValue := reflect.ValueOf(target)
	sValue := sPtrValue.Elem()

	switch sValue.Kind() {
	case reflect.Struct:
		for i := 0; i < sValue.Type().NumField(); i++ {
			size := sValue.Type().Field(i).Type.Size()
			if sValue.Type().Field(i).Name == "_" {
				index += size
				continue
			}
			raw := memory[index : index+size]
			if err := decodeAnonymous(sValue.Field(i), raw); err != nil {
				return err
			}
			index += size
		}
	default:
		return errors.New("target must be a pointer to a struct")
	}

	return nil
}

func decodeAnonymous(target reflect.Value, raw []byte) error {

	if target.Kind() == reflect.Ptr {
		target.Set(reflect.New(target.Type().Elem()))
		target = target.Elem()
	}

	if !target.CanSet() {
		return fmt.Errorf("cannot set %s", target.Type())
	}

	switch target.Kind() {
	case reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8, reflect.Uint, reflect.Uintptr:
		target.SetUint(decodeUint(raw))
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		target.SetInt(decodeInt(raw))
	case reflect.String:
		target.SetString(string(raw))
	case reflect.Struct:
		if err := decodeStruct(raw, target.Addr().Interface()); err != nil {
			return err
		}
	case reflect.Array, reflect.Slice:
		var index uintptr
		for i := 0; i < target.Len(); i++ {
			memory := raw[index : index+target.Type().Elem().Size()]
			if err := decodeAnonymous(target.Index(i), memory); err != nil {
				return err
			}
			index += target.Type().Elem().Size()
		}
	default:
		return fmt.Errorf("unsupported kind for field '%s': %s", target.String(), target.Kind().String())
	}
	return nil
}

func decodeUint(raw []byte) uint64 {
	var output uint64
	for i := 0; i < len(raw); i++ {
		output |= uint64(raw[i]) << uint(i*8)
	}
	return output
}
func decodeInt(raw []byte) int64 {
	var output int64
	for i := 0; i < len(raw); i++ {
		output |= int64(raw[i]) << uint(i*8)
	}
	return output
}
