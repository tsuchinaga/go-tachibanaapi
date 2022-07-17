package main

import (
	"context"
	"log"

	"gitlab.com/tsuchinaga/go-tachibanaapi/examples"

	tachibana "gitlab.com/tsuchinaga/go-tachibanaapi"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client := tachibana.NewClient(tachibana.EnvironmentProduction, tachibana.ApiVersionLatest)

	// ログイン
	var session *tachibana.Session
	{
		res, err := client.Login(context.Background(), tachibana.LoginRequest{
			UserId:   examples.UserId,
			Password: examples.Password,
		})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v\n", res)
		if res.ResultCode != "0" || res.UnreadDocument {
			log.Fatalf("ResultCode: %s, ResultText: %s, UnreadDocument: %v\n", res.ResultCode, res.ResultText, res.UnreadDocument)
			return
		}

		session, err = res.Session()
		if err != nil {
			log.Fatalln(err)
			return
		}
		log.Printf("%+v\n", session)
	}

	// ストリーム
	{
		resCh, errCh := client.Stream(context.Background(), session, tachibana.StreamRequest{
			ColumnNumber:      nil,
			IssueCodes:        nil,
			MarketCodes:       nil,
			StartStreamNumber: 0,
			StreamEventTypes:  []tachibana.EventType{},
		})

		for {
			select {
			case err, ok := <-errCh:
				if ok {
					log.Printf("err: %+v\n", err)
				}
				return
			case res, ok := <-resCh:
				if !ok {
					return
				}
				log.Printf("res: %+v\n", res)
			}
		}
	}
}
