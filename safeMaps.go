package utils

import "sync"

type SafeByteInt64MapStruct struct {
	M map[byte]int64
	S sync.Mutex
}
type SafeByteIntMapStruct struct {
	M map[byte]int
	S sync.Mutex
}
type SafeByteBoolMapStruct struct {
	M map[byte]bool
	S sync.Mutex
}
type SafeByteStringMapStruct struct {
	M map[byte]string
	S sync.Mutex
}

type SafeStringInt64MapStruct struct {
	M map[string]int64
	S sync.Mutex
}
type SafeStringIntMapStruct struct {
	M map[string]int
	S sync.Mutex
}
type SafeStringBoolMapStruct struct {
	M map[string]bool
	S sync.Mutex
}
type SafeStringStringMapStruct struct {
	M map[string]string
	S sync.Mutex
}

type SafeIntInt64MapStruct struct {
	M map[int]int64
	S sync.Mutex
}
type SafeIntIntMapStruct struct {
	M map[int]int
	S sync.Mutex
}
type SafeIntBoolMapStruct struct {
	M map[int]bool
	S sync.Mutex
}
type SafeIntStringMapStruct struct {
	M map[int]string
	S sync.Mutex
}

type SafeInt64Int64MapStruct struct {
	M map[int64]int64
	S sync.Mutex
}
type SafeInt64IntMapStruct struct {
	M map[int64]int
	S sync.Mutex
}
type SafeInt64BoolMapStruct struct {
	M map[int64]bool
	S sync.Mutex
}
type SafeInt64StringMapStruct struct {
	M map[int64]string
	S sync.Mutex
}

type SafeCounter struct {
	Value int
	s     sync.Mutex
}

func (c *SafeCounter) Increase() {
	c.s.Lock()
	c.Value++
	c.s.Unlock()
}

func (c *SafeCounter) Decrease() {
	c.s.Lock()
	c.Value--
	c.s.Unlock()
}
