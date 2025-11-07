# Checklist de Verificación - Properties API

Este documento contiene una lista de verificación para validar que todos los componentes de la API funcionan correctamente.

---

## ✅ 1. Compilación

Verificar que el código compila sin errores.

### Comando

```bash
cd backend/properties-api
go build -o main .
```

### Resultado Esperado

- ✅ Sin errores de compilación
- ✅ Binario `main` generado exitosamente

### Si falla

- Verificar que todas las dependencias estén instaladas: `go mod tidy`
- Revisar errores de sintaxis en el código
- Verificar imports faltantes

---

## ✅ 2. Tests

Verificar que todos los tests unitarios pasan.

### Comando

```bash
cd backend/properties-api
go test ./... -v -cover
```

### Resultado Esperado

- ✅ Todos los tests pasan
- ✅ Cobertura de código reportada
- ✅ Sin errores ni fallos

### Tests Incluidos

- `TestCreateProperty_Success`
- `TestCreateProperty_UserNotFound`
- `TestGetPropertyByID_Success`
- `TestGetPropertyByID_NotFound`
- `TestUpdateProperty_Unauthorized`
- `TestDeleteProperty_Success`
- `TestDeleteProperty_Unauthorized`

### Si falla

- Revisar mensajes de error en los tests
- Verificar que los mocks estén correctamente configurados
- Verificar que las dependencias estén disponibles

---

## ✅ 3. Docker

Verificar que docker-compose levanta todos los servicios correctamente.

### Comando

```bash
cd /ruta/al/proyecto/raiz
docker-compose up --build -d
```

### Resultado Esperado

- ✅ Todos los contenedores se construyen exitosamente
- ✅ Todos los contenedores están en estado "Up"
- ✅ Sin errores en los logs

### Verificar Servicios

```bash
docker-compose ps
```

**Servicios esperados:**
- `properties-api` - Estado: Up
- `mongodb` - Estado: Up
- `rabbitmq` - Estado: Up
- `users-api` - Estado: Up

### Si falla

- Verificar que Docker esté ejecutándose
- Revisar logs: `docker-compose logs`
- Verificar que los puertos no estén en uso
- Verificar que el Dockerfile esté correcto

---

## ✅ 4. Conectividad

Verificar que MongoDB y RabbitMQ son accesibles.

### MongoDB

#### Comando

```bash
# Verificar conexión a MongoDB
docker-compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

#### Resultado Esperado

```json
{ ok: 1 }
```

#### Test de Conectividad

```bash
# Verificar que MongoDB responde
curl -s http://localhost:27017 || echo "MongoDB no accesible directamente (normal)"
```

### RabbitMQ

#### Comando

```bash
# Verificar estado de RabbitMQ
docker-compose exec rabbitmq rabbitmqctl status
```

#### Resultado Esperado

- ✅ Status de RabbitMQ mostrado sin errores
- ✅ Nodo RabbitMQ funcionando

#### Test de Conectividad

```bash
# Verificar panel web de RabbitMQ
curl -s -u guest:guest http://localhost:15672/api/overview | head -20
```

#### Resultado Esperado

- ✅ JSON con información de RabbitMQ
- ✅ Sin errores de autenticación

### Si falla

- Verificar que los contenedores estén ejecutándose: `docker-compose ps`
- Revisar logs: `docker-compose logs mongodb` o `docker-compose logs rabbitmq`
- Verificar que los puertos no estén bloqueados por firewall

---

## ✅ 5. Endpoints - POST /properties

Verificar que el endpoint de creación funciona correctamente.

### Comando

```bash
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad de prueba",
    "description": "Descripción de prueba para verificar el endpoint",
    "price": 100000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool"],
    "capacity": 4,
    "available": true
  }' | jq .
```

### Resultado Esperado

**Status Code:** `201 Created`

**Response Body:**

```json
{
  "success": true,
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "title": "Propiedad de prueba",
    "description": "Descripción de prueba para verificar el endpoint",
    "price": 121270.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool"],
    "capacity": 4,
    "available": true,
    "createdAt": "2024-01-15T10:30:00Z",
    "updatedAt": "2024-01-15T10:30:00Z"
  },
  "message": "Property created successfully"
}
```

### Verificaciones

- ✅ Status code es 201
- ✅ Campo `success` es `true`
- ✅ Campo `data` contiene la propiedad creada
- ✅ `id` está presente y no está vacío
- ✅ `price` fue calculado con concurrencia (mayor al precio base)
- ✅ `createdAt` y `updatedAt` están presentes

### Guardar ID para pruebas posteriores

```bash
# Guardar el ID de la propiedad creada
PROPERTY_ID=$(curl -s -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad de prueba",
    "description": "Descripción de prueba",
    "price": 100000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool"],
    "capacity": 4,
    "available": true
  }' | jq -r '.data.id')

echo "Property ID: $PROPERTY_ID"
```

### Si falla

- Verificar que la API esté ejecutándose: `curl http://localhost:8081/health`
- Revisar logs: `docker-compose logs properties-api`
- Verificar que users-api esté disponible (para validar el ownerId)
- Verificar formato JSON del request

