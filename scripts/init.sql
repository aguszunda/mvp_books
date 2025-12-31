-- Script de Inicialización de Base de Datos para mvp_books
-- Basado en la entidad 'internal/domain/Book' y convenciones de GORM.

CREATE DATABASE IF NOT EXISTS `books_db`;
USE `books_db`;

-- Tabla 'books' derivada de struct internal/domain/Book
-- struct Book {
--     ID     uint   `gorm:"primaryKey"`
--     Title  string
--     Author string
-- }

CREATE TABLE IF NOT EXISTS `books` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `title` longtext,
    `author` longtext,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Datos de ejemplo (Opcional, para testing manual)
INSERT INTO `books` (`title`, `author`) VALUES 
('The Go Programming Language', 'Alan A. A. Donovan'),
('Clean Architecture', 'Robert C. Martin');
