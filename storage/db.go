package storage

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/porthos-rpc/porthos-playground/models"
)

// DBStorage holds a db connection.
type DBStorage struct {
	db *sqlx.DB
}

// SetMaxIdleConns in the db connection
func (s *DBStorage) SetMaxIdleConns(n int) {
	s.db.SetMaxIdleConns(n)
}

// SetMaxOpenConns in the db connection
func (s *DBStorage) SetMaxOpenConns(n int) {
	s.db.SetMaxOpenConns(n)
}

// Ping the db connection
func (s *DBStorage) Ping() bool {
	_, err := s.db.Exec("SELECT 1")

	if err != nil {
		fmt.Errorf("Postgres Ping error: %s", err)
	}

	return err == nil
}

// SaveServiceSpecs persists the given specs.
func (s *DBStorage) SaveServiceSpecs(serviceSpecs *models.ServiceSpecs) {
	fmt.Println("Inserting service specs", serviceSpecs.Service, serviceSpecs.Specs)

	stmtDelete, _ := s.db.Prepare(`DELETE FROM specs WHERE service = ?`)
	defer stmtDelete.Close()

	_, err := stmtDelete.Exec(serviceSpecs.Service)

	if err != nil {
		fmt.Errorf("Error executing sql statement: %s", err)
		return
	}

	stmtInsert, _ := s.db.Prepare(`INSERT INTO specs(service, specs) VALUES(?, ?)`)
	defer stmtInsert.Close()

	_, err = stmtInsert.Exec(serviceSpecs.Service, serviceSpecs.Specs)

	if err != nil {
		fmt.Errorf("Error executing sql statement: %s", err)
	}
}

// GetSpecs returns all persisted specs.
func (s *DBStorage) GetSpecs() ([]*models.ServiceSpecs, error) {
	specs := []*models.ServiceSpecs{}

	err := s.db.Select(&specs, `SELECT service, specs FROM specs`)

	if err != nil {
		return nil, err
	}

	return specs, nil
}

// NewDb creates a new DB connection.
func NewDb(driver, url string) *sqlx.DB {
	db, err := sqlx.Connect(driver, url)

	if err != nil {
		panic(err)
	}

	_, err = createSchemaIfNotExists(db)

	if err != nil {
		panic(err)
	}

	return db
}

func createSchemaIfNotExists(db *sqlx.DB) (sql.Result, error) {
	schema := `CREATE TABLE IF NOT EXISTS specs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		service VARCHAR NOT NULL,
		specs TEXT NOT NULL
	);`

	// execute a query on the server
	return db.Exec(schema)
}

// NewStorage creates a new DB.
func NewStorage(db *sqlx.DB) Storage {
	return &DBStorage{db}
}
