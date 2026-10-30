package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/JoelTeoGom/login-service/internal/domain"
	"github.com/JoelTeoGom/login-service/internal/ports"
	"github.com/JoelTeoGom/login-service/internal/service"
)

// In-memory fakes: the service is tested without MySQL, bcrypt or JWT.

type fakeRepo struct {
	byEmail map[string]*domain.User
	nextID  int64
}

func newFakeRepo() *fakeRepo { return &fakeRepo{byEmail: map[string]*domain.User{}} }

func (r *fakeRepo) Create(_ context.Context, u *domain.User) error {
	if _, ok := r.byEmail[u.Email]; ok {
		return domain.ErrEmailTaken
	}
	r.nextID++
	u.ID = r.nextID
	r.byEmail[u.Email] = u
	return nil
}

func (r *fakeRepo) FindByEmail(_ context.Context, email string) (*domain.User, error) {
	if u, ok := r.byEmail[email]; ok {
		return u, nil
	}
	return nil, domain.ErrNotFound
}

func (r *fakeRepo) FindByID(context.Context, int64) (*domain.User, error) {
	return nil, domain.ErrNotFound
}
func (r *fakeRepo) Update(context.Context, *domain.User) error { return nil }
func (r *fakeRepo) Delete(context.Context, int64) error        { return nil }

type plainHasher struct{}

func (plainHasher) Hash(p string) (string, error) { return "h:" + p, nil }
func (plainHasher) Matches(h, p string) bool      { return h == "h:"+p }

type fakeTokens struct{}

func (fakeTokens) Issue(int64) (string, error)  { return "token", nil }
func (fakeTokens) Verify(string) (int64, error) { return 1, nil }

func newAuth() *service.AuthService {
	return service.NewAuthService(newFakeRepo(), plainHasher{}, fakeTokens{})
}

func TestRegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	auth := newAuth()

	u, err := auth.Register(ctx, ports.RegisterInput{Email: " Joel@Example.com ", Name: "Joel", Password: "supersecret"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if u.Email != "joel@example.com" {
		t.Errorf("email not normalized: %q", u.Email)
	}

	token, err := auth.Login(ctx, ports.LoginInput{Email: "joel@example.com", Password: "supersecret"})
	if err != nil || token == "" {
		t.Fatalf("login: token=%q err=%v", token, err)
	}
}

func TestRegisterDuplicateEmail(t *testing.T) {
	ctx := context.Background()
	auth := newAuth()
	in := ports.RegisterInput{Email: "a@b.com", Name: "A", Password: "supersecret"}

	if _, err := auth.Register(ctx, in); err != nil {
		t.Fatal(err)
	}
	if _, err := auth.Register(ctx, in); !errors.Is(err, domain.ErrEmailTaken) {
		t.Fatalf("want ErrEmailTaken, got %v", err)
	}
}

func TestLoginWrongPasswordAndUnknownEmailLookTheSame(t *testing.T) {
	ctx := context.Background()
	auth := newAuth()
	_, _ = auth.Register(ctx, ports.RegisterInput{Email: "a@b.com", Name: "A", Password: "supersecret"})

	_, errWrong := auth.Login(ctx, ports.LoginInput{Email: "a@b.com", Password: "nope-nope"})
	_, errUnknown := auth.Login(ctx, ports.LoginInput{Email: "x@y.com", Password: "supersecret"})

	if !errors.Is(errWrong, domain.ErrInvalidCredentials) || !errors.Is(errUnknown, domain.ErrInvalidCredentials) {
		t.Fatalf("want ErrInvalidCredentials for both, got %v / %v", errWrong, errUnknown)
	}
}

func TestRegisterValidation(t *testing.T) {
	auth := newAuth()
	_, err := auth.Register(context.Background(), ports.RegisterInput{Email: "not-an-email", Name: "A", Password: "supersecret"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("want ErrInvalidInput, got %v", err)
	}
}
