package ranks

type Administrator struct{}

func (admin Administrator) GetId() uint8 {
	return 3
}

func (admin Administrator) GetName() string {
	return "Administrator"
}

/* Global permissions */

func (admin Administrator) CanView() bool {
	return true
}

func (admin Administrator) CanReportUser() bool {
	return true
}

func (admin Administrator) CanViewAndReplyToUserReports() bool {
	return true
}

func (admin Administrator) CanBanUser() bool {
	return true
}

func (admin Administrator) CanDeleteUsersAvatar() bool {
	return true
}

func (admin Administrator) CanViewUserProfile() bool {
	return false
}

func (admin Administrator) CanDeleteUser() bool {
	return true
}

func (admin Administrator) CanPromoteDemoteUser() bool {
	return true
}

/* Posts permissions */

func (admin Administrator) CanReact() bool {
	return true
}

func (admin Administrator) CanCreatePost() bool {
	return true
}

func (admin Administrator) CanModifyOwnPost() bool {
	return true
}

func (admin Administrator) CanDeleteOwnPost() bool {
	return true
}

func (admin Administrator) CanReportPost() bool {
	return true
}

func (admin Administrator) CanViewAndReplyToPostReports() bool {
	return true
}

func (admin Administrator) CanModifyUsersPosts() bool {
	return true
}

func (admin Administrator) CanDeleteUsersPosts() bool {
	return true
}

/* Topics permissions */

func (admin Administrator) CanCreateTopic() bool {
	return true
}

func (admin Administrator) CanModifyOwnTopic() bool {
	return true
}

func (admin Administrator) CanDeleteOwnTopic() bool {
	return true
}

func (admin Administrator) CanLockOwnTopic() bool {
	return true
}

func (admin Administrator) CanLockUsersTopic() bool {
	return true
}

func (admin Administrator) CanReportTopic() bool {
	return true
}

func (admin Administrator) CanViewAndReplyToTopicReports() bool {
	return true
}

func (admin Administrator) CanModifyUsersTopics() bool {
	return true
}

func (admin Administrator) CanDeleteUsersTopics() bool {
	return true
}

/* Categories permissions */

func (admin Administrator) CanCreateCategory() bool {
	return true
}

func (admin Administrator) CanModifyCategory() bool {
	return true
}

func (admin Administrator) CanDeleteCategory() bool {
	return true
}
