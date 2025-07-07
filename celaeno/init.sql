-- init.sql

CREATE DATABASE IF NOT EXISTS ctf;
USE ctf;

CREATE USER 'ctf_user'@'%' IDENTIFIED BY 'ctf_password';
GRANT ALL PRIVILEGES ON ctf.* TO 'ctf_user'@'%';

CREATE TABLE IF NOT EXISTS users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL
);

-- 可选：插入示例数据
INSERT INTO users (username, password) VALUES
('admin', '__flag__');
