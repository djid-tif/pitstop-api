package users

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"pitstop-api/src/config"
	"pitstop-api/src/utils"
	"runtime"
	"strconv"
)

func updateAvatar(r *http.Request) (error, int) {
	_ = r.ParseForm()
	err := r.ParseMultipartForm(5 << 20) // For 5MB max file size
	if err != nil {
		return err, http.StatusInternalServerError
	}

	avatarFile, header, err := r.FormFile("avatar")
	if err != nil {
		return err, http.StatusInternalServerError
	}
	defer avatarFile.Close()

	if header.Size > 5000000 { // 5000000 = 5MB
		return errors.New("the image is too big"), http.StatusBadRequest
	}

	img, format, err := image.Decode(avatarFile)
	if err != nil {
		return err, http.StatusBadRequest
	}
	if format != "png" {
		return fmt.Errorf("type d'image invalide: %s", format), http.StatusBadRequest
	}

	user, found := LoadUserFromRequest(r)
	if !found {
		return errors.New("id invalide"), http.StatusUnauthorized
	}

	avatarPath := path.Join(config.AvatarsFolder, idToString(user.GetId())+".png")
	targetFile, err := os.Create(avatarPath)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	defer targetFile.Close()

	err = png.Encode(targetFile, img)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	return nil, http.StatusOK
}

func printAvatar(w http.ResponseWriter, id string) {
	file, err := os.Open(path.Join(config.AvatarsFolder, id+".png"))
	if err != nil {
		utils.Prettier(w, "aucun avatar trouvé", nil, http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, file)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, "échec de l'envoi de l'avatar", nil, http.StatusInternalServerError)
		return
	}
}

func deleteAvatar(id uint) error {
	filePath := path.Join(config.AvatarsFolder, idToString(id)+".png")
	err := os.Remove(filePath)
	if err != nil {
		if err == os.ErrNotExist {
			return errors.New("aucun avatar trouvé")
		}
		if runtime.GOOS == "windows" {
			return errors.New("error due to windows")
		}
		return err
	}
	return nil
}

func idToString(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}
