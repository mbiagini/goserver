package dto

type CreateUserRequestDTO struct {
	// The name of the user
	Name string `json:"name" validate:"required" example:"Martín"`
	// The surname of the user
	Surname string `json:"surname" validate:"required" example:"Biagini"`
	// The age of the user
	Age int `json:"age" validate:"required,number,gte=0,lte=130" example:"27"`
	// The email of the user
	Email string `json:"email,omitempty" validate:"email" example:"martinbiagini@gmail.com"`
} // @name CreateUserRequestDTO