# 🔄 Cambios Realizados: Normalización de Roles

## 📋 Resumen

Se ha normalizado el campo de roles en el sistema de usuarios, cambiando de `UserType` (con valores "normal", "admin") a `Role` (con valores "USER", "ADMIN", "OWNER"). El role ahora se incluye como claim en el JWT.

---

## ✅ Cambios Implementados

### 1. **Modelo de Usuario** (`domain/user.go`)

**Antes:**
```go
type UserType string
const (
    UserTypeNormal UserType = "normal"
    UserTypeAdmin  UserType = "admin"
)
type User struct {
    UserType string `gorm:"default:'normal'"`
}
```

**Después:**
```go
type Role string
const (
    RoleUser  Role = "USER"   // Usuario común
    RoleAdmin Role = "ADMIN"  // Administrador
    RoleOwner Role = "OWNER"  // Propietario
)
type User struct {
    Role string `gorm:"default:'USER';type:varchar(20)"`
}
```

**Funciones agregadas:**
- `Role.IsValid()`: Valida que un role sea válido (USER, ADMIN, OWNER)

---

### 2. **JWT Claims** (`utils/jwt.go`)

**Antes:**
```go
type Claims struct {
    UserType string `json:"user_type"`
}
func GenerateToken(userID uint, username, userType string) (string, error)
```

**Después:**
```go
type Claims struct {
    Role string `json:"role"` // USER, ADMIN, OWNER - incluido como claim en JWT
}
func GenerateToken(userID uint, username, role string) (string, error)
```

**Funciones agregadas:**
- `IsAdmin(role string) bool`: Verifica si el role es ADMIN
- `IsOwner(role string) bool`: Verifica si el role es OWNER
- `IsUser(role string) bool`: Verifica si el role es USER

---

### 3. **Servicio de Usuario** (`services/user_service.go`)

**Cambios:**
- `CreateUser`: Ahora asigna `RoleUser` por defecto en lugar de "normal"
- `Login`: Genera token con `user.Role` en lugar de `user.UserType`
- `UpdateUser`: Ahora permite actualizar el role (solo ADMIN puede cambiar roles)
- `toDTO`: Retorna `Role` en lugar de `UserType`

---

### 4. **DTOs** (`dto/user_dto.go`)

**Antes:**
```go
type UserResponse struct {
    UserType string `json:"userType"`
}
type UpdateUserRequest struct {
    // Sin campo para role
}
```

**Después:**
```go
type UserResponse struct {
    Role string `json:"role"` // USER, ADMIN, OWNER
}
type UpdateUserRequest struct {
    Role *string `json:"role"` // USER, ADMIN, OWNER - solo ADMIN puede cambiar
}
```

---

### 5. **Middleware** (`middleware/auth_middleware.go`)

**Antes:**
```go
c.Set("user_type", claims.UserType)
// ...
if userType != "admin" {
    // Error
}
```

**Después:**
```go
c.Set("role", claims.Role) // Role incluido en el contexto desde el JWT
// ...
if role != "ADMIN" {
    // Error
}
```

---

### 6. **Controlador** (`controllers/user_controller.go`)

**Cambios:**
- `UpdateUser`: Agrega validación para que solo ADMIN pueda cambiar roles

```go
// Validar que solo ADMIN puede cambiar roles
if req.Role != nil {
    role, exists := c.Get("role")
    if !exists || role != "ADMIN" {
        c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "solo ADMIN puede cambiar roles"})
        return
    }
}
```

---

### 7. **Tests** (`services/user_service_test.go`)

**Cambios:**
- Actualizado para usar `Role` en lugar de `UserType`
- Mensajes de error actualizados a español para coincidir con el código
- Test de contraseña corregido para obtener usuario del repositorio

---

## 📚 Documentación

Se creó el archivo **`ROLES_Y_PERMISOS.md`** que documenta:

1. **Roles disponibles:**
   - **USER**: Leer y reservar propiedades
   - **ADMIN**: Todas las acciones de USER + crear/editar/eliminar propiedades y cambiar roles
   - **OWNER**: Todas las acciones de USER + gestionar sus propias propiedades

2. **Inclusión en JWT:**
   - El role se incluye como claim `role` en el JWT
   - Valores posibles: `USER`, `ADMIN`, `OWNER`

3. **Validación de roles:**
   - Ejemplos de cómo verificar roles en endpoints
   - Uso del middleware `AdminMiddleware`

4. **Cambio de roles:**
   - Solo ADMIN puede cambiar roles de otros usuarios
   - Validaciones implementadas

5. **Matriz de permisos:**
   - Tabla comparativa de permisos por rol

---

## 🔄 Migración de Datos

Si ya tienes usuarios en la base de datos con `user_type` antiguo, necesitas migrarlos:

```sql
-- Migrar usuarios existentes
UPDATE users SET role = 'USER' WHERE user_type = 'normal' OR role IS NULL;
UPDATE users SET role = 'ADMIN' WHERE user_type = 'admin';
```

**Nota:** Después de la migración, puedes eliminar la columna `user_type` si ya no la necesitas.

---

## 🧪 Pruebas

Todos los tests han sido actualizados y deberían pasar correctamente:

```bash
cd backend/users-api
go test ./services/...
```

---

## ⚠️ Notas Importantes

1. **Por defecto:** Todos los usuarios nuevos se crean con role `USER`
2. **JWT Claims:** El role siempre se incluye en el JWT como claim `role`
3. **Validación:** El role se valida en el middleware y está disponible en el contexto
4. **Cambio de Roles:** Solo ADMIN puede cambiar roles de otros usuarios
5. **Backward Compatibility:** Si hay usuarios con `user_type` antiguo, necesitan migración

---

## 📝 Archivos Modificados

1. `backend/users-api/domain/user.go` - Modelo y constantes de roles
2. `backend/users-api/utils/jwt.go` - Claims y funciones de validación
3. `backend/users-api/services/user_service.go` - Lógica de negocio
4. `backend/users-api/dto/user_dto.go` - DTOs actualizados
5. `backend/users-api/middleware/auth_middleware.go` - Middleware actualizado
6. `backend/users-api/controllers/user_controller.go` - Validación de roles
7. `backend/users-api/services/user_service_test.go` - Tests actualizados
8. `backend/users-api/ROLES_Y_PERMISOS.md` - **NUEVO** - Documentación de roles

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el proyecto:**
   ```bash
   cd backend/users-api
   go build
   ```

2. **Ejecutar tests:**
   ```bash
   go test ./services/...
   ```

3. **Probar el login:**
   - El JWT generado debe incluir el claim `role`
   - Verificar con [jwt.io](https://jwt.io) que el token contiene `"role": "USER"` (o "ADMIN"/"OWNER")

4. **Probar cambio de roles:**
   - Solo usuarios con role `ADMIN` pueden cambiar roles
   - Intentar cambiar un role sin ser ADMIN debe retornar 403 Forbidden

---

## 🎯 Próximos Pasos

1. Actualizar `properties-api` para usar el claim `role` del JWT
2. Implementar validaciones de ownership en `properties-api` para el role `OWNER`
3. Actualizar el frontend para mostrar/ocultar funcionalidades según el role del usuario

