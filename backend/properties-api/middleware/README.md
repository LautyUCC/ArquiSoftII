# 🔐 Middleware de Autorización - Properties API

Este documento describe los helpers de autorización disponibles en el servicio properties-api.

---

## 📋 Helpers Disponibles

### 1. **RequireAuth(jwtSecret)** - Autenticación Requerida

Valida que el request tenga un token JWT válido.

**Parámetros:**
- `jwtSecret`: String con el secreto JWT para validar el token

**Comportamiento:**
- ✅ Si el token es válido: permite continuar y guarda la info del usuario en el contexto
- ❌ Si no hay token: responde **401 Unauthorized**
- ❌ Si el token es inválido o expirado: responde **401 Unauthorized**

**Información guardada en el contexto:**
- `userID` / `user_id`: ID del usuario
- `username`: Username del usuario
- `role`: Role del usuario (USER, ADMIN, OWNER)

**Uso:**
```go
jwtSecret := getEnv("JWT_SECRET", "default-secret")

protected := router.Group("/api")
protected.Use(middleware.RequireAuth(jwtSecret))
{
    protected.POST("/properties", propertyController.CreateProperty)
    protected.PUT("/properties/:id", propertyController.UpdateProperty)
}
```

---

### 2. **RequireAdmin(jwtSecret)** - Requiere Token Válido Y Rol ADMIN

Valida que el request tenga un token JWT válido **Y** que el usuario tenga rol `ADMIN`.

**Parámetros:**
- `jwtSecret`: String con el secreto JWT para validar el token

**Comportamiento:**
- ✅ Si el token es válido Y el role es ADMIN: permite continuar
- ❌ Si no hay token o es inválido: responde **401 Unauthorized**
- ❌ Si hay token válido pero el role NO es ADMIN: responde **403 Forbidden**

**Información guardada en el contexto:**
- `userID` / `user_id`: ID del usuario
- `username`: Username del usuario
- `role`: Role del usuario (siempre será "ADMIN" si pasa la validación)

**Uso:**
```go
jwtSecret := getEnv("JWT_SECRET", "default-secret")

admin := router.Group("/api/admin")
admin.Use(middleware.RequireAdmin(jwtSecret))
{
    admin.GET("/properties", propertyController.GetAllProperties)
}
```

---

## 🔄 Compatibilidad

Los middlewares antiguos siguen funcionando pero están marcados como deprecated:

- `AuthMiddleware(jwtSecret)` → Usar `RequireAuth(jwtSecret)` en su lugar
- `AdminRequired()` → Usar `RequireAdmin(jwtSecret)` en su lugar

**Nota:** `AdminRequired()` sin parámetros todavía funciona pero requiere que `RequireAuth()` se haya ejecutado antes. Se recomienda usar `RequireAdmin()` que incluye ambas validaciones.

---

## 📝 Ejemplos de Uso

### Ejemplo 1: Ruta Protegida Simple

```go
jwtSecret := getEnv("JWT_SECRET", "default-secret")

// Requiere autenticación (cualquier usuario autenticado)
protected := router.Group("/api/properties")
protected.Use(middleware.RequireAuth(jwtSecret))
{
    protected.POST("/", propertyController.CreateProperty)
    protected.PUT("/:id", propertyController.UpdateProperty)
    protected.DELETE("/:id", propertyController.DeleteProperty)
}
```

### Ejemplo 2: Ruta Solo para Administradores

```go
jwtSecret := getEnv("JWT_SECRET", "default-secret")

// Requiere autenticación Y rol ADMIN
admin := router.Group("/api/admin")
admin.Use(middleware.RequireAdmin(jwtSecret))
{
    admin.GET("/properties", propertyController.GetAllProperties)
}
```

### Ejemplo 3: Verificar Role en un Endpoint

```go
func (ctrl *PropertyController) UpdateProperty(c *gin.Context) {
    userID, _ := c.Get("user_id")
    role, _ := c.Get("role")
    
    // Obtener propiedad
    property := getProperty(id)
    
    // Solo OWNER de la propiedad o ADMIN pueden actualizar
    if property.OwnerID != userID && role != "ADMIN" {
        c.JSON(http.StatusForbidden, gin.H{"error": "owner or admin privileges required"})
        return
    }
    
    // ... lógica de actualización
}
```

---

## ⚠️ Notas Importantes

1. **JWT Secret:** Debe ser el mismo secreto usado en `users-api` para generar los tokens.

2. **Orden de Middlewares:** `RequireAdmin()` ya incluye la validación del token, no necesitas usar `RequireAuth()` antes.

3. **Códigos de Estado:**
   - **401 Unauthorized**: No hay token o el token es inválido
   - **403 Forbidden**: Hay token válido pero no tiene permisos suficientes

4. **Contexto:** Después de pasar por `RequireAuth()` o `RequireAdmin()`, puedes acceder a:
   - `c.Get("userID")` o `c.Get("user_id")` - ID del usuario
   - `c.Get("username")` - Username del usuario
   - `c.Get("role")` - Role del usuario (USER, ADMIN, OWNER)

5. **Actualización de Roles:** El middleware ahora usa `role` (USER, ADMIN, OWNER) en lugar de `user_type` (normal, admin). Los tokens antiguos con `user_type` pueden no funcionar correctamente.

---

## 🧪 Testing

Para probar los middlewares:

```bash
# Test con token válido
curl -H "Authorization: Bearer <token>" http://localhost:8082/api/properties

# Test sin token (debe retornar 401)
curl http://localhost:8082/api/properties

# Test con token de usuario normal en ruta admin (debe retornar 403)
curl -H "Authorization: Bearer <token_user>" http://localhost:8082/api/admin/properties
```

---

## 🔄 Migración desde Middlewares Antiguos

Si tienes código que usa los middlewares antiguos:

**Antes:**
```go
protected.Use(middleware.AuthMiddleware(jwtSecret))
admin.Use(middleware.AuthMiddleware(jwtSecret))
admin.Use(middleware.AdminRequired())
```

**Después:**
```go
protected.Use(middleware.RequireAuth(jwtSecret))
admin.Use(middleware.RequireAdmin(jwtSecret))
```

