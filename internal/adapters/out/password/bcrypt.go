package password

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/JoelTeoGom/login-service/internal/ports"
)

// Bcrypt implements ports.PasswordHasher.
type Bcrypt struct {
	cost int
}

var _ ports.PasswordHasher = (*Bcrypt)(nil)

func NewBcrypt(cost int) *Bcrypt {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &Bcrypt{cost: cost}
}

func (b *Bcrypt) Hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(h), err
}

func (b *Bcrypt) Matches(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
