CREATE TABLE fund_class_prices (
    id SERIAL PRIMARY KEY,
    fund_class_id INT REFERENCES fund_classes(id) ON DELETE CASCADE,
    price_date DATE NOT NULL,
    nav DECIMAL(12,2) NOT NULL,
    scraped_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(fund_class_id, price_date)
);

CREATE INDEX idx_fund_class_prices_fund_class_id ON fund_class_prices(fund_class_id);
CREATE INDEX idx_fund_class_prices_price_date ON fund_class_prices(price_date);
CREATE INDEX idx_fund_class_prices_scraped_at ON fund_class_prices(scraped_at);