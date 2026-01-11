package atomic

import (
	"reflect"
	"strings"
)

type Reflected struct {
	reflect.Type
	GetSetAny
}

func ReflectType(tp reflect.Type) (Reflected, bool) {
	if tp.Kind() != reflect.Ptr {
		return Reflected{}, false
	}
	return ReflectValue(reflect.New(tp).Elem())
}

func ReflectValue(reflectValue reflect.Value) (Reflected, bool) {
	if !reflectValue.IsValid() || reflectValue.IsZero() {
		return Reflected{}, false
	}

	tp := reflectValue.Type()

	if tp.Kind() != reflect.Ptr {
		return Reflected{}, false
	}

	elemType := tp.Elem()
	if !strings.HasPrefix(elemType.Name(), "Value[") ||
		elemType.PkgPath() != "github.com/bitmagnet-io/bitmagnet/internal/atomic" {
		return Reflected{}, false
	}

	var getSetAny GetSetAny
	if reflectValue.IsZero() {
		getSetAny = reflect.New(tp.Elem()).Interface().(GetSetAny)
	} else {
		getSetAny = reflectValue.Interface().(GetSetAny)
	}

	value := getSetAny.GetAny()

	return Reflected{
		Type:      reflect.TypeOf(value),
		GetSetAny: getSetAny,
	}, true
}
