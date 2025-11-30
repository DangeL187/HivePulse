package domain

type Flusher[T any] interface {
	Flush(batch []*T)
}
