package nats_consumer

import (
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/stan.go"
	"github.com/noah-blockchain/coin-price-backend/internal/usecase"
	"github.com/noah-blockchain/coinExplorer-tools"
	"github.com/noah-blockchain/coinExplorer-tools/helpers"
)

type handlers struct {
	app usecase.Usecase
}

// StartConsumer starting consumer
func StartConsumer(sc stan.Conn, app usecase.Usecase) {

	h := &handlers{app}

	_, _ = sc.QueueSubscribe(
		helpers.CoinCreatedSubject,
		helpers.CoinCreatedSubject+"Queue",
		h.coinCreatedMessage,
		stan.DurableName(helpers.CoinCreatedSubject+"Name"),
		stan.StartWithLastReceived(),
	)
}

func (h *handlers) coinCreatedMessage(msg *stan.Msg) {
	log.Println("NEW COIN MESSAGE")
	eventStore := coin_extender.Coin{}
	err := proto.Unmarshal(msg.Data, &eventStore)
	if err != nil {
		return
	}
	log.Println("COIN NAME", eventStore.Symbol)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = h.app.CreateCoinInfo(ctx, eventStore); err != nil {
		log.Println(err)
	}
}
