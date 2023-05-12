package mockclient

type CreateUserCDO struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Age     int    `json:"age"`
	Email   string `json:"email"`
}