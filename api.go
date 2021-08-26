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
	http.HandleFunc("/API/get_users", returnAllUsers)
	http.HandleFunc("/API/get_user", returnSingleUser)
	http.HandleFunc("/API/get_history_by_tg", returnSingleUserHistory)
	http.HandleFunc("/API/delete_history_by_tg", deleteHistoryField)
	http.HandleFunc("/API/"+os.Getenv("TELEGRAM_BOT_TOKEN")+"/add_admin",
		addAdmin)
}

func addAdmin(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	user_id, ok_id := parseUserId(w, r)
	if !ok_id {
		return
	}
	db, ok_conn := api_connectToDB(w, r)
	if !ok_conn {
		return
	}
	user, ok_user := api_giveUserById(db, user_id, w, r)
	if !ok_user {
		return
	}
	user.UserRole = "admin"
	db.Save(&user)
	w.WriteHeader(http.StatusOK)
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	db, ok_conn := api_connectToDB(w, r)
	if !ok_conn {
		return
	}
	users, ok_users := api_giveUsers(db, w, r)
	if !ok_users {
		return
	}
	sendData(users, w, r)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	user_id, ok_id := parseUserId(w, r)
	if !ok_id {
		return
	}
	db, ok_conn := api_connectToDB(w, r)
	if !ok_conn {
		return
	}
	user, ok_user := api_giveUserById(db, user_id, w, r)
	if !ok_user {
		return
	}
	sendData(user, w, r)
}

func returnSingleUserHistory(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	user_id, ok_id := parseUserId(w, r)
	if !ok_id {
		return
	}
	db, ok_conn := api_connectToDB(w, r)
	if !ok_conn {
		return
	}
	info_list, ok_info := api_giveUserHistoryRet(db, user_id, w, r)
	if !ok_info {
		return
	}
	sendData(info_list, w, r)
}

func deleteHistoryField(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete history field")
}

func sendData(data interface{}, w http.ResponseWriter, r *http.Request) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't marshal json"}`))
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func isMethodGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "method not found"}`))
		return false
	}
	return true
}

func parseUserId(w http.ResponseWriter, r *http.Request) (int, bool) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "arguments params is missing"}`))
		return 0, false
	}
	user_id, err := strconv.Atoi(keys[0])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "can't pars id"}`))
		return 0, false
	}
	return user_id, true
}

func api_connectToDB(w http.ResponseWriter, r *http.Request) (*gorm.DB, bool) {
	db, err_conn := connectToDB()
	if err_conn != nil {
		log.Error(err_conn)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't connect to data base"}`))
		return nil, false
	}
	return db, true
}

func api_giveUsers(db *gorm.DB, w http.ResponseWriter,
	r *http.Request) ([]User, bool) {

	users, err := giveUsers(db)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return []User{}, false
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return []User{}, false
	}
	return users, true
}

func api_giveUserById(db *gorm.DB, user_id int, w http.ResponseWriter,
	r *http.Request) (*User, bool) {

	user, err := giveUserByID(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unrecognized user"}`))
		return nil, false
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return nil, false
	}
	return user, true
}

func api_giveUserHistoryRet(db *gorm.DB, user_id int, w http.ResponseWriter,
	r *http.Request) ([]InfoIP, bool) {

	info_list, err := giveUserHistoryRet(db, user_id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		return info_list, true
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return nil, false
	}
	return info_list, true
}
