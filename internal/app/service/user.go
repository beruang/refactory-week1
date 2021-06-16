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
	"refactory/notes/internal/mail"
	"refactory/notes/internal/security/token"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.UserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, username, password string) (*model.LoginResponse, error)
	VerifyCode(ctx context.Context, session token.Token, code int) error
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

	u := model.NewUser(req.FirstName, req.LastName, req.Email, req.Username, string(pass), "", UserRole.Int())

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
		return app.UnauthenticateError
	}

	if err := a.repo.VerifyUser(ctx, session.Username); nil != err {
		return err
	}

	return nil
}
