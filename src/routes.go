package main

import (
	"github.com/gorilla/mux"
	"pitstop-api/src/categories"
	"pitstop-api/src/categories/topics"
	"pitstop-api/src/categories/topics/posts"
	"pitstop-api/src/reports"
	"pitstop-api/src/users"
	"pitstop-api/src/users/auth"
	"pitstop-api/src/users/friends"
)

func initRoutes() *mux.Router {
	r := mux.NewRouter()

	// auth
	r.HandleFunc("/api/auth/login", auth.Login).Methods("POST")
	r.HandleFunc("/api/auth/logout", auth.AuthRequired(auth.Logout)).Methods("POST")
	r.HandleFunc("/api/auth/register", auth.Register).Methods("POST")
	r.HandleFunc("/api/auth/forgot", auth.Forgot).Methods("POST")
	r.HandleFunc("/api/auth/activate/{token}", auth.Activate).Methods("POST")
	r.HandleFunc("/api/auth/refresh", auth.Refresh).Methods("POST")
	r.HandleFunc("/api/auth/create", auth.Create).Methods("POST")      // temp
	r.HandleFunc("/api/auth/token", auth.DecryptToken).Methods("POST") // temp

	// users //

	//avatar
	r.HandleFunc("/api/user/avatar", auth.AuthRequired(users.GetOwnAvatar)).Methods("GET")
	r.HandleFunc("/api/user/{id}/avatar", users.GetAvatar).Methods("GET")
	r.HandleFunc("/api/user/avatar", auth.AuthRequired(users.UpdateOwnAvatar)).Methods("PUT")
	r.HandleFunc("/api/user/avatar", auth.AuthRequired(users.DeleteOwnAvatar)).Methods("DELETE")
	r.HandleFunc("/api/user/{id}/avatar", auth.AuthRequired(users.DeleteAvatar)).Methods("DELETE")

	// update user
	r.HandleFunc("/api/user", auth.AuthRequired(users.UpdateUser)).Methods("PUT")

	// password
	r.HandleFunc("/api/user/password", auth.AuthRequired(users.UpdateOwnPassword)).Methods("PUT")
	r.HandleFunc("/api/user/password/{token}", users.UpdatePassword).Methods("PUT")
	r.HandleFunc("/api/user/valid-password", auth.AuthRequired(users.ValidPassword)).Methods("POST")

	// get user
	r.HandleFunc("/api/user", auth.AuthRequired(users.GetCurrentUser)).Methods("GET")
	r.HandleFunc("/api/user/{id}", users.HandleUser).Methods("GET")
	r.HandleFunc("/api/user/search/{username}", users.SearchUserByUsername).Methods("GET")

	// promote demote
	r.HandleFunc("/api/user/{id}/promote", auth.AuthRequired(users.Promote)).Methods("POST")
	r.HandleFunc("/api/user/{id}/demote", auth.AuthRequired(users.Demote)).Methods("POST")

	// ban unban
	r.HandleFunc("/api/user/{id}/ban", auth.AuthRequired(users.BanUser)).Methods("POST")
	r.HandleFunc("/api/user/{id}/unban", auth.AuthRequired(users.UnBanUser)).Methods("POST")

	//permission
	r.HandleFunc("/api/user/permission/{permission}", auth.AuthRequired(users.HasPermission)).Methods("GET")

	// settings
	r.HandleFunc("/api/user/settings/{setting}", auth.AuthRequired(users.UpdateSettings)).Methods("PUT")

	// report
	r.HandleFunc("/api/report/{id}", auth.AuthRequired(reports.GetReport)).Methods("GET")
	r.HandleFunc("/api/report/{id}", auth.AuthRequired(reports.ReplyToReport)).Methods("POST")
	r.HandleFunc("/api/report/message/{id}", auth.AuthRequired(reports.GetReportMessage)).Methods("GET")
	r.HandleFunc("/api/report/{id}", auth.AuthRequired(reports.CloseReport)).Methods("DELETE")

	r.HandleFunc("/api/reports/user", auth.AuthRequired(reports.ListUserReports)).Methods("GET")
	r.HandleFunc("/api/report/user/{id}", auth.AuthRequired(reports.ReportUser)).Methods("POST")

	r.HandleFunc("/api/reports/topic", auth.AuthRequired(reports.ListTopicReports)).Methods("GET")
	r.HandleFunc("/api/report/topic/{id}", auth.AuthRequired(reports.ReportTopic)).Methods("POST")

	r.HandleFunc("/api/reports/post", auth.AuthRequired(reports.ListPostReports)).Methods("GET")
	r.HandleFunc("/api/report/post/{id}", auth.AuthRequired(reports.ReportPost)).Methods("POST")

	// friends
	r.HandleFunc("/api/friend/{id}", auth.AuthRequired(friends.AddFriend)).Methods("PUT")
	r.HandleFunc("/api/friend/{id}", auth.AuthRequired(friends.RemoveFriend)).Methods("DELETE")
	r.HandleFunc("/api/friends", auth.AuthRequired(friends.ListFriends)).Methods("GET")

	// categories
	r.HandleFunc("/api/categories", categories.ListCategories).Methods("GET")
	r.HandleFunc("/api/category/{id}", categories.GetCategory).Methods("GET")
	r.HandleFunc("/api/category", auth.AuthRequired(categories.CreateCategory)).Methods("POST")
	r.HandleFunc("/api/category/{id}", auth.AuthRequired(categories.UpdateCategory)).Methods("PUT")
	r.HandleFunc("/api/category/{id}", auth.AuthRequired(categories.DeleteCategory)).Methods("DELETE")

	// topics
	r.HandleFunc("/api/category/{id}/topics", categories.GetTopicsOfCategory).Methods("GET")
	r.HandleFunc("/api/topic/{id}", topics.GetTopic).Methods("GET")
	r.HandleFunc("/api/category/{id}/topic", auth.AuthRequired(categories.CreateTopicInCategory)).Methods("POST")
	r.HandleFunc("/api/topic/{id}", auth.AuthRequired(topics.UpdateTopic)).Methods("PUT")
	r.HandleFunc("/api/topic/{id}/pin", auth.AuthRequired(topics.PinTopic)).Methods("POST")
	r.HandleFunc("/api/topic/{id}/unpin", auth.AuthRequired(topics.UnpinTopic)).Methods("POST")
	r.HandleFunc("/api/topic/{id}/lock", auth.AuthRequired(topics.LockTopic)).Methods("POST")
	r.HandleFunc("/api/topic/{id}/unlock", auth.AuthRequired(topics.UnlockTopic)).Methods("POST")
	r.HandleFunc("/api/topic/{id}", auth.AuthRequired(topics.DeleteTopic)).Methods("DELETE")

	// posts
	r.HandleFunc("/api/topic/{id}/posts", topics.GetPostsOfTopic).Methods("GET")
	r.HandleFunc("/api/post/{id}", posts.GetPost).Methods("GET")
	r.HandleFunc("/api/topic/{id}/post", auth.AuthRequired(topics.CreatePostInTopic)).Methods("POST")
	r.HandleFunc("/api/post/{id}", auth.AuthRequired(posts.UpdatePost)).Methods("PUT")
	r.HandleFunc("/api/post/{id}/like", auth.AuthRequired(posts.LikePost)).Methods("POST")
	r.HandleFunc("/api/post/{id}/dislike", auth.AuthRequired(posts.DislikePost)).Methods("POST")
	r.HandleFunc("/api/post/{id}", auth.AuthRequired(posts.DeletePost)).Methods("DELETE")

	r.HandleFunc("/api/trends/post", posts.GetMostLikedPosts).Methods("GET")

	return r
}
