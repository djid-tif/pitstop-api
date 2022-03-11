package users

import (
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/users/auth"
	"pitstop-api/src/users/ranks"
	"pitstop-api/src/users/settings"
	"pitstop-api/src/users/stats"
	"pitstop-api/src/utils"
	"time"
)

func LoadUser(id uint) (User, bool) {
	dbUser, found := database.GetUserById(id)
	if !found {
		return User{}, false
	}
	user := User{
		id:                  dbUser.Id,
		username:            dbUser.Username,
		email:               dbUser.Email,
		password:            dbUser.Password,
		registrationDate:    time.Unix(dbUser.RegistrationDate, 0),
		rank:                ranks.GetRankById(dbUser.Rank),
		lastPasswordRefresh: time.Unix(dbUser.LastPasswordRefresh, 0),
		friends:             dbUser.Friends,
		settings:            settings.LoadSettings(dbUser.Settings),
		stats:               stats.LoadStats(dbUser.Stats),
	}
	return user, true
}

func LoadUserFromRequest(r *http.Request) (User, bool) {
	// get id with token
	id, errorToken := auth.ExtractIdFromRequest(r)
	if errorToken != nil {
		return User{}, false
	}

	return LoadUser(id)
}

func (user User) ToPublic() schemas.UserPublic {
	return schemas.UserPublic{
		Id:               user.id,
		Username:         user.username,
		RegistrationDate: utils.TimeToIso(user.registrationDate),
		Rank:             user.rank.GetName(),
		Posts:            database.GetUserPosts(user.id),
	}
}

func (user User) ToPrivate() schemas.UserPrivate {
	return schemas.UserPrivate{
		Id:               user.id,
		Username:         user.username,
		Email:            user.email,
		RegistrationDate: utils.TimeToIso(user.registrationDate),
		Rank:             user.rank.GetName(),
		Friends:          user.friends,
		Settings:         user.settings.ToSchema(),
		Stats:            user.stats.ToSchema(),
		Posts:            database.GetUserPosts(user.id),
	}
}
