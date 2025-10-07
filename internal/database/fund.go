package database

import "github.com/Alexandervanderleek/FundFinderZA/internal/models"

func (db *DB) SaveFunds(funds []*models.Fund) error {

	query := `
		INSERT INTO funds (trust_no, name, secondary_name, manager_id)
		VALUES (:trust_no, :name, :secondary_name, :manager_id)
		ON CONFLICT (trust_no) DO UPDATE
		SET name = EXCLUDE.name,
		 	secondary_name = EXCLUDE.secondary_name,
			manager_id = EXCLUDE.manager_id
		`

	_, err := db.conn.NamedExec(query, funds)

	return err

}
