package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sahilsnghai/bookstore/pkg/models"
	"github.com/sahilsnghai/bookstore/pkg/utils"
)

var NewBook models.Book
var contentType string = "content-type"
var pkglication string = "pkglication/json"

func GetBook(w http.ResponseWriter, r *http.Request) {
	newBook := models.GetAllBooks()
	res, _ := json.Marshal(&newBook)
	w.Header().Set(contentType, pkglication)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["Id"]
	BId, err := strconv.ParseInt(idStr, 0, 0)
	if err != nil {
		log.Fatal("Error while parsing", err)
	}

	bookDetails, _ := models.GetBookById(BId)
	res, _ := json.Marshal(bookDetails)
	w.Header().Set(contentType, pkglication)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	CreateBook := &models.Book{}
	utils.ParscBody(r, CreateBook)
	book := CreateBook.CreateBook()
	res, _ := json.Marshal(book)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParscBody(r, updateBook)
	vars := mux.Vars(r)
	idStr := vars["Id"]
	BId, err := strconv.ParseInt(idStr, 0, 0)
	if err != nil {
		log.Fatal("Error while parsing", err)
	}

	bookDetails, db := models.GetBookById(BId)
	if updateBook.Name != "" {
		bookDetails.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		bookDetails.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		bookDetails.Publication = updateBook.Publication
	}
	db.Save(&bookDetails)

	res, _ := json.Marshal(bookDetails)
	w.Header().Set(contentType, pkglication)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["Id"]
	BId, err := strconv.ParseInt(idStr, 0, 0)
	if err != nil {
		log.Fatal("Error while parsing", err)
	}
	book := models.DeleteBook(BId)
	res, _ := json.Marshal(book)
	w.Header().Set(contentType, pkglication)
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
