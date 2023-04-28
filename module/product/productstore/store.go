package productstore

import (
	"context"
	"database/sql"
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

func (s *store) Select(ctx context.Context, offset, limit int) ([]*productmodel.Product, error) {
	ctx = tracing.StartSpan(ctx, "productstore.Select")
	defer tracing.EndSpan(ctx)

	end := metric.Store().Start(name, "Select")
	defer end()

	cmd := fmt.Sprintf(`SELECT id, name FROM %s LIMIT ? OFFSET ?`, s.tbl)
	var rows []*productRow
	if err := s.db.SelectContext(ctx, &rows, cmd, limit, offset); err != nil {
		return nil, err
	}

	var results []*productmodel.Product
	for _, r := range rows {
		results = append(results, r.toModel())
	}
	return results, nil
}

func (s *store) Update(ctx context.Context, product *productmodel.Product) error {
	ctx = tracing.StartSpan(ctx, "productstore.Update")
	defer tracing.EndSpan(ctx)

	end := metric.Store().Start(name, "Update")
	defer end()

	var (
		stm = fmt.Sprintf(`UPDATE %v SET name = :name WHERE id = :id;`, s.tbl)
		row = toProductRow(product)
	)
	_, err := s.db.NamedExecContext(ctx, stm, row)
	if err != nil {
		return err
	}

	return nil
}
func (s *store) FindById(ctx context.Context, id int64) (*productmodel.Product, error) {
	ctx = tracing.StartSpan(ctx, "productstore.FindById")
	defer tracing.EndSpan(ctx)

	end := metric.Store().Start(name, "FindById")
	defer end()

	var row productRow
	cmd := fmt.Sprintf(`SELECT id, name FROM %s WHERE id = ?;`, s.tbl)
	err := s.db.GetContext(ctx, &row, cmd, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return row.toModel(), nil
}
func (s *store) Delete(ctx context.Context, id int64) error {
	ctx = tracing.StartSpan(ctx, "productstore.Delete")
	defer tracing.EndSpan(ctx)

	end := metric.Store().Start(name, "Delete")
	defer end()

	cmd := fmt.Sprintf(`DELETE FROM %s WHERE id = ?;`, s.tbl)
	if _, err := s.db.ExecContext(ctx, cmd, id); err != nil {
		return err
	}

	return nil
}
