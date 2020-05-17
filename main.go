package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"unicode/utf8"

	_ "github.com/mattn/go-sqlite3"
)

type Phonebook struct {
	ID          int64
	Name        string
	PhoneNumber string
}

func validInput(str string) error {
	if str == "" {
		return fmt.Errorf("Input Error: %w", errors.New("Unexpected Input."))
	}

	if utf8.ValidString(str) && utf8.RuneCountInString(str) != len(str) {
		return fmt.Errorf("Input Error: %w", errors.New("Not use ASCII."))
	}
	return nil
}

func inputElement(strName string) (string, error) {
	var str string
	var s = bufio.NewScanner(os.Stdin)
	fmt.Print(strName, ">>")
	if s.Scan() {
		str = s.Text()
	}
	if err := validInput(str); err != nil {
		return "", err
	}
	return str, nil
}

func createTable(db *sql.DB) error {
	const sql = `
	CREATE TABLE IF NOT EXISTS phonebook (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		phoneNumber TEXT NOT NULL
	);
	`
	if _, err := db.Exec(sql); err != nil {
		return err
	}

	return nil
}

func showTable(db *sql.DB) error {
	fmt.Println("----------------------------------------------")
	rows, err := db.Query("SELECT * FROM phonebook")
	if err != nil {
		return err
	}

	for rows.Next() {
		var p Phonebook
		if err := rows.Scan(&p.ID, &p.Name, &p.PhoneNumber); err != nil {
			return err
		}
		fmt.Println("ID:", p.ID, "Name:", p.Name, "PhoneNumber:", p.PhoneNumber)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	fmt.Println("----------------------------------------------")

	return nil
}

func insertTable(db *sql.DB) error {
	var name, phoneNumber string
	for {
		if inputStr, err := inputElement("name"); err != nil {
			fmt.Println(err)
			continue
		} else {
			name = inputStr
		}
		break
	}
	for {
		if inputStr, err := inputElement("phoneNumber"); err != nil {
			fmt.Println(err)
			continue
		} else {
			phoneNumber = inputStr
		}
		break
	}

	var sql = "INSERT INTO phonebook(name, phoneNumber) values (?,?)"
	r, err := db.Exec(sql, name, phoneNumber)
	if err != nil {
		return err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	fmt.Println("INSERT:", id)
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "phonebook.db")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := createTable(db); err != nil {
			log.Fatal(err)
		}

		for {
			if err := showTable(db); err != nil {
				log.Fatal(err)
			}

			if err := insertTable(db); err != nil {
				log.Fatal(err)
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println(" Ctrl+Cを検知しました.プログラムを終了します.")
}
