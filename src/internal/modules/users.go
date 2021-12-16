package modules

import (
	"api/internal/database"
	"api/internal/types"
)

// SaveUser ...
func SaveUser(user types.User) error {
	err := database.SaveUser(user)

	if err != nil {
		return err
	}

	return nil
}

// GetUser ...
func GetUser(publicKey string) (types.User, error) {
	user, err := database.GetUser(publicKey)

	// If 404 we need to create a new user in the database
	// because there's no user with the provided publicKey
	if err != nil && err.Error() == "404" {
		err := SaveUser(types.User{
			PublicKey: publicKey,
			Username:  "",
			Role:      "user",
		})

		if err != nil {
			return types.User{}, err
		}
	} else if err != nil {
		return types.User{}, err
	}

	return user, nil
}
