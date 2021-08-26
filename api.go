package main

import (
	"encoding/json"
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
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "method not found"}`))
	}
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "arguments params is missing"}`))
		return
	}
	user_id, err_id := strconv.Atoi(keys[0])
	if err_id != nil {
		log.Error(err_id)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "can't pars id"}`))
		return
	}
	db, err_conn := connectToDB()
	if err_conn != nil {
		log.Error(err_conn)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't connect to data base"}`))
		return
	}
	user, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unrecognized user"}`))
		return
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return
	}
	user.UserRole = "admin"
	db.Save(&user)
	w.WriteHeader(http.StatusOK)
}

func homePageAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(`{"message": "Welcome to the HomePageAPI!}`))
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "method not found"}`))
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "method not found"}`))
	}
	db, err_conn := connectToDB()
	if err_conn != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't connect to data base"}`))
		return
	}
	users, err := giveUsers(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return
	}
	data, err_m := json.Marshal(users)
	if err_m != nil {
		log.Error(err_m)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't marshal json"}`))
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "method not found"}`))
	}
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "arguments params is missing"}`))
		return
	}
	user_id, err_id := strconv.Atoi(keys[0])
	if err_id != nil {
		log.Error(err_id)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "can't pars id"}`))
		return
	}
	db, err_conn := connectToDB()
	if err_conn != nil {
		log.Error(err_conn)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't connect to data base"}`))
		return
	}
	user, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unrecognized user"}`))
		return
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return
	}
	data, err_m := json.Marshal(user)
	if err_m != nil {
		log.Error(err_m)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't marshal json"}`))
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)
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
