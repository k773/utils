package utils

import (
	"errors"
	"math/rand"
	"sync"
)

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
type SafeStringBoolMap struct {
	M map[string]bool
	S sync.Mutex
}

func (m *SafeStringBoolMap) Set(k string, v bool) {
	m.S.Lock()
	defer m.S.Unlock()
	m.M[k] = v
}

func (m *SafeStringBoolMap) Get(k string) bool {
	m.S.Lock()
	defer m.S.Unlock()
	return m.M[k]
}

func (m *SafeStringBoolMap) GetOk(k string) (bool, bool) {
	m.S.Lock()
	defer m.S.Unlock()
	a, b := m.M[k]
	return a, b
}

type SafeStringStringMap struct {
	M map[string]string
	S sync.Mutex
}

func (m *SafeStringStringMap) Set(k, v string) {
	m.S.Lock()
	defer m.S.Unlock()
	m.M[k] = v
}

func (m *SafeStringStringMap) Get(k string) string {
	m.S.Lock()
	defer m.S.Unlock()
	return m.M[k]
}

func (m *SafeStringStringMap) GetOk(k string) (string, bool) {
	m.S.Lock()
	defer m.S.Unlock()
	a, b := m.M[k]
	return a, b
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

type SafeStringArray struct {
	A []string
	sync.Mutex
}

func (s *SafeStringArray) Get(i int) string {
	s.Lock()
	defer s.Unlock()
	return s.A[i]
}

func (s *SafeStringArray) GetRandom() (string, error) {
	s.Lock()
	defer s.Unlock()
	if len(s.A) == 0 {
		return "", errors.New("GetRandom() for an empty array is not supported")
	}
	return s.A[rand.Intn(len(s.A))], nil
}

func (s *SafeStringArray) Fill(a []string) {
	s.Lock()
	defer s.Unlock()
	s.A = a
}

func (s *SafeStringArray) Copy() []string {
	s.Lock()
	defer s.Unlock()
	var a = make([]string, len(s.A))
	copy(a, s.A)
	return a
}

func (s *SafeStringArray) Marshal() []byte {
	s.Lock()
	defer s.Unlock()
	return Marshal(s.A)
}

func (s *SafeStringArray) Len() int {
	s.Lock()
	defer s.Unlock()
	return len(s.A)
}

func (s *SafeStringArray) Append(s1 string) {
	s.Lock()
	defer s.Unlock()
	s.A = append(s.A, s1)
}

func (s *SafeStringArray) AppendNoLock(s1 string) {
	s.A = append(s.A, s1)
}

func (s *SafeStringArray) RemoveByValue(s1 string) {
	for i, v := range s.A {
		if v == s1 {
			s.A[i] = s.A[len(s.A)-1]
			s.A = s.A[:len(s.A)-1]
			break
		}
	}
}

func (s *SafeStringArray) GenerateValueToNothingMap() map[string]struct{} {
	s.Lock()
	defer s.Unlock()

	var m = map[string]struct{}{}
	for _, v := range s.A {
		m[v] = struct{}{}
	}
	return m
}

func (s *SafeStringArray) GenerateValueToNothingMapExcludeEmpty() map[string]struct{} {
	s.Lock()
	defer s.Unlock()

	var m = map[string]struct{}{}
	for _, v := range s.A {
		if v != "" {
			m[v] = struct{}{}
		}
	}
	return m
}
