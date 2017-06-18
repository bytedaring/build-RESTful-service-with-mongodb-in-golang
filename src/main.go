package main

import (
	"log"
	"net/http"

	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
	Age   int    `json:"age,omitempty"`
}

var session *mgo.Session

func init() {
	s, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal(err)
	}
	session = s
}

func main() {
	defer session.Close()
	ensureIndex(session)

	session.SetMode(mgo.Monotonic, true)

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/users", allUsers)
	e.GET("/user/:id", getUser)
	e.PUT("/user", updateUser)
	e.DELETE("/user/:id", deleteUser)
	e.POST("/user", saveUser)

	e.Logger.Fatal(e.Start(":1424"))
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("store").C("users")

	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func saveUser(e echo.Context) error {
	u := new(User)
	if err := e.Bind(u); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	s := session.Copy()
	defer s.Close()

	c := s.DB("store").C("users")
	err := c.Insert(u)
	if err != nil {
		log.Println("Failed insert user", u)
		if mgo.IsDup(err) {
			return e.JSON(http.StatusBadRequest, "User with this id alread exists.")
		}

		return e.JSON(http.StatusInternalServerError, "Database error")
	}

	return e.JSON(http.StatusCreated, "SUCCESS")
}

func getUser(e echo.Context) error {
	s := session.Copy()
	defer s.Clone()

	c := s.DB("store").C("users")
	var u User
	id := e.Param("id")
	fmt.Println("userid", id)
	err := c.Find(bson.M{"_id": id}).One(&u)
	if err != nil {
		log.Println("Failed get user", err)
		return e.JSON(http.StatusNotFound, "Database error")
	}

	return e.JSON(http.StatusOK, u)
}

func updateUser(e echo.Context) error {
	u := new(User)
	if err := e.Bind(u); err != nil {
		return e.JSON(http.StatusBadRequest, err)
	}

	s := session.Copy()
	defer s.Close()

	c := s.DB("store").C("users")
	err := c.Update(bson.M{"_id": u.ID}, &u)
	if err != nil {
		switch err {
		default:
			log.Fatalln("Failed update user: ", err)
			return e.JSON(http.StatusInternalServerError, "Database error")
		case mgo.ErrNotFound:
			return e.JSON(http.StatusNotFound, "Not found")
		}
	}
	return e.JSON(http.StatusOK, u)
}

func deleteUser(e echo.Context) error {
	s := session.Copy()
	defer s.Close()

	id := e.Param("id")

	c := s.DB("store").C("users")
	err := c.Remove(bson.M{"_id": id})
	if err != nil {
		switch err {
		default:
			e.JSON(http.StatusInternalServerError, "Database error")
			log.Fatalln("Failed delete user: ", err)
			return err
		case mgo.ErrNotFound:
			e.JSON(http.StatusInternalServerError, "User not found")
			return err
		}
	}

	return e.JSON(http.StatusOK, "Sucess")
}

func allUsers(e echo.Context) error {
	s := session.Copy()
	defer s.Close()

	c := s.DB("store").C("users")

	var users []User
	err := c.Find(bson.M{}).All(&users)
	if err != nil {
		e.JSON(http.StatusInternalServerError, "Database Error")
		return err
	}

	return e.JSON(http.StatusOK, users)
}

