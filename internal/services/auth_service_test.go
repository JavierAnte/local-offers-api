package services

import (
	"errors"
	"testing"

	"github.com/JavierAnte/local-offers-api/internal/dto"
	"github.com/JavierAnte/local-offers-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type fakeUserRepository struct {
	created            *models.User
	user               *models.User
	createErr, findErr error
}

func (f *fakeUserRepository) Create(user *models.User) error           { f.created = user; return f.createErr }
func (f *fakeUserRepository) FindByEmail(string) (*models.User, error) { return f.user, f.findErr }

func TestAuthServiceRegister(t *testing.T) {
	repo := &fakeUserRepository{}
	response, err := NewAuthService(repo, "secret").Register(dto.RegisterRequest{Name: "  Local User  ", Email: " USER@EXAMPLE.COM ", Password: "password"})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if repo.created.Name != "Local User" || repo.created.Email != "user@example.com" {
		t.Fatalf("created = %#v", repo.created)
	}
	if bcrypt.CompareHashAndPassword([]byte(repo.created.PasswordHash), []byte("password")) != nil {
		t.Fatal("password was not hashed correctly")
	}
	if response.Token == "" || response.User.Email != "user@example.com" {
		t.Fatalf("response = %#v", response)
	}
}

func TestAuthServiceRegisterErrors(t *testing.T) {
	tests := []dto.RegisterRequest{
		{Name: "", Email: "user@example.com", Password: "password"},
		{Name: "User", Email: "invalid", Password: "password"},
		{Name: "User", Email: "user@example.com", Password: "short"},
	}
	for _, req := range tests {
		if _, err := NewAuthService(&fakeUserRepository{}, "secret").Register(req); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("Register(%#v) = %v", req, err)
		}
	}
	if _, err := NewAuthService(&fakeUserRepository{createErr: gorm.ErrDuplicatedKey}, "secret").Register(dto.RegisterRequest{Name: "User", Email: "user@example.com", Password: "password"}); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("duplicate error = %v", err)
	}
}

func TestAuthServiceLogin(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	user := &models.User{Name: "User", Email: "user@example.com", PasswordHash: string(hash)}
	response, err := NewAuthService(&fakeUserRepository{user: user}, "secret").Login(dto.LoginRequest{Email: " USER@EXAMPLE.COM ", Password: "password"})
	if err != nil || response.Token == "" {
		t.Fatalf("Login() = %#v, %v", response, err)
	}
	if _, err := NewAuthService(&fakeUserRepository{user: user}, "secret").Login(dto.LoginRequest{Email: "user@example.com", Password: "wrong"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("wrong password = %v", err)
	}
	if _, err := NewAuthService(&fakeUserRepository{findErr: gorm.ErrRecordNotFound}, "secret").Login(dto.LoginRequest{Email: "missing@example.com", Password: "password"}); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("missing user = %v", err)
	}
	databaseErr := errors.New("database unavailable")
	if _, err := NewAuthService(&fakeUserRepository{findErr: databaseErr}, "secret").Login(dto.LoginRequest{Email: "user@example.com", Password: "password"}); !errors.Is(err, databaseErr) {
		t.Fatalf("database error = %v", err)
	}
}
