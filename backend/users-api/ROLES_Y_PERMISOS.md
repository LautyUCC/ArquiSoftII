# 🔐 Roles y Permisos - Users API

Este documento describe los roles disponibles en el sistema y los permisos asociados a cada uno.

---

## 📋 Roles Disponibles

### 1. **USER** (Rol por defecto)
**Valor en BD:** `USER`

**Permisos:**
- ✅ Leer propiedades (GET /api/properties/:id)
- ✅ Buscar propiedades (GET /api/search)
- ✅ Crear reservas/bookings (POST /api/bookings)
- ✅ Ver sus propias reservas (GET /api/bookings/user/:userId)
- ✅ Ver sus propias propiedades (GET /api/properties/user/:userId)

**Restricciones:**
- ❌ No puede crear propiedades
- ❌ No puede editar/eliminar propiedades de otros
- ❌ No puede cambiar roles de usuarios
- ❌ No puede acceder a endpoints de administración

---

### 2. **ADMIN** (Administrador)
**Valor en BD:** `ADMIN`

**Permisos:**
- ✅ **Todos los permisos de USER**
- ✅ Crear propiedades (POST /api/properties)
- ✅ Editar cualquier propiedad (PUT /api/properties/:id)
- ✅ Eliminar cualquier propiedad (DELETE /api/properties/:id)
- ✅ Ver todas las propiedades (GET /api/admin/properties)
- ✅ Ver todos los usuarios (GET /admin/users)
- ✅ Actualizar usuarios (PUT /admin/users/:id) - **incluyendo cambiar roles**
- ✅ Eliminar usuarios (DELETE /admin/users/:id)

**Restricciones:**
- ❌ No puede eliminar su propio usuario (si se implementa esta validación)

---

### 3. **OWNER** (Propietario)
**Valor en BD:** `OWNER`

**Permisos:**
- ✅ **Todos los permisos de USER**
- ✅ Crear propiedades (POST /api/properties)
- ✅ Editar sus propias propiedades (PUT /api/properties/:id) - validación de ownership
- ✅ Eliminar sus propias propiedades (DELETE /api/properties/:id) - validación de ownership
- ✅ Ver sus propias propiedades (GET /api/properties/user/:userId)

**Restricciones:**
- ❌ No puede editar/eliminar propiedades de otros
- ❌ No puede cambiar roles de usuarios
- ❌ No puede acceder a endpoints de administración

---

## 🔑 Inclusión en JWT

El **role** se incluye como claim en el JWT token:

```json
{
  "user_id": 1,
  "username": "admin",
  "role": "ADMIN",
  "exp": 1234567890,
  "iat": 1234567890
}
```

**Claim name:** `role`

**Valores posibles:** `USER`, `ADMIN`, `OWNER`

---

## 🛡️ Validación de Roles

### En el Middleware

El middleware `AuthMiddleware` extrae el `role` del JWT y lo guarda en el contexto:

```go
c.Set("role", claims.Role)
```

### Verificar Rol en Endpoints

```go
// Verificar si es ADMIN
role, _ := c.Get("role")
if role != "ADMIN" {
    c.JSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
    c.Abort()
    return
}

// Verificar si es OWNER o ADMIN
role, _ := c.Get("role")
if role != "OWNER" && role != "ADMIN" {
    c.JSON(http.StatusForbidden, gin.H{"error": "owner or admin privileges required"})
    c.Abort()
    return
}
```

---

## 📝 Cambio de Roles

**Solo ADMIN puede cambiar roles de usuarios:**

```bash
PUT /admin/users/:id
{
  "role": "ADMIN"  // Solo ADMIN puede hacer esto
}
```

**Validaciones:**
- El role debe ser válido: `USER`, `ADMIN`, o `OWNER`
- Solo usuarios con role `ADMIN` pueden cambiar roles
- Un usuario no puede cambiar su propio role (recomendado)

---

## 🔄 Migración de Datos

Si ya tienes usuarios con `user_type = "normal"` o `user_type = "admin"`, necesitas migrarlos:

```sql
-- Migrar usuarios existentes
UPDATE users SET role = 'USER' WHERE user_type = 'normal' OR role IS NULL;
UPDATE users SET role = 'ADMIN' WHERE user_type = 'admin';
```

---

## 📊 Matriz de Permisos

| Acción | USER | OWNER | ADMIN |
|--------|------|-------|-------|
| Leer propiedades | ✅ | ✅ | ✅ |
| Buscar propiedades | ✅ | ✅ | ✅ |
| Crear reservas | ✅ | ✅ | ✅ |
| Ver propias reservas | ✅ | ✅ | ✅ |
| Crear propiedades | ❌ | ✅ | ✅ |
| Editar propias propiedades | ❌ | ✅ | ✅ |
| Editar cualquier propiedad | ❌ | ❌ | ✅ |
| Eliminar propias propiedades | ❌ | ✅ | ✅ |
| Eliminar cualquier propiedad | ❌ | ❌ | ✅ |
| Ver todas las propiedades | ❌ | ❌ | ✅ |
| Ver todos los usuarios | ❌ | ❌ | ✅ |
| Cambiar roles | ❌ | ❌ | ✅ |
| Eliminar usuarios | ❌ | ❌ | ✅ |

---

## 🔍 Ejemplos de Uso

### Verificar si es ADMIN en un endpoint:

```go
func (ctrl *PropertyController) DeleteProperty(c *gin.Context) {
    role, _ := c.Get("role")
    if role != "ADMIN" {
        c.JSON(http.StatusForbidden, gin.H{"error": "admin privileges required"})
        return
    }
    // ... lógica de eliminación
}
```

### Verificar si es OWNER o ADMIN:

```go
func (ctrl *PropertyController) UpdateProperty(c *gin.Context) {
    role, _ := c.Get("role")
    userID, _ := c.Get("user_id")
    
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

1. **Por defecto:** Todos los usuarios nuevos se crean con role `USER`
2. **JWT Claims:** El role siempre se incluye en el JWT como claim `role`
3. **Validación:** El role se valida en el middleware y está disponible en el contexto
4. **Cambio de Roles:** Solo ADMIN puede cambiar roles de otros usuarios
5. **Backward Compatibility:** Si hay usuarios con `user_type` antiguo, necesitan migración

