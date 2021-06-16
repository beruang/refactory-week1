package model

type User struct {
	Id         int
	FirstName  string
	LastName   string
	Email      string
	Username   string
	Password   string
	Photo      string
	Role       int
	IsVerified bool
	IsActive   bool
}

func NewUser(firstName string, lastName string, email string, username string, password string, photo string, role int) *User {
	return &User{FirstName: firstName, LastName: lastName, Email: email, Username: username, Password: password, Photo: photo, Role: role}
}

type UserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,alphanum,max=128"`
	Password  string `json:"password" validate:"required,min=6,max=128"`
	Photo     string `json:"photo"`
}

type UserResponse struct {
	Id        int    `json:"id_user"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Photo     string `json:"photo"`
	Role      string `json:"role"`
	Token     string `json:"token"`
}

func NewUserResponse(id int, firstName string, lastName string, email string, username string, password string, photo string, role string, token string) *UserResponse {
	return &UserResponse{Id: id, FirstName: firstName, LastName: lastName, Email: email, Username: username, Password: password, Photo: photo, Role: role, Token: token}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func NewLoginResponse(username string, token string) *LoginResponse {
	return &LoginResponse{Username: username, Token: token}
}

type Session struct {
	UserId     int    `json:"user_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Code       int    `json:"code"`
	RoleId     int    `json:"role_id"`
	IsVerified bool   `json:"is_verified"`
	IsSent     bool   `json:"is_sent"`
	IsActive   bool   `json:"is_active"`
}

type VerifyRequest struct {
	Code int `json:"code" validate: "required"`
}