---

## ✅ 6. Validación - POST con Owner Inválido

Verificar que la validación de usuario funciona correctamente.

### Comando

```bash
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad con owner inválido",
    "description": "Esta propiedad no debería crearse",
    "price": 100000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "nonexistent-user-999",
    "amenities": ["wifi"],
    "capacity": 2,
    "available": true
  }' | jq .
```

### Resultado Esperado

**Status Code:** `404 Not Found`

**Response Body:**

```json
{
  "error": "User not found",
  "message": "usuario owner con ID 'nonexistent-user-999' no existe"
}
```

### Verificaciones

- ✅ Status code es 404
- ✅ Campo `error` contiene "User not found"
- ✅ Mensaje indica que el usuario no existe
- ✅ La propiedad NO fue creada en la base de datos

### Test Adicional - Campo Requerido Faltante

```bash
# Test con campo requerido faltante
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad sin precio",
    "description": "Esta propiedad no debería crearse",
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "capacity": 2
  }' | jq .
```

**Resultado Esperado:** `400 Bad Request` con mensaje de validación

### Si falla

- Verificar que users-api esté ejecutándose
- Revisar logs de properties-api para ver el error específico
- Verificar que la validación se esté ejecutando correctamente

---

## ✅ 7. Eventos - RabbitMQ

Verificar que RabbitMQ recibe mensajes cuando se crean/actualizan/eliminan propiedades.

### Verificar Cola en RabbitMQ

#### Comando

```bash
# Verificar que la cola existe
curl -s -u guest:guest http://localhost:15672/api/queues | jq '.[] | select(.name == "property_events")'
```

### Crear Propiedad y Verificar Evento

#### Paso 1: Crear Propiedad

```bash
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad para test de eventos",
    "description": "Verificar que se publica evento en RabbitMQ",
    "price": 150000.00,
    "location": "Medellín, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool"],
    "capacity": 3,
    "available": true
  }' | jq -r '.data.id'
```

#### Paso 2: Verificar Mensajes en la Cola

```bash
# Verificar mensajes en la cola property_events
curl -s -u guest:guest http://localhost:15672/api/queues/%2F/property_events | jq '.messages'
```

#### Resultado Esperado

- ✅ La cola `property_events` existe
- ✅ El número de mensajes es mayor a 0
- ✅ Mensaje contiene `operation: "create"` y `propertyId`

### Verificar Evento de Actualización

```bash
# Guardar ID de propiedad
PROPERTY_ID="507f1f77bcf86cd799439011"

# Actualizar propiedad (requiere user_id en contexto, usar mock si es necesario)
curl -X PUT http://localhost:8081/properties/$PROPERTY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{
    "title": "Título actualizado"
  }' | jq .
```

Luego verificar que el mensaje de "update" fue publicado.

### Si falla

- Verificar que RabbitMQ esté ejecutándose: `docker-compose ps rabbitmq`
- Revisar logs de RabbitMQ: `docker-compose logs rabbitmq`
- Verificar que la conexión a RabbitMQ esté configurada correctamente
- Revisar logs de properties-api para errores de publicación

---

## ✅ 8. Concurrencia - Cálculo de Precio

Verificar que el precio se calcula correctamente usando goroutines.

### Comando

```bash
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad para test de precio",
    "description": "Verificar cálculo de precio con concurrencia",
    "price": 100000.00,
    "location": "Cali, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool", "parking"],
    "capacity": 4,
    "available": true
  }' | jq '.data.price'
```

### Cálculo Esperado

**Precio base:** $100,000.00

**Desglose:**
- Precio base con impuestos (21%): $100,000 × 1.21 = **$121,000**
- Costo por amenidades (3 amenidades): 3 × $50 = **$150**
- Costo por capacidad (4 personas): 4 × $30 = **$120**

**Precio total esperado:** $121,000 + $150 + $120 = **$121,270**

### Resultado Esperado

```json
121270.00
```

### Verificaciones

- ✅ El precio retornado es exactamente $121,270.00
- ✅ El precio es mayor al precio base (incluye impuestos y extras)
- ✅ El cálculo se realiza correctamente usando goroutines

### Test con Diferentes Valores

```bash
# Test con precio base 200000, 5 amenidades, capacidad 6
curl -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad grande",
    "description": "Test de cálculo con valores mayores",
    "price": 200000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi", "pool", "parking", "kitchen", "tv"],
    "capacity": 6,
    "available": true
  }' | jq '.data.price'
```

**Cálculo esperado:**
- Precio base con impuestos: $200,000 × 1.21 = $242,000
- Amenidades (5): 5 × $50 = $250
- Capacidad (6): 6 × $30 = $180
- **Total: $242,430**

### Si falla

- Verificar que la función `CalculatePriceWithConcurrency` esté implementada correctamente
- Revisar logs de la API para errores en el cálculo
- Verificar que las goroutines se estén ejecutando correctamente

