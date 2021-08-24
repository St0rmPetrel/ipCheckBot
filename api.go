package main

import (
	"fmt"
	"log"
	"net/http"
)

func initBackendApi() {
	http.HandleFunc("/API/", homePageAPI)
	http.HandleFunc("/API/get_users", returnAllUsers)
	http.HandleFunc("/API/get_user", returnSingleUser)
	http.HandleFunc("/API/get_history_by_tg", returnSingleUserHistory)
	http.HandleFunc("/API/delete_history_by_tg", deleteHistoryField)
}

func homePageAPI(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePageAPI!")
}
func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Return All Users")
}
func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	user_id := keys[0]
	fmt.Fprintf(w, fmt.Sprintf("Return a User id = %v", user_id))
}
func returnSingleUserHistory(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}
	user_id := keys[0]
	fmt.Fprintf(w, fmt.Sprintf("Return History for User with id = %v", user_id))
}
func deleteHistoryField(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete history field")
}
