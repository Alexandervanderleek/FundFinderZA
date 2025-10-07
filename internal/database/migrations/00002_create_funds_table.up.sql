CREATE TABLE funds (
    trust_no INT PRIMARY KEY,
    name VARCHAR(500) NOT NULL,
    secondary_name VARCHAR(500),
    manager_id INT NOT NULL,
    FOREIGN KEY (manager_id) REFERENCES cisManagers(id)
)