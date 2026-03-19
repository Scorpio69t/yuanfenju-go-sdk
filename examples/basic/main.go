package main

import (
	"context"
	"fmt"
	"log"
	"time"

	yuanfenju "github.com/Scorpio69t/yuanfenju-go-sdk"
)

func main() {
	client, err := yuanfenju.NewClient(yuanfenju.Config{
		APIKey:  "replace_with_your_api_key",
		Timeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	merchant, err := client.Free.QueryMerchant(ctx)
	if err != nil {
		log.Fatalf("query merchant failed: %v", err)
	}
	fmt.Printf("merchant type: %s\n", merchant.Data.MerchantType)

	meiri, err := client.Divination.Meiri(ctx, yuanfenju.MeiriRequest{Lang: "zh-cn"})
	if err != nil {
		log.Fatalf("meiri failed: %v", err)
	}
	fmt.Printf("meiri gua: %s (#%d)\n", meiri.Data.GuaMing, meiri.Data.Number)
}
