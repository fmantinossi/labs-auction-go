package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

const defaultAuctionDuration = 60

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	ar.startAuctionExprationWatcher(ctx, *auctionEntityMongo)

	return nil
}

func (ar *AuctionRepository) startAuctionExprationWatcher(ctx context.Context, auction AuctionEntityMongo) {
	go func() {
		duration := time.Duration(defaultAuctionDuration) * time.Second
		if envDuration, err := strconv.Atoi(os.Getenv("ACTION_DURATION_SECONDS")); err != nil {
			duration = time.Duration(envDuration) * time.Second
		}

		timer := time.NewTimer(duration)
		<-timer.C

		var currentAuction AuctionEntityMongo
		err := ar.Collection.FindOne(ctx, map[string]interface{}{"_id": auction.Id}).Decode(&currentAuction)
		if err != nil || currentAuction.Status != auction_entity.Active {
			return
		}

		_, err = ar.Collection.UpdateOne(ctx,
			map[string]interface{}{"_id": auction.Id},
			map[string]interface{}{"$set": map[string]interface{}{"status": auction_entity.Completed}},
		)
		if err != nil {
			logger.Error("erro ao fechar leilão automaticamente", err)
		}
	}()

}
