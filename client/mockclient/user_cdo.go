package mockclient

import (
	"errors"
	"goserver/model"
	"strconv"
)

type UserCDO struct {
	Id      string `json:"id"      validate:"required"`
	Name    string `json:"name"    validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Age     int    `json:"age"     validate:"required,number,gte=0,lte=130"`
	Email   string `json:"email,omitempty" validate:"omitempty,email"`
}

func (cdo *UserCDO) ToModel() (m *model.User, err error) {
	id, err := strconv.Atoi(cdo.Id)
	if err != nil {
		return nil, errors.New("could not map UserCDO.Id field to User.Id: could not convert value to string")
	}
	return &model.User{
		Id:      id,
		Name:    cdo.Name,
		Surname: cdo.Surname,
		Age: 	 cdo.Age,
		Email:   cdo.Email,
	}, nil
}