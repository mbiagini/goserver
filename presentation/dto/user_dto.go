package dto

type UserDTO struct {
	// The ID to uniquely identify a user
	Id int `json:"id" example:"1"`
	// The name of the user
	Name string `json:"name,omitempty" example:"Martín"`
	// The surname of the user
	Surname string `json:"surname,omitempty" example:"Biagini"`
	// The age of the user
	Age int `json:"age" example:"27"`
	// The email of the user
	Email string `json:"email,omitempty" example:"martinbiagini@gmail.com"`
} // @name UserDTO
