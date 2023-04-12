package book

import book_list "github.com/samlitowitz/protoc-gen-crud/examples/book-list"

type Repository interface {
	Create([]*book_list.Book) ([]*book_list.Book, error)
	Read()
	Update([]*book_list.Book) ([]*book_list.Book, error)
	Delete([]string) error
}
