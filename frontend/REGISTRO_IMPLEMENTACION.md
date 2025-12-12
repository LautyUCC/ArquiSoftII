# 📝 Implementación de Registro - Frontend

Este documento describe la implementación de la ruta de registro en el frontend.

---

## ✅ Funcionalidades Implementadas

### 1. **Ruta /register**

**Componente:** `src/pages/Register.jsx`

**Características:**
- ✅ Formulario de registro con campos: username, email, password, firstName, lastName
- ✅ Validación en tiempo real de todos los campos
- ✅ Validación del lado del cliente antes de enviar
- ✅ Manejo de errores de validación del backend
- ✅ Login automático después del registro exitoso
- ✅ Redirección a /login si el login automático falla
- ✅ Enlace para ir a /login si ya tienes cuenta
- ✅ Diseño consistente con la página de Login

---

## 📋 Campos del Formulario

| Campo | Validación Frontend | Validación Backend |
|-------|-------------------|-------------------|
| **username** | Mínimo 3 caracteres, máximo 50 | required, min=3, max=50 |
| **email** | Formato de email válido | required, email |
| **password** | Mínimo 6 caracteres | required, min=6 |
| **firstName** | Requerido | required |
| **lastName** | Requerido | required |

---

## 🔄 Flujo de Registro

1. **Usuario completa el formulario**
   - Validación en tiempo real mientras escribe
   - Indicadores visuales de errores

2. **Usuario hace clic en "Registrarse"**
   - Validación completa de todos los campos
   - Si hay errores, se muestran y no se envía el request

3. **Request al backend**
   - POST `/api/users` a través de Nginx
   - Nginx hace proxy a `users-api`

4. **Manejo de respuesta:**
   - **Éxito (201):** Intenta login automático
   - **Error (400):** Muestra errores de validación del backend
   - **Error (409/422):** Muestra error si usuario/email ya existe
   - **Error de red:** Muestra mensaje de error de conexión

5. **Login automático:**
   - Si el registro es exitoso, intenta loguear automáticamente
   - Si el login automático es exitoso, guarda token y redirige a `/search`
   - Si el login automático falla, redirige a `/login`

---

## 🎨 Diseño

El componente sigue el mismo diseño que `Login.jsx`:
- Mismo estilo de tarjeta blanca con sombra
- Mismo esquema de colores (primary, secondary)
- Mismos estilos de inputs y botones
- Misma estructura de layout

**Diferencias:**
- Título: "Crea tu cuenta" en lugar de "Propiedades de lujo"
- Formulario con más campos (5 en lugar de 2)
- Grid de 2 columnas para firstName y lastName
- Enlace al final: "¿Ya tienes cuenta? Inicia sesión"

---

## 🔐 Manejo de Errores

### Errores de Validación del Backend

El componente maneja diferentes tipos de errores:

1. **400 Bad Request:**
   - Errores de validación de campos
   - Intenta extraer el campo específico del mensaje de error
   - Muestra el error en el campo correspondiente

2. **409 Conflict / 422 Unprocessable Entity:**
   - Usuario o email ya existe
   - Muestra el error en el campo correspondiente (username o email)

3. **ERR_NETWORK:**
   - Error de conexión
   - Muestra mensaje genérico de error de conexión

4. **Otros errores:**
   - Muestra el mensaje de error del backend o un mensaje genérico

### Ejemplo de Errores Manejados

```javascript
// Error de validación del backend
if (err.response?.status === 400) {
  const errorMessage = err.response?.data?.error || 'Error de validación';
  setError(errorMessage);
  
  // Intentar extraer errores de campos específicos
  const errorText = errorMessage.toLowerCase();
  if (errorText.includes('username')) {
    setFieldErrors(prev => ({...prev, username: errorMessage}));
  }
  // ... más campos
}
```

---

## 🧪 Testing

### Test 1: Registro Exitoso

1. Ir a `/register`
2. Completar todos los campos correctamente
3. Hacer clic en "Registrarse"
4. **Resultado esperado:** Login automático y redirección a `/search`

### Test 2: Validación de Campos

1. Dejar campos vacíos
2. Hacer clic en "Registrarse"
3. **Resultado esperado:** Errores de validación en cada campo

### Test 3: Usuario/Email Duplicado

1. Intentar registrar con username o email que ya existe
2. **Resultado esperado:** Error 409/422 con mensaje específico

### Test 4: Errores del Backend

1. Registrar con datos inválidos (ej: email mal formateado)
2. **Resultado esperado:** Error 400 con mensaje del backend

### Test 5: Navegación

1. Desde `/register`, hacer clic en "Inicia sesión"
2. **Resultado esperado:** Redirección a `/login`
3. Desde `/login`, hacer clic en "Regístrate"
4. **Resultado esperado:** Redirección a `/register`

---

## 📊 Código HTTP y Respuestas

| Código | Descripción | Acción del Frontend |
|--------|-------------|---------------------|
| 201 Created | Usuario creado exitosamente | Intenta login automático |
| 400 Bad Request | Error de validación | Muestra errores en campos |
| 409 Conflict | Usuario/email ya existe | Muestra error específico |
| 422 Unprocessable Entity | Error de validación | Muestra error específico |
| ERR_NETWORK | Error de conexión | Muestra mensaje de error de conexión |

---

## 🔗 Integración con Backend

### Endpoint

**POST `/api/users`** → Nginx → `users-api:8081/users`

### Request Body

```json
{
  "username": "usuario123",
  "email": "usuario@email.com",
  "password": "password123",
  "firstName": "Juan",
  "lastName": "Pérez"
}
```

### Response Success (201)

```json
{
  "message": "Usuario creado exitosamente"
}
```

### Response Error (400)

```json
{
  "error": "el username ya existe"
}
```

---

## ✅ Verificación

Para verificar que todo funciona correctamente:

1. **Compilar el frontend:**
   ```bash
   cd frontend && npm run build
   ```

2. **Probar el registro:**
   - Ir a `http://localhost/register`
   - Completar el formulario
   - Verificar que se registra y loguea automáticamente

3. **Probar validaciones:**
   - Intentar registrar con campos vacíos
   - Intentar registrar con email inválido
   - Intentar registrar con usuario/email duplicado

4. **Probar navegación:**
   - Ir de `/register` a `/login` y viceversa

---

## 📊 Resumen de Archivos

1. `frontend/src/pages/Register.jsx` - **NUEVO** - Componente de registro
2. `frontend/src/App.jsx` - Actualizado para incluir ruta `/register`
3. `frontend/src/services/api.js` - Ya tenía el endpoint de registro configurado
4. `frontend/REGISTRO_IMPLEMENTACION.md` - Este documento

---

## 🎯 Características Destacadas

1. **Validación en Tiempo Real:** Los campos se validan mientras el usuario escribe
2. **Manejo de Errores del Backend:** Extrae y muestra errores específicos de campos
3. **Login Automático:** Después del registro exitoso, intenta loguear automáticamente
4. **Fallback a Login:** Si el login automático falla, redirige a `/login`
5. **Diseño Consistente:** Sigue el mismo estilo que la página de Login
6. **UX Mejorada:** Indicadores visuales de errores y validación

