package models

type Fund struct {
	TrustNo       int    `db:"trust_no"`
	Name          string `db:"name"`
	SecondaryName string `db:"secondary_name"`
	ManagerID     int    `db:"manager_id"`
}
