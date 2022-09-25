package utils

import (
	"bytes"
	"encoding/gob"
	"errors"
	"math/rand"
	"sync"
)

type unexportedMutex struct {
	s sync.RWMutex
}

func (u *unexportedMutex) RLock() {
	u.s.RLock()
}
func (u *unexportedMutex) RUnlock() {
	u.s.RUnlock()
}
func (u *unexportedMutex) Lock() {
	u.s.Lock()
}
func (u *unexportedMutex) TryLock() bool {
	return u.s.TryLock()
}
func (u *unexportedMutex) Unlock() {
	u.s.Unlock()
}

type SafeValue[T any] struct {
	V T
	unexportedMutex
}

func NewSafeValue[T any]() *SafeValue[T] {
	return &SafeValue[T]{}
}

func NewSafeValueFrom[T any](a T) *SafeValue[T] {
	return &SafeValue[T]{V: a}
}

type SafeValueTools[T comparable] struct {
	SafeValue[T]
}

func NewSafeValueTools[T comparable]() *SafeValueTools[T] {
	return &SafeValueTools[T]{}
}

func NewSafeValueToolsFrom[T comparable](a T) *SafeValueTools[T] {
	return &SafeValueTools[T]{SafeValue: SafeValue[T]{V: a}}
}

func (s *SafeValueTools[T]) Get() T {
	s.RLock()
	defer s.RUnlock()
	return s.V
}

func (s *SafeValueTools[T]) Set(v T) {
	s.Lock()
	defer s.Unlock()
	s.V = v
}

func (s *SafeValueTools[T]) SetIfEquals(v, ifEquals T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V == ifEquals
	if r {
		s.V = v
	}
	return r
}

type SafeNumericValueTools[T Ints | Uints | Floats] struct {
	SafeValue[T]
}

func NewSafeNumericValueTools[T Ints | Uints | Floats]() *SafeNumericValueTools[T] {
	return &SafeNumericValueTools[T]{}
}

func (s *SafeNumericValueTools[T]) Add(v T) {
	s.RLock()
	defer s.RUnlock()
	s.V += v
}

func (s *SafeNumericValueTools[T]) Get() T {
	s.RLock()
	defer s.RUnlock()
	return s.V
}

func (s *SafeNumericValueTools[T]) Set(v T) {
	s.Lock()
	defer s.Unlock()
	s.V = v
}

func (s *SafeNumericValueTools[T]) SetIfEquals(v, ifEquals T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V == ifEquals
	if r {
		s.V = v
	}
	return r
}

func (s *SafeNumericValueTools[T]) SetIfGreater(v, ifGreaterThanThis T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V > ifGreaterThanThis
	if r {
		s.V = v
	}
	return r
}

func (s *SafeNumericValueTools[T]) SetIfGreaterOrEquals(v, ifGreaterOrEqualsThanThis T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V >= ifGreaterOrEqualsThanThis
	if r {
		s.V = v
	}
	return r
}

func (s *SafeNumericValueTools[T]) SetIfLesser(v, ifLesserThanThis T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V < ifLesserThanThis
	if r {
		s.V = v
	}
	return r
}

func (s *SafeNumericValueTools[T]) SetIfLesserOrEquals(v, ifLesserOrEqualsThanThis T) bool {
	s.Lock()
	defer s.Unlock()
	r := s.V <= ifLesserOrEqualsThanThis
	if r {
		s.V = v
	}
	return r
}

type SafeMap[K comparable, V any] struct {
	M map[K]V
	unexportedMutex
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{M: make(map[K]V)}
}

func NewSafeMapFrom[K, V comparable](m map[K]V) *SafeMap[K, V] {
	return &SafeMap[K, V]{M: m}
}

type SafeMapGetSetHas[K, V comparable] struct {
	SafeMap[K, V]
}

func NewSafeMapGetSetHas[K, V comparable]() *SafeMapGetSetHas[K, V] {
	return &SafeMapGetSetHas[K, V]{SafeMap[K, V]{M: make(map[K]V)}}
}

func NewSafeMapGetSetHasFrom[K, V comparable](m map[K]V) *SafeMapGetSetHas[K, V] {
	return &SafeMapGetSetHas[K, V]{SafeMap[K, V]{M: m}}
}

func (s *SafeMapGetSetHas[K, V]) Get(k K, externalLock bool) V {
	if !externalLock {
		s.s.RLock()
		defer s.s.RUnlock()
	}

	return s.M[k]
}

func (s *SafeMapGetSetHas[K, V]) GetHas(k K, externalLock bool) (V, bool) {
	if !externalLock {
		s.s.RLock()
		defer s.s.RUnlock()
	}

	v, h := s.M[k]
	return v, h
}

func (s *SafeMapGetSetHas[K, V]) Has(k K, externalLock bool) bool {
	if !externalLock {
		s.s.RLock()
		defer s.s.RUnlock()
	}

	_, h := s.M[k]
	return h
}

func (s *SafeMapGetSetHas[K, V]) Set(k K, v V, externalLock bool) {
	if !externalLock {
		s.s.Lock()
		defer s.s.Unlock()
	}

	s.M[k] = v
}

func (s *SafeMapGetSetHas[K, V]) Delete(k K, externalLock bool) {
	if !externalLock {
		s.s.Lock()
		defer s.s.Unlock()
	}

	delete(s.M, k)
}

type SafeMapL2[K, K2, V comparable] struct {
	M map[K]map[K2]V
	unexportedMutex
}

