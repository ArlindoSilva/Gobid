package services

import (
	"context"
	"errors"

	"github.com/ArlindoSilva/gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidsService struct {
	pool *pgxpool.Pool
	Queries *pgstore.Queries
}

func NewBidsService (pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool: pool,
		Queries: pgstore.New(pool),
	}
}

var ErrBidIsTooLow = errors.New("the bid value is too low")

func (bs *BidsService) Placebid(ctx context.Context, product_id, bidder_id uuid.UUID, amount float64) (pgstore.Bid, error) {
	//Amount > previous_amount
	//Amount > baseprice

	product, err := bs.Queries.GetProductById(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	highestBid, err := bs.Queries.GetHighestBidByProductId(ctx, product_id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}

	if product.Baseprice >= amount || highestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}

	highestBid, err = bs.Queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: product_id,
		BidderID: bidder_id,
		BidAmount: amount,
	})

	if err != nil {
		return pgstore.Bid{}, err
	}

	return highestBid, nil
}