package signercache

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	signedCachePrefix = "signed-v6-"
	txMapPrefix       = "tx-map-v6-"
	vaultCachePrefix  = "vault-v6-"
)

// CacheStore manage the key value store used to store what tx out items have been signed before
type CacheStore struct {
	logger zerolog.Logger
	db     *leveldb.DB
}

// NewCacheStore create a new instance of CacheStore
func NewCacheStore(db *leveldb.DB) *CacheStore {
	return &CacheStore{
		db:     db,
		logger: log.With().Str("module", "signer-cache").Logger(),
	}
}

// SetSigned update key value store to set the given height and hash as signed
func (s *CacheStore) SetSigned(hash string) error {
	key := s.getSignedKey(hash)
	s.logger.Debug().Msgf("key:%s set to signed", key)
	return s.db.Put([]byte(key), []byte{1}, nil)
}

func (s *CacheStore) getSignedKey(hash string) string {
	return fmt.Sprintf("%s%s", signedCachePrefix, hash)
}

func (s *CacheStore) getMapKey(txHash string) string {
	return fmt.Sprintf("%s%s", txMapPrefix, txHash)
}

func (s *CacheStore) getVaultKey(vaultKey string) string {
	return fmt.Sprintf("%s%s", vaultCachePrefix, vaultKey)
}

// HasSigned check whether the given height and hash has been signed before or not
func (s *CacheStore) HasSigned(hash string) bool {
	key := s.getSignedKey(hash)
	exist, _ := s.db.Has([]byte(key), nil)
	s.logger.Debug().Msgf("key:%s has signed: %t", key, exist)
	return exist
}

// RemoveSigned removes the corresponding TxOutItem from the signer cache. The provided
// transaction hash should be for the broadcast transaction - it is internally mapped to
// the cache key for the TxOutItem.
func (s *CacheStore) RemoveSigned(transactionHash string) error {
	mapKey := s.getMapKey(transactionHash)
	value, err := s.db.Get([]byte(mapKey), nil)
	if err != nil {
		// bifrost didn't sign this tx , so it is fine
		if errors.Is(err, leveldb.ErrNotFound) {
			return nil
		}
		s.logger.Err(err).Msg("fail to check map key exist")
		return err
	}
	key := s.getSignedKey(string(value))
	if err = s.db.Delete([]byte(key), nil); err != nil {
		s.logger.Error().Err(err).Msgf("fail to remove %s from signed cache", string(value))
		return fmt.Errorf("fail to remove signed cache, err: %w", err)
	}
	return nil
}

// SetTransactionHashMap map a transaction hash to a tx out item hash
func (s *CacheStore) SetTransactionHashMap(txOutItemHash, transactionHash string) error {
	key := s.getMapKey(transactionHash)
	return s.db.Put([]byte(key), []byte(txOutItemHash), nil)
}

// SetLatestRecordedTx map a vault and transaction inbound or outbound to transaction hash
func (s *CacheStore) SetLatestRecordedTx(vaultKey, transactionHash string) error {
	key := s.getVaultKey(vaultKey)
	return s.db.Put([]byte(key), []byte(transactionHash), nil)
}

func (s *CacheStore) GetLatestRecordedTx(vaultKey string) (string, error) {
	key := s.getVaultKey(vaultKey)
	hash, err := s.db.Get([]byte(key), nil)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Close underlying db
func (s *CacheStore) Close() error {
	return s.db.Close()
}
