package utils

func RefPointer[T any | interface{} | string | bool | int | int8 | int16 | int32 | int64 | float32 | float64](val T) *T {
	return &val
}

func DerefPointer[T any | interface{} | string | bool | int | int8 | int16 | int32 | int64 | float32 | float64](val *T) T {
	return *val
}
