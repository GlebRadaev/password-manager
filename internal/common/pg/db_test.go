package pg

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

type mockDatabase struct{}

func (m *mockDatabase) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}
func (m *mockDatabase) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func (m *mockDatabase) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return nil
}

func (m *mockDatabase) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return pgconn.NewCommandTag(""), nil
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDB := NewMockDatabase(ctrl)
	dbInstance := New(mockDB)
	assert.NotNil(t, dbInstance)
	assert.IsType(t, &db{}, dbInstance)
	assert.Equal(t, mockDB, dbInstance.(*db).Database)
}
