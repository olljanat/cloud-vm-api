package auth

import (
	"encoding/base64"
	"strings"

	"github.com/pkg/errors"
)

type Credentials struct {
	AccessKey string
	Secret    string
}

func DecodeCredentials(token string) (*Credentials, error) {
	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode token")
	}

	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	return &Credentials{
		AccessKey: parts[0],
		Secret:    parts[1],
	}, nil
}
