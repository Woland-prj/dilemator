package user_entity

type Profile struct {
	Name    *string
	Surname *string
	Avatar  *string
}

func NewProfile(name, surname, avatar *string) *Profile {
	return &Profile{
		Name:    name,
		Surname: surname,
		Avatar:  avatar,
	}
}
