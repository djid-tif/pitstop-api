package ranks

type Banned struct{}

func (banned Banned) GetId() uint8 {
	return 0
}

func (banned Banned) GetName() string {
	return "Banned"
}

/* Global permissions */

func (banned Banned) CanView() bool {
	return true
}

func (banned Banned) CanReportUser() bool {
	return false
}

func (banned Banned) CanViewAndReplyToUserReports() bool {
	return false
}

func (banned Banned) CanBanUser() bool {
	return false
}

func (banned Banned) CanDeleteUsersAvatar() bool {
	return false
}

func (banned Banned) CanViewUserProfile() bool {
	return false
}

func (banned Banned) CanDeleteUser() bool {
	return false
}

func (banned Banned) CanPromoteDemoteUser() bool {
	return false
}

/* Posts permissions */

func (banned Banned) CanReact() bool {
	return false
}

func (banned Banned) CanCreatePost() bool {
	return false
}

func (banned Banned) CanModifyOwnPost() bool {
	return false
}

func (banned Banned) CanDeleteOwnPost() bool {
	return false
}

func (banned Banned) CanReportPost() bool {
	return false
}

func (banned Banned) CanViewAndReplyToPostReports() bool {
	return false
}

func (banned Banned) CanModifyUsersPosts() bool {
	return false
}

func (banned Banned) CanDeleteUsersPosts() bool {
	return false
}

/* Topics permissions */

func (banned Banned) CanCreateTopic() bool {
	return false
}

func (banned Banned) CanModifyOwnTopic() bool {
	return false
}

func (banned Banned) CanDeleteOwnTopic() bool {
	return false
}

func (banned Banned) CanLockOwnTopic() bool {
	return false
}

func (banned Banned) CanLockUsersTopic() bool {
	return false
}

func (banned Banned) CanReportTopic() bool {
	return false
}

func (banned Banned) CanViewAndReplyToTopicReports() bool {
	return false
}

func (banned Banned) CanModifyUsersTopics() bool {
	return false
}

func (banned Banned) CanDeleteUsersTopics() bool {
	return false
}

/* Categories permissions */

func (banned Banned) CanCreateCategory() bool {
	return false
}

func (banned Banned) CanModifyCategory() bool {
	return false
}

func (banned Banned) CanDeleteCategory() bool {
	return false
}
