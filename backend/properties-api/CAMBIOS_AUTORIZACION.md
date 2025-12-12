# 🔐 Cambios de Autorización - Properties API

Este documento resume los cambios realizados para restringir las operaciones de creación, actualización y eliminación de propiedades solo a usuarios con rol ADMIN.

---

## ✅ Cambios Implementados

### 1. **Rutas Actualizadas** (`main.go`)

**Antes:**
```go
// Rutas protegidas (requieren autenticación - cualquier rol)
protected := router.Group("/api")
protected.Use(middleware.RequireAuth(jwtSecret))
{
    protected.POST("/properties", propertyController.CreateProperty)
    protected.PUT("/properties/:id", propertyController.UpdateProperty)
    protected.DELETE("/properties/:id", propertyController.DeleteProperty)
}
```

**Después:**
```go
// Rutas de administrador (solo ADMIN puede crear/editar/eliminar)
admin := router.Group("/api")
admin.Use(middleware.RequireAdmin(jwtSecret)) // Valida token Y rol ADMIN
{
    admin.POST("/properties", propertyController.CreateProperty)
    admin.PUT("/properties/:id", propertyController.UpdateProperty)
    admin.DELETE("/properties/:id", propertyController.DeleteProperty)
}
```

**Endpoints de Lectura (sin cambios):**
```go
// Rutas públicas (lectura - cualquier usuario puede ver)
public := router.Group("/api")
{
    public.GET("/properties/:id", propertyController.GetPropertyByID)
    public.GET("/properties/user/:userId", propertyController.GetUserProperties)
}
```

---

### 2. **Controladores Actualizados** (`controllers/properties_controller.go`)

#### CreateProperty
- **Antes:** Cualquier usuario autenticado podía crear propiedades
- **Después:** Solo ADMIN puede crear propiedades (validado por middleware `RequireAdmin`)

#### UpdateProperty
- **Antes:** Cualquier usuario autenticado podía actualizar propiedades (con validación de owner)
- **Después:** Solo ADMIN puede actualizar propiedades
- **Validación de owner:** Se mantiene en el servicio para consistencia de datos, pero ADMIN puede editar cualquier propiedad sin importar el owner

**Cambios técnicos:**
- Actualizado para usar `role` del contexto en lugar de `isAdmin`
- Mejorada la conversión de `userID` de `uint` a `string` usando `fmt.Sprintf`

#### DeleteProperty
- **Antes:** Cualquier usuario autenticado podía eliminar propiedades (con validación de owner)
- **Después:** Solo ADMIN puede eliminar propiedades
- **Validación de owner:** Se mantiene en el servicio para consistencia de datos, pero ADMIN puede eliminar cualquier propiedad sin importar el owner

**Cambios técnicos:**
- Actualizado para usar `role` del contexto en lugar de `isAdmin`
- Mejorada la conversión de `userID` de `uint` a `string` usando `fmt.Sprintf`

---

### 3. **Servicio Actualizado** (`services/properties_service.go`)

Los métodos del servicio mantienen la validación de owner para consistencia de datos, pero ahora solo reciben requests de usuarios con rol ADMIN (validado en el middleware).

**Comentarios actualizados:**
- `CreateProperty`: Documentado que solo ADMIN puede crear propiedades
- `UpdateProperty`: Documentado que solo ADMIN puede actualizar, pero la validación de owner se mantiene
- `DeleteProperty`: Documentado que solo ADMIN puede eliminar, pero la validación de owner se mantiene

---

## 📋 Matriz de Permisos

| Endpoint | Método | Antes | Después |
|----------|--------|-------|---------|
| `/api/properties` | POST | Cualquier autenticado | **Solo ADMIN** |
| `/api/properties/:id` | PUT | Cualquier autenticado | **Solo ADMIN** |
| `/api/properties/:id` | DELETE | Cualquier autenticado | **Solo ADMIN** |
| `/api/properties/:id` | GET | Público | Público (sin cambios) |
| `/api/properties/user/:userId` | GET | Público | Público (sin cambios) |
| `/api/admin/properties` | GET | Solo ADMIN | Solo ADMIN (sin cambios) |

---

## 🔐 Validación de Owner

