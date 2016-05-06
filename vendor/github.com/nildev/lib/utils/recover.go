package utils

import "fmt"

// PanicRecover is handler for panic, wraps panic in error
func PanicRecover() error {
	if r := recover(); r != nil {
		var ok bool
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("pkg: %v", r)
		}
		return err
	}

	return nil
}
