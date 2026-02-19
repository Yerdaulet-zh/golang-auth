package config

import (
	"crypto/rsa"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/awnumar/memguard"
	"github.com/golang-auth/internal/core/ports" // Verify this matches your module path
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type JWTTokenKeys struct {
	kid      string
	public   *memguard.Enclave
	private  *memguard.Enclave
	logger   ports.Logger
	Duration time.Duration
}

func NewJWTConfig(logger ports.Logger) (*JWTTokenKeys, error) {
	KID := viper.GetString("jwt.KID")
	pubPathRel := viper.GetString("jwt.pubPath")
	privPathRel := viper.GetString("jwt.privPath")
	minute := viper.GetInt("jwt.durationInMinute")
	duration := time.Duration(minute) * time.Minute

	// Resolve absolute paths based on current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	pubPath := filepath.Join(wd, pubPathRel)
	privPath := filepath.Join(wd, privPathRel)

	// Read key files into temporary heap memory
	privRaw, err := os.ReadFile(privPath)
	if err != nil {
		return nil, err
	}
	pubRaw, err := os.ReadFile(pubPath)
	if err != nil {
		return nil, err
	}

	// Move keys into secure memguard enclaves
	keys := &JWTTokenKeys{
		kid:      KID,
		public:   memguard.NewEnclave(pubRaw),
		private:  memguard.NewEnclave(privRaw),
		logger:   logger,
		Duration: duration,
	}

	// Zero out the plain-text slices immediately
	scrubByteSlice(privRaw)
	scrubByteSlice(pubRaw)

	return keys, nil
}

func (j *JWTTokenKeys) SignToken(jti string) (string, error) {
	claims := jwt.MapClaims{
		"jti": jti,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(j.Duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = j.kid

	// Access the protected buffer
	lockedBuf, err := j.private.Open()
	if err != nil {
		j.logger.Error("Failed to open private key enclave", "error", err)
		return "", err
	}
	// Wipes the PEM bytes as soon as signing is finished
	defer lockedBuf.Destroy()

	// Parse RSA Private Key into heap
	key, err := jwt.ParseRSAPrivateKeyFromPEM(lockedBuf.Bytes())
	if err != nil {
		j.logger.Error("Failed to parse private key from PEM", "error", err)
		return "", err
	}

	// Generate signed string
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	// Wipe the RSA object's math bits from the heap
	scrubRSAPrivateKey(key)

	return tokenString, nil
}

// --- Internal Security Helpers ---

func scrubByteSlice(b []byte) {
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(b)
}

func scrubRSAPrivateKey(k *rsa.PrivateKey) {
	if k == nil {
		return
	}
	zeroBigInt(k.D)
	for _, p := range k.Primes {
		zeroBigInt(p)
	}
	zeroBigInt(k.Precomputed.Dp)
	zeroBigInt(k.Precomputed.Dq)
	zeroBigInt(k.Precomputed.Qinv)
	runtime.KeepAlive(k)
}

func zeroBigInt(n *big.Int) {
	if n == nil {
		return
	}
	words := n.Bits()
	for i := range words {
		words[i] = 0
	}
}
