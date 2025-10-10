package database

import "strings"

func (db *DB) FuzzyMatchFundName(fundName string) (int, string, error) {
	fundId, err := db.GetFundByName(fundName)
	if err != nil {
		return 0, "", err
	}

	if fundId != 0 {
		return fundId, fundName, nil
	}

	if !strings.HasSuffix(fundName, "Fund") {
		fundId, err = db.GetFundByName(fundName + " Fund")
		if err != nil {
			return 0, "", err
		}
		if fundId != 0 {
			return fundId, fundName + " Fund", nil
		}
	}

	if strings.HasSuffix(fundName, " Fund") {
		trimmedName := strings.TrimSuffix(fundName, " Fund")
		fundId, err = db.GetFundByName(trimmedName)
		if err != nil {
			return 0, "", err
		}
		if fundId != 0 {
			return fundId, trimmedName, nil
		}
	}

	query := `
		SELECT trust_no, name 
		FROM funds 
		WHERE LOWER(name) = LOWER($1)
		LIMIT 1
	`

	var result struct {
		TrustNo int    `db:"trust_no"`
		Name    string `db:"name"`
	}

	err = db.conn.Get(&result, query, fundName)
	if err == nil {
		return result.TrustNo, result.Name, nil
	}

	query = `
		SELECT trust_no, name 
		FROM funds 
		WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
		ORDER BY LENGTH(name)
		LIMIT 1
	`

	err = db.conn.Get(&result, query, fundName)
	if err == nil {
		return result.TrustNo, result.Name, nil
	}

	return 0, "", nil
}

func NormalizeFundName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Join(strings.Fields(name), " ")

	replacements := map[string]string{
		"&":  "and",
		"  ": " ",
	}

	for old, new := range replacements {
		name = strings.ReplaceAll(name, old, new)
	}

	return name
}
