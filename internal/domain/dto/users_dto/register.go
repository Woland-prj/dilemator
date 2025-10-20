package users_dto

type RegisterDto struct {
	Email    string
	Password string
	Name     *string
	Surname  *string
}
