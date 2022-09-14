package models

type User struct {
	BaseModel
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Token    string `json:"token"`
}

func (User) TableName() string {
	return "users"
}

type GoplayRepo interface {
	Login(email string) (User, error)
	UpdateToken(email string, token string) error 
}
