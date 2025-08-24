package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type ReplicationRoutingDataSource struct {
	masterDB  *sql.DB
	slaveDBs  []*sql.DB
	currentDB *sql.DB
	dbMutex   sync.RWMutex
	randomize bool
}

func NewReplicationRoutingDataSource(masterURL, slaveURLs []string, randomize bool) (*ReplicationRoutingDataSource, error) {
	masterDB, err := openDB(masterURL[0])
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}

	slaveDBs := make([]*sql.DB, len(slaveURLs))
	for i, url := range slaveURLs {
		db, err := openDB(url)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to slave database %s: %w", url, err)
		}
		slaveDBs[i] = db
	}

	return &ReplicationRoutingDataSource{
		masterDB:  masterDB,
		slaveDBs:  slaveDBs,
		currentDB: masterDB, // Начинаем с мастера
		dbMutex:   sync.RWMutex{},
		randomize: randomize,
	}, nil
}

func openDB(url string) (*sql.DB, error) {

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	// Важно:  Установите таймаут для подключения.
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(500)
	db.SetConnMaxLifetime(5 * time.Minute)
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func (r *ReplicationRoutingDataSource) GetDBMaster(ctx context.Context) (*sql.DB, error) {
	r.dbMutex.Lock()
	defer r.dbMutex.Unlock()

	if err := r.checkDB(r.masterDB); err != nil {
		return nil, fmt.Errorf("no available master databases", err)
	}

	return r.masterDB, nil
}

func (r *ReplicationRoutingDataSource) GetDB(ctx context.Context) (*sql.DB, error) {
	r.dbMutex.Lock()
	defer r.dbMutex.Unlock()

	// Проверяем доступность мастера
	if err := r.checkDB(r.masterDB); err != nil {
		// Если мастер недоступен, выбираем случайного слейва
		r.currentDB = r.ChooseSlave()
		if r.currentDB == nil {
			return nil, fmt.Errorf("no available slave databases")
		}
	}

	return r.currentDB, nil
}

func (r *ReplicationRoutingDataSource) ChooseSlave() *sql.DB {
	if len(r.slaveDBs) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())
	indexDbSlave := rand.Intn(2)

	err := r.slaveDBs[indexDbSlave].PingContext(context.Background())
	if err != nil {

		if indexDbSlave == 0 {
			indexDbSlave = 1
		} else {
			indexDbSlave = 0
		}
		return r.slaveDBs[indexDbSlave]
	}

	return r.slaveDBs[indexDbSlave] // Возвращаем первый слейв, если не randomize
}

func (r *ReplicationRoutingDataSource) checkDB(db *sql.DB) error {
	err := db.PingContext(context.Background())
	return err
}

func (r *ReplicationRoutingDataSource) Close() {
	r.dbMutex.Lock()
	defer r.dbMutex.Unlock()
	if r.masterDB != nil {
		r.masterDB.Close()
	}
	for _, db := range r.slaveDBs {
		if db != nil {
			db.Close()
		}
	}
}
