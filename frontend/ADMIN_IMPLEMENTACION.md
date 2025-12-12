# 🔐 Implementación de Vista Admin - Frontend

Este documento describe la implementación de la vista `/admin` en el frontend, accesible solo para usuarios con rol ADMIN.

---

## ✅ Funcionalidades Implementadas

### 1. **Ruta /admin**

**Componente:** `src/pages/Admin.jsx`

**Características:**
- ✅ Accesible solo para rol ADMIN (validado por `AdminRoute`)
- ✅ Dos secciones: Usuarios y Propiedades
- ✅ Sistema de tabs para cambiar entre secciones
- ✅ Mensajes de éxito/error con auto-ocultado
- ✅ Loaders durante operaciones
- ✅ Diseño consistente con el resto de la aplicación

---

## 📋 Sección Usuarios

### Funcionalidades

**Listar Usuarios:**
- Consume `GET /api/admin/users` desde users-api
- Muestra tabla con: ID, Email, Rol Actual, Cambiar Rol
- Badges de colores para cada rol (ADMIN, USER, OWNER)

**Cambiar Roles:**
- Selector dropdown para cada usuario
- Permite cambiar entre USER y ADMIN
- Consume `PUT /api/admin/users/:id` con campo `role`
- Muestra mensaje de éxito/error
- Recarga automática de la lista después del cambio

**Validaciones:**
- Solo ADMIN puede cambiar roles (validado en backend)
- No permite cambiar un ADMIN a USER (protección en frontend)

---

## 🏠 Sección Propiedades

### Funcionalidades

**Listar Propiedades:**
- Consume `GET /api/admin/properties` desde properties-api
- Muestra tabla con: Título, Ubicación, Precio, Capacidad, Disponible, Acciones

**Crear Propiedad:**
- Botón "Crear Propiedad" abre modal
- Formulario con campos: título, descripción, precio, ubicación, ownerId, capacidad, disponible
- Consume `POST /api/properties` (solo ADMIN)
- Muestra mensaje de éxito/error
- Recarga automática de la lista

**Editar Propiedad:**
- Botón "Editar" en cada fila abre modal
- Formulario pre-llenado con datos actuales
- Solo actualiza campos modificados
- Consume `PUT /api/properties/:id` (solo ADMIN)
- Muestra mensaje de éxito/error
- Recarga automática de la lista

**Eliminar Propiedad:**
- Botón "Eliminar" en cada fila
- Modal de confirmación antes de eliminar
- Consume `DELETE /api/properties/:id` (solo ADMIN)
- Muestra mensaje de éxito/error
- Recarga automática de la lista

---

## 🔐 Protección de Ruta

### AdminRoute Component

**Componente:** `src/components/AdminRoute.jsx`

**Validaciones:**
1. Verifica que hay token en localStorage
2. Verifica que hay usuario en localStorage
3. Verifica que el usuario tiene rol `ADMIN`
4. Si no cumple alguna condición, redirige apropiadamente

**Comportamiento:**
- Sin token → Redirige a `/login`
- Token válido pero no ADMIN → Redirige a `/search`
- Token válido y ADMIN → Permite acceso

---

## 📊 Endpoints Consumidos

### Usuarios (users-api)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/admin/users` | Listar todos los usuarios (solo ADMIN) |
| PUT | `/api/admin/users/:id` | Actualizar usuario, incluyendo role (solo ADMIN) |

### Propiedades (properties-api)

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| GET | `/api/admin/properties` | Listar todas las propiedades (solo ADMIN) |
| POST | `/api/properties` | Crear propiedad (solo ADMIN) |
| PUT | `/api/properties/:id` | Actualizar propiedad (solo ADMIN) |
| DELETE | `/api/properties/:id` | Eliminar propiedad (solo ADMIN) |

---

## 🎨 Diseño

### Tabs

- Dos tabs: "Usuarios" y "Propiedades"
- Tab activo con borde inferior y color primary
- Iconos: Users y Home

### Tabla de Usuarios

**Columnas:**
- ID: Número de usuario
- Email: Email del usuario
- Rol Actual: Badge con color según rol
- Cambiar Rol: Selector dropdown

