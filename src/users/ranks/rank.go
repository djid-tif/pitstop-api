package ranks

type Rank interface {
	GetId() uint8
	GetName() string

	/* Global permissions */

	CanView() bool
	CanReportUser() bool
	CanViewAndReplyToUserReports() bool
	CanBanUser() bool
	CanDeleteUsersAvatar() bool
	CanViewUserProfile() bool
	CanDeleteUser() bool
	CanPromoteDemoteUser() bool

	/* Posts permissions */

	CanReact() bool
	CanCreatePost() bool
	CanModifyOwnPost() bool
	CanDeleteOwnPost() bool
	CanReportPost() bool
	CanViewAndReplyToPostReports() bool
	CanModifyUsersPosts() bool
	CanDeleteUsersPosts() bool

	/* Topics permissions */

	CanCreateTopic() bool
	CanModifyOwnTopic() bool
	CanDeleteOwnTopic() bool
	CanLockOwnTopic() bool
	CanLockUsersTopic() bool
	CanReportTopic() bool
	CanViewAndReplyToTopicReports() bool
	CanModifyUsersTopics() bool
	CanDeleteUsersTopics() bool

	/* Categories permissions */

	CanCreateCategory() bool
	CanModifyCategory() bool
	CanDeleteCategory() bool
}
