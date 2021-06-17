package service

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/repository"
	"refactory/notes/internal/config"
	"refactory/notes/internal/mail"
	"refactory/notes/internal/security/token"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.UserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, username, password string) (*model.LoginResponse, error)
	VerifyCode(ctx context.Context, session token.Token, code int) error
	ListUser(ctx context.Context) ([]*model.UserResponse, error)
	DetailUser(ctx context.Context, Id int) (*model.UserResponse, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.UserResponse, error)
	DeleteUser(ctx context.Context, id int) error
	ActiveUser(ctx context.Context, id int) error
}

type role int

const (
	UserRole role = iota + 1
	AdminRole
)

func (r role) Int() int {
	return int(r)
}

func (r role) String() string {
	return [...]string{"user", "admin"}[r-1]
}

type userService struct {
	repo   repository.UserRepository
	mailer mail.Mailer
}

func NewUserService(repo repository.UserRepository) *userService {
	return &userService{repo: repo, mailer: mail.NewMailer(repo)}
}

func (a *userService) CreateUser(ctx context.Context, req model.UserRequest) (*model.UserResponse, error) {
	// encrypt password
	pass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if nil != err {
		return nil, err
	}

	u := model.NewUser(0, req.FirstName, req.LastName, req.Email, req.Username, string(pass), "", UserRole.Int())

	session := model.Session{
		UserId:     u.Id,
		Email:      u.Email,
		Username:   u.Username,
		RoleId:     u.Role,
		Code:       rand.Intn(999999),
		IsVerified: false,
		IsSent:     false,
		IsActive:   true,
	}

	// record data to database
	token, err := a.repo.Create(ctx, u, session)
	if nil != err {
		// check if user already exist
		if vErr, ok := err.(*pq.Error); ok && vErr.Code == "23505" {
			return nil, app.Error{Code: app.DuplicateCode.Int(), Message: fmt.Sprintf("duplicate value for field %s", vErr.Column)}
		}

		return nil, err
	}

	// adding user to mail verification queue
	a.mailer.Add(session)

	return model.NewUserResponse(u.Id, u.FirstName, u.LastName, u.Email, u.Username, u.Password, u.Photo, UserRole.String(), token), nil
}

func (a *userService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// finding user
	u, err := a.repo.FindUser(ctx, username)
	if nil != err {
		return nil, err
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); nil != err {
		return nil, app.Error{Code: app.UnauthenticatedCode.Int(), Message: "Unauthenticated"}
	}

	// check if user is verified or active
	if !u.IsVerified || !u.IsActive {
		return nil, app.Error{Code: app.UnauthorizedCode.Int(), Message: "Unauthorized"}
	}

	session := model.Session{
		UserId:     u.Id,
		Email:      u.Email,
		Username:   u.Username,
		RoleId:     u.Role,
		Code:       0,
		IsVerified: u.IsVerified,
		IsSent:     false,
		IsActive:   u.IsActive,
	}

	// generate new token and update session
	t, err := token.GenerateToken(session)
	if err := a.repo.UpdateSession(ctx, session); nil != err {
		return nil, app.Error{Code: app.InternalCode.Int(), Message: "Internal Server Error"}
	}

	return model.NewLoginResponse(u.Username, t), nil
}

func (a *userService) VerifyCode(ctx context.Context, session token.Token, code int) error {
	sess, err := a.repo.FindSession(ctx, session.Username)
	if nil != err {
		return app.Error{Code: app.InternalCode.Int(), Message: fmt.Sprintf("finding session for user: %s", session.Username)}
	}

	if sess.Code != code {
		return app.InvalidCodeError
	}

	if err := a.repo.VerifyUser(ctx, session.Username); nil != err {
		return err
	}

	return nil
}

func (a *userService) ListUser(ctx context.Context) ([]*model.UserResponse, error) {
	var response []*model.UserResponse
	result, err := a.repo.ListUser(ctx)
	if nil != err {
		return nil, err
	}

	for _, r := range result {
		u := model.NewUserResponse(r.Id, r.FirstName, r.LastName, r.Email, r.Username, r.Password, fmt.Sprintf("%s/api/media/%d", config.Cfg().WebAddress, r.MediaId), "", "")
		switch r.Role {
		case UserRole.Int():
			u.Role = UserRole.String()
		case AdminRole.Int():
			u.Role = AdminRole.String()
		}
		response = append(response, u)
	}

	return response, nil
}

func (a *userService) DetailUser(ctx context.Context, Id int) (*model.UserResponse, error) {
	result, err := a.repo.DetailUser(ctx, Id)
	if nil != err {
		return nil, err
	}

	u := model.NewUserResponse(result.Id, result.FirstName, result.LastName, result.Email, result.Username, result.Password, fmt.Sprintf("%s/api/media/%d", config.Cfg().WebAddress, result.MediaId), "", "")
	switch result.Role {
	case UserRole.Int():
		u.Role = UserRole.String()
	case AdminRole.Int():
		u.Role = AdminRole.String()
	}

	return u, nil
}

func (a *userService) UpdateUser(ctx context.Context, user *model.User) (*model.UserResponse, error) {
	if err := a.repo.UpdateUser(ctx, user); nil != err {
		return nil, err
	}

	u := model.NewUserResponse(user.Id, user.FirstName, user.LastName, user.Email, user.Username, user.Password, fmt.Sprintf("%s/api/media/%d", config.Cfg().WebAddress, user.MediaId), "", "")
	switch user.Role {
	case UserRole.Int():
		u.Role = UserRole.String()
	case AdminRole.Int():
		u.Role = AdminRole.String()
	}

	return u, nil
}

func (a *userService) DeleteUser(ctx context.Context, id int) error {
	if err := a.repo.DeleteUser(ctx, id); nil != err {
		return err
	}
	return nil
}

func (a *userService) ActiveUser(ctx context.Context, id int) error {
	return a.repo.ActiveUser(ctx, id)
}
