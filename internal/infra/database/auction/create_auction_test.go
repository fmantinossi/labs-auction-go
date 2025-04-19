package auction_test

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionAutoClose(t *testing.T) {
	os.Setenv("AUCTION_DURATION_SECONDS", "5")

	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		t.Fatalf("erro ao conectar no MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("testdb")
	repo := auction.NewAuctionRepository(db)

	auctionEntity := auction_entity.Auction{
		Id:          "teste123",
		ProductName: "produto teste",
		Category:    "categoria teste",
		Description: "descrição teste",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	err = repo.CreateAuction(context.Background(), &auctionEntity)
	if err != nil {
		t.Fatalf("erro ao criar leilão: %v", err)
	}

	time.Sleep(6 * time.Second)

	var result auction.AuctionEntityMongo
	err = repo.Collection.FindOne(context.Background(), map[string]interface{}{"_id": auctionEntity.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("erro ao buscar leilão: %v", err)
	}

	if result.Status != auction_entity.Completed {
		t.Errorf("esperado status Closed, mas retornou %v", result.Status)
	}

}
