package mail

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed mail_credentials.json
var jsonCredentials []byte

//go:embed templates/*.tmpl
var templates embed.FS

type mailCredentials struct {
	Host     string `json:"SMTP_HOST"`
	Port     int    `json:"SMTP_PORT"`
	User     string `json:"SMTP_USER"`
	Password string `json:"SMTP_PASSWORD"`
}

var credentials mailCredentials

func init() {
	err := json.Unmarshal(jsonCredentials, &credentials)
	if err != nil {
		panic(fmt.Errorf("échec du parsing des identifiants de messagerie: %v\n", err))
	}

	/*dir, err := templates.ReadDir(".")
	if err != nil {
		panic(err)
	}
	fmt.Println("Templates found:")
	printDir(dir)*/
}

func printDir(dir []fs.DirEntry) {
	for _, dirEntry := range dir {
		if dirEntry.IsDir() {
			fmt.Println("=>", dirEntry.Name())
			subDir, err := templates.ReadDir(dirEntry.Name())
			if err != nil {
				panic(err)
			}
			printDir(subDir)
		} else {
			fmt.Println("==>", dirEntry.Name())
		}
	}
}

type writeReceiver struct {
	buffer []byte
}

func (r *writeReceiver) Write(p []byte) (n int, err error) {
	r.buffer = append(r.buffer, p...)
	return len(p), nil
}

func (r writeReceiver) Export() string {
	return string(r.buffer)
}
