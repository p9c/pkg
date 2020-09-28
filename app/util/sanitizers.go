package util

import (
	"errors"
	"path"
	"path/filepath"
	"time"
)

// Example is a dummy that returns always no error, use as a prototype
func Example() HookFunc {
	cu := func(in interface{}) (err error) {
		return nil
	}
	return &cu
}

// CheckURL cleans the URL
func CheckURL() HookFunc {
	cu := func(in interface{}) (err error) {
		if pointer, ok := in.(*string); ok {
			*pointer = path.Clean(*pointer)
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &cu
}

// CheckPath cleans the filesystem path of special characters
func CheckPath() HookFunc {
	cp := func(in interface{}) (err error) {
		if pointer, ok := in.(*string); ok {
			*pointer = filepath.Clean(*pointer)
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &cp
}

// IntBounds allows setting a floor and ceiling on a value
func IntBounds(min, max int) HookFunc {
	if min > max {
		panic("bad parameters for bounds not in correct relation")
	}
	fn := func(in interface{}) (err error) {
		if pointer, ok := in.(*int); ok {
			val := *pointer
			if val < min {
				val = min
			}
			if val > max {
				val = max
			}
			*pointer = val
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &fn
}

// UintBounds allows setting a floor and ceiling on a value
func UintBounds(min, max uint) HookFunc {
	if min > max {
		panic("bad parameters for bounds not in correct relation")
	}
	fn := func(in interface{}) (err error) {
		if pointer, ok := in.(*uint); ok {
			val := *pointer
			if val < min {
				val = min
			}
			if val > max {
				val = max
			}
			*pointer = val
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &fn
}

// Float64Bounds allows setting a floor and ceiling on a value
func Float64Bounds(min, max float64) HookFunc {
	if min > max {
		panic("bad parameters for bounds not in correct relation")
	}
	fn := func(in interface{}) (err error) {
		if pointer, ok := in.(*float64); ok {
			val := *pointer
			if val < min {
				val = min
			}
			if val > max {
				val = max
			}
			*pointer = val
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &fn
}

// DurationBounds allows setting a floor and ceiling on a value
func DurationBounds(min, max time.Duration) HookFunc {
	if min > max {
		panic("bad parameters for bounds not in correct relation")
	}
	fn := func(in interface{}) (err error) {
		if pointer, ok := in.(*time.Duration); ok {
			val := *pointer
			if val < min {
				val = min
			}
			if val > max {
				val = max
			}
			*pointer = val
		} else {
			err = errors.New("no Value was found in the pointer")
		}
		return
	}
	return &fn
}
