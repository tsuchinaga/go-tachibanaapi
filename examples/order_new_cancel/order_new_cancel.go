package main

import (
	"context"
	"log"
	"time"

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

	// 新規注文
	var orderNumber string
	var executionDate time.Time
	{
		res, err := client.NewOrder(context.Background(), session, tachibana.NewOrderRequest{
			AccountType:         tachibana.AccountTypeSpecific,
			DeliveryAccountType: tachibana.DeliveryAccountTypeUnused,
			IssueCode:           "1475",
			Exchange:            tachibana.ExchangeToushou,
			Side:                tachibana.SideBuy,
			ExecutionTiming:     tachibana.ExecutionTimingClosing,
			OrderPrice:          0,
			OrderQuantity:       1,
			TradeType:           tachibana.TradeTypeStock,
			ExpireDate:          time.Time{},
			ExpireDateIsToday:   true,
			StopOrderType:       tachibana.StopOrderTypeNormal,
			TriggerPrice:        0,
			StopOrderPrice:      0,
			ExitPositionType:    tachibana.ExitPositionTypeUnused,
			SecondPassword:      examples.SecondPassword,
			ExitPositions:       nil,
		})

		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", res)
		if res.ResultCode != "0" {
			return
		}

		orderNumber = res.OrderNumber
		executionDate = res.ExecutionDate
	}

	// 取消注文
	{
		res, err := client.CancelOrder(context.Background(), session, tachibana.CancelOrderRequest{
			OrderNumber:    orderNumber,
			ExecutionDate:  executionDate,
			SecondPassword: examples.SecondPassword,
		})

		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", res)
		if res.ResultCode != "0" {
			return
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
