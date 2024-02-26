package client

import (
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/jsondb"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/local"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/local/storage"
)

type dataStoreFactoryImpl struct {
	cfg *config.Config
}

var _ persistence.DataStoreFactory = (*dataStoreFactoryImpl)(nil)

func NewDataStoreFactory(cfg *config.Config) persistence.DataStoreFactory {
	return &dataStoreFactoryImpl{
		cfg: cfg,
	}
}

func (f dataStoreFactoryImpl) NewHistoryStore() persistence.HistoryStore {
	// TODO: Add support for other data stores (e.g. sqlite, postgres, etc.)
	return jsondb.New(f.cfg.DataDir, f.cfg.DAGs)
}

func (f dataStoreFactoryImpl) NewDAGStore() persistence.DAGStore {
	return local.NewDAGStore(f.cfg.DAGs)
}

func (f dataStoreFactoryImpl) NewFlagStore() persistence.FlagStore {
	s := storage.NewStorage(f.cfg.SuspendFlagsDir)
	return local.NewFlagStore(s)
}
