package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func api() {
	fmt.Println("Api initializing...")
	r := mux.NewRouter()
	r.HandleFunc("/", healthCheck)
	r.HandleFunc("/users/login", UserLogin)
	r.HandleFunc("/users/logout", UserLogout)
	http.ListenAndServe(":8080", r)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.Method + " Sistem aktif."))
}

//urlden değişken alma /* vars := mux.Vars(r) id := vars["id"] */
//byte dönme /* w.Write([]byte(a)) */

func UserLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		userControl := Login(r.FormValue("username"), r.FormValue("password"))
		jsonString, err := json.Marshal(userControl)
		checkError(err)
		fmt.Fprintf(w, string(jsonString))

	} else {
		w.Write([]byte("Desteklenmeyen method: " + r.Method))
	}

}
func UserLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	if r.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		userId, err := strconv.Atoi(r.FormValue("id"))
		checkError(err)
		userControl := Logout(userId)
		var result string
		if userControl {
			result = "Başarılı çıkış"
		} else {
			result = "Başarısız çıkış"
		}
		jsonString, err := json.Marshal(result)
		checkError(err)
		fmt.Fprintf(w, string(jsonString))

	} else {
		w.Write([]byte("Desteklenmeyen method: " + r.Method))
	}

}
