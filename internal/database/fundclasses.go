package database

import (
	"fmt"

	"github.com/Alexandervanderleek/FundFinderZA/internal/models"
)

func (db *DB) SaveFundClass(fundClass *models.FundClass) error {
	query := `
		INSERT INTO fund_classes (fund_id, class_name, target_market, add_fee, max_init_fee, category)
		VALUES (:fund_id, :class_name, :target_market, :add_fee, :max_init_fee, :category)
		ON CONFLICT (fund_id, class_name) DO UPDATE
		SET target_market = EXCLUDED.target_market,
			add_fee = EXCLUDED.add_fee,
			max_init_fee = EXCLUDED.max_init_fee,
			category = EXCLUDED.category
	`

	_, err := db.conn.NamedExec(query, fundClass)

	return err
}

func (db *DB) SaveFundClassCosts(fundClassCost *models.FundClassCost) error {
	query := `
		INSERT INTO fund_class_costs (fund_class_id, tic_date, ter_perf_comp, ter, tc, tic)
		VALUES (:fund_class_id, :tic_date, :ter_perf_comp, :ter, :tc, :tic)
		ON CONFLICT (fund_class_id, tic_date) DO UPDATE
		SET ter_perf_comp = EXCLUDED.ter_perf_comp,
			ter = EXCLUDED.ter,
			tc = EXCLUDED.tc,
			tic = EXCLUDED.tic
	`

	_, err := db.conn.NamedExec(query, fundClassCost)

	return err
}

func (db *DB) SaveFundClassPrice(fundClassPrice *models.FundClassPrice) error {
	query := `
		INSERT INTO fund_class_prices (fund_class_id, price_date, nav)
		VALUES (:fund_class_id, :price_date, :nav)
		ON CONFLICT (fund_class_id, price_date) DO UPDATE
		SET nav = EXCLUDED.nav
	`
	_, err := db.conn.NamedExec(query, fundClassPrice)

	return err
}

func (db *DB) GetFundByName(name string) (int, error) {
	var fundID int

	query := `
		SELECT trust_no
		FROM funds
		WHERE name = $1
		LIMIT 1
	`

	err := db.conn.Get(&fundID, query, name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 0, nil
		}
		return 0, err
	}

	return fundID, nil
}

func (db *DB) GetAllFundNames() (map[int]string, error) {
	query := `SELECT trust_no, name FROM funds`

	rows, err := db.conn.Query(query)

	if err != nil {
		return nil, fmt.Errorf("error getting fund names: %s", err)
	}
	defer rows.Close()

	fundNames := make(map[int]string)
	for rows.Next() {
		var trustNo int
		var name string

		if err := rows.Scan(&trustNo, &name); err != nil {
			return nil, fmt.Errorf("error scaning in fund: %s", err)
		}

		fundNames[trustNo] = name
	}

	return fundNames, nil
}
