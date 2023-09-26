package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/Jfreyer99/gobank/types"
)

// ----------------------------------------------------------------------------------------------------------
// Contains struct PostgresStore which implements Storage interface
type PostgresStore struct {
	db *sql.DB
	mu sync.Mutex
}

//	Postgres Specific Setup
//----------------------------------------------------------------------------------------------------------

func NewPostgresStore() (*PostgresStore, error) {
	//connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	connStr := "postgresql://postgres:gobank@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
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

func (s *PostgresStore) Init() error {

	if err := s.CreateAccountNumberSequence(); err != nil {
		return err
	}
	if err := s.CreateUserAccountTable(); err != nil {
		return err
	}
	return s.CreateAccountTable()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS ACCOUNT (
		account_id INT NOT NULL,
		last_name VARCHAR(30) NOT NULL,
		first_name VARCHAR(30) NOT NULL,
		account_number integer DEFAULT nextval('account_account_number_seq'),
		balance	REAL,
		created_at TIMESTAMP NOT NULL,
		constraint PK_ACCOUNT_TABLE PRIMARY KEY (account_id, account_number),
		CONSTRAINT fk_user FOREIGN KEY(account_id) REFERENCES USERACCOUNT(account_id)
	);`
	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateUserAccountTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS USERACCOUNT (
		account_id INT GENERATED ALWAYS AS IDENTITY,
		email VARCHAR(70) NOT NULL,
		passhash VARCHAR(100) NOT NULL,
		salthash VARCHAR(100) NOT NULL,
		PRIMARY KEY(account_id)
	);`
	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) DropTableAccount() error {

	stmt, err := s.db.Prepare(`DROP TABLE IF EXISTS ACCOUNT`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) CreateAccountNumberSequence() error {

	stmt, err := s.db.Prepare(
		`CREATE SEQUENCE IF NOT EXISTS account_account_number_seq
	START 10000;`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) printAccountTable() error {
	rows, err := s.db.Query("SELECT * FROM ACCOUNT")
	if err != nil {
		return err
	}
	defer rows.Close()

	accounts := make([]*types.Account, 0)

	for rows.Next() {
		account := &types.Account{}

		if err := rows.Scan(
			&account.ID, &account.LastName,
			&account.FirstName, &account.Number,
			&account.Balance, &account.CreatedAt); err != nil {
			return err
		}

		accounts = append(accounts, account)
	}

	rerr := rows.Close()
	if rerr != nil {
		log.Fatal(rerr)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(accounts); i++ {
		fmt.Printf("%+v\n", accounts[i])
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------

// Implementation of Account Methods of Interface Storage
//----------------------------------------------------------------------------------------------------------

func (s *PostgresStore) CreateAccount(a *types.Account) error {

	stmt, err := s.db.Prepare(
		`INSERT INTO ACCOUNT(account_id, last_name, first_name, balance, created_at) 
		VALUES($1,$2,$3,$4) RETURNING account_id, account_number`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	reerr := stmt.QueryRow(
		a.ID,
		a.LastName,
		a.FirstName,
		a.Balance,
		a.CreatedAt,
	).Scan(&a.ID, &a.Number)

	if reerr != nil {
		return reerr
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id, number int) error {

	stmt, err := s.db.Prepare(
		"DELETE FROM ACCOUNT WHERE account_id = $1 AND account_number = $2")
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Exec(id, number)

	if err != nil {
		return err
	}

	numDeleted, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if numDeleted <= 0 {
		return fmt.Errorf("now rows were deleted")
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(a *types.Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByIDAndNumber(id, number int) (*types.Account, error) {

	//Cache the stmt somewhere
	stmt, err := s.db.Prepare("SELECT * FROM ACCOUNT WHERE account_id = $1 AND account_number = $2")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	account := &types.Account{}
	reerr := stmt.QueryRow(id, number).Scan(
		&account.ID,
		&account.LastName,
		&account.FirstName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if reerr != nil {
		return nil, reerr
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {

	stmt, err := s.db.Prepare("SELECT * FROM ACCOUNT")

	if err != nil {
		return nil, err
	}

	rows, errQ := stmt.Query()

	if err != nil {
		return nil, errQ
	}

	accounts := make([]*types.Account, 0)

	for rows.Next() {

		account := &types.Account{}

		if err := rows.Scan(
			&account.ID,
			&account.LastName,
			&account.FirstName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)

	}

	return accounts, nil
}

//----------------------------------------------------------------------------------

// Implementation of Account Methods of Interface Storage
//----------------------------------------------------------------------------------------------------------

func (s *PostgresStore) CreateUserAccount(u *types.UserAccount) error {

	stmt, err := s.db.Prepare("INSERT INTO USERACCOUNT(email, passhash, salthash) VALUES($1,$2,$3) RETURNING account_id")

	if err != nil {
		return err
	}

	row := stmt.QueryRow(u.Email, u.PassHash, u.SaltHash).Scan(&u.ID)

	if row != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteUserAccount(id int) error {

	//Needs to cascade delete because of references to account table, testing this scenario is nessecary
	stmt, err := s.db.Prepare("DELETE FROM USERACCOUNT WHERE account_id = $1")

	if err != nil {
		return err
	}

	rerr := stmt.QueryRow(id).Scan()

	if rerr != nil {
		return rerr
	}

	return nil
}

func (s *PostgresStore) UpdateUserAccount(*types.UserAccount) error {
	return fmt.Errorf("not implemented yet")
}

func (s *PostgresStore) GetUserAccountByID(id int) (*types.UserAccount, error) {

	stmt, err := s.db.Prepare("SELECT * FROM USERACCOUNT WHERE account_id = $1")

	if err != nil {
		return nil, err
	}

	userAccount := &types.UserAccount{}

	err = stmt.QueryRow(id).Scan(
		&userAccount.ID,
		&userAccount.Email,
		&userAccount.PassHash,
		&userAccount.SaltHash,
	)

	if err != nil {
		return nil, err
	}

	return userAccount, nil
}

func (s *PostgresStore) GetUserAccounts() ([]*types.UserAccount, error) {

	stmt, err := s.db.Prepare("SELECT * FROM USERACCOUNT")

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}

	userAccounts := make([]*types.UserAccount, 0)

	for rows.Next() {

		userAccount := &types.UserAccount{}

		err := rows.Scan(
			&userAccount.ID,
			&userAccount.Email,
			&userAccount.PassHash,
			&userAccount.SaltHash,
		)

		if err != nil {
			return nil, err
		}

		userAccounts = append(userAccounts, userAccount)
	}

	return userAccounts, nil
}

//----------------------------------------------------------------------------------------------------------
