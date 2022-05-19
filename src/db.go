package main

import (
	"database/sql"
	"fmt"
	"math/rand"

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
	db.Exec("CREATE TABLE if not exists users( Id INTEGER NOT NULL, Username text UNIQUE NOT NULL, Password text NOT NULL, IsOnline INTEGER NOT NULL, PRIMARY KEY(Id AUTOINCREMENT) )")
	db.Exec("CREATE TABLE if not exists wallet( UserId INTEGER NOT NULL, CardNo	TEXT NOT NULL, Balance	INTEGER NOT NULL, PIN	TEXT NOT NULL, PRIMARY KEY(CardNo), UNIQUE(UserId))")
}

type User struct {
	Id       int
	Username string
	IsOnline int
}
type wallet struct {
	UserId  int
	CardNo  string
	Balance int
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
			var isOnline int
			err = rows.Scan(&id, &isOnline)
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
		} else if len(userList) > 0 {
			user = userList[0]
		}
	} else {
		fmt.Println(err)
	}
	return user
}

func Wallets() []wallet {
	var walletList []wallet
	walletList = nil
	stmt, err := db.Prepare("select UserId,CardNo,Balance from wallet")
	rows, err := stmt.Query()
	if err == nil {
		for rows.Next() {
			var userId int
			var cardNo string
			var balance int
			err = rows.Scan(&userId, &cardNo, &balance)
			if err == nil {
				walletList = append(walletList, wallet{UserId: userId, CardNo: cardNo, Balance: balance})
			} else {
				fmt.Println(err)
			}
		}
		rows.Close()
	} else {
		fmt.Println(err)
	}
	return walletList
}

func Logout(id int) bool {
	stmt, err := db.Prepare("update users set IsOnline=0 where Id = ? and IsOnline = 1")
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
func CreateUser(username string, password string) int64 {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO users(Username, Password, IsOnline) values (?,?,0)")
	checkError(err)
	res, err := stmt.Exec(username, password)
	checkError(err)
	id, err := res.LastInsertId()
	checkError(err)
	return id
}
func CreateWallet(id int64, pin string) (string, error) {
	stmt, err := db.Prepare("INSERT OR IGNORE INTO wallet(UserId,CardNo,Balance,PIN ) VALUES(?,?,100,?)")
	checkError(err)
	fmt.Println(id)
	res, err := stmt.Exec(id, rand.Intn(999999), pin)
	checkError(err)
	_, err = res.RowsAffected()
	checkError(err)
	if err == nil {
		return "Kart oluşturuldu. PIN = " + pin, err
	}
	return "Kart oluşturulamadı", err
}
func BalanceTransferToCardNo(id int64, cardNoSender string, pin string, cardNorecipient string, amount int64) bool {
	stmt, err := db.Prepare("update wallet set Balance = Balance - ? where UserId = ? and CardNo = ? and PIN = ? and Balance >= ?")
	checkError(err)
	res, err := stmt.Exec(amount, id, cardNoSender, pin, amount)
	checkError(err)
	affectedRow, err := res.RowsAffected()
	checkError(err)
	if err == nil && affectedRow > 0 {
		stmt2, err2 := db.Prepare("update wallet set Balance = Balance + ? where CardNo = ? ")
		checkError(err2)
		fmt.Println(id)
		res2, err2 := stmt2.Exec(amount, cardNorecipient)
		checkError(err2)
		affectedRow2, err2 := res2.RowsAffected()
		checkError(err2)
		if err2 == nil && affectedRow2 > 0 {
			return true
		}
	}
	return false
}
