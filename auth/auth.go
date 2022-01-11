package auth

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/term"
)

const authFileName string = ".auth"

type gitAuth struct {
	accessToken string
	id          string
}

func GetAuth() (string, string) {
	auth, err := loadAuth()
	checkErr(err)

	return auth.accessToken, auth.id
}

func inputAuth() gitAuth {
	var auth gitAuth

	fmt.Print("Github user id: ")
	fmt.Scanln(&auth.id)

	fmt.Print("Github access token: ")
	token, err := term.ReadPassword(int(syscall.Stdin))
	checkErr(err)
	auth.accessToken = string(token[:])

	return auth
}

func saveAuth() gitAuth {
	authFile, err := os.Create(authFileName)
	checkErr(err)
	defer authFile.Close()

	auth := inputAuth()

	fmt.Fprintln(authFile, auth.accessToken, auth.id)

	return auth
}

func loadAuth() (gitAuth, error) {
	authFile, err := os.Open(authFileName)
	defer authFile.Close()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			auth := saveAuth()
			return auth, nil
		} else {
			return gitAuth{}, err
		}
	}

	var accessToken string
	var id string

	fmt.Fscanf(authFile, "%s %s", &accessToken, &id)

	return gitAuth{accessToken: accessToken, id: id}, nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
