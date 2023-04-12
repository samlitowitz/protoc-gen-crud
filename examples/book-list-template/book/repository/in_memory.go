package repository

import book_list "github.com/samlitowitz/protoc-gen-crud/examples/book-list"

type InMemory struct {
	booksByUNID1 map[string]*book_list.Book
	booksByUNID2 map[string]*book_list.Book
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

func NewInMemory() *InMemory {
	return &InMemory{
		booksByUNID1: make(map[string]*book_list.Book),
		booksByUNID2: make(map[string]*book_list.Book),
	}
}
