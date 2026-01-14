package models

type Response[T any] struct {
	Data  T
	Error error
}

func (r Response[T]) HasData() bool {
	data := &r.Data
	return data != nil
}

func (r Response[T]) HasError() bool {
	return r.Error != nil
}
