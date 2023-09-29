package types

type UserAccount struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	PassHash string `json:"passHash"`
}

func NewUserAccount(email, passhash string) *UserAccount {
	return &UserAccount{
		Email:    email,
		PassHash: passhash,
	}
}

type CreateUserAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
