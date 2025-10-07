package database

import "github.com/Alexandervanderleek/FundFinderZA/internal/models"

func (db *DB) SaveFunds(funds []*models.Fund) error {

	query := `
		INSERT INTO funds (trust_no, name, secondary_name, manager_id)
		VALUES (:trust_no, :name, :secondary_name, :manager_id)
		ON CONFLICT (trust_no) DO UPDATE
		SET name = EXCLUDED.name,
		 	secondary_name = EXCLUDED.secondary_name,
			manager_id = EXCLUDED.manager_id
		`

	_, err := db.conn.NamedExec(query, funds)

	return err

}
