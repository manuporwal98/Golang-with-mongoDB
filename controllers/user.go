package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/manuporwal98/Golang-with-mongoDB/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController { //returns pointer to struct
	return &UserController{s} //return address of struct
}

//struct method, checkout video from tech with tim

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) { ///need to checkout on the struct method
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
	}
	oid := bson.ObjectIdHex(id) ///need to check this function
	u := models.User{}
	err := uc.session.DB("mongo-golang").C("users").FindId(oid).One(&u)
	if err != nil {
		w.WriteHeader(404)
		return
	} //database is mongo-golang and C is user    //need to complete this later
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	u.Id = bson.NewObjectId()
	uc.session.DB("mongo-golang").C("users").Insert(u) //inserting u in the collection users
	uj, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)

}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}
	oid := bson.ObjectIdHex(id)
	err := uc.session.DB("mongo-golang").C("users").Remove(oid)
	if err != nil {
		w.WriteHeader(404)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted user", oid, "\n")
}
