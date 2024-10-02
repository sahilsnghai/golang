package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sahilsnghai/mongo-golang/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	fmt.Println("Start Running at port 8008")

	http.ListenAndServe(":8008", r)

}

func getSession() *mgo.Session {

	session, err := mgo.Dial("mongodb://root:mypass@localhost:27017/go-mongo")
	if err != nil {
		panic(err)
	}
	return session
}
