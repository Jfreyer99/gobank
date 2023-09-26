package storage

import (
	"github.com/Jfreyer99/gobank/types"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// -------------------------------------------------------------------------------------------------------
// Contains only defintion of Storage interface
// Could use Generics interface{} or any with the reflection api to determine fields of a passed in struct that represents the relation in postgres
type Storage interface {
	AccountStorage
	UserAccountStorage
}

type AccountStorage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(id, number int) error
	UpdateAccount(*types.Account) error
	GetAccountByIDAndNumber(id, number int) (*types.Account, error)
	GetAccounts() ([]*types.Account, error)
}

type UserAccountStorage interface {
	CreateUserAccount(*types.UserAccount) error
	DeleteUserAccount(int) error
	UpdateUserAccount(*types.UserAccount) error
	GetUserAccountByID(int) (*types.UserAccount, error)
	GetUserAccounts() ([]*types.UserAccount, error)
}
