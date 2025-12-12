#!/bin/bash

# Script para crear una propiedad de prueba
# Requiere que el usuario admin exista y que los servicios estén corriendo

set -e

echo "========================================"
echo "Creando Propiedad de Prueba"
echo "========================================"
echo ""

# Verificar que los servicios estén corriendo
echo "[1/4] Verificando que la API esté disponible..."
if ! curl -s -f http://localhost/api/users/login > /dev/null; then
  echo "ERROR: La API no está disponible en http://localhost/api"
  echo "Asegúrate de que los servicios estén ejecutándose con docker-compose up"
  exit 1
fi
echo "OK: API disponible"
echo ""

# Hacer login como admin
echo "[2/4] Haciendo login como admin..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"usernameOrEmail":"admin","password":"admin123"}')

# Extraer el token del JSON (usando jq si está disponible, o grep/sed como fallback)
if command -v jq &> /dev/null; then
  TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
else
  TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "ERROR: No se pudo obtener el token. Verificando respuesta:"
  echo "$LOGIN_RESPONSE"
  exit 1
fi

echo "OK: Token obtenido"
echo ""

# Crear propiedad de prueba
echo "[3/4] Creando propiedad de prueba..."
PROPERTY_RESPONSE=$(curl -s -X POST http://localhost/api/properties \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Apartamento Moderno en el Centro",
    "description": "Hermoso apartamento completamente amueblado en el corazón de la ciudad. Cerca de restaurantes, transporte público y centros comerciales. Ideal para turistas o profesionales.",
    "price": 120000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "1",
    "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning"],
    "capacity": 4,
    "available": true,
    "images": []
  }')

# Verificar si se creó correctamente
if echo "$PROPERTY_RESPONSE" | grep -q '"id"'; then
  echo "OK: Propiedad creada exitosamente"
  echo ""
  echo "[4/4] Respuesta del servidor:"
  echo "$PROPERTY_RESPONSE" | jq '.' 2>/dev/null || echo "$PROPERTY_RESPONSE"
else
  echo "ERROR: No se pudo crear la propiedad. Respuesta:"
  echo "$PROPERTY_RESPONSE"
  exit 1
fi

echo ""
echo "========================================"
echo "Propiedad de prueba creada exitosamente"
echo "========================================"
echo ""
echo "Puedes verificar la propiedad en el frontend:"
echo "http://localhost"
echo ""

