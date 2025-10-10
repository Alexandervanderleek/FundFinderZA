CREATE TYPE target_market_type AS ENUM ('Retail', 'Institutional');

CREATE TABLE fund_classes (
    id SERIAL PRIMARY KEY,
    fund_id INT NOT NULL REFERENCES funds(trust_no),
    class_name VARCHAR(50) NOT NULL,
    add_fee BOOLEAN,
    target_market target_market_type,
    max_init_fee DECIMAL(5,2),
    category VARCHAR(255),
    UNIQUE(fund_id, class_name)
);

CREATE INDEX idx_fund_classes_fund_id ON fund_classes(fund_id);
CREATE INDEX idx_fund_classes_category ON fund_classes(category);
CREATE INDEX idx_fund_classes_target_market ON fund_classes(target_market);