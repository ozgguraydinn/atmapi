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
	db.Exec("CREATE TABLE if not exists users ( Id INTEGER NOT NULL, Username text UNIQUE, Password text, IsOnline INTEGER, PRIMARY KEY(Id AUTOINCREMENT) )")
	db.Exec("CREATE TABLE sqlite_sequence(name,seq)")
	db.Exec("INSERT INTO sqlite_sequence VALUES(users,0)")

}

type User struct {
	Id       int
	Username string
	IsOnline int
}

func Login(username string, password string) User {
	user := User{}
	var userList []User
	userList = nil
	stmt, err := db.Prepare("select Id,IsOnline from users where Username = ? and Password = ?")
	rows, err := stmt.Query(username, password)
	if err == nil {
		for rows.Next() {
			var id int
			var username string
			var isOnline int
			err = rows.Scan(&id, &username, &isOnline)
			if err == nil {
				userList = append(userList, User{Id: id, Username: username, IsOnline: isOnline})
			} else {
				fmt.Println(err)
			}
		}
		rows.Close()
		if len(userList) > 0 && userList[0].IsOnline != 1 {
			stmt, err := db.Prepare("update users set IsOnline=1 where Username = ?")
			checkError(err)
			fmt.Println(username)
			res, err := stmt.Exec(username)
			checkError(err)
			_, err = res.RowsAffected()
			checkError(err)
			if err == nil {
				user = userList[0]
			}
		}
	} else {
		fmt.Println(err)
	}
	return user
}
func Logout(id int) bool {
	stmt, err := db.Prepare("update users set IsOnline=0 where Id = ?")
	checkError(err)
	fmt.Println(id)
	res, err := stmt.Exec(id)
	checkError(err)
	_, err = res.RowsAffected()
	checkError(err)
	if err == nil {
		return true
	}
	return false
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
