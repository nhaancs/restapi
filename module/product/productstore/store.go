package productstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"restapi/module/product/productmodel"
	"restapi/pkg/metric"
	"restapi/pkg/tracing"
)

var ErrNotFound = errors.New("product not found")

type productRow struct {
	Id   int64  `db:"id,primarykey"`
	Name string `db:"name"`
}

func toProductRow(m *productmodel.Product) *productRow {
	return &productRow{
		Id:   m.Id,
		Name: m.Name,
	}
}

func (r *productRow) toModel() *productmodel.Product {
	return &productmodel.Product{
		Id:   r.Id,
		Name: r.Name,
	}
}

type store struct {
	db  *sqlx.DB
	tbl string
}

func New(db *sqlx.DB) *store {
	return &store{
		db:  db,
		tbl: "products",
	}
}

const name = "product.store"

func (s *store) Insert(ctx context.Context, product *productmodel.Product) error {
	ctx = tracing.StartSpan(ctx, "productstore.Insert")
	defer tracing.EndSpan(ctx)

	end := metric.Store().Start(name, "Insert")
	defer end()

	cmd := fmt.Sprintf(`INSERT INTO %s (name) VALUES (:name)`, s.tbl)
	row := toProductRow(product)
	if _, err := s.db.NamedExecContext(ctx, cmd, row); err != nil {
		return err
	}
	return nil
}
