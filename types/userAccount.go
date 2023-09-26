package types

type UserAccount struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	PassHash string `json:"passHash"`
	SaltHash string `json:"saltHash"`
}

func NewUserAccount(email, passhash, salthash string) *UserAccount {
	return &UserAccount{
		Email:    email,
		PassHash: passhash,
		SaltHash: salthash,
	}
}

type CreateUserAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
