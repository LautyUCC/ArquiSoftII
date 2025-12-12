#!/bin/bash

# Script para crear usuarios de prueba usando la API
# Este script crea los usuarios y luego actualiza uno a admin usando SQL directo

echo "🔧 Creando usuarios de prueba..."

# URL base de la API (a través de nginx)
API_URL="http://localhost/api/users"

# 1. Crear usuario ADMINISTRADOR
echo ""
echo "📝 Creando usuario administrador..."
ADMIN_RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@spotly.com",
    "password": "admin123",
    "firstName": "Admin",
    "lastName": "Usuario"
  }')

echo "Respuesta: $ADMIN_RESPONSE"

# 2. Crear usuario NORMAL
echo ""
echo "📝 Creando usuario normal..."
USER_RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "usuario",
    "email": "usuario@spotly.com",
    "password": "usuario123",
    "firstName": "Usuario",
    "lastName": "Normal"
  }')

echo "Respuesta: $USER_RESPONSE"

# 3. Actualizar el usuario admin a tipo "admin" usando SQL directo
echo ""
echo "🔧 Actualizando usuario admin a tipo 'admin'..."
docker-compose exec -T mysql mysql -u spotly -pspotlypass spotly_users <<EOF
UPDATE users SET user_type = 'admin' WHERE username = 'admin';
SELECT id, username, email, user_type FROM users WHERE username IN ('admin', 'usuario');
EOF

echo ""
echo "✅ Usuarios creados exitosamente!"
echo ""
echo "📋 Credenciales:"
echo "   ADMIN:"
echo "     Username: admin"
echo "     Password: admin123"
echo "     Email: admin@spotly.com"
echo ""
echo "   USUARIO NORMAL:"
echo "     Username: usuario"
echo "     Password: usuario123"
echo "     Email: usuario@spotly.com"
echo ""

