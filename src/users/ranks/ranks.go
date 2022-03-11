package ranks

import (
	"fmt"
	"pitstop-api/src/utils"
	"strings"
)

func GetRankById(id uint8) Rank {
	switch id {
	case 0: // Banned
		return Banned{}
	case 1: // User
		return User{}
	case 2: // Mod
		return Moderator{}
	case 3: // Admin
		return Administrator{}
	default:
		utils.PrintError(fmt.Errorf("rank with id %d not found", id))
		return nil
	}
}

func HasPermission(permission string, rank Rank) bool {
	switch strings.ToLower(permission) {
	case "view":
		return rank.CanView()
	case "report_user":
		return rank.CanReportUser()
	case "view_reply_to_user_reports":
		return rank.CanViewAndReplyToUserReports()
	case "ban_user":
		return rank.CanBanUser()
	case "delete_user_avatar":
		return rank.CanDeleteUsersAvatar()
	case "view_user_profile":
		return rank.CanViewUserProfile()
	case "delete_user":
		return rank.CanDeleteUser()
	case "promote_demote_user":
		return rank.CanPromoteDemoteUser()
	case "react":
		return rank.CanReact()
	case "create_post":
		return rank.CanCreatePost()
	case "modify_own_post":
		return rank.CanModifyOwnPost()
	case "delete_own_post":
		return rank.CanDeleteOwnPost()
	case "report_post":
		return rank.CanReportPost()
	case "view_reply_to_post_reports":
		return rank.CanViewAndReplyToPostReports()
	case "modify_users_posts":
		return rank.CanModifyUsersPosts()
	case "delete_users_posts":
		return rank.CanDeleteUsersPosts()
	case "create_topic":
		return rank.CanCreateTopic()
	case "modify_own_topic":
		return rank.CanModifyOwnTopic()
	case "delete_own_topic":
		return rank.CanDeleteOwnTopic()
	case "lock_own_topic":
		return rank.CanLockOwnTopic()
	case "lock_users_topics":
		return rank.CanLockUsersTopic()
	case "report_topic":
		return rank.CanReportTopic()
	case "view_reply_to_topic_reports":
		return rank.CanViewAndReplyToTopicReports()
	case "modify_users_topics":
		return rank.CanModifyUsersTopics()
	case "delete_users_topics":
		return rank.CanDeleteUsersTopics()
	case "create_category":
		return rank.CanCreateCategory()
	case "modify_category":
		return rank.CanModifyCategory()
	case "delete_category":
		return rank.CanDeleteCategory()
	default:
		return false
	}
}
