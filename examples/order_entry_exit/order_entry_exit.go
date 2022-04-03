package main

import (
	"context"
	"log"
	"time"

	tachibana "gitlab.com/tsuchinaga/go-tachibanaapi"
)

// ログイン -> エントリー注文(成行) -> 約定確認 -> 約定情報取得 -> エグジット注文(OCO) -> 約定確認 -> ログアウト
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	userId := "user-id"
	password := "password"
	secondPassword := "password"

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
			ExecutionTiming:     tachibana.ExecutionTimingNormal,
			OrderPrice:          0,
			OrderQuantity:       1,
			TradeType:           tachibana.TradeTypeStock,
			ExpireDate:          time.Time{},
			ExpireDateIsToday:   true,
			StopOrderType:       tachibana.StopOrderTypeNormal,
			TriggerPrice:        0,
			StopOrderPrice:      0,
			ExitPositionType:    tachibana.ExitPositionTypeUnused,
			SecondPassword:      secondPassword,
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

	// 約定確認
	// 10sに1回確認する
	{
		var isDone bool
		var status tachibana.OrderStatus
		var contractedQuantity float64

		for {
			res, err := client.OrderList(context.Background(), session, tachibana.OrderListRequest{
				IssueCode:          "",
				ExecutionDate:      time.Time{},
				OrderInquiryStatus: tachibana.OrderInquiryStatusUnspecified,
			})
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("%+v\n", res)
			if res.ResultCode != "0" {
				return
			}

			for _, o := range res.Orders {
				// 違う注文は無視
				if o.OrderNumber != orderNumber {
					continue
				}

				if contractedQuantity != o.ContractQuantity || status != o.OrderStatus {
					contractedQuantity = o.ContractQuantity
					status = o.OrderStatus
					isDone = o.ContractStatus == tachibana.ContractStatusDone
				}
			}

			if isDone {
				log.Println("新規注文の約定確認完了")
				break
			}
			<-time.After(10 * time.Second)
		}
	}

	// 注文詳細から約定値を取得
	var price float64
	{
		res, err := client.OrderListDetail(context.Background(), session, tachibana.OrderListDetailRequest{
			OrderNumber:   orderNumber,
			ExecutionDate: executionDate,
		})
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", res)
		if res.ResultCode != "0" {
			return
		}

		price = res.Contracts[0].Price
	}

	// エグジット
	{
		res, err := client.NewOrder(context.Background(), session, tachibana.NewOrderRequest{
			AccountType:         tachibana.AccountTypeSpecific,
			DeliveryAccountType: tachibana.DeliveryAccountTypeUnused,
			IssueCode:           "1475",
			Exchange:            tachibana.ExchangeToushou,
			Side:                tachibana.SideSell,
			ExecutionTiming:     tachibana.ExecutionTimingNormal,
			OrderPrice:          price + 3,
			OrderQuantity:       1,
			TradeType:           tachibana.TradeTypeStock,
			ExpireDate:          time.Time{},
			ExpireDateIsToday:   true,
			StopOrderType:       tachibana.StopOrderTypeOCO,
			TriggerPrice:        price - 3,
			StopOrderPrice:      0,
			ExitPositionType:    tachibana.ExitPositionTypeUnused,
			SecondPassword:      secondPassword,
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

	// 約定確認
	// 10sに1回確認する
	{
		var isDone bool
		var contractedQuantity float64

		for {
			res, err := client.OrderList(context.Background(), session, tachibana.OrderListRequest{
				IssueCode:          "",
				ExecutionDate:      time.Time{},
				OrderInquiryStatus: tachibana.OrderInquiryStatusUnspecified,
			})
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("%+v\n", res)
			if res.ResultCode != "0" {
				return
			}

			for _, o := range res.Orders {
				// 違う注文は無視
				if o.OrderNumber != orderNumber {
					continue
				}

				// 注文の状態が変わっているかをチェック
				if o.ContractQuantity != contractedQuantity {
					contractedQuantity = o.ContractQuantity
				}

				isDone = o.ContractStatus == tachibana.ContractStatusDone
			}

			if isDone {
				log.Println("決済注文の約定確認完了")
				break
			}
			<-time.After(10 * time.Second)
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
