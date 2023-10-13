package id

import "fmt"

// NewSequentialGenerator creates a new transaction.IDGenerator implementation.
func NewSequentialGenerator() *SequentialGenerator {
	return &SequentialGenerator{}
}

// SequentialGenerator generates sequential, predicatable ids to enable integration testing.
type SequentialGenerator struct {
	previousID int
}

// NewID returns a new id, or an error if there was a problem creating one.
func (g *SequentialGenerator) NewID() (string, error) {
	g.previousID++
	return fmt.Sprintf("sequentialID-%d", g.previousID), nil
}
