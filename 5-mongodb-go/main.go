package main

import (
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"mongo-golang/controllers"
	"net/http"
)

func main() {
	r := httprouter.New()
	uc := controllers.NewUserController(getSession())
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe(":8000", r)
}

func getSession() *mgo.Session {
	dial, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}
	return dial
}
