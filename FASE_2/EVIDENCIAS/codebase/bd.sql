CREATE TABLE "user" (
  id TEXT PRIMARY KEY,
  name VARCHAR,
  email VARCHAR UNIQUE,
  rut VARCHAR UNIQUE,
  uid TEXT,
  type VARCHAR, 
  created_at TIMESTAMP
);

CREATE TABLE customer (
  id TEXT PRIMARY KEY REFERENCES "user"(id),
  type VARCHAR 
);

CREATE TABLE employee (
  id TEXT PRIMARY KEY REFERENCES "user"(id),
  role VARCHAR 
);

CREATE TABLE parking (
  id SERIAL PRIMARY KEY,
  code VARCHAR,
  location TEXT,
  zone TEXT,
  is_active BOOLEAN
);

CREATE TABLE vehicle (
  id SERIAL PRIMARY KEY,
  plate VARCHAR UNIQUE,
  brand VARCHAR,
  model VARCHAR,
  vehicle_type VARCHAR, 
  customer_id TEXT REFERENCES customer(id),
  created_at TIMESTAMP,
);

CREATE TABLE reservation (
  id SERIAL PRIMARY KEY,
  customer_id TEXT REFERENCES customer(id),
  parking_id INT REFERENCES parking(id),
  vehicle_id INT REFERENCES vehicle(id),
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  status VARCHAR, 
  created_at TIMESTAMP
);

CREATE TABLE customer_schedule (
  id SERIAL PRIMARY KEY,
  customer_id TEXT REFERENCES customer(id),
  start_time TIMESTAMP,
  end_time TIMESTAMP,
  note TEXT, 
  created_at TIMESTAMP
);

CREATE TABLE parking_usage (
  id SERIAL PRIMARY KEY,
  reservation_id INT REFERENCES reservation(id),
  vehicle_id INT REFERENCES vehicle(id),
  parking_id INT REFERENCES parking(id),
  entry_time TIMESTAMP,
  exit_time TIMESTAMP,
  ocr_plate VARCHAR,
  qr_scanned BOOLEAN,
  registered_by TEXT REFERENCES employee(id),
  manual_entry BOOLEAN DEFAULT FALSE,
  visitor_name VARCHAR,
  visitor_rut VARCHAR,
  visitor_contact VARCHAR,
  zone TEXT
);

CREATE TABLE feedback (
  id SERIAL PRIMARY KEY,
  from_user_id TEXT REFERENCES "user"(id),
  comment TEXT,
  response_comment TEXT,
  rating INT, 
  created_at TIMESTAMP,
  response_at TIMESTAMP
);

CREATE TABLE notification_template (
  id SERIAL PRIMARY KEY,
  title VARCHAR,
  message TEXT,
  created_at TIMESTAMP
);

CREATE TABLE user_notification (
  id SERIAL PRIMARY KEY,
  user_id TEXT REFERENCES "user"(id),
  notification_template_id INT REFERENCES notification_template(id),
  is_read BOOLEAN DEFAULT FALSE,
  read_at TIMESTAMP
);

CREATE TABLE parking_access_log (
  id SERIAL PRIMARY KEY,
  plate VARCHAR,
  parking_id INT REFERENCES parking(id),
  detected_at TIMESTAMP,
  type VARCHAR, 
  status VARCHAR, 
  reservation_id INT REFERENCES reservation(id)
);


INSERT INTO parking (code, location, zone, is_active) VALUES
('P0001', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0002', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0003', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0004', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0005', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0006', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0007', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0008', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0009', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0010', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0011', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0012', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0013', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0014', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0015', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0016', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0017', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0018', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0019', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0020', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0021', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0022', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0023', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0024', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0025', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0026', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0027', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0028', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0029', 'Zona Estudiantes', 'Frontis Externo', TRUE),
('P0030', 'Zona Estudiantes', 'Frontis Externo', TRUE);


INSERT INTO parking (code, location, zone, is_active) VALUES
('P0031', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0032', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0033', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0034', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0035', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0036', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0037', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0038', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0039', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0040', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0041', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0042', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0043', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0044', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0045', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0046', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0047', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0048', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0049', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0050', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0051', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0052', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0053', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0054', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0055', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0056', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0057', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0058', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0059', 'Zona Profesores', 'Frontis Externo', TRUE),
('P0060', 'Zona Profesores', 'Frontis Externo', TRUE);


CREATE OR REPLACE VIEW active_parking_usages AS
SELECT  * FROM 
  parking_usage pu
WHERE 
  pu.exit_time IS NULL;
  
  CREATE OR REPLACE VIEW user_details AS
SELECT
  u.id AS user_id,
  u.name,
  u.email,
  u.rut,
  u.uid,
  u.type AS user_type,
  u.created_at,
  c.type AS customer_type,
  e.role AS employee_role
FROM
  miappdb.public.user u
LEFT JOIN customer c ON u.id = c.id
LEFT JOIN employee e ON u.id = e.id;
