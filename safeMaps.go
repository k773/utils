package utils

import "sync"

type safeByteInt64MapStruct struct {
	m map[byte]int64
	s sync.Mutex
}
type safeByteIntMapStruct struct {
	m map[byte]int
	s sync.Mutex
}
type safeByteBoolMapStruct struct {
	m map[byte]bool
	s sync.Mutex
}
type safeByteStringMapStruct struct {
	m map[byte]string
	s sync.Mutex
}

type safeStringInt64MapStruct struct {
	m map[string]int64
	s sync.Mutex
}
type safeStringIntMapStruct struct {
	m map[string]int
	s sync.Mutex
}
type safeStringBoolMapStruct struct {
	m map[string]bool
	s sync.Mutex
}
type safeStringStringMapStruct struct {
	m map[string]string
	s sync.Mutex
}

type safeIntInt64MapStruct struct {
	m map[int]int64
	s sync.Mutex
}
type safeIntIntMapStruct struct {
	m map[int]int
	s sync.Mutex
}
type safeIntBoolMapStruct struct {
	m map[int]bool
	s sync.Mutex
}
type safeIntStringMapStruct struct {
	m map[int]string
	s sync.Mutex
}

type safeInt64Int64MapStruct struct {
	m map[int64]int64
	s sync.Mutex
}
type safeInt64IntMapStruct struct {
	m map[int64]int
	s sync.Mutex
}
type safeInt64BoolMapStruct struct {
	m map[int64]bool
	s sync.Mutex
}
type safeInt64StringMapStruct struct {
	m map[int64]string
	s sync.Mutex
}