**Badges de Rol:**
- **ADMIN:** Morado (bg-purple-100 text-purple-800)
- **USER:** Azul (bg-blue-100 text-blue-800)
- **OWNER:** Verde (bg-green-100 text-green-800)

### Tabla de Propiedades

**Columnas:**
- Título: Nombre de la propiedad
- Ubicación: Dirección/ciudad
- Precio: Formato monetario ($XXX.XX)
- Capacidad: Número de huéspedes
- Disponible: Badge Sí/No
- Acciones: Botones Editar y Eliminar

### Modales

**Modal Crear Propiedad:**
- Formulario completo con validación
- Botones: "Crear Propiedad" y "Cancelar"
- Se cierra automáticamente después de crear exitosamente

**Modal Editar Propiedad:**
- Formulario pre-llenado
- Solo actualiza campos modificados
- Botones: "Guardar Cambios" y "Cancelar"

**Modal Confirmar Eliminación:**
- Mensaje de confirmación
- Botones: "Eliminar" (rojo) y "Cancelar"

---

## 💬 Mensajes de Éxito/Error

### Mensajes de Éxito

- Color: Verde (bg-green-50 border-green-200)
- Icono: Check
- Auto-ocultado: Desaparece después de 3 segundos
- Ejemplos:
  - "Rol del usuario actualizado a ADMIN exitosamente"
  - "Propiedad creada exitosamente"
  - "Propiedad actualizada exitosamente"
  - "Propiedad eliminada exitosamente"

### Mensajes de Error

- Color: Rojo (bg-red-50 border-red-200)
- Icono: AlertCircle
- Persisten hasta que se resuelva el error o se realice otra acción
- Ejemplos:
  - "No tienes permisos para cambiar roles. Se requiere rol ADMIN."
  - "Error al crear la propiedad"
  - "Sesión expirada. Por favor inicia sesión nuevamente."

---

## ⏳ Loaders

### Estados de Carga

**Cargando Usuarios:**
- Spinner animado
- Texto: "Cargando usuarios..."

**Cargando Propiedades:**
- Spinner animado
- Texto: "Cargando propiedades..."

**Durante Operaciones:**
- Los botones se deshabilitan
- Los modales muestran estado de carga

---

## 🔄 Flujo de Funcionamiento

### Cambiar Rol de Usuario

1. ADMIN selecciona nuevo rol en el dropdown
2. Se envía `PUT /api/admin/users/:id` con `{ role: "ADMIN" }` o `{ role: "USER" }`
3. Si es exitoso:
   - Muestra mensaje de éxito
   - Actualiza el badge de rol inmediatamente
   - Recarga la lista después de 1 segundo
4. Si falla:
   - Muestra mensaje de error
   - Mantiene el rol anterior

### Crear Propiedad

1. ADMIN hace clic en "Crear Propiedad"
2. Se abre modal con formulario
3. ADMIN completa los campos
4. Se envía `POST /api/properties` con los datos
5. Si es exitoso:
   - Muestra mensaje de éxito
   - Cierra el modal
   - Recarga la lista de propiedades
6. Si falla:
   - Muestra mensaje de error
   - Mantiene el modal abierto

### Editar Propiedad

1. ADMIN hace clic en "Editar" en una fila
2. Se abre modal con formulario pre-llenado
3. ADMIN modifica los campos deseados
4. Se envía `PUT /api/properties/:id` solo con los campos modificados
5. Si es exitoso:
   - Muestra mensaje de éxito
   - Cierra el modal
   - Recarga la lista de propiedades
6. Si falla:
   - Muestra mensaje de error
   - Mantiene el modal abierto

### Eliminar Propiedad

1. ADMIN hace clic en "Eliminar" en una fila
2. Se abre modal de confirmación
3. ADMIN confirma la eliminación
4. Se envía `DELETE /api/properties/:id`
5. Si es exitoso:
   - Muestra mensaje de éxito
   - Cierra el modal
   - Recarga la lista de propiedades
6. Si falla:
   - Muestra mensaje de error
   - Cierra el modal de confirmación

---

## 🧪 Testing

### Test 1: Acceso como ADMIN

1. Login como ADMIN
2. Ir a `/admin`
3. **Resultado esperado:** Vista admin cargada con tabs Usuarios y Propiedades

### Test 2: Acceso como USER

