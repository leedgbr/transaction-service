package transaction_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/date"
	"transaction-service/internal/id"
	"transaction-service/internal/transaction"
)

var (
	repository *transaction.InMemoryRepository
)

func TestRepository(t *testing.T) {
	t.Run("store", func(t *testing.T) {
		t.Run("success - should return entity with new id", func(t *testing.T) {
			setUpRepository()

			tcs := []struct {
				name       string
				entity     transaction.Entity
				wantEntity transaction.Entity
			}{
				{
					name:   "empty entity",
					entity: transaction.Entity{},
					wantEntity: transaction.Entity{
						ID: "sequentialID-1",
					},
				},
				{
					name: "complete entity (less ID)",
					entity: transaction.Entity{
						Description:     "*description*",
						TransactionDate: date.NewInUTC(2023, time.January, 23),
						AmountInCents:   5432,
					},
					wantEntity: transaction.Entity{
						ID:              "sequentialID-2",
						Description:     "*description*",
						TransactionDate: date.NewInUTC(2023, time.January, 23),
						AmountInCents:   5432,
					},
				},
			}
			for _, tc := range tcs {
				t.Run(tc.name, func(t *testing.T) {
					entity, err := repository.Save(tc.entity)
					assert.Nil(t, err)
					assert.Equal(t, tc.wantEntity, entity)
				})
			}

		})
		t.Run("failure - should return empty transaction and the error details", func(t *testing.T) {
			repository = transaction.NewInMemoryRepository(&alwaysErrorIDGenerator{})

			entity := transaction.Entity{
				Description:     "*description*",
				TransactionDate: date.NewInUTC(2023, time.January, 23),
				AmountInCents:   5432,
			}

			entity, err := repository.Save(entity)
			assert.Equal(t, errors.New("problem"), err)
			assert.Equal(t, transaction.Entity{}, entity)
		})
	})

	t.Run("find by id", func(t *testing.T) {
		t.Run("should return previously stored entity with id supplied", func(t *testing.T) {
			setUpRepository()
			repository.Save(transaction.Entity{
				Description:     "*description*",
				TransactionDate: date.NewInUTC(2023, time.January, 23),
				AmountInCents:   5432,
			})

			entity := repository.FindByID("sequentialID-1")
			wantEntity := transaction.Entity{
				ID:              "sequentialID-1",
				Description:     "*description*",
				TransactionDate: date.NewInUTC(2023, time.January, 23),
				AmountInCents:   5432,
			}
			assert.Equal(t, wantEntity, entity)
		})
		t.Run("should return empty entity when nothing has been stored with the provided id", func(t *testing.T) {
			setUpRepository()

			entity := repository.FindByID("sequentialID-1")
			wantEntity := transaction.Entity{}
			assert.Equal(t, wantEntity, entity)
		})
	})
}

func setUpRepository() {
	repository = transaction.NewInMemoryRepository(id.NewSequentialGenerator())
}

type alwaysErrorIDGenerator struct{}

func (id *alwaysErrorIDGenerator) NewID() (string, error) {
	return "", errors.New("problem")
}
