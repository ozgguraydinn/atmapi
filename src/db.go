package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func databaseConnect() {
	fmt.Println("database bağlanıyor...")
	var err error
	db, err = sql.Open("sqlite3", "../atm.db")
	checkError(err)
	_ = db
	tableCreate(db)
}

func tableCreate(db *sql.DB) {
	db.Exec("create table if not exists users (username text, password text)")
}

type User struct {
	Username string
	Password string
}

func GetUserList(username string, password string) []User {

	var userList []User
	userList = nil
	fmt.Println(username)
	fmt.Println(password)
	stmt, err := db.Prepare("select * from users where username = ? and password = ?")
	rows, err := stmt.Query(username, password)
	if err == nil {
		fmt.Println(rows.Columns())
		for rows.Next() {
			var username string
			var password string
			err = rows.Scan(&username, &password)
			if err == nil {
				userList = append(userList, User{Username: username, Password: password})
			} else {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(err)
	}
	rows.Close()
	return userList
}

func addDictionary(word, wordMean string) {

	stmt, err := db.Prepare("INSERT INTO tblSozluk(kelime, kelimeAnlam) values (?,?)")
	checkError(err)
	res, err := stmt.Exec(word, wordMean)
	checkError(err)

	id, err := res.LastInsertId()
	checkError(err)

	fmt.Println("Son eklenen kayıt: ", id)

}

func deleteDictionary(wordID int) {

	stmt, err := db.Prepare("DELETE FROM TBLsOZLUK where kelimeID=?")
	checkError(err)
	res, err := stmt.Exec(wordID)
	checkError(err)

	_, err = res.RowsAffected()
	checkError(err)

	fmt.Println("Kayıt silindi")

}

func updateDictionary(wordID int, word, wordMean string) {

	stmt, err := db.Prepare("update tblSoz set kelime=? .....................")
	checkError(err)

	res, err := stmt.Exec(word, wordMean, wordID)
	checkError(err)

	_, err = res.RowsAffected()
	checkError(err)

	fmt.Println("Veri güncellendi")

}