1. Login como USER
2. Intentar ir a `/admin`
3. **Resultado esperado:** Redirección a `/search`

### Test 3: Cambiar Rol de Usuario

1. Como ADMIN, ir a `/admin` → tab Usuarios
2. Cambiar rol de un usuario de USER a ADMIN
3. **Resultado esperado:** Mensaje de éxito, badge actualizado, lista recargada

### Test 4: Crear Propiedad

1. Como ADMIN, ir a `/admin` → tab Propiedades
2. Hacer clic en "Crear Propiedad"
3. Completar formulario y enviar
4. **Resultado esperado:** Mensaje de éxito, modal cerrado, propiedad en la lista

### Test 5: Editar Propiedad

1. Como ADMIN, ir a `/admin` → tab Propiedades
2. Hacer clic en "Editar" en una propiedad
3. Modificar campos y guardar
4. **Resultado esperado:** Mensaje de éxito, modal cerrado, cambios reflejados

### Test 6: Eliminar Propiedad

1. Como ADMIN, ir a `/admin` → tab Propiedades
2. Hacer clic en "Eliminar" en una propiedad
3. Confirmar eliminación
4. **Resultado esperado:** Mensaje de éxito, propiedad eliminada de la lista

---

## ⚠️ Manejo de Errores

### Errores Manejados

1. **401 Unauthorized:**
   - Sesión expirada
   - Muestra mensaje y redirige a `/login` después de 2 segundos

2. **403 Forbidden:**
   - Usuario no es ADMIN
   - Muestra mensaje de permisos insuficientes

3. **400 Bad Request:**
   - Error de validación
   - Muestra mensaje de error del backend

4. **404 Not Found:**
   - Propiedad no encontrada (al eliminar)
   - Muestra mensaje de error

5. **ERR_NETWORK:**
   - Error de conexión
   - Muestra mensaje de error de conexión

---

## 📝 Formato de Datos

### Request Body - Actualizar Usuario

```json
{
  "role": "ADMIN"
}
```

### Request Body - Crear Propiedad

```json
{
  "title": "Apartamento en el centro",
  "description": "Descripción de la propiedad",
  "price": 150.00,
  "location": "Bogotá, Colombia",
  "ownerId": "123",
  "capacity": 4,
  "available": true,
  "amenities": ["wifi", "pool"],
  "images": []
}
```

### Request Body - Actualizar Propiedad

```json
{
  "title": "Nuevo título",
  "price": 200.00,
  "available": false
}
```

Solo se envían los campos que se modificaron (opcionales).

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el frontend:**
   ```bash
   cd frontend && npm run build
   ```

2. **Probar acceso:**
   - Login como ADMIN → debe poder acceder a `/admin`
   - Login como USER → debe redirigir a `/search` si intenta acceder a `/admin`

3. **Probar gestión de usuarios:**
   - Cambiar rol de un usuario
   - Verificar que se actualiza correctamente

4. **Probar gestión de propiedades:**
   - Crear propiedad
   - Editar propiedad
   - Eliminar propiedad
   - Verificar que todas las operaciones funcionan

---

## 📊 Resumen de Archivos

1. `frontend/src/pages/Admin.jsx` - **NUEVO** - Componente principal de administración
2. `frontend/src/components/AdminRoute.jsx` - **NUEVO** - Componente de protección de ruta para ADMIN
3. `frontend/src/services/api.js` - Actualizado con endpoints de admin
4. `frontend/src/App.jsx` - Ruta `/admin` agregada
5. `frontend/ADMIN_IMPLEMENTACION.md` - Este documento

---

## 🎯 Características Destacadas

1. **Protección de Ruta:** Solo ADMIN puede acceder, validado en frontend y backend
2. **Dos Secciones:** Usuarios y Propiedades con tabs para navegación fácil
3. **Gestión Completa:** Crear, editar y eliminar propiedades
4. **Cambio de Roles:** Selector dropdown para cambiar roles de usuarios
5. **Mensajes Claros:** Éxito/error con auto-ocultado y colores distintivos
6. **Loaders:** Indicadores de carga durante operaciones
7. **Modales:** Interfaz limpia para crear/editar/eliminar
8. **Validaciones:** Validación en frontend y backend
9. **UX Mejorada:** Confirmación antes de eliminar, recarga automática, estados visuales

