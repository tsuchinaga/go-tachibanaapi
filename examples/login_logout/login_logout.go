package main

import (
	"context"
	"log"

	tachibana "gitlab.com/tsuchinaga/go-tachibanaapi"
)

func main() {
	userId := "user-id"
	password := "password"

	client := tachibana.NewClient(tachibana.EnvironmentProduction, tachibana.ApiVersionLatest)

	// ログイン
	var session *tachibana.Session
	{
		res, err := client.Login(context.Background(), tachibana.LoginRequest{
			UserId:   userId,
			Password: password,
		})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v\n", res)

		session, err = res.Session()
		if err != nil {
			log.Fatalln(err)
		}
	}

	// ログアウト
	{
		res, err := client.Logout(context.Background(), session, tachibana.LogoutRequest{})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v\n", res)
	}
}
