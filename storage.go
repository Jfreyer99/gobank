package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
)

//-------------------------------------------------------------------------------------------------------
// Contains only defintion of Storage interface
type Storage interface{
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

// ---------------------------------------------------------------------------------------------------
// Contains struct PostgresStore which implements Storage interface
type PostgresStore struct{
	db *sql.DB
	mu sync.Mutex
}

func NewPostgresStore() (*PostgresStore, error) {
	//connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	connStr := "postgresql://postgres:gobank@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
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

func (s *PostgresStore) Close() {
    s.mu.Lock()
    defer s.mu.Unlock()
    if s.db != nil {
        s.db.Close()
        s.db = nil
    }
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

	for rows.Next() {
		account := &Account{}

		if err := rows.Scan(
			&account.ID, &account.LastName,
			&account.FirstName, &account.Number,
			&account.Balance, &account.CreatedAt); err != nil {
			log.Fatal(err)
		}

		accounts = append(accounts, account)
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
		fmt.Printf("%+v\n", accounts[i])
	}
	return nil
}

func (s *PostgresStore) CreateAccount(a *Account) error{

	stmt, err := s.db.Prepare(
	"INSERT INTO ACCOUNT(last_name, first_name, balance, created_at) VALUES($1,$2,$3,$4) RETURNING account_id, account_number")
	 if err != nil{
		return err
	 }

	 defer stmt.Close()

	 reerr := stmt.QueryRow(
		a.LastName,
		a.FirstName,
		a.Balance,
		a.CreatedAt,
	 ).Scan(&a.ID, &a.Number)

	 if reerr != nil{
		return reerr
	 }

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error{

	stmt, err := s.db.Prepare(
	"DELETE FROM ACCOUNT WHERE account_id = $1 RETURNING account_id")
	if err != nil{
		return err
	}
	
	defer stmt.Close()

	var account_id int

	reerr := stmt.QueryRow(id).Scan(&account_id)

	if reerr != nil{
		return reerr
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error{
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error){

	stmt, err := s.db.Prepare("SELECT * FROM ACCOUNT WHERE account_id = $1")

	if err != nil{
		return nil, err
	}

	defer stmt.Close()

	account := &Account{}

	reerr := stmt.QueryRow(id).Scan(
		&account.ID,
		&account.LastName,
		&account.FirstName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if reerr != nil{
		return nil, reerr
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error){
	return nil, fmt.Errorf("Not implemented yet")
}

//----------------------------------------------------------------------------------