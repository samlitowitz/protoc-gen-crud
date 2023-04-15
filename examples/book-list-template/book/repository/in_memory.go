package repository

import book_list "github.com/samlitowitz/protoc-gen-crud/examples/book-list"

type NextField1 func() (int, error)
type NextField2 func() (string, error)

type InMemory struct {
	booksByUID1 map[string]*book_list.Book
	booksByUID2 map[string]*book_list.Book

	nextField1 NextField1
	nextField2 NextField2
}

func NewInMemory(nextField1 NextField1, nextField2 NextField2) *InMemory {
	return &InMemory{
		booksByUID1: make(map[string]*book_list.Book),
		booksByUID2: make(map[string]*book_list.Book),
		nextField1:  nextField1,
		nextField2:  nextField2,
	}
}

func (repo *InMemory) Create(books []*book_list.Book) ([]*book_list.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *InMemory) Read() {
	//TODO implement me
	panic("implement me")
}

func (repo *InMemory) Update(books []*book_list.Book) ([]*book_list.Book, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *InMemory) Delete(strings []string) error {
	//TODO implement me
	panic("implement me")
}
