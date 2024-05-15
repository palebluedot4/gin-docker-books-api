package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-docker-books-api/pkg/model"
)

func GetBooksHandler(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, model.Books)
}

func FindBookByISBN(isbn string) *model.Book {
	for _, book := range model.Books {
		if book.ISBN == isbn {
			return book
		}
	}
	return nil
}

func GetBookByISBNHandler(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	book := FindBookByISBN(isbn)
	if book == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, book)
}

func CreateBookHandler(ctx *gin.Context) {
	var newBook *model.Book

	if err := ctx.BindJSON(&newBook); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid book data"})
		return
	}

	model.Books = append(model.Books, newBook)
	ctx.IndentedJSON(http.StatusOK, newBook)
}

func DeleteBookHandler(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	index := -1

	for i, book := range model.Books {
		if book.ISBN == isbn {
			index = i
			break
		}
	}

	if index == -1 {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	model.Books = append(model.Books[:index], model.Books[index+1:]...)
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "book deleted"})
}

func UpdateBookHandler(ctx *gin.Context) {
	isbn := ctx.Param("isbn")
	book := FindBookByISBN(isbn)
	if book == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found", "isbn": isbn})
		return
	}

	var updatedBook model.Book
	if err := ctx.BindJSON(&updatedBook); err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid update data"})
		return
	}

	book.ISBN = updatedBook.ISBN
	book.Title = updatedBook.Title
	book.Author = updatedBook.Author
	book.Price = updatedBook.Price
	book.Stock = updatedBook.Stock

	ctx.IndentedJSON(http.StatusOK, book)
}

func CheckoutBookHandler(ctx *gin.Context) {
	isbn := ctx.Query("isbn")
	if isbn == "" {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "missing isbn query parameter"})
		return
	}

	book := FindBookByISBN(isbn)
	if book == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "Book not found", "isbn": isbn})
		return
	}

	if book.Stock <= 0 {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"error": "book not available"})
		return
	}
	book.Stock -= 1

	ctx.IndentedJSON(http.StatusOK, book)
}
