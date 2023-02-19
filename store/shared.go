package store

type EventWrapper[T any] struct {
	EventTime int64 // nanos
	Event     T
}
