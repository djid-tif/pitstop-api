package users

import (
	"encoding/json"
	"pitstop-api/src/database"
	"pitstop-api/src/users/ranks"
	"pitstop-api/src/users/settings"
	"pitstop-api/src/users/stats"
	"pitstop-api/src/utils"
	"time"
)

// User represents a user
type User struct {
	id                  uint
	username            string
	email               string
	password            string
	registrationDate    time.Time
	rank                ranks.Rank
	lastPasswordRefresh time.Time
	friends             []uint
	settings            settings.UserSettings
	stats               stats.UserStats
}

// String method override
func (user User) String() string {
	return "User: " + user.username
}

func (user User) GetId() uint {
	return user.id
}

func (user User) GetUsername() string {
	return user.username
}

func (user *User) SetUsername(username string) error {
	// already exist check here or before
	err := database.UpdateUserById(user.id, "username", username)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	user.username = username
	return nil
}

func (user User) GetEmail() string {
	return user.email
}

func (user *User) SetEmail(email string) error {
	// already exist check here or before
	err := database.UpdateUserById(user.id, "email", email)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	user.email = email
	return nil
}

func (user User) GetPassword() string {
	return user.password
}

func (user *User) SetPassword(password string) error {
	err := database.UpdateUserById(user.id, "password", password)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	err = database.UpdateUserById(user.id, "last_password_refresh", time.Now().Unix())
	if err != nil {
		utils.PrintError(err)
		return err
	}

	user.stats.ResetPassword(user.GetId())

	user.password = password
	return nil
}

func (user User) GetRegistrationDate() time.Time {
	return user.registrationDate
}

func (user User) GetRank() ranks.Rank {
	return user.rank
}

func (user *User) SetRank(rank ranks.Rank) error {
	err := database.UpdateUserById(user.id, "rank", rank.GetId())
	if err != nil {
		utils.PrintError(err)
		return err
	}
	user.rank = rank
	return nil
}

func (user User) GetFriendsList() []uint {
	return user.friends
}

func (user *User) SetFriendsList(friendsList []uint) error {
	jsonList, err := json.Marshal(friendsList)
	if err != nil {
		return err
	}
	err = database.UpdateUserById(user.id, "friends", string(jsonList))
	if err != nil {
		utils.PrintError(err)
		return err
	}
	user.friends = friendsList
	return nil
}

func (user *User) Settings() *settings.UserSettings {
	return &user.settings
}

//func (user *User) SetSettings(settings settings.UserSettings) error {
//	jsonSettings, err := json.Marshal(settings)
//	if err != nil {
//		return err
//	}
//	err = database.UpdateUserById(user.id, "settings", string(jsonSettings))
//	if err != nil {
//		utils.PrintError(err)
//		return err
//	}
//	user.settings = settings
//	return nil
//}

func (user *User) Stats() *stats.UserStats {
	return &user.stats
}
