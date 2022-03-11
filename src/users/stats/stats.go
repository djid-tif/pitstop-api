package stats

import (
	"encoding/json"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
)

type UserStats struct {
	postsCreated   uint
	postsModified  uint
	topicsCreated  uint
	topicsModified uint
	topicsLocked   uint
	reactionsAdded uint
	passwordResets uint
	reportedPosts  uint
	reportedUsers  uint
}

func (stats *UserStats) CreatePost(id uint) {
	stats.postsCreated++
	stats.save(id)
}

func (stats *UserStats) ModifyPost(id uint) {
	stats.postsModified++
	stats.save(id)
}

func (stats *UserStats) CreateTopic(id uint) {
	stats.topicsCreated++
	stats.save(id)
}

func (stats *UserStats) ModifyTopic(id uint) {
	stats.topicsModified++
	stats.save(id)
}

func (stats *UserStats) LockTopic(id uint) {
	stats.topicsLocked++
	stats.save(id)
}

func (stats *UserStats) AddReaction(id uint) {
	stats.reactionsAdded++
	stats.save(id)
}

func (stats *UserStats) ResetPassword(id uint) {
	stats.passwordResets++
	stats.save(id)
}

func (stats *UserStats) ReportPost(id uint) {
	stats.reportedPosts++
	stats.save(id)
}

func (stats *UserStats) ReportUser(id uint) {
	stats.reportedUsers++
	stats.save(id)
}

func (stats *UserStats) save(id uint) {
	marshaledStats, err := json.Marshal(stats.ToSchema())
	if err != nil {
		utils.PrintError(err)
		return
	}
	err = database.UpdateUserById(id, "stats", string(marshaledStats))
	if err != nil {
		utils.PrintError(err)
	}
}

func (stats *UserStats) ToSchema() schemas.UserStatsSchema {
	return schemas.UserStatsSchema{
		PostsCreated:   stats.postsCreated,
		PostsModified:  stats.postsModified,
		TopicsCreated:  stats.topicsCreated,
		TopicsModified: stats.topicsModified,
		TopicsLocked:   stats.topicsLocked,
		ReactionsAdded: stats.reactionsAdded,
		PasswordResets: stats.passwordResets,
		ReportedPosts:  stats.reportedPosts,
		ReportedUsers:  stats.reportedUsers,
	}
}

func LoadStats(statsSchema schemas.UserStatsSchema) UserStats {
	return UserStats{
		postsCreated:   statsSchema.PostsCreated,
		postsModified:  statsSchema.PostsModified,
		topicsCreated:  statsSchema.TopicsCreated,
		topicsModified: statsSchema.TopicsModified,
		topicsLocked:   statsSchema.TopicsLocked,
		reactionsAdded: statsSchema.ReactionsAdded,
		passwordResets: statsSchema.PasswordResets,
		reportedPosts:  statsSchema.ReportedPosts,
		reportedUsers:  statsSchema.ReportedUsers,
	}
}
