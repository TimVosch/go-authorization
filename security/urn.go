package security

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidURN = errors.New("invalid URN")
)

type URN []string

func (u URN) Resource() string {
	return u[0]
}

func (u URN) Identifier() string {
	return u[1]
}

func OrganisationIDToURN(userID int) URN {
	return URN{"organisation", strconv.Itoa(userID)}
}
func URNToOrganisationID(urn URN) (int, error) {
	if urn.Resource() != "organisation" {
		return 0, ErrInvalidURN
	}
	return strconv.Atoi(urn.Identifier())
}

func UserIDToURN(userID int) URN {
	return URN{"user", strconv.Itoa(userID)}
}
func URNToUserID(urn URN) (int, error) {
	if urn.Resource() != "user" {
		return 0, ErrInvalidURN
	}
	return strconv.Atoi(urn.Identifier())
}
