package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const cryptCost = bcrypt.DefaultCost

var privKey *rsa.PrivateKey
var keySet jwk.Set

const expirationRefreshToken = time.Hour * 168
const expirationAccessToken = time.Hour * 3
const expirationResetToken = time.Hour * 24
const expirationActivateToken = time.Hour * 168

func init() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("échec de la génération de la clé privée private key: %s\n", err)
		return
	}
	privKey = key

	pubKey, err := jwk.New(privKey.PublicKey)
	if err != nil {
		fmt.Printf("échec de la création du JWK: %s\n", err)
		return
	}
	_ = pubKey.Set(jwk.AlgorithmKey, jwa.RS256)

	// This JWKS can *only* have 1 key.
	set := jwk.NewSet()
	set.Add(pubKey)
	keySet = set
}
