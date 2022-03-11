package reports

import (
	"encoding/json"
	"errors"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"time"
)

const (
	userReport  = "user"
	topicReport = "topic"
	postReport  = "post"
)

type Report struct {
	id           uint
	reportType   string
	targetId     uint
	creationDate time.Time
	authorId     uint
	messages     []uint
}

type ReportMessage struct {
	id           uint
	parentId     uint
	creationDate time.Time
	authorId     uint
	content      string
}

type reportCreated struct {
	ReportId uint `json:"report_id"`
}

type reportList struct {
	Reports []uint `json:"reports"`
}

type reportMessageCreated struct {
	ReportMessageId uint `json:"report_message_id"`
}

func (report *Report) createReportMessage(authorId uint, content string) (uint, error) {
	messageId, err := database.CreateReportMessage(report.id, authorId, content)
	if err != nil {
		return 0, err
	}

	report.messages = append(report.messages, messageId)
	jsonList, err := json.Marshal(report.messages)
	if err != nil {
		return 0, err
	}

	return messageId, database.UpdateReportById(report.id, "messages", string(jsonList))
}

func (report Report) hasAccess(user users.User) bool {
	switch report.reportType {
	case userReport:
		return user.GetRank().CanViewAndReplyToUserReports()
	case topicReport:
		return user.GetRank().CanViewAndReplyToTopicReports()
	case postReport:
		return user.GetRank().CanViewAndReplyToPostReports()
	default:
		return false
	}
}

func (report Report) ToSchema() schemas.ReportSchema {
	return schemas.ReportSchema{
		Id:           report.id,
		ReportType:   report.reportType,
		TargetId:     report.targetId,
		CreationDate: utils.TimeToIso(report.creationDate),
		AuthorId:     report.authorId,
		Messages:     report.messages,
	}
}

func (report Report) close() error {
	return database.DeleteReport(report.id)
}

func (message ReportMessage) hasAccess(user users.User) bool {
	parentReport, found := loadReportFromId(message.parentId)
	if !found {
		return false
	}
	return parentReport.hasAccess(user)
}

func (message ReportMessage) ToSchema() schemas.ReportMessageSchema {
	return schemas.ReportMessageSchema{
		Id:           message.id,
		ParentId:     message.parentId,
		CreationDate: utils.TimeToIso(message.creationDate),
		AuthorId:     message.authorId,
		Content:      message.content,
	}
}

func isReportTypeValid(reportType string) bool {
	return reportType == userReport || reportType == topicReport || reportType == postReport
}

func getAllReports(reportType string) ([]uint, error) {
	if !isReportTypeValid(reportType) {
		return []uint{}, errors.New("type de report invalide: " + reportType)
	}
	return database.GetAllReports(reportType)
}

func loadReportFromId(id uint) (Report, bool) {
	dbReport, err := database.GetReport(id)
	if err != nil {
		return Report{}, false
	}

	return Report{
		id:           dbReport.Id,
		reportType:   dbReport.ReportType,
		targetId:     dbReport.TargetId,
		creationDate: time.Unix(dbReport.CreationDate, 0),
		authorId:     dbReport.AuthorId,
		messages:     dbReport.Messages,
	}, true
}

func loadReportMessageFromId(id uint) (ReportMessage, bool) {
	dbReport, err := database.GetReportMessage(id)
	if err != nil {
		return ReportMessage{}, false
	}

	return ReportMessage{
		id:           dbReport.Id,
		parentId:     dbReport.ParentId,
		creationDate: time.Unix(dbReport.CreationDate, 0),
		authorId:     dbReport.AuthorId,
		content:      dbReport.Content,
	}, true
}

func createReport(reportType string, targetId uint, authorId uint) (uint, error) {
	if !isReportTypeValid(reportType) {
		return 0, errors.New("type de report invalide: " + reportType)
	}
	return database.CreateReport(reportType, targetId, authorId)
}
