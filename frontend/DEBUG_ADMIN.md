# 🔍 Debug: Verificar si el Usuario es ADMIN

## Pasos para Verificar

### 1. Abre la Consola del Navegador (F12)

### 2. Verifica el Usuario en localStorage

Ejecuta en la consola:

```javascript
// Ver el usuario guardado
const user = JSON.parse(localStorage.getItem('user'));
console.log('Usuario:', user);
console.log('Role:', user?.role);
console.log('Es ADMIN?', user?.role === 'ADMIN');
```

### 3. Verifica el Token JWT

```javascript
// Ver el token
const token = localStorage.getItem('token');
console.log('Token:', token ? 'Existe' : 'No existe');

// Decodificar el token (solo para ver, no usar en producción)
if (token) {
  const parts = token.split('.');
  const payload = JSON.parse(atob(parts[1]));
  console.log('JWT Payload:', payload);
  console.log('Role en JWT:', payload.role);
}
```

### 4. Verifica el Estado en React

Si tienes React DevTools instalado, puedes ver el estado del componente Search:
- Busca el componente `Search`
- Verifica el estado `isAdmin`
- Debería ser `true` si el usuario es ADMIN

## Problemas Comunes

### El role no es "ADMIN"

Si el role es diferente (por ejemplo, "admin" en minúsculas o "ADMINISTRATOR"):

1. **Verifica en la base de datos:**
   ```sql
   SELECT id, username, email, role FROM users WHERE username = 'admin';
   ```

2. **Actualiza el role si es necesario:**
   ```sql
   UPDATE users SET role = 'ADMIN' WHERE username = 'admin';
   ```

3. **Vuelve a hacer login** para que se actualice el localStorage

### El role está en minúsculas

El código busca `role === 'ADMIN'` (mayúsculas). Si en la BD está en minúsculas, actualiza:

```sql
UPDATE users SET role = 'ADMIN' WHERE role = 'admin';
```

### El usuario no tiene role

Si el campo `role` está NULL o vacío:

```sql
UPDATE users SET role = 'ADMIN' WHERE username = 'admin' AND (role IS NULL OR role = '');
```

## Solución Rápida

Si quieres forzar que el usuario sea ADMIN temporalmente para probar:

```javascript
// En la consola del navegador (F12)
const user = JSON.parse(localStorage.getItem('user'));
user.role = 'ADMIN';
localStorage.setItem('user', JSON.stringify(user));
// Recargar la página
window.location.reload();
```

**NOTA:** Esto es solo para pruebas. El role real debe venir del backend.