/*
	Safe array
*/

type SafeArray[T any] struct {
	L []T
	unexportedMutex
}

func NewSafeArray[T any]() *SafeArray[T] {
	return new(SafeArray[T])
}

func NewSafeArrayFrom[T any](arr []T) *SafeArray[T] {
	return &SafeArray[T]{L: arr}
}

/*
	Safe array tools
*/

type SafeArrayTools[T comparable] struct {
	SafeArray[T]
}

func NewSafeArrayTools[T comparable]() *SafeArrayTools[T] {
	return new(SafeArrayTools[T])
}

func NewSafeArrayToolsFrom[T comparable](arr []T) *SafeArrayTools[T] {
	return &SafeArrayTools[T]{SafeArray: SafeArray[T]{L: arr}}
}

func (s *SafeArrayTools[T]) Get(i int) T {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	return s.SafeArray.L[i]
}

func (s *SafeArrayTools[T]) GetRandom() (T, error) {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	if len(s.SafeArray.L) == 0 {
		var a T
		return a, errors.New("GetRandom() for an empty array is not supported")
	}
	return s.SafeArray.L[rand.Intn(len(s.SafeArray.L))], nil
}

func (s *SafeArrayTools[T]) Set(a []T) {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	s.SafeArray.L = a
}

func (s *SafeArrayTools[T]) Copy() []T {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	return Copy(s.L)
}

func (s *SafeArrayTools[T]) Marshal() []byte {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	return Marshal(s.SafeArray.L)
}

func (s *SafeArrayTools[T]) Len() int {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	return len(s.SafeArray.L)
}

func (s *SafeArrayTools[T]) Append(s1 T) {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()
	s.SafeArray.L = append(s.SafeArray.L, s1)
}

func (s *SafeArrayTools[T]) AppendNoLock(s1 T) {
	s.SafeArray.L = append(s.SafeArray.L, s1)
}

func (s *SafeArrayTools[T]) RemoveByValue(s1 T, single bool) {
	for i, v := range s.SafeArray.L {
		if v == s1 {
			s.SafeArray.L[i] = s.SafeArray.L[len(s.SafeArray.L)-1]
			s.SafeArray.L = s.SafeArray.L[:len(s.SafeArray.L)-1]
			if single {
				break
			}
		}
	}
}

func (s *SafeArrayTools[T]) ToHasMap() map[T]struct{} {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()

	return Slice2HasMap(s.L)
}

func (s *SafeArrayTools[T]) ToHasMapExcludeEmpty() map[T]struct{} {
	s.SafeArray.Lock()
	defer s.SafeArray.Unlock()

	return Slice2HasMapExcludeEmpty(s.L)
}

type SafeUniqueArray[T comparable] struct {
	unexportedMutex
	m        map[T]struct{}
	onUpdate func()
}

func NewSafeUniqueArray[T comparable](onUpdate func()) *SafeUniqueArray[T] {
	return &SafeUniqueArray[T]{m: map[T]struct{}{}, onUpdate: onUpdate}
}

func (s *SafeUniqueArray[T]) CallOnUpdate(doLock bool) {
	if s.onUpdate != nil {
		if doLock {
			s.Lock()
			defer s.Unlock()
		}
		s.onUpdate()
	}
}

func (s *SafeUniqueArray[T]) Add(a T) {
	s.Lock()
	defer s.Unlock()
	if _, h := s.m[a]; !h {
		s.m[a] = struct{}{}
		s.CallOnUpdate(false)
	}
}

func (s *SafeUniqueArray[T]) Set(a []T) {
	s.Lock()
	defer s.Unlock()

	s.m = make(map[T]struct{}, len(a))
	for _, v := range a {
		s.m[v] = struct{}{}
	}
	s.CallOnUpdate(false)
}

func (s *SafeUniqueArray[T]) Remove(a T) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, a)
	s.CallOnUpdate(false)
}

func (s *SafeUniqueArray[T]) GetList() (ret []T) {
	s.Lock()
	defer s.Unlock()

	ret = make([]T, len(s.m))
	i := 0
	for k := range s.m {
		ret[i] = k
		i++
	}
	return
}

func (s *SafeUniqueArray[T]) GobEncode() ([]byte, error) {
	s.Lock()
	defer s.Unlock()

	var buf = bytes.NewBuffer(nil)
	e := gob.NewEncoder(buf).Encode(&s.m)
	return buf.Bytes(), e
}

func (s *SafeUniqueArray[T]) GobDecode(data []byte) error {
	s.Lock()
	defer s.Unlock()

	return gob.NewDecoder(bytes.NewReader(data)).Decode(&s.m)
}

type SafeKeyValue[K, V any] struct {
	unexportedMutex
	K K
	V V
}

func NewSafeKeyValue[K, V any]() *SafeKeyValue[K, V] {
	return &SafeKeyValue[K, V]{}
}

func NewSafeKeyValueFrom[K, V any](k K, v V) *SafeKeyValue[K, V] {
	return &SafeKeyValue[K, V]{K: k, V: v}
}

type KeyValue[K, V any] struct {
	K K
	V V
}

func NewKeyValue[K, V any]() *KeyValue[K, V] {
	return &KeyValue[K, V]{}
}

func NewKeyValueFrom[K, V any](k K, v V) *KeyValue[K, V] {
	return &KeyValue[K, V]{K: k, V: v}
}
