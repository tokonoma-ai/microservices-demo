package api

// service.go contains the definition and implementation (business logic) of the
// user service. Everything here is agnostic to the transport (HTTP).

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/microservices-demo/user/db"
	"github.com/microservices-demo/user/users"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

// Service is the user service, providing operations for users to login, register, and retrieve customer information.
type Service interface {
	Login(ctx context.Context, username, password string) (users.User, error) // GET /login
	Register(ctx context.Context, username, password, email, first, last string) (string, error)
	GetUsers(ctx context.Context, id string) ([]users.User, error)
	PostUser(ctx context.Context, u users.User) (string, error)
	GetAddresses(ctx context.Context, id string) ([]users.Address, error)
	PostAddress(ctx context.Context, u users.Address, userid string) (string, error)
	GetCards(ctx context.Context, id string) ([]users.Card, error)
	PostCard(ctx context.Context, u users.Card, userid string) (string, error)
	Delete(ctx context.Context, entity, id string) error
	Health(ctx context.Context) []Health // GET /health
}

// NewFixedService returns a simple implementation of the Service interface,
func NewFixedService() Service {
	return &fixedService{}
}

type fixedService struct{}

type Health struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

func (s *fixedService) Login(_ context.Context, username, password string) (users.User, error) {
	u, err := db.GetUserByName(username)
	if err != nil {
		return users.New(), err
	}
	if u.Password != calculatePassHash(password, u.Salt) {
		return users.New(), ErrUnauthorized
	}
	db.GetUserAttributes(&u)
	u.MaskCCs()
	return u, nil

}

func (s *fixedService) Register(_ context.Context, username, password, email, first, last string) (string, error) {
	u := users.New()
	u.Username = username
	u.Password = calculatePassHash(password, u.Salt)
	u.Email = email
	u.FirstName = first
	u.LastName = last
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s *fixedService) GetUsers(_ context.Context, id string) ([]users.User, error) {
	if id == "" {
		us, err := db.GetUsers()
		for k, u := range us {
			u.AddLinks()
			us[k] = u
		}
		return us, err
	}
	u, err := db.GetUser(id)
	u.AddLinks()
	return []users.User{u}, err
}

func (s *fixedService) PostUser(_ context.Context, u users.User) (string, error) {
	u.NewSalt()
	u.Password = calculatePassHash(u.Password, u.Salt)
	err := db.CreateUser(&u)
	return u.UserID, err
}

func (s *fixedService) GetAddresses(_ context.Context, id string) ([]users.Address, error) {
	if id == "" {
		as, err := db.GetAddresses()
		for k, a := range as {
			a.AddLinks()
			as[k] = a
		}
		return as, err
	}
	a, err := db.GetAddress(id)
	a.AddLinks()
	return []users.Address{a}, err
}

func (s *fixedService) PostAddress(_ context.Context, add users.Address, userid string) (string, error) {
	err := db.CreateAddress(&add, userid)
	return add.ID, err
}

func (s *fixedService) GetCards(_ context.Context, id string) ([]users.Card, error) {
	if id == "" {
		cs, err := db.GetCards()
		for k, c := range cs {
			c.AddLinks()
			cs[k] = c
		}
		return cs, err
	}
	c, err := db.GetCard(id)
	c.AddLinks()
	return []users.Card{c}, err
}

func (s *fixedService) PostCard(_ context.Context, card users.Card, userid string) (string, error) {
	err := db.CreateCard(&card, userid)
	return card.ID, err
}

func (s *fixedService) Delete(_ context.Context, entity, id string) error {
	return db.Delete(entity, id)
}

func (s *fixedService) Health(_ context.Context) []Health {
	var health []Health
	dbstatus := "OK"

	err := db.Ping()
	if err != nil {
		dbstatus = "err"
	}

	app := Health{"user", "OK", time.Now().String()}
	db := Health{"user-db", dbstatus, time.Now().String()}

	health = append(health, app)
	health = append(health, db)

	return health
}

func calculatePassHash(pass, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt)
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}
