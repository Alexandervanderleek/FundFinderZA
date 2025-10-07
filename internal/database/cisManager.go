package database

import (
	"fmt"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
)

func (db *DB) GetAllCISManagers() ([]*models.CISManager, error) {
	var cisMangers []*models.CISManager

	err := db.conn.Select(&cisMangers, "SELECT * FROM cisManagers")

	if err != nil {
		return nil, fmt.Errorf("failed to select all cisManagers: %w", err)
	}

	return cisMangers, nil
}

func (db *DB) SaveCISManagers(cisManager []*models.CISManager) error {
	query := `
		INSERT INTO cisManagers (id, name)
		VALUES (:id, :name)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name
	`
	_, err := db.conn.NamedExec(query, cisManager)

	return err
}
