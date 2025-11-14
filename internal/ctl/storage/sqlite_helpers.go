package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/mattn/go-sqlite3"
)

const sqliteMainDB = "main"

func openInMemoryDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, err
}

func deserializeInMemoryDBFromBytes(ctx context.Context, data []byte) (*sql.DB, error) {
	db, err := openInMemoryDB()
	if err != nil {
		return nil, err
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database during error handling: %v", closeErr)
		}
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	err = conn.Raw(func(driverConn any) error {
		sqliteConn, ok := driverConn.(*sqlite3.SQLiteConn)
		if !ok {
			return fmt.Errorf("driver connection is not SQLiteConn")
		}

		if err := sqliteConn.Deserialize(data, sqliteMainDB); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database during error handling: %v", closeErr)
		}
		return nil, err
	}

	return db, nil
}

func serializeInMemoryDBToBytes(ctx context.Context, db *sql.DB) ([]byte, error) {
	var bytes []byte

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

	err = conn.Raw(func(driverConn any) error {
		sqliteConn, ok := driverConn.(*sqlite3.SQLiteConn)
		if !ok {
			return fmt.Errorf("driver connection is not SQLiteConn")
		}

		data, err := sqliteConn.Serialize(sqliteMainDB)
		if err != nil {
			return err
		}
		bytes = data
		return nil
	})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
