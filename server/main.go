package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "./database"
	handler "./handler"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	db.ConnectDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/getUsers", handler.GetUsers)
	r.HandleFunc("/addUser", handler.AddUser)
	r.HandleFunc("/modifyUser", handler.ModifyUser)
	r.HandleFunc("/deleteUser", handler.DeleteUser)
	r.HandleFunc("/getUnreservedItems", handler.GetUnreservedItems)
	r.HandleFunc("/getReservedItems", handler.GetReservedItems)
	r.HandleFunc("/addItem", handler.AddItem)
	r.HandleFunc("/modifyItem", handler.ModifyItem)
	r.HandleFunc("/deleteItem", handler.DeleteItem)
	r.HandleFunc("/getLastNBookings", handler.GetLastNBookings)
	r.HandleFunc("/getUserDebts", handler.GetUserDebts)
	r.HandleFunc("/checkout", handler.Checkout)
	r.HandleFunc("/pay", handler.Pay)
	r.HandleFunc("/deleteBookEntry", handler.DeleteBookEntry)
	r.HandleFunc("/updateFavoriteItems", handler.UpdateFavoriteItems)
	r.HandleFunc("/getFavoriteItemIDs", handler.GetFavoriteItemIDs)
	r.HandleFunc("/deleteUserFromFavoriteItems", handler.DeleteUserFromFavoriteItems)
	r.HandleFunc("/getFeedback", handler.GetFeedback)
	r.HandleFunc("/addFeedback", handler.AddFeedback)
	r.HandleFunc("/deleteFeedback", handler.DeleteFeedback)
	r.HandleFunc("/login", handler.Login)
	r.HandleFunc("/changeAdminPassword", handler.ChangeAdminPassword)

	serveIndexHTML := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		http.ServeFile(w, r, "./../client/dist/index.html")
	}
	r.PathPrefix("/").Handler(isAuthorized(serveIndexHTML))
	// r.PathPrefix("/").Handler(handler.CustomFileServer(http.Dir("./../client/dist"), serveIndexHTML))

	server := &http.Server{
		Addr:           ":8081",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(server.ListenAndServeTLS("server.crt", "server.key"))
}

var mySigningKey = []byte("captainjacksparrowsayshi")

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
