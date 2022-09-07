package types

import (
	"errors"

	"github.com/spf13/cast"
)

var Purposes = []string{
	"KYC",
}
var (
	PurposeError = errors.New("purpose is invalid")
)

type Purpose string

func (a Purpose) Validate() error {
	for _, purpose := range Purposes {
		if purpose == string(a) {
			return nil
		}
	}
	return PurposeError
}

func (a Purpose) String() string {
	return string(a)
}

var IsPurpose = &isPurpose{}

type isPurpose struct{}

func (ia *isPurpose) Validate(value interface{}) error {
	a, err := cast.ToStringE(value)
	if err != nil {
		return err
	}
	purpose := Purpose(a)
	return purpose.Validate()
}
