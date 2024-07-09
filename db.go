package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Student struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type myDb struct {
	conn *sql.DB
}

func (student Student) String() string {
	return fmt.Sprintf("ID: %d, Name: %s, Age: %s", student.Id, student.Name, student.Age)
}

// const (
// 	host     = os.Getenv("DB_HOST")
// 	port, _  = strconv.Atoi(os.Getenv("DB_PORT"))
// 	user     = os.Getenv("DB_USER")
// 	password = os.Getenv("DB_PASSWORD")
// 	dbName   = os.Getenv("DB_NAME")
// )

func DBConn() *myDb {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")
	return &myDb{conn: db}
}

func (db *myDb) maybeCreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS Student (
		id SERIAL NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		age TEXT,
		grade TEXT,
		gender TEXT
	);`
	_, err := db.conn.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Table created or already exists")
}

func (db *myDb) GetById(id int) (Student, error) {
	var student Student
	query := "SELECT id, name, age FROM Student WHERE id = $1"
	err := db.conn.QueryRow(query, id).Scan(&student.Id, &student.Name, &student.Age)
	if err != nil {
		return Student{}, err
	}
	return student, nil
}

func (db *myDb) getAll() ([]Student, error) {
	query := "SELECT id, name, age FROM Student"
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (db *myDb) addUser(name, age string) error {
	query := "INSERT INTO Student (name, age) VALUES ($1, $2)"
	_, err := db.conn.Exec(query, name, age)
	return err
}

func (db *myDb) editUser(id int, name, age string) error {
	query := "UPDATE Student SET name=$1, age=$2 WHERE id=$3"
	_, err := db.conn.Exec(query, name, age, id)
	return err
}

func (db *myDb) deleteUser(id int) error {
	query := "DELETE FROM Student WHERE id=$1"
	_, err := db.conn.Exec(query, id)
	return err
}
