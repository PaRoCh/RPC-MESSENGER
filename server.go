package main

import (
	"context"
	"encoding/json"
	"fmt"

	//"github.com/jinzhu/gorm"
	"log"
	"os"
	//"strconv"
	guuid "github.com/google/uuid"
	//"math/rand"
	"net/http"
	//"strconv"
	//"database/sql"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx"
	//"github.com/rs/cors"
	//_ "github.com/mattn/go-sqlite3"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func genUUID() string {
	id := guuid.New()
	return id.String()
}

//Author *Author `json:"author"`

//var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("hello Roma Lox")
}

func login(w http.ResponseWriter, r *http.Request) {
	type Person struct {
		NicknameOrMail 	string `json:"nicknameormail"`
		Password     	string `json:"password"`
	}
	type PersonData struct {
		Id			string `json:"id"`
		Nickname    string `json:"nickname"`
		Mail    	string `json:"mail"`
		Password 	string `json:"password"`
	}
	w.Header().Set("Content-Type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	log.Println(person)

	res, err := db.Exec(context.Background(),`select * from users where (nickname = $1 and password = $2) or (mail = $1 and password = $2)`,
		person.NicknameOrMail, person.Password)
	if err!= nil {
		fmt.Fprintf(os.Stderr, "Unable to get person: %v\n", err)
		json.NewEncoder(w).Encode(err)
	}else if res.RowsAffected()==1{
		fmt.Fprintf(os.Stderr, "Result: %v\n", res)
		res, err := db.Query(context.Background(),`select * from users where (nickname = $1 and password = $2) or (mail = $1 and password = $2)`,
			person.NicknameOrMail, person.Password)
		if err!= nil {
			fmt.Fprintf(os.Stderr, "Unable to get person: %v\n", err)
			json.NewEncoder(w).Encode(err)
		}else{
			var pd PersonData
			for res.Next() {
				err := res.Scan(&pd.Id,&pd.Nickname,&pd.Mail,&pd.Password)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to get personn: %v\n", err)
				}
			}
			json.NewEncoder(w).Encode(&pd)
			fmt.Printf("Succesful\n")
		}
	}else{
		json.NewEncoder(w).Encode("nety takogo")
	}
}
func auth(w http.ResponseWriter, r *http.Request) {
	type Person struct {
		NicknameOrMail 	string `json:"nicknameormail"`
		Password     	string `json:"password"`
	}
	type PersonData struct {
		Id			string `json:"id"`
		Nickname    string `json:"nickname"`
		Mail    	string `json:"mail"`
		Password 	string `json:"password"`
	}
	w.Header().Set("Content-Type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	log.Println(person)

	res, err := db.Exec(context.Background(),`select * from users where (nickname = $1 and password = $2) or (mail = $1 and password = $2)`,
		person.NicknameOrMail, person.Password)
	if err!= nil {
		fmt.Fprintf(os.Stderr, "Unable to get person: %v\n", err)
		json.NewEncoder(w).Encode(err)
	}else if res.RowsAffected()==1{
		fmt.Fprintf(os.Stderr, "Result: %v\n", res)
		res, err := db.Query(context.Background(),`select * from users where (nickname = $1 and password = $2) or (mail = $1 and password = $2)`,
			person.NicknameOrMail, person.Password)
		if err!= nil {
			fmt.Fprintf(os.Stderr, "Unable to get person: %v\n", err)
			json.NewEncoder(w).Encode(err)
		}else{
			var pd PersonData
			for res.Next() {
				err := res.Scan(&pd.Id,&pd.Nickname,&pd.Mail,&pd.Password)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Unable to get personn: %v\n", err)
				}
			}
			json.NewEncoder(w).Encode(&pd)
			fmt.Printf("Succesful\n")
		}
	}else{
		json.NewEncoder(w).Encode("nety takogo")
	}
}
func register(w http.ResponseWriter, r *http.Request) {
	type Person struct {
		Id			string `json:"id"`
		Nickname    string `json:"nickname"`
		Mail    	string `json:"mail"`
		Password 	string `json:"password"`
		Token 		string `json:"token"`
	}
	w.Header().Set("Content-Type", "application/json")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)
	person.Id=genUUID()
	person.Token=genUUID()
	log.Println(person)
	_, err := db.Exec(context.Background(),`INSERT INTO users (id, nickname, mail, password) VALUES ($1, $2, $3, $4)`,
		person.Id,person.Nickname, person.Mail, person.Password)
	if err!= nil {
		fmt.Fprintf(os.Stderr, "Unable to create person: %v\n", err)
		json.NewEncoder(w).Encode(err)
	}else {
		_, err = db.Exec(context.Background(),`INSERT INTO tokens (id, token) VALUES ($1,$2)`,
			person.Id,person.Token)
		if err!= nil {
			fmt.Fprintf(os.Stderr, "Unable to create token: %v\n", err)
			json.NewEncoder(w).Encode(err)
		}else{
			json.NewEncoder(w).Encode(&person)
			fmt.Printf("Succesful\n")
		}
	}
}
var db = setDB()
func setDB() *pgx.Conn {
	db, err := pgx.Connect(context.Background(),"postgresql://postgres:2058@localhost:5432/messenger")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
func configDB() {
	_, err := db.Exec(context.Background(),`create table users (
    id text primary key not null unique,
    nickname text not null,
	mail text not null unique,
    password text not null
	);`)
	//first_name text not null,
	//	last_name text not null);
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create users table: %v\n", err)
		//os.Exit(1)
	}
	fmt.Printf("Successfully created users table\n")

	_, err = db.Exec(context.Background(),`create table tokens (
    id text not null unique,
    token text not null unique
	);`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create tokens table: %v\n", err)
		//os.Exit(1)
	}
	fmt.Printf("Successfully created tokens table\n")
}
func main() {
	//url := "Server=192.168.1.163;Port=5432;Database=postgres;User Id=postgres;Password=root"

	//configDB()



	//db.Close(context.Background())
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
		//log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/", getBooks).Methods("GET")
	r.HandleFunc("/log", login).Methods("POST")
	r.HandleFunc("/reg", register).Methods("POST")
	r.HandleFunc("/auth", auth).Methods("POST")
	// Start server

	log.Fatal(http.ListenAndServe(":"+port, r))
}
