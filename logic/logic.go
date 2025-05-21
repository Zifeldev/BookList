package logic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type Booklist []Book

var books Booklist

func (b *Booklist) Load() {
	data, err := os.ReadFile("db.json")
	if err != nil {
		fmt.Println("db not loaded...")
		return

	}
	err = json.Unmarshal(data, &b)
	if err != nil {
		panic(err)
	}
}

func (b *Booklist) Save() {
	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("db.json", data, 2)
	if err != nil {
		panic(err)
	}
}

func GetBook(c *gin.Context) {
	books.Load()
	c.JSON(http.StatusOK, books)
}

func PostBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	newBook.ID = len(books) + 1
	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
	books.Save()
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")

	for _, i := range books {
		id, err := strconv.Atoi(id)
		if err != nil {
			return
		}
		if i.ID == id {
			c.JSON(http.StatusOK, i)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
}

func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for index, book := range books {
		if book.ID == id {
			books[index].Title = updatedBook.Title
			books[index].Author = updatedBook.Author
			books.Save()
			c.JSON(http.StatusOK, books[index])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}

func DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	for index, book := range books {
		if book.ID == id {

			books = append(books[:index], books[index+1:]...)

			file, err := os.OpenFile("db.json", os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
				return
			}
			defer file.Close()

			newData, err := json.Marshal(books)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal data"})
				return
			}

			file.Truncate(0)
			file.Seek(0, 0)
			_, err = file.Write(newData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write data"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}
