package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func api() {
	fmt.Println("Api initializing...")

	r := mux.NewRouter()
	r.HandleFunc("/", healthCheck)
	r.HandleFunc("/healthCheck", healthCheck)
	r.HandleFunc("/users", Users)
	http.ListenAndServe(":8080", r)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	success := "true"
	msg := "Health Check Başarılı."
	a := `{"success":` + success + ` , "message":` + msg + `}`
	w.Write([]byte(a))
}

//urlden değişken alma /* vars := mux.Vars(r) id := vars["id"] */
//byte dönme /* w.Write([]byte(a)) */

func Users(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)

	if r.Method == "POST" {
		fmt.Println("---Post---")
		log.Println("kelime:", r.FormValue("kelime"))
		fmt.Fprintf(w, "yeap")

	} else if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		userControl := GetUserList(r.FormValue("username"), r.FormValue("password"))
		var result string
		if len(userControl) > 0 {
			result = "Başarılı giriş"
		} else {
			result = "Başarısız giriş"
		}
		jsonString, err := json.Marshal(result)
		checkError(err)
		fmt.Fprintf(w, string(jsonString))
	}

}
