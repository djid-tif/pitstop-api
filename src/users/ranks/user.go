package ranks

type User struct{}

func (user User) GetId() uint8 {
	return 1
}

func (user User) GetName() string {
	return "User"
}

/* Global permissions */

func (user User) CanView() bool {
	return true
}

func (user User) CanReportUser() bool {
	return true
}

func (user User) CanViewAndReplyToUserReports() bool {
	return false
}

func (user User) CanBanUser() bool {
	return false
}

func (user User) CanDeleteUsersAvatar() bool {
	return false
}

func (user User) CanViewUserProfile() bool {
	return false
}

func (user User) CanDeleteUser() bool {
	return false
}

func (user User) CanPromoteDemoteUser() bool {
	return false
}

/* Posts permissions */

func (user User) CanReact() bool {
	return true
}

func (user User) CanCreatePost() bool {
	return true
}

func (user User) CanModifyOwnPost() bool {
	return true
}

func (user User) CanDeleteOwnPost() bool {
	return true
}

func (user User) CanReportPost() bool {
	return true
}

func (user User) CanViewAndReplyToPostReports() bool {
	return false
}

func (user User) CanModifyUsersPosts() bool {
	return false
}

func (user User) CanDeleteUsersPosts() bool {
	return false
}

/* Topics permissions */

func (user User) CanCreateTopic() bool {
	return true
}

func (user User) CanModifyOwnTopic() bool {
	return true
}

func (user User) CanDeleteOwnTopic() bool {
	return true
}

func (user User) CanLockOwnTopic() bool {
	return true
}

func (user User) CanLockUsersTopic() bool {
	return false
}

func (user User) CanReportTopic() bool {
	return false
}

func (user User) CanViewAndReplyToTopicReports() bool {
	return false
}

func (user User) CanModifyUsersTopics() bool {
	return false
}

func (user User) CanDeleteUsersTopics() bool {
	return false
}

/* Categories permissions */

func (user User) CanCreateCategory() bool {
	return false
}

func (user User) CanModifyCategory() bool {
	return false
}

func (user User) CanDeleteCategory() bool {
	return false
}
