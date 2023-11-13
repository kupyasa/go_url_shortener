package models

type UserRegisterForm struct {
	Name     string `param:"name" query:"name" form:"name" json:"name"`
	Email    string `param:"email" query:"email" form:"email" json:"email"`
	Password string `param:"password" query:"password" form:"password" json:"-"`
}

type UserPresentation struct {
	Id      uint   `gorm:"primarykey" param:"user_id" query:"user_id" form:"user_id" json:"user_id"`
	Name    string `param:"name" query:"name" form:"name" json:"name"`
	Email   string `gorm:"unique" param:"email" query:"email" form:"email" json:"email"`
	IsAdmin bool   `param:"is_admin" query:"is_admin" form:"is_admin" json:"is_admin"`
}

type User struct {
	Id       uint   `gorm:"primarykey" param:"user_id" query:"user_id" form:"user_id" json:"user_id"`
	Name     string `param:"name" query:"name" form:"name" json:"name"`
	Email    string `gorm:"unique" param:"email" query:"email" form:"email" json:"email"`
	Password string `param:"password" query:"password" form:"password" json:"password"`
	IsAdmin  bool   `param:"is_admin" query:"is_admin" form:"is_admin" json:"is_admin"`
}
