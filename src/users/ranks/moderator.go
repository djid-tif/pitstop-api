package ranks

type Moderator struct{}

func (mod Moderator) GetId() uint8 {
	return 2
}

func (mod Moderator) GetName() string {
	return "Moderator"
}

/* Global permissions */

func (mod Moderator) CanView() bool {
	return true
}

func (mod Moderator) CanReportUser() bool {
	return true
}

func (mod Moderator) CanViewAndReplyToUserReports() bool {
	return true
}

func (mod Moderator) CanBanUser() bool {
	return true
}

func (mod Moderator) CanDeleteUsersAvatar() bool {
	return true
}

func (mod Moderator) CanViewUserProfile() bool {
	return false
}

func (mod Moderator) CanDeleteUser() bool {
	return false
}

func (mod Moderator) CanPromoteDemoteUser() bool {
	return false
}

/* Posts permissions */

func (mod Moderator) CanReact() bool {
	return true
}

func (mod Moderator) CanCreatePost() bool {
	return true
}

func (mod Moderator) CanModifyOwnPost() bool {
	return true
}

func (mod Moderator) CanDeleteOwnPost() bool {
	return true
}

func (mod Moderator) CanReportPost() bool {
	return true
}

func (mod Moderator) CanViewAndReplyToPostReports() bool {
	return true
}

func (mod Moderator) CanModifyUsersPosts() bool {
	return true
}

func (mod Moderator) CanDeleteUsersPosts() bool {
	return true
}

/* Topics permissions */

func (mod Moderator) CanCreateTopic() bool {
	return true
}

func (mod Moderator) CanModifyOwnTopic() bool {
	return true
}

func (mod Moderator) CanDeleteOwnTopic() bool {
	return true
}

func (mod Moderator) CanLockOwnTopic() bool {
	return true
}

func (mod Moderator) CanLockUsersTopic() bool {
	return true
}

func (mod Moderator) CanReportTopic() bool {
	return true
}

func (mod Moderator) CanViewAndReplyToTopicReports() bool {
	return true
}

func (mod Moderator) CanModifyUsersTopics() bool {
	return true
}

func (mod Moderator) CanDeleteUsersTopics() bool {
	return true
}

/* Categories permissions */

func (mod Moderator) CanCreateCategory() bool {
	return false
}

func (mod Moderator) CanModifyCategory() bool {
	return false
}

func (mod Moderator) CanDeleteCategory() bool {
	return false
}
