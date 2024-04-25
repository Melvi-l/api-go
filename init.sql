SHOW DATABASES;
CREATE DATABASE mydb;
USE mydb;
CREATE TABLE comment (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO comment (username, content) VALUES ('test_user', 'This is a comment.');
SELECT * FROM comment;
