-- Script SQL para crear usuarios de prueba
-- NOTA: Este script requiere que los usuarios ya existan (creados vía API)
-- Luego actualiza el usuario 'admin' a tipo 'admin'
-- 
-- Ejecutar después de crear usuarios con la API:
-- docker-compose exec mysql mysql -u spotly -pspotlypass spotly_users < scripts/create_test_users.sql

USE spotly_users;

-- Actualizar usuario admin a tipo 'admin' (si ya existe)
UPDATE users SET user_type = 'admin' WHERE username = 'admin';

-- Verificar usuarios creados
SELECT id, username, email, user_type, created_at FROM users WHERE username IN ('admin', 'usuario');

