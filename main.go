package main

import(
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	_"github.com/go-sql-driver/mysql"
	// "github.com/kataras/go-sessions"

)

var db *sql.DB
var err error

type user struct{
	ID int
	Username string
	Password string
	Email string
}

type response struct{
	Status int `json:"status"`
	Message string `json:"message"`
	Data []user
}

func connect_db(){
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/slivth-login")

	if err != nil{
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil{
		log.Fatalln(err)
	}
}

func checkErr(w http.ResponseWriter, r *http.Request, err error) bool {
	if err != nil {
		fmt.Println(r.Host + r.URL.Path)
		   
   		http.Redirect(w, r, r.Host+r.URL.Path, 301)
		return false
	}
  	 return true
   }

func QueryUser(email string) user{
	var users = user{}
	fmt.Println(users)
	err = db.QueryRow(`	
		SELECT username,
		email,
		password
		FROM users WHERE email=?
		`, email).
	Scan(
		&users.Username,
		&users.Password,
		&users.Email,
	)
	return users
	// fmt.Println(users)
}

func register(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST"{
		var username = r.FormValue("username")
		var password = r.FormValue("password")
		var email = r.FormValue("email")
		var response response

		users := QueryUser(username)
		if (user{}) == users{
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
			if len(hashedPassword) != 0 && checkErr(w, r, err) {
				stmt, err := db.Prepare("INSERT INTO users SET username=?, password=?, email=?")
				
				if err == nil{
					_, err := stmt.Exec(&username, &hashedPassword, &email)
					if err != nil{
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					
					}
					response.Status = 1
					response.Message = "Success Register"
					// response.Data = arr_user

					json.NewEncoder(w).Encode(response)
					} 
				} 
			}
		}
	}

func routes(){
	http.HandleFunc("/register", register)
}

func main(){
	connect_db()
	routes()

	defer db.Close()

	fmt.Println("Server running on port:8080")
	http.ListenAndServe(":8080", nil)
}