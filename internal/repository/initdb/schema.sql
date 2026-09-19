CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    "status" TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS departments (
    id SERIAL PRIMARY KEY,
    department TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id SERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    department_id INT NOT NULL,
    status_id INT NOT NULL,

    CONSTRAINT status_id_fk FOREIGN KEY (status_id) REFERENCES statuses (id),
    CONSTRAINT department_id_fk FOREIGN KEY (department_id) REFERENCES departments (id),
    UNIQUE (first_name, last_name) 
);

CREATE TABLE IF NOT EXISTS pending_employees (
    id SERIAL PRIMARY KEY,
    creating_time TIMESTAMPTZ DEFAULT NOW(),
    employee_id INT,

    CONSTRAINT employee_id_fk FOREIGN KEY (employee_id) REFERENCES employees (id)
);

CREATE INDEX IF NOT EXISTS department_idx
ON departments USING HASH (department)