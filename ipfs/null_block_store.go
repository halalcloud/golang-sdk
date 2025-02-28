package ipfs

import (
	"context"

	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
)

type NullBlockstore struct {
}

func (s *NullBlockstore) DeleteBlock(ctx context.Context, req cid.Cid) error {

	return nil
}

// Has(context.Context, cid.Cid) (bool, error)
func (s *NullBlockstore) Has(ctx context.Context, req cid.Cid) (bool, error) {

	return false, nil
}

// Get(context.Context, cid.Cid) (blocks.Block, error)
func (s *NullBlockstore) Get(ctx context.Context, req cid.Cid) (blocks.Block, error) {

	// data, err := storagePool.Read(oid, )
	return nil, nil
}

// GetSize(context.Context, cid.Cid) (int, error)
func (s *NullBlockstore) GetSize(ctx context.Context, req cid.Cid) (int, error) {

	return 0, nil
}

// Put(context.Context, blocks.Block) error
func (s *NullBlockstore) Put(ctx context.Context, block blocks.Block) error {

	return nil
}

// PutMany(context.Context, []blocks.Block) error
func (s *NullBlockstore) PutMany(ctx context.Context, blocks []blocks.Block) error {
	return nil
}

// AllKeysChan(ctx context.Context) (<-chan cid.Cid, error)
func (s *NullBlockstore) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	return nil, nil
}

// HashOnRead(enabled bool)
func (s *NullBlockstore) HashOnRead(enabled bool) {
}
