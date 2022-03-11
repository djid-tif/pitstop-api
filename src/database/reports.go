package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"pitstop-api/src/schemas"
	"time"
)

func GetAllReports(reportType string) ([]uint, error) {
	db, err := connect()
	if err != nil {
		return []uint{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM reports WHERE type = ?", reportType)
	if err != nil {
		return []uint{}, err
	}
	defer rows.Close()

	reports := []uint{}
	for rows.Next() {
		var reportId uint
		err = rows.Scan(&reportId)
		if err != nil {
			return []uint{}, err
		}
		reports = append(reports, reportId)
	}

	return reports, nil
}

func GetReport(id uint) (schemas.DbReport, error) {
	db, err := connect()
	if err != nil {
		return schemas.DbReport{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, type, target_id, creation_date, author_id, messages FROM reports WHERE id = ?", id)
	if err != nil {
		return schemas.DbReport{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return schemas.DbReport{}, errors.New("report not found")
	}

	report := new(schemas.DbReport)
	var jsonMessages string
	err = rows.Scan(&report.Id, &report.ReportType, &report.TargetId, &report.CreationDate, &report.AuthorId, &jsonMessages)
	if err != nil {
		return schemas.DbReport{}, err
	}

	err = json.Unmarshal([]byte(jsonMessages), &report.Messages)
	if err != nil {
		return schemas.DbReport{}, err
	}

	return *report, nil
}

func CreateReport(reportType string, targetId uint, authorId uint) (uint, error) {
	creationDate := time.Now().Unix()

	db, err := connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO reports(type, target_id, creation_date, author_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(reportType, targetId, creationDate, authorId)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), tx.Commit()
}

func UpdateReportById(reportId uint, field string, value interface{}) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE reports SET %s = ? WHERE id = ?", field))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, reportId)
	return err
}

func DeleteReport(reportId uint) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM reports WHERE id = ?", reportId)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM report_messages WHERE parent_id = ?", reportId)
	return err
}

func GetReportMessage(id uint) (schemas.DbReportMessage, error) {
	db, err := connect()
	if err != nil {
		return schemas.DbReportMessage{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, parent_id, creation_date, author_id, content FROM report_messages WHERE id = ?", id)
	if err != nil {
		return schemas.DbReportMessage{}, err
	}
	defer rows.Close()

	if !rows.Next() {
		return schemas.DbReportMessage{}, errors.New("message not found")
	}

	message := new(schemas.DbReportMessage)
	err = rows.Scan(&message.Id, &message.ParentId, &message.CreationDate, &message.AuthorId, &message.Content)
	if err != nil {
		return schemas.DbReportMessage{}, err
	}

	return *message, nil
}

func CreateReportMessage(reportId uint, authorId uint, content string) (uint, error) {
	creationDate := time.Now().Unix()

	db, err := connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO report_messages(parent_id, creation_date, author_id, content) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(reportId, creationDate, authorId, content)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), tx.Commit()
}
