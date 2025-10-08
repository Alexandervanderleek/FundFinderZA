CREATE TABLE fund_class_costs (
    id SERIAL PRIMARY KEY,
    fund_class_id INT NOT NULL REFERENCES fund_classes(id) ON DELETE CASCADE,
    tic_date DATE,
    ter_perf_comp DECIMAL(5,2),
    ter DECIMAL(5,2),
    tc DECIMAL(5,2),
    tic DECIMAL(5,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fund_class_id, tic_date)
);

CREATE INDEX idx_fund_class_costs_fund_class_id ON fund_class_costs(fund_class_id);
CREATE INDEX idx_fund_class_costs_tic_date ON fund_class_costs(tic_date);