package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func initBackendApi() {
	http.HandleFunc("/API/", homePageAPI)
	http.HandleFunc("/API/get_users", returnAllUsers)
	http.HandleFunc("/API/get_user", returnSingleUser)
	http.HandleFunc("/API/get_history_by_tg", returnSingleUserHistory)
	http.HandleFunc("/API/delete_history_by_tg", deleteHistoryField)
	http.HandleFunc("/API/"+os.Getenv("TELEGRAM_BOT_TOKEN")+"/add_admin",
		addAdmin)
}

func addAdmin(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		log.Error("Argument params is missing")
		fmt.Fprintf(w, "Error: argument params is missing\n")
		return
	}
	user_id, err_id := strconv.Atoi(keys[0])
	if err_id != nil {
		log.Error(err_id)
		fmt.Fprintf(w, "Error: while parse id\n")
		return
	}
	db, err_conn := connectToDB()
	if err_conn != nil {
		log.Error(err_conn)
		fmt.Fprintf(w, "Error: while connecting to data base\n")
		return
	}
	user, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Fprintf(w, "Error: unrecognized user\n")
		return
	} else if err != nil {
		log.Error(err)
		fmt.Fprintf(w, "Error\n")
		return
	}
	user.UserRole = "admin"
	db.Save(&user)
	fmt.Fprintf(w, fmt.Sprintf("Success\n"))
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
