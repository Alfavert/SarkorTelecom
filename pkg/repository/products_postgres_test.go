package repository

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestProducts_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()
	r := NewProductPostgres(db)

	type args struct {
		product   Products
		productId int
	}
	type mockBehaviour func(args args, id int)

	tests := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		id            int
		wantErr       bool
	}{{
		name: "success",
		args: args{
			productId: 1,
			product: Products{
				Name:        "Product 1",
				Description: "Product Description 1",
				Price:       10,
				Quantity:    10,
			},
		},
		id: 1,
		mockBehaviour: func(args args, id int) {
			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
			mock.ExpectQuery(`INSERT INTO product \(name, description, price, quantity\) values \(\$1, \$2, \$3, \$4\) RETURNING id`).
				WithArgs(args.product.Name, args.product.Description, args.product.Price, args.product.Quantity).
				WillReturnRows(rows)

			mock.ExpectCommit()
		},
	}, {
		name: "Empty Fields",
		args: args{
			productId: 1,
			product: Products{
				Name:        "",
				Description: "",
				Price:       1,
				Quantity:    1,
			},
		},
		mockBehaviour: func(args args, id int) {
			mock.ExpectBegin()
			rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
			mock.ExpectQuery(`INSERT INTO product \(name, description, price, quantity\) values \(\$1, \$2, \$3, \$4\) RETURNING id`).
				WithArgs(args.product.Name, args.product.Description, args.product.Price, args.product.Quantity).
				WillReturnRows(rows)

			mock.ExpectRollback()
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args, tt.id)

			got, err := r.CreateProduct(tt.args.product)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, got)
			}
		})
	}
}
func TestProducts_Update(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		logrus.Fatal(err)
	}
	defer db.Close()
	r := NewProductPostgres(db)

	type args struct {
		productId int
		input     UpdateProducts
	}

	tests := []struct {
		name          string
		args          args
		mockBehaviour func(args args)
		wantErr       bool
	}{
		{
			name: "success",
			args: args{
				productId: 1,
				input: UpdateProducts{
					Name:        strPtr("Updated Product"),
					Description: strPtr("Updated Description"),
					Price:       float64Ptr(15.5),
					Quantity:    intPtr(50),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE product SET name = \$1, description = \$2, price = \$3, quantity = \$4 WHERE id = \$5`).
					WithArgs(args.input.Name, args.input.Description, args.input.Price, args.input.Quantity, args.productId).
					WillReturnResult(sqlmock.NewResult(1, 1)) // Возвращаем успешный результат

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "transaction begin error",
			args: args{
				productId: 1,
				input: UpdateProducts{
					Name:        strPtr("Updated Product"),
					Description: strPtr("Updated Description"),
					Price:       float64Ptr(15.5),
					Quantity:    intPtr(50),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
			},
			wantErr: true,
		},
		{
			name: "update query error",
			args: args{
				productId: 1,
				input: UpdateProducts{
					Name:        strPtr("Updated Product"),
					Description: strPtr("Updated Description"),
					Price:       float64Ptr(15.5),
					Quantity:    intPtr(50),
				},
			},
			mockBehaviour: func(args args) {
				mock.ExpectBegin()

				mock.ExpectExec(`UPDATE product SET name = \$1, description = \$2, price = \$3, quantity = \$4 WHERE id = \$5`).
					WithArgs(args.input.Name, args.input.Description, args.input.Price, args.input.Quantity, args.productId).
					WillReturnError(errors.New("update query error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour(tt.args)

			err := r.Update(tt.args.productId, tt.args.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func TestAuthPostgres_GetAll(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	r := &AuthPostgres{db: db}

	tests := []struct {
		name          string
		mockBehaviour func()
		want          []Products
		wantErr       bool
	}{
		{
			name: "success",
			mockBehaviour: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "quantity"}).
					AddRow(1, "Product 1", "Description 1", 10.0, 100).
					AddRow(2, "Product 2", "Description 2", 20.0, 200)

				mock.ExpectQuery(`SELECT id, name, description, price, quantity FROM product`).
					WillReturnRows(rows)

				mock.ExpectCommit()
			},
			want: []Products{
				{Id: 1, Name: "Product 1", Description: "Description 1", Price: 10.0, Quantity: 100},
				{Id: 2, Name: "Product 2", Description: "Description 2", Price: 20.0, Quantity: 200},
			},
			wantErr: false,
		},
		{
			name: "transaction begin error",
			mockBehaviour: func() {
				mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "query error",
			mockBehaviour: func() {
				mock.ExpectBegin()

				mock.ExpectQuery(`SELECT id, name, description, price, quantity FROM product`).
					WillReturnError(errors.New("query error"))

				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "rows scan error",
			mockBehaviour: func() {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "quantity"}).
					AddRow(1, "Product 1", "Description 1", "invalid_price", 100)

				mock.ExpectQuery(`SELECT id, name, description, price, quantity FROM product`).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			got, err := r.GetAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func TestAuthPostgres_Delete(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	r := &AuthPostgres{db: db}

	tests := []struct {
		name          string
		mockBehaviour func()
		productId     int
		wantErr       bool
	}{
		{
			name:      "success",
			productId: 1,
			mockBehaviour: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM product WHERE id = \$1`). // Используем правильный формат для PostgreSQL
											WithArgs(1).
											WillReturnResult(sqlmock.NewResult(1, 1)) // Один результат с одной строкой, которая удалена

				mock.ExpectCommit() // Ожидание коммита транзакции
			},
			wantErr: false,
		},
		{
			name:      "transaction begin error",
			productId: 1,
			mockBehaviour: func() {
				mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
			},
			wantErr: true,
		},
		{
			name:      "delete query error",
			productId: 1,
			mockBehaviour: func() {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM product WHERE id = \$1`). // Используем правильный формат для PostgreSQL
											WithArgs(1).
											WillReturnError(errors.New("delete query error"))

				mock.ExpectRollback() // Ожидание отката транзакции
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			err := r.Delete(tt.productId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
func TestAuthPostgres_GetById(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	require.NoError(t, err)
	defer db.Close()

	r := &AuthPostgres{db: db}

	tests := []struct {
		name          string
		mockBehaviour func()
		productId     int
		want          Products
		wantErr       bool
	}{
		{
			name:      "success",
			productId: 1,
			want: Products{
				Id:          1,
				Name:        "Product 1",
				Description: "Product Description 1",
				Price:       10.0,
				Quantity:    100,
			},
			mockBehaviour: func() {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "quantity"}).
					AddRow(1, "Product 1", "Product Description 1", 10.0, 100)
				mock.ExpectQuery(`SELECT id, name, description, price, quantity FROM product WHERE id = \$1`).
					WithArgs(1).WillReturnRows(rows)
				mock.ExpectCommit() // Ожидание коммита транзакции
			},
			wantErr: false,
		},
		{
			name:      "transaction begin error",
			productId: 1,
			mockBehaviour: func() {
				mock.ExpectBegin().WillReturnError(errors.New("transaction begin error"))
			},
			want:    Products{},
			wantErr: true,
		},
		{
			name:      "query error",
			productId: 1,
			mockBehaviour: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT id, name, description, price, quantity FROM product WHERE id = \$1`).
					WithArgs(1).WillReturnError(errors.New("query error"))
				mock.ExpectRollback() // Ожидание отката транзакции
			},
			want:    Products{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehaviour()

			got, err := r.GetById(tt.productId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
