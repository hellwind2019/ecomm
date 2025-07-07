package storer

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestCreateProduct(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database %v", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	st := NewMySQLStorer(db)

	p := &Product{
		Name:         "Test Product",
		Image:        "test_image.jpg",
		Category:     "Test Category",
		Description:  "This is a test product",
		Rating:       5,
		NumReviews:   10,
		Price:        99.99,
		CountInStock: 100,
	}
	mock.ExpectExec(`
    INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`).WillReturnResult(sqlmock.NewResult(1, 1))
	cp, err := st.CreateProduct(context.Background(), p)
	require.NoError(t, err)
	require.Equal(t, int64(1), cp.ID)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)

}
