package transaction

import (
	"sync"

	"github.com/google/uuid"
)

// IDGenerator is the expected interface to be used when generating ids for stored transactions.
type IDGenerator interface {
	NewID() (string, error)
}

// NewUUIDGenerator creates an IDGenerator that generates UUIDs.
func NewUUIDGenerator() UUIDGenerator {
	return UUIDGenerator{}
}

// UUIDGenerator generates UUIDs
type UUIDGenerator struct {
}

// NewID returns a new UUID or an error (if one occurred)
func (f UUIDGenerator) NewID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

// NewInMemoryRepository creates a new in memory repository with the supplied id generator.
func NewInMemoryRepository(idGenerator IDGenerator) *InMemoryRepository {
	return &InMemoryRepository{
		data:        make(map[string]Entity),
		idGenerator: idGenerator,
	}
}

// InMemoryRepository stores transactions in an in memory map.  It generates a new id for each transaction as it is
// stored, using its configured IDGenerator.  This is intended a very simple way of storing transactions.  These
// transactions do no persist once the application is shut down.  In a production environment a repository such as
// this would manage communication with a real database to persist transactions long term.
type InMemoryRepository struct {
	data        map[string]Entity
	mu          sync.RWMutex
	idGenerator IDGenerator
}

// Save stores the transaction in the repository, or returns an error if one occurred.  It performs locking to ensure
// safe access for concurrent operations.
func (r *InMemoryRepository) Save(txn Entity) (Entity, error) {
	id, err := r.idGenerator.NewID()
	if err != nil {
		return Entity{}, err
	}
	txn.ID = id
	r.mu.Lock()
	r.data[txn.ID] = txn
	r.mu.Unlock()
	return txn, nil
}

// FindByID fetches the transaction with the provided id from the store.  An empty Entity will be returned if a
// transaction with the supplied id is not found.  It performs locking to ensure safe access for concurrent operations.
func (r *InMemoryRepository) FindByID(id string) Entity {
	r.mu.RLock()
	txn := r.data[id]
	r.mu.RUnlock()
	return txn
}
