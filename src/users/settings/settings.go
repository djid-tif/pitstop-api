package settings

import (
	"encoding/json"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
)

type UserSettings struct {
	darkTheme           bool
	notifications       bool
	blockFriendRequests bool
	sellUserData        bool
}

func (settings *UserSettings) ToggleDarkTheme(id uint) {
	settings.darkTheme = !settings.darkTheme
	settings.save(id)
}

func (settings *UserSettings) IsDarkThemeActive() bool {
	return settings.darkTheme
}

func (settings *UserSettings) ToggleNotifications(id uint) {
	settings.notifications = !settings.notifications
	settings.save(id)
}

func (settings *UserSettings) AreNotificationsActive() bool {
	return settings.notifications
}

func (settings *UserSettings) ToggleBlockFriendRequests(id uint) {
	settings.blockFriendRequests = !settings.blockFriendRequests
	settings.save(id)
}

func (settings *UserSettings) AreFriendRequestsBlocked() bool {
	return settings.blockFriendRequests
}

func (settings *UserSettings) save(id uint) {
	marshaledSettings, err := json.Marshal(settings.ToSchema())
	if err != nil {
		utils.PrintError(err)
		return
	}
	err = database.UpdateUserById(id, "settings", marshaledSettings)
	if err != nil {
		utils.PrintError(err)
	}
}

func (settings UserSettings) ToSchema() schemas.UserSettingsSchema {
	return schemas.UserSettingsSchema{
		DarkTheme:           settings.darkTheme,
		Notifications:       settings.notifications,
		BlockFriendRequests: settings.blockFriendRequests,
	}
}

func LoadSettings(settingsSchema schemas.UserSettingsSchema) UserSettings {
	return UserSettings{
		darkTheme:           settingsSchema.DarkTheme,
		notifications:       settingsSchema.Notifications,
		blockFriendRequests: settingsSchema.BlockFriendRequests,
	}
}
