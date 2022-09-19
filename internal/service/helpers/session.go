package helpers

import (
	"errors"
	"net/http"
)

func Authorization(r *http.Request, resourceOwner string) (bool, error) {
	doorman := DoormanConnector(r)

	token, err := doorman.GetAuthToken(r)
	if err != nil {
		return false, errors.New("invalid token")
	}
	permission, err := doorman.CheckPermission(resourceOwner, token)
	if err != nil || !permission {
		return false, errors.New("user does not have permission")
	}
	return true, nil
}
