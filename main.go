package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/gorilla/mux"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

//setting up conection with database

var db *sql.DB

func getMySQLDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/Practice2") // sql.Open("database name = mysql","root:@(server appdress = 127.0.0.1:portNo = 3306)/databsae name in which table is created")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type studentInfo struct {
	Sid    string
	Name   string
	Course string
}

func init() {
	tmpl = template.Must(template.ParseFiles("crude.html"))
}

func main() {
	rout := mux.NewRouter()

	defer db.Close()

	// setting the static folder

	fs := http.FileServer(http.Dir("assets"))
	rout.PathPrefix("/assets/").Handler(http.StripPrefix("/assets", fs))

	http.HandleFunc("/", homeHandler)

	http.ListenAndServe(":9999", nil)

}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	db = getMySQLDB()

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	student := studentInfo{
		Sid:    r.FormValue("Sid"),
		Name:   r.FormValue("name"),
		Course: r.FormValue("course"),
	}

	if r.FormValue("submit") == "Insert" {
		sid, _ := strconv.Atoi(student.Sid)

		_, err := db.Exec("insert into studentinfo1(Sid,Name,Course) values(?,?,?)", sid, student.Name, student.Course)

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: "Record Inserted"})
		}

	} else if r.FormValue("submit") == "Read" {

		data := []string{}

		rows, err := db.Query("select * from studentinfo1")

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {

			s := studentInfo{}

			data = append(data, "<table border=1>")                                                           // creating a table in slice
			data = append(data, "<tr><th> Student id </th><th>Student Name</td><th>Student Course</th></tr>") // defining column structure

			for rows.Next() {
				rows.Scan(&s.Sid, &s.Name, &s.Course)
				data = append(data, fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td>", s.Sid, s.Name, s.Course)) // assigning value to table
			}
			data = append(data, " </table>") // closing table
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: strings.Trim(fmt.Sprint(data), "[]")})

		}

	} else if r.FormValue("submit") == "Update" {

		sid, _ := strconv.Atoi(student.Sid)
		result, err := db.Exec("update studentinfo1 set name=?,course=? where sid=?", student.Name, student.Course, sid)

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			_, er := result.RowsAffected()

			if er != nil {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "record not pdated"})
			} else {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Welcome to Update"})
			}

		}

	} else if r.FormValue("submit") == "Delet" {
		sid, _ := strconv.Atoi(student.Sid)
		result, err := db.Exec("delete from studentinfo1  where sid=?",sid)

		if err != nil {
			tmpl.Execute(w, struct {
				Success bool
				Message string
			}{Success: true, Message: err.Error()})
		} else {
			_, er := result.RowsAffected()

			if er != nil {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "record not deleted"})
			} else {
				tmpl.Execute(w, struct {
					Success bool
					Message string
				}{Success: true, Message: "Record Deleted"})
			}

		}
	}
	// tmpl.Execute(w, struct {
	// 	Success bool
	// 	Message string
	// }{Success: true, Message: "Welcome"})
	fmt.Println(student)

}
