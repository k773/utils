package utils

import "sync"

type SafeByteInt64MapStruct struct {
	m map[byte]int64
	s sync.Mutex
}
type SafeByteIntMapStruct struct {
	m map[byte]int
	s sync.Mutex
}
type SafeByteBoolMapStruct struct {
	m map[byte]bool
	s sync.Mutex
}
type SafeByteStringMapStruct struct {
	m map[byte]string
	s sync.Mutex
}

type SafeStringInt64MapStruct struct {
	m map[string]int64
	s sync.Mutex
}
type SafeStringIntMapStruct struct {
	m map[string]int
	s sync.Mutex
}
type SafeStringBoolMapStruct struct {
	m map[string]bool
	s sync.Mutex
}
type SafeStringStringMapStruct struct {
	m map[string]string
	s sync.Mutex
}

type SafeIntInt64MapStruct struct {
	m map[int]int64
	s sync.Mutex
}
type SafeIntIntMapStruct struct {
	m map[int]int
	s sync.Mutex
}
type SafeIntBoolMapStruct struct {
	m map[int]bool
	s sync.Mutex
}
type SafeIntStringMapStruct struct {
	m map[int]string
	s sync.Mutex
}

type SafeInt64Int64MapStruct struct {
	m map[int64]int64
	s sync.Mutex
}
type SafeInt64IntMapStruct struct {
	m map[int64]int
	s sync.Mutex
}
type SafeInt64BoolMapStruct struct {
	m map[int64]bool
	s sync.Mutex
}
type SafeInt64StringMapStruct struct {
	m map[int64]string
	s sync.Mutex
}
