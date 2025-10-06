package database

import (
	"fmt"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
)

func (db *DB) GetAllCisManagers() ([]*models.CISManager, error) {
	var cisMangers []*models.CISManager

	err := db.conn.Select(&cisMangers, "SELECT * FROM cisManagers")

	if err != nil {
		return nil, fmt.Errorf("failed to select all cisManagers: %w", err)
	}

	return cisMangers, nil
}

func (db *DB) SaveCisManagers(cisManager []*models.CISManager) error {
	query := `INSERT INTO cismanagers (id, name)
			  VALUES (:id, :name)`
	_, err := db.conn.NamedExec(query, cisManager)

	return err
}
