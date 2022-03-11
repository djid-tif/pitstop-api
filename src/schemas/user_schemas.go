package schemas

type UserPublic struct {
	Id               uint   `json:"id"`
	Username         string `json:"username"`
	RegistrationDate string `json:"registration_date"`
	Rank             string `json:"rank"`
	Posts            []uint `json:"posts"`
}

type UserPrivate struct {
	Id               uint               `json:"id"`
	Username         string             `json:"username"`
	Email            string             `json:"email"`
	RegistrationDate string             `json:"registration_date"`
	Rank             string             `json:"rank"`
	Friends          []uint             `json:"friends"`
	Settings         UserSettingsSchema `json:"settings"`
	Stats            UserStatsSchema    `json:"stats"`
	Posts            []uint             `json:"posts"`
}

type UserSettingsSchema struct {
	DarkTheme           bool `json:"dark_theme"`
	Notifications       bool `json:"notifications"`
	BlockFriendRequests bool `json:"block_friend_requests"`
}

type UserStatsSchema struct {
	PostsCreated   uint `json:"posts_created"`
	PostsModified  uint `json:"posts_modified"`
	TopicsCreated  uint `json:"topics_created"`
	TopicsModified uint `json:"topics_modified"`
	TopicsLocked   uint `json:"topics_locked"`
	ReactionsAdded uint `json:"reactions_added"`
	PasswordResets uint `json:"password_resets"`
	ReportedPosts  uint `json:"reported_posts"`
	ReportedUsers  uint `json:"reported_users"`
}
