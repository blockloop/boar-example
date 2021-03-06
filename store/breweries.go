package store

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/blockloop/beer_ratings/models"
	"github.com/blockloop/scan"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

var (
	selectBreweriesCols = scan.Columns(new(models.Brewery))
	selectBreweries     = sq.Select(selectBreweriesCols...).From("breweries")

	insertBreweriesCols = scan.Columns(new(models.Brewery), "id")
	insertBreweries     = sq.Insert("breweries").Columns(insertBreweriesCols...)
)

func NewBreweries(db *sql.DB) Breweries {
	return &breweries{
		db: db,
	}
}

type breweries struct {
	db *sql.DB
}

func (b *breweries) Create(ctx context.Context, userID int64, brewery models.Brewery) (*models.Brewery, error) {
	now := time.Now()

	brewery.Created = now
	brewery.Modified = now
	brewery.UUID = uuid.New()

	vals := scan.Values(insertBreweriesCols, &brewery)
	res, err := insertBreweries.
		Values(vals...).
		RunWith(b.db).
		ExecContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert brewery into db")
	}

	brewery.ID, err = res.LastInsertId()
	if err != nil {
		return &brewery, errors.Wrap(err, "failed to determine last inserted ID")
	}

	return &brewery, nil
}

func (b *breweries) Search(context.Context, string) ([]*models.Brewery, *models.Pagination, error) {
	panic("not implemented")
}

func (b *breweries) ForBeer(ctx context.Context, beerID int64) (*models.Brewery, error) {
	panic("not implemented")
}

func (b *breweries) Get(ctx context.Context, id int64) (*models.Brewery, error) {
	rows, err := selectBreweries.
		Where(sq.Eq{"id": id}).
		Limit(1).
		RunWith(b.db).
		QueryContext(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to execute query")
	}

	var brewery models.Brewery
	if err := scan.Row(&brewery, rows); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to scan brewery")
	}

	return &brewery, nil
}

func (b *breweries) Update(context.Context, models.Brewery) error {
	panic("not implemented")
}
