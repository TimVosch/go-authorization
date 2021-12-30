package security

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMalformedCredentialString = errors.New("invalid credential string")

	enc = base64.StdEncoding.WithPadding(base64.NoPadding)
)

// Credential ...
type Credential struct {
	Type   AuthType
	Secret string
}

func (c Credential) String() string {
	return strings.Join([]string{
		enc.EncodeToString([]byte(c.Type)),
		enc.EncodeToString([]byte(c.Secret)),
	}, ".",
	)
}

func UnmarshalCredentialString(str string) (Credential, error) {
	var c Credential

	parts := strings.Split(str, ".")
	if len(parts) != 2 {
		return c, ErrMalformedCredentialString
	}

	typ, err := enc.DecodeString(parts[0])
	if err != nil {
		return c, fmt.Errorf("%w: %v", ErrMalformedCredentialString, err)
	}
	secret, err := enc.DecodeString(parts[1])
	if err != nil {
		return c, fmt.Errorf("%w: %v", ErrMalformedCredentialString, err)
	}

	c.Type = AuthType(typ)
	c.Secret = string(secret)

	return c, nil
}
