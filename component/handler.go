package component

import "reflect"

type Handler struct {
	Receiver reflect.Value  // receiver of method
	Method   reflect.Method // method stub
}
