package netchan

import (
	"reflect"
)

// Convert channel pointer to reflect.Value.
// If ch cannot convert, this func rises an panic.
func valueChan(ch interface{}) reflect.Value {

	// assert pointer
	v := reflect.ValueOf(ch)
	if v.Kind() != reflect.Ptr {
		panic("ch must be pointer of channel.")
	}

	// assert chan
	pv := reflect.Indirect(v)
	if pv.Kind() != reflect.Chan {
		panic("ch must be pointer of channel.")
	}

	return pv
}
