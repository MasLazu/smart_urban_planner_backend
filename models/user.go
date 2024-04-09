package models

type UserRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *UserRegister) ToUser() *User {
	return &User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
	}
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u *UserLogin) ToUser() *User {
	return &User{
		Email:    u.Email,
		Password: u.Password,
	}
}

type User struct {
	ID       string   `json:"id" gorm:"primary_key" param:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email" gorm:"unique"`
	Password string   `json:"-"`
	Reports  []Report `json:"reports,omitempty" gorm:"foreignKey:AuthorID"`
}