La validación de owner se mantiene en el servicio para:
1. **Consistencia de datos:** Asegurar que solo el owner o ADMIN pueden modificar propiedades
2. **Futuras extensiones:** Si en el futuro se permite que OWNER también pueda editar sus propiedades, la validación ya está implementada
3. **Logging y auditoría:** Registrar quién modificó qué propiedad

**Comportamiento actual:**
- ADMIN puede crear/editar/eliminar cualquier propiedad sin importar el owner
- La validación de owner se ejecuta pero siempre pasa porque `isAdmin` es `true` cuando llega al servicio

---

## 🧪 Testing

### Test 1: Crear Propiedad sin Token (debe retornar 401)

```bash
curl -X POST http://localhost:8082/api/properties \
  -H "Content-Type: application/json" \
  -d '{"title": "Test", "ownerID": "1", "price": 100}'
# Response: {"error": "authorization header required"}
```

### Test 2: Crear Propiedad con Token de USER (debe retornar 403)

```bash
curl -X POST http://localhost:8082/api/properties \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{"title": "Test", "ownerID": "1", "price": 100}'
# Response: {"error": "admin privileges required"}
```

### Test 3: Crear Propiedad con Token de ADMIN (debe retornar 201)

```bash
curl -X POST http://localhost:8082/api/properties \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{"title": "Test", "ownerID": "1", "price": 100}'
# Response: 201 Created con la propiedad creada
```

### Test 4: Leer Propiedad sin Token (debe funcionar - público)

```bash
curl http://localhost:8082/api/properties/123
# Response: 200 OK con los datos de la propiedad
```

### Test 5: Actualizar Propiedad con Token de USER (debe retornar 403)

```bash
curl -X PUT http://localhost:8082/api/properties/123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_user>" \
  -d '{"title": "Updated"}'
# Response: {"error": "admin privileges required"}
```

### Test 6: Actualizar Propiedad con Token de ADMIN (debe retornar 200)

```bash
curl -X PUT http://localhost:8082/api/properties/123 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token_admin>" \
  -d '{"title": "Updated"}'
# Response: {"message": "Propiedad actualizada exitosamente"}
```

---

## ⚠️ Notas Importantes

1. **Middleware RequireAdmin:** Las rutas POST, PUT y DELETE ahora usan `RequireAdmin()` que valida:
   - Token JWT válido (401 si no hay token o es inválido)
   - Rol ADMIN en el claim del JWT (403 si el token es válido pero no es ADMIN)

2. **Endpoints de Lectura:** Siguen siendo públicos y no requieren autenticación. Cualquier usuario puede ver propiedades.

3. **Validación de Owner:** Se mantiene en el servicio pero siempre pasa porque solo ADMIN puede llegar a esos métodos.

4. **Conversión de userID:** Mejorada para manejar correctamente `uint`, `int` y `string` usando `fmt.Sprintf`.

5. **Compatibilidad:** El código es compatible con el nuevo sistema de roles (USER, ADMIN, OWNER) implementado en users-api.

---

## 🔄 Migración

Si tienes código que usa estos endpoints:

**Antes:**
- Cualquier usuario autenticado podía crear/editar/eliminar propiedades
- Solo el owner o admin podía editar/eliminar propiedades específicas

**Después:**
- Solo ADMIN puede crear/editar/eliminar propiedades
- Los endpoints de lectura siguen siendo públicos
- La validación de owner se mantiene pero siempre pasa porque solo ADMIN puede hacer estas operaciones

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el servicio:**
   ```bash
   cd backend/properties-api && go build
   ```

2. **Probar con diferentes roles:**
   - Crear un usuario con rol USER y intentar crear una propiedad → debe retornar 403
   - Crear un usuario con rol ADMIN y crear una propiedad → debe retornar 201
   - Leer una propiedad sin token → debe funcionar (público)

3. **Verificar códigos de estado:**
   - Sin token → 401
   - Token válido de USER → 403
   - Token válido de ADMIN → 200/201

---

## 📊 Resumen de Archivos Modificados

1. `backend/properties-api/main.go` - Rutas actualizadas para usar `RequireAdmin`
2. `backend/properties-api/controllers/properties_controller.go` - Controladores actualizados para usar `role` del contexto
3. `backend/properties-api/services/properties_service.go` - Comentarios actualizados
4. `backend/properties-api/CAMBIOS_AUTORIZACION.md` - Este documento