---

## ✅ 9. Permisos - Update/Delete Validan Ownership

Verificar que solo el propietario puede actualizar o eliminar su propiedad.

### Test Update con Usuario No Autorizado

#### Paso 1: Crear Propiedad con Owner "user123"

```bash
PROPERTY_ID=$(curl -s -X POST http://localhost:8081/properties \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propiedad del usuario 123",
    "description": "Esta propiedad pertenece a user123",
    "price": 100000.00,
    "location": "Bogotá, Colombia",
    "ownerId": "user123",
    "amenities": ["wifi"],
    "capacity": 2,
    "available": true
  }' | jq -r '.data.id')

echo "Property ID: $PROPERTY_ID"
```

#### Paso 2: Intentar Actualizar con Usuario Diferente

```bash
# Intentar actualizar con user456 (no es el owner)
curl -X PUT http://localhost:8081/properties/$PROPERTY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "X-User-ID: user456" \
  -d '{
    "title": "Intento de actualización no autorizada"
  }' | jq .
```

**Nota:** Si el middleware obtiene user_id del contexto, ajustar según la implementación real.

#### Resultado Esperado

**Status Code:** `403 Forbidden`

**Response Body:**

```json
{
  "error": "Forbidden",
  "message": "usuario con ID 'user456' no tiene permisos para actualizar propiedad '...' (owner: 'user123')"
}
```

### Test Delete con Usuario No Autorizado

```bash
# Intentar eliminar con usuario diferente al owner
curl -X DELETE http://localhost:8081/properties/$PROPERTY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "X-User-ID: user456" \
  | jq .
```

#### Resultado Esperado

**Status Code:** `403 Forbidden`

### Test Update con Owner Correcto

```bash
# Actualizar con el owner correcto (user123)
curl -X PUT http://localhost:8081/properties/$PROPERTY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "X-User-ID: user123" \
  -d '{
    "title": "Título actualizado por el owner"
  }' | jq .
```

#### Resultado Esperado

**Status Code:** `200 OK`

**Response Body:**

```json
{
  "success": true,
  "message": "Property updated successfully"
}
```

### Test Delete con Owner Correcto

```bash
# Eliminar con el owner correcto
curl -X DELETE http://localhost:8081/properties/$PROPERTY_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -H "X-User-ID: user123" \
  | jq .
```

#### Resultado Esperado

**Status Code:** `200 OK`

**Response Body:**

```json
{
  "success": true,
  "message": "Property deleted successfully"
}
```

### Verificaciones

- ✅ Update con usuario no autorizado retorna 403
- ✅ Delete con usuario no autorizado retorna 403
- ✅ Update con owner correcto retorna 200
- ✅ Delete con owner correcto retorna 200
- ✅ Mensajes de error son descriptivos

### Si falla

- Verificar que el middleware de autenticación esté configurado correctamente
- Revisar cómo se obtiene el `user_id` del contexto
- Verificar que la validación de ownership se ejecute en el servicio
- Revisar logs de la API para errores específicos

---

## 📋 Resumen de Verificación

### Checklist Rápido

```bash
# 1. Compilación
go build -o main . && echo "✅ Compilación OK"

# 2. Tests
go test ./... -v && echo "✅ Tests OK"

# 3. Docker
docker-compose ps | grep -q "Up" && echo "✅ Docker OK"

# 4. Health Check
curl -s http://localhost:8081/health | grep -q "ok" && echo "✅ API OK"

# 5. MongoDB
docker-compose exec -T mongodb mongosh --eval "db.adminCommand('ping')" | grep -q "ok.*1" && echo "✅ MongoDB OK"

# 6. RabbitMQ
docker-compose exec rabbitmq rabbitmqctl status > /dev/null 2>&1 && echo "✅ RabbitMQ OK"
```

### Comandos de Verificación Rápida

```bash
# Verificar todos los servicios
./scripts/start.sh

# Ejecutar tests
./scripts/test.sh

# Crear datos de prueba
./scripts/seed.sh

# Ver logs
docker-compose logs -f properties-api
```

---

## 🐛 Troubleshooting

### Problemas Comunes

1. **Error de conexión a MongoDB**
   - Verificar que el contenedor esté ejecutándose
   - Verificar la URI de conexión en el código

2. **Error de conexión a RabbitMQ**
   - Verificar que el contenedor esté ejecutándose
   - Verificar credenciales y URL

3. **Error de validación de usuario**
   - Verificar que users-api esté ejecutándose
   - Verificar la URL de users-api en la configuración

4. **Tests fallan**
   - Ejecutar `go mod tidy`
   - Verificar que los mocks estén correctamente implementados

---

## 📝 Notas

- Todos los comandos curl asumen que la API está en `http://localhost:8081`
- Los tokens de autenticación deben ser proporcionados según la implementación real
- El `user_id` se obtiene del contexto; ajustar los headers según la implementación del middleware
- Los IDs de propiedades son ObjectIDs de MongoDB (24 caracteres hexadecimales)

