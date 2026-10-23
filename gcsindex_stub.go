//go:build no_store_gcs

package desync

import (
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type GCIndexStore struct {
	GCStoreBase
}

func NewGCIndexStore(location *url.URL, opt StoreOptions) (s GCIndexStore, e error) {
	return GCIndexStore{}, errors.New("no_store_gcs")
}

func (s GCIndexStore) GetIndexReader(name string) (r io.ReadCloser, err error) {
	return nil, errors.New("no_store_gcs")
}

func (s GCIndexStore) GetIndex(name string) (i Index, e error) {
	return Index{}, errors.New("no_store_gcs")
}

func (s GCIndexStore) StoreIndex(name string, idx Index) error {
	errors.New("no_store_gcs")
	return nil
}
