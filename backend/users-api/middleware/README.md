# 🔐 Middleware de Autorización - Users API

Este documento describe los helpers de autorización disponibles en el servicio users-api.

---

## 📋 Helpers Disponibles

### 1. **RequireAuth()** - Autenticación Requerida

Valida que el request tenga un token JWT válido.

**Comportamiento:**
- ✅ Si el token es válido: permite continuar y guarda la info del usuario en el contexto
- ❌ Si no hay token: responde **401 Unauthorized**
- ❌ Si el token es inválido o expirado: responde **401 Unauthorized**

**Información guardada en el contexto:**
- `user_id`: ID del usuario
- `username`: Username del usuario
- `role`: Role del usuario (USER, ADMIN, OWNER)

**Uso:**
```go
protected := router.Group("/api")
protected.Use(middleware.RequireAuth())
{
    protected.GET("/profile", userController.GetProfile)
    protected.PUT("/profile", userController.UpdateProfile)
}
```

---

### 2. **RequireAdmin()** - Requiere Token Válido Y Rol ADMIN

Valida que el request tenga un token JWT válido **Y** que el usuario tenga rol `ADMIN`.

**Comportamiento:**
- ✅ Si el token es válido Y el role es ADMIN: permite continuar
- ❌ Si no hay token o es inválido: responde **401 Unauthorized**
- ❌ Si hay token válido pero el role NO es ADMIN: responde **403 Forbidden**

**Información guardada en el contexto:**
- `user_id`: ID del usuario
- `username`: Username del usuario
- `role`: Role del usuario (siempre será "ADMIN" si pasa la validación)

**Uso:**
```go
admin := router.Group("/admin")
admin.Use(middleware.RequireAdmin())
{
    admin.GET("/users", userController.GetAllUsers)
    admin.PUT("/users/:id", userController.UpdateUser)
    admin.DELETE("/users/:id", userController.DeleteUser)
}
```

---

## 🔄 Compatibilidad

Los middlewares antiguos siguen funcionando pero están marcados como deprecated:

- `AuthMiddleware()` → Usar `RequireAuth()` en su lugar
- `AdminMiddleware()` → Usar `RequireAdmin()` en su lugar

---

## 📝 Ejemplos de Uso

### Ejemplo 1: Ruta Protegida Simple

```go
// Requiere autenticación (cualquier usuario autenticado)
protected := router.Group("/api/properties")
protected.Use(middleware.RequireAuth())
{
    protected.POST("/", propertyController.CreateProperty)
    protected.PUT("/:id", propertyController.UpdateProperty)
}
```

### Ejemplo 2: Ruta Solo para Administradores

```go
// Requiere autenticación Y rol ADMIN
admin := router.Group("/admin")
admin.Use(middleware.RequireAdmin())
{
    admin.GET("/users", userController.GetAllUsers)
    admin.PUT("/users/:id", userController.UpdateUser)
}
```

### Ejemplo 3: Verificar Role en un Endpoint

```go
func (ctrl *UserController) SomeEndpoint(c *gin.Context) {
    role, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "role not found"})
        return
    }
    
    if role == "ADMIN" {
        // Lógica solo para admin
    } else if role == "OWNER" {
        // Lógica solo para owner
    } else {
        // Lógica para USER
    }
}
```

---

## ⚠️ Notas Importantes

1. **Orden de Middlewares:** `RequireAdmin()` ya incluye la validación del token, no necesitas usar `RequireAuth()` antes.

2. **Códigos de Estado:**
   - **401 Unauthorized**: No hay token o el token es inválido
   - **403 Forbidden**: Hay token válido pero no tiene permisos suficientes

3. **Contexto:** Después de pasar por `RequireAuth()` o `RequireAdmin()`, puedes acceder a:
   - `c.Get("user_id")` - ID del usuario
   - `c.Get("username")` - Username del usuario
   - `c.Get("role")` - Role del usuario (USER, ADMIN, OWNER)

---

## 🧪 Testing

Para probar los middlewares:

```bash
# Test con token válido
curl -H "Authorization: Bearer <token>" http://localhost:8081/api/properties

# Test sin token (debe retornar 401)
curl http://localhost:8081/api/properties

# Test con token de usuario normal en ruta admin (debe retornar 403)
curl -H "Authorization: Bearer <token_user>" http://localhost:8081/admin/users
```

