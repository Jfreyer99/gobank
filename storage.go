package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	_ "github.com/lib/pq"
)

//-------------------------------------------------------------------------------------------------------
// Contains only defintion of Storage interface
type Storage interface{
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

// ---------------------------------------------------------------------------------------------------
// Contains struct PostgresStore which implements Storage interface
type PostgresStore struct{
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error){
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil{
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error{
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error{
	query := `
		CREATE TABLE IF NOT EXISTS ACCOUNT (
		account_id SERIAL PRIMARY KEY,
		last_name VARCHAR(30) NOT NULL,
		first_name VARCHAR(30) NOT NULL,
		account_number SERIAL NOT NULL,
		balance	REAL,
		created_at TIMESTAMP NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) printAccountTable() error {
	rows, err := s.db.Query("SELECT * FROM ACCOUNT")
	if err != nil{
		return err
	}
	defer rows.Close()

	accounts := make([]*Account, 0)
	timestamps := make([]string, 0)

	for rows.Next() {
		var first_name string
		var last_name string
		var account_number int64
		var account_id int
		var balance float64
		var created_at string

		if err := rows.Scan(&account_id, &last_name, &first_name, &account_number, &balance, &created_at); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}

		accounts = append(accounts, &Account{
			ID: account_id,
			FirstName: first_name,
			LastName: last_name,
			Number: account_number,
			Balance: balance,
		})

		timestamps = append(timestamps, created_at)
	}

	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for i:=0; i < len(accounts); i++{
		fmt.Printf("%+v%s\n", accounts[i], timestamps[i])
	}

	return nil
}


func (s *PostgresStore) CreateAccount(a *Account) error{
	_, err := s.db.Exec("INSERT INTO ACCOUNT(last_name, first_name, balance, created_at) VALUES($1,$2,$3,$4)",
	a.LastName, a.FirstName, a.Balance, time.Now())
	 if err != nil{
		return err
	 }
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error{
	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error{
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error){
	return nil, nil
}

//----------------------------------------------------------------------------------