#!/bin/bash

# Script para crear datos de prueba en la API
# Crea 3 propiedades de ejemplo usando curl

set -e

echo "🌱 Creando datos de prueba en Properties API..."
echo ""

# URL base de la API
API_URL="http://localhost:8081"

# Verificar que la API esté disponible
echo "🔍 Verificando que la API esté disponible..."
if ! curl -s -f "$API_URL/health" > /dev/null; then
  echo "❌ Error: La API no está disponible en $API_URL"
  echo "   Asegúrate de que los servicios estén ejecutándose con ./scripts/start.sh"
  exit 1
fi
echo "✅ API disponible"
echo ""

# Propiedad 1: Apartamento moderno en el centro
echo "📝 Creando Propiedad 1: Apartamento moderno en el centro..."
PROPERTY_1=$(curl -s -X POST "$API_URL/properties" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Apartamento moderno en el centro",
    "description": "Hermoso apartamento completamente amueblado en el corazón de la ciudad. Cerca de restaurantes, transporte público y centros comerciales. Ideal para turistas o profesionales.",
    "price": 120000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user001",
    "amenities": ["wifi", "pool", "parking", "kitchen", "air-conditioning"],
    "capacity": 3,
    "available": true
  }')

if echo "$PROPERTY_1" | grep -q "success"; then
  PROPERTY_1_ID=$(echo "$PROPERTY_1" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
  echo "   ✅ Propiedad creada con ID: $PROPERTY_1_ID"
else
  echo "   ❌ Error creando propiedad 1:"
  echo "$PROPERTY_1" | head -5
fi
echo ""

# Propiedad 2: Casa con jardín
echo "📝 Creando Propiedad 2: Casa con jardín en las afueras..."
PROPERTY_2=$(curl -s -X POST "$API_URL/properties" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Casa con jardín en las afueras",
    "description": "Casa espaciosa con jardín privado, perfecta para familias. Incluye cocina completa, sala de estar amplia y espacios al aire libre. Ideal para estancias largas.",
    "price": 200000.00,
    "location": "Medellín, Colombia",
    "ownerId": "user002",
    "amenities": ["wifi", "kitchen", "parking", "garden", "tv"],
    "capacity": 6,
    "available": true
  }')

if echo "$PROPERTY_2" | grep -q "success"; then
  PROPERTY_2_ID=$(echo "$PROPERTY_2" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
  echo "   ✅ Propiedad creada con ID: $PROPERTY_2_ID"
else
  echo "   ❌ Error creando propiedad 2:"
  echo "$PROPERTY_2" | head -5
fi
echo ""

# Propiedad 3: Loft industrial
echo "📝 Creando Propiedad 3: Loft industrial en zona comercial..."
PROPERTY_3=$(curl -s -X POST "$API_URL/properties" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Loft industrial en zona comercial",
    "description": "Loft moderno con diseño industrial, ubicado en zona comercial y de entretenimiento. Perfecto para parejas o profesionales jóvenes. Incluye todas las comodidades modernas.",
    "price": 95000.00,
    "location": "Cali, Colombia",
    "ownerId": "user003",
    "amenities": ["wifi", "kitchen", "air-conditioning", "tv", "workspace"],
    "capacity": 2,
    "available": true
  }')

if echo "$PROPERTY_3" | grep -q "success"; then
  PROPERTY_3_ID=$(echo "$PROPERTY_3" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
  echo "   ✅ Propiedad creada con ID: $PROPERTY_3_ID"
else
  echo "   ❌ Error creando propiedad 3:"
  echo "$PROPERTY_3" | head -5
fi
echo ""

echo "✅ Datos de prueba creados exitosamente"
echo ""
echo "📋 Resumen de propiedades creadas:"
echo "   - Propiedad 1: Apartamento moderno en el centro (Bogotá)"
echo "   - Propiedad 2: Casa con jardín en las afueras (Medellín)"
echo "   - Propiedad 3: Loft industrial en zona comercial (Cali)"
echo ""
echo "💡 Puedes verificar las propiedades con:"
echo "   curl http://localhost:8081/properties/<property-id>"

