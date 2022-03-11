package database

import (
	"encoding/json"
	"pitstop-api/src/schemas"
)

var defaultUserSettings, _ = json.Marshal(schemas.UserSettingsSchema{
	DarkTheme:           true,
	Notifications:       true,
	BlockFriendRequests: false,
})

var defaultUserStats, _ = json.Marshal(schemas.UserStatsSchema{})
