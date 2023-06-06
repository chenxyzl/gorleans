package component

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

var typeOfProtoMsg = reflect.TypeOf(new(proto.Message)).Elem()

func IsExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

// isHandlerMethod decide a method is suitable handler method
func isHandlerMethod(method reflect.Method, p reflect.Type) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs two or three ins: receiver, context.Context and optional []byte or pointer.
	if mt.NumIn() != 3 || mt.NumOut() != 2 {
		return false
	}

	//mn := method.Name

	//匹配参数1的类型
	if t1 := mt.In(1); !t1.AssignableTo(p) && !t1.AssignableTo(p.Elem()) {
		return false
	}
	//匹配参数2的类型 必须是proto 且名字为{mn}Req
	if t1 := mt.In(2); t1.Kind() != reflect.Ptr || !t1.Implements(typeOfProtoMsg) || !strings.HasSuffix(t1.Elem().Name(), "Req") {
		return false
	}

	//匹配返回值1的类型 必须是proto 且名字必须是{mn}Rsp
	if t1 := mt.Out(0); t1.Kind() != reflect.Ptr || !t1.Implements(typeOfProtoMsg) || !strings.HasSuffix(t1.Elem().Name(), "Rsp") {
		return false
	}
	//匹配返回值2的类型 必须是code
	if t1 := mt.Out(1); t1.Kind() != reflect.Int32 || t1.Name() != "Code" {
		return false
	}

	return true
}

func SuitableHandlerMethods[T any, C any](nameFunc func(string) string) map[string]*Handler {
	var t = new(T)
	var typ = reflect.TypeOf(t)
	var c = reflect.TypeOf(new(C))
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mn := method.Name
		if isHandlerMethod(method, c) {
			// rewrite handler name
			if nameFunc != nil {
				mn = nameFunc(mn)
			}
			handler := &Handler{
				Receiver: reflect.ValueOf(t),
				Method:   method,
			}
			if _, ok := methods[mn]; ok {
				err := fmt.Errorf("repeated handler, %v", mn)
				panic(err)
			}
			methods[mn] = handler
		}
	}
	return methods
}
