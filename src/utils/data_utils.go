package utils

import (
	"errors"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

func IsEmailValid(email string) error {

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if len(email) == 0 {
		return errors.New("aucun email fourni")
	}

	if len(email) < 5 && len(email) > 254 {
		return errors.New("la longueur de l'email doit être comprise entre 5 et 254 caractères")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("l'email ne respecte pas le format demandé")
	}

	return nil
}

func IsUserNameValid(username string) error {

	userNameRegex := regexp.MustCompile(`^([[:alnum:]]|[_,-]){3,16}$`)

	if len(username) == 0 {
		return errors.New("aucun nom d'utilisateur fourni")
	}

	if len(username) < 3 || len(username) > 16 {
		return errors.New("la longueur du nom d'utilisateur doit être comprise entre 3 et 16 caractères")
	}

	if !userNameRegex.MatchString(username) {
		return errors.New("le nom d'utilisateur ne respecte pas le format demandé")
	}

	return nil
}

func IsPasswordValid(password string) error {

	if !regexp.MustCompile(`(?i)^.{8,64}$`).MatchString(password) {
		return errors.New("le mot de passe ne respecte pas le format demandé")
	}

	if !regexp.MustCompile(`^.*[A-Z]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins une lettre majuscule")
	}

	if !regexp.MustCompile(`^.*[a-z]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins une lettre miniscule")
	}

	if !regexp.MustCompile(`(?i)^.*[0-9]+.*$`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins un chiffre")
	}

	if !regexp.MustCompile(`^.*([ -/]|[:-@]|[\[-\x60]|[{-~])+`).MatchString(password) {
		return errors.New("le mot de passe doit contenir au moins un caractère spécial")
	}

	return nil
}

func IsTitleValid(title string) bool {
	return len(title) >= 1 && len(title) <= 255
}

func Slugify(name string) (slug string) {
	slug = strings.ToLower(name)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, err := transform.String(t, slug)
	if err == nil {
		slug = output
	}
	re := regexp.MustCompile(`(?i)[^0-9a-z]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
