package models

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r *User) SetId(id string) {
	r.Id = id
}

type UserCredentials struct {
	Name     string `json:"name" validate:"required,max=20"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserTokenResponse struct {
	Token   string `json:"token"`
	LocalId string `json:"localId"`
}
