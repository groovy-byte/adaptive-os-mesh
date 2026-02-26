package mirror

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ContextRecord struct {
	ID        int    `json:"id"`
	Account   string `json:"account"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp string `json:"timestamp"`
}

// MirrorService watches SQLite and prepares data for Cloud Upload
type MirrorService struct {
	db         *sql.DB
	lastSeenID int
	storePath  string
}

func NewMirrorService(sqlitePath string, storePath string) (*MirrorService, error) {
	db, err := sql.Open("sqlite3", sqlitePath)
	if err != nil {
		return nil, err
	}
	return &MirrorService{
		db:        db,
		storePath: storePath,
	}, nil
}

// SyncOnce reads new records and writes them to a JSON file for File Search indexing
func (m *MirrorService) SyncOnce() (int, error) {
	rows, err := m.db.Query("SELECT id, account, key, value, timestamp FROM shared_context WHERE id > ?", m.lastSeenID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var records []ContextRecord
	maxID := m.lastSeenID

	for rows.Next() {
		var r ContextRecord
		if err := rows.Scan(&r.ID, &r.Account, &r.Key, &r.Value, &r.Timestamp); err != nil {
			return 0, err
		}
		records = append(records, r)
		if r.ID > maxID {
			maxID = r.ID
		}
	}

	if len(records) == 0 {
		return 0, nil
	}

	// Write to a temporary JSON file for cloud indexing
	// In Task 2.2, we will use the actual file_search_upload tool
	filename := fmt.Sprintf("%s/sync_%d.json", m.storePath, time.Now().Unix())
	data, _ := json.MarshalIndent(records, "", "  ")
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return 0, err
	}

	log.Printf("[Mirror] Synced %d new records to %s", len(records), filename)
	m.lastSeenID = maxID
	return len(records), nil
}

func (m *MirrorService) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	log.Printf("[Mirror] Started SQLite watch on interval %v", interval)
	for range ticker.C {
		_, err := m.SyncOnce()
		if err != nil {
			log.Printf("[Mirror] Error during sync: %v", err)
		}
	}
}
