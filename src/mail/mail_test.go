package mail

import (
	"pitstop-api/src/config"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"testing"
	"time"
)

func TestMail(t *testing.T) {
	t.Run("mail_test", testMail)
	t.Run("template_test", testTemplate)
}

func testMail(t *testing.T) {
	const username = "Michel"
	const email = "michel@email.net"

	err := SendTemplate(username, email, "Test", "test", struct {
		Username    string
		SendingDate string
	}{
		Username:    username,
		SendingDate: time.Now().Format("02/01/2006 15:04:05"),
	})
	if err != nil {
		utils.PrintError(err)
		t.Errorf("échec de l'envoi du mail: %v", err)
	}
}

func testTemplate(t *testing.T) {
	const username = "Michel"
	const email = "michel@email.net"

	t.Run("reset_password", func(t *testing.T) {
		err := SendTemplate(username, email, "Modification du mot de passe", "reset_password", schemas.ResetPassword{
			Username:    username,
			ResetLink:   config.OriginServerFront + "/api/user/password?token=TOKEN_HERE",
			PitStopLink: config.OriginServerFront,
			SendingDate: time.Now().Format("02/01/2006 15:04:05"),
		})
		if err != nil {
			utils.PrintError(err)
			t.Error(err)
		}
	})

	t.Run("confirm_email", func(t *testing.T) {
		err := SendTemplate(username, email, "Confirmation de votre adresse email", "confirm_email", schemas.ConfirmEmail{
			Username:    username,
			ConfirmLink: config.OriginServerFront + "/api/activate?token=TOKEN_HERE",
			PitStopLink: config.OriginServerFront,
			SendingDate: time.Now().Format("02/01/2006 15:04:05"),
		})
		if err != nil {
			utils.PrintError(err)
			t.Error(err)
		}
	})
}
