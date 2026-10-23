//go:build no_store_gcs

package desync

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
)

var _ WriteStore = GCStore{}

type GCStoreBase struct {
}

type GCStore struct {
	GCStoreBase
}

func NewGCStoreBase(u *url.URL, opt StoreOptions) (GCStoreBase, error) {

	return GCStoreBase{}, errors.New("no_store_gcs")
}

func (s GCStoreBase) String() string {
	return ""
}

func (s GCStoreBase) Close() error { return nil }

func NewGCStore(location *url.URL, opt StoreOptions) (s GCStore, e error) {
	return GCStore{}, errors.New("no_store_gcs")
}

func (s GCStore) GetChunk(id ChunkID) (*Chunk, error) {

	return nil, errors.New("no_store_gcs")
}

func (s GCStore) StoreChunk(chunk *Chunk) error {
	return errors.New("no_store_gcs")
}

func (s GCStore) HasChunk(id ChunkID) (bool, error) {
	return false, errors.New("no_store_gcs")
}

func (s GCStore) RemoveChunk(id ChunkID) error {
	return errors.New("no_store_gcs")
}

func (s GCStore) Prune(ctx context.Context, ids map[ChunkID]struct{}) error {
	return errors.New("no_store_gcs")
}
