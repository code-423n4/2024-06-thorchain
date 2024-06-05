package signer

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"gitlab.com/thorchain/thornode/bifrost/blockscanner"
	"gitlab.com/thorchain/thornode/bifrost/db"
	"gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/config"
)

const (
	DefaultSignerLevelDBFolder = "signer_data"
	txOutPrefix                = "txout-v4-"
)

type TxStatus int

const (
	TxUnknown TxStatus = iota
	TxAvailable
	TxUnavailable
	TxSpent
)

type TxOutStoreItem struct {
	TxOutItem    types.TxOutItem
	Status       TxStatus
	Height       int64
	Index        int64
	Round7Retry  bool
	Checkpoint   []byte
	SignedTx     []byte
	RetrievalKey string `json:"-"`
	// RetrievalKey is to ensure consistent KV overwrite/deletion after iterator retrieval;
	// the json "-" tag is to not store it in the KVStore.
}

func NewTxOutStoreItem(height int64, item types.TxOutItem, idx int64) TxOutStoreItem {
	return TxOutStoreItem{
		TxOutItem: item,
		Height:    height,
		Status:    TxAvailable,
		Index:     idx,
	}
}

func (s *TxOutStoreItem) Key() string {
	// If this is a retrieved item then refer to the same key-value pair
	// for overwriting/deletion, else newly derive it.
	if len(s.RetrievalKey) != 0 {
		return s.RetrievalKey
	}

	buf, _ := json.Marshal(struct {
		TxOutItem types.TxOutItem
		Height    int64
		Index     int64
	}{
		s.TxOutItem,
		s.Height,
		s.Index,
	})
	sha256Bytes := sha256.Sum256(buf)
	return fmt.Sprintf("%s%s", txOutPrefix, hex.EncodeToString(sha256Bytes[:]))
}

type SignerStorage interface {
	Set(item TxOutStoreItem) error
	Batch(items []TxOutStoreItem) error
	Get(key string) (TxOutStoreItem, error)
	Has(key string) bool
	Remove(item TxOutStoreItem) error
	List() []TxOutStoreItem
	OrderedLists() map[string][]TxOutStoreItem
	Close() error
}

type SignerStore struct {
	*blockscanner.LevelDBScannerStorage
	logger     zerolog.Logger
	db         *leveldb.DB
	passphrase string
}

// NewSignerStore create a new instance of SignerStore. If no folder is given,
// an in memory implementation is used.
func NewSignerStore(levelDbFolder string, opts config.LevelDBOptions, passphrase string) (*SignerStore, error) {
	ldb, err := db.NewLevelDB(levelDbFolder, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create level db: %w", err)
	}

	levelDbStorage, err := blockscanner.NewLevelDBScannerStorage(ldb)
	if err != nil {
		return nil, fmt.Errorf("failed to create scanner storage: %w", err)
	}
	return &SignerStore{
		LevelDBScannerStorage: levelDbStorage,
		logger:                log.With().Str("module", "signer-storage").Logger(),
		db:                    ldb,
		passphrase:            passphrase,
	}, nil
}

func (s *SignerStore) Set(item TxOutStoreItem) error {
	key := item.Key()
	buf, err := json.Marshal(item)
	if err != nil {
		s.logger.Error().Err(err).Msg("fail to marshal to txout store item")
		return err
	}
	if err = s.db.Put([]byte(key), buf, nil); err != nil {
		s.logger.Error().Err(err).Msg("fail to set txout item")
		return err
	}
	return nil
}

func (s *SignerStore) Batch(items []TxOutStoreItem) error {
	batch := new(leveldb.Batch)
	for _, item := range items {
		key := item.Key()
		buf, err := json.Marshal(item)
		if err != nil {
			s.logger.Error().Err(err).Msg("fail to marshal to txout store item")
			return err
		}
		batch.Put([]byte(key), buf)
	}
	return s.db.Write(batch, nil)
}

func (s *SignerStore) Get(keyString string) (item TxOutStoreItem, err error) {
	key := []byte(keyString)

	ok, err := s.db.Has(key, nil)
	if !ok || err != nil {
		return
	}
	buf, _ := s.db.Get(key, nil)
	if err = json.Unmarshal(buf, &item); err != nil {
		s.logger.Error().Err(err).Msg("fail to unmarshal to txout store item")
		return item, err
	}
	// Record the key so not needing to successfully rederive to overwrite/delete the key-value pair.
	item.RetrievalKey = keyString
	return
}

// Has check whether the given key exist in key value store
func (s *SignerStore) Has(key string) (ok bool) {
	ok, _ = s.db.Has([]byte(key), nil)
	return
}

// Remove remove the given item from key values store
func (s *SignerStore) Remove(item TxOutStoreItem) error {
	return s.db.Delete([]byte(item.Key()), nil)
}

// List send back tx out to retry depending on arg failed only
func (s *SignerStore) List() []TxOutStoreItem {
	iterator := s.db.NewIterator(util.BytesPrefix([]byte(txOutPrefix)), nil)
	defer iterator.Release()
	var results []TxOutStoreItem
	for iterator.Next() {
		buf := iterator.Value()
		if len(buf) == 0 {
			continue
		}

		var item TxOutStoreItem
		if err := json.Unmarshal(buf, &item); err != nil {
			s.logger.Error().Err(err).Msg("fail to unmarshal to txout store item")
			continue
		}

		// ignore already spent items
		if item.Status == TxSpent {
			continue
		}

		// Record the key so not needing to successfully rederive to overwrite/delete the key-value pair.
		item.RetrievalKey = string(iterator.Key())

		results = append(results, item)
	}

	// Ensure that we sort our list by block height (lowest to highest), then
	// by Hash. This makes best efforts to ensure that each node is iterating
	// through their list of items as closely as possible
	sort.SliceStable(results, func(i, j int) bool { return results[i].TxOutItem.Hash() < results[j].TxOutItem.Hash() })
	sort.SliceStable(results, func(i, j int) bool { return results[i].Height < results[j].Height })
	return results
}

// OrderedLists
func (s *SignerStore) OrderedLists() map[string][]TxOutStoreItem {
	lists := make(map[string][]TxOutStoreItem)
	for _, item := range s.List() {
		key := fmt.Sprintf("%s-%s", item.TxOutItem.Chain.String(), item.TxOutItem.VaultPubKey.String())
		if _, ok := lists[key]; !ok {
			lists[key] = make([]TxOutStoreItem, 0)
		}
		lists[key] = append(lists[key], item)
	}
	return lists
}

// Close underlying db
func (s *SignerStore) Close() error {
	return s.db.Close()
}

func (s *SignerStore) GetInternalDb() *leveldb.DB {
	return s.db
}
