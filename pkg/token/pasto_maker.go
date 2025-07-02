package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

// import (
// 	"fmt"
// 	"time"
// )

type PasetoMarker struct {
	paseto       *paseto.V2
	synmetricKey []byte
}

func NewPasetoMaker(synmetricKey string) (Maker, error) {
	if len(synmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalide key size: must be exaxtly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMarker{
		paseto:       paseto.NewV2(),
		synmetricKey: []byte(synmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username anf duration
func (maker *PasetoMarker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}

	return maker.paseto.Encrypt(maker.synmetricKey, payload, nil)
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMarker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.synmetricKey, payload, nil)
	if err != nil{
		return nil, ErrInvalideToken
	}

	err = payload.Valid()
	if err != nil{
		return nil, err
	}

	return payload, nil
}
