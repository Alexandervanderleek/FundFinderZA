package models

type FundClass struct {
	ID           int      `db:"id"`
	FundID       int      `db:"fund_id"`
	ClassName    string   `db:"class_name"`
	AddFee       bool     `db:"add_fee"`
	MaxInitFee   *float64 `db:"max_init_fee"`
	Category     string   `db:"category"`
	TargetMarket string   `db:"target_market"`

	FundName string `db:"-"`
}

type FundClassCost struct {
	ID          int      `db:"id"`
	FundClassID int      `db:"fund_class_id"`
	TICDate     *string  `db:"tic_date"`
	TERPerfComp *float64 `db:"ter_perf_comp"`
	TER         *float64 `db:"ter"`
	TC          *float64 `db:"tc"`
	TIC         *float64 `db:"tic"`
}

type FundClassPrice struct {
	ID          int      `db:"id"`
	FundClassID int      `db:"fund_class_id"`
	PriceDate   *string  `db:"price_date"`
	NAV         *float64 `db:"nav"`
}
