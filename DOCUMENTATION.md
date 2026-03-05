# Documentación - Prendeluz ERP

## Autenticación y Seguridad

### Sistema de Autenticación

El sistema utiliza autenticación basada en tokens Bearer para proteger los endpoints de la API.

#### Login
**Endpoint:** `POST /login`

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "contraseña"
}
```

**Response:**
```json
{
  "token": "token_generado",
  "id": 123,
  "role": 4
}
```

#### Proceso de Autenticación

1. El usuario envía credenciales (email y password)
2. El sistema verifica las credenciales usando bcrypt para comparar contraseñas hasheadas
3. Si las credenciales son válidas, se genera un token aleatorio de 240 caracteres
4. El token se almacena en la tabla `user_token_erps` con el prefijo "Bearer "
5. Se retorna el token (sin prefijo), ID de usuario y rol

#### Almacenamiento de Tokens

Los tokens se almacenan en la tabla `user_token_erps` con la siguiente estructura:
- `id`: ID único del token
- `user_id`: ID del usuario asociado
- `token`: Token con prefijo "Bearer "
- `valid`: Estado del token (activo/inactivo)
- `created_at` / `updated_at`: Timestamps

### Middlewares de Seguridad

#### 1. Auth Middleware
**Ubicación:** `internal/middlewares/api_autentication.go`

Verifica que el token proporcionado en el header `Authorization` sea válido:
- Extrae el token del header
- Verifica su existencia y validez en la base de datos
- Si es inválido, retorna 401 Unauthorized
- Si es válido, permite continuar con la petición

#### 2. Check Roles Middleware
**Ubicación:** `internal/middlewares/check_roles.go`

Verifica que el usuario tenga los permisos necesarios según su rol.

## Sistema de Roles y Permisos

### Roles Definidos

```go
const (
    StoreManager    = 4  // Gerente de almacén
    StoreWorker     = 5  // Trabajador de almacén
    StoreSupervisor = 6  // Supervisor de almacén
)
```

### Niveles de Acceso

#### AllStoreUsers
Permite acceso a: StoreManager, StoreSupervisor y StoreWorker
- Todos los usuarios del almacén pueden acceder

#### AdminStoreUsers
Permite acceso a: StoreManager y StoreSupervisor
- Solo usuarios con permisos administrativos

### Asignación de Roles

Los roles se almacenan en la tabla `model_has_roles` con la estructura:
- `model_type`: 'App\\Models\\User'
- `model_id`: ID del usuario
- `role_id`: ID del rol asignado

## Funcionalidades por Módulo

### 1. Módulo de Órdenes (Orders)

#### Endpoints para Todos los Usuarios (AllStoreUsers)
- `GET /order` - Obtener listado de órdenes
- `GET /order/status` - Obtener estados de órdenes
- `GET /order/type` - Obtener tipos de órdenes
- `GET /order/supplierOrders` - Obtener órdenes de proveedores
- `PATCH /order/closeOrders` - Cerrar líneas de órdenes
- `PATCH /order/openOrders` - Abrir líneas de órdenes
- `POST /order/editOrders` - Actualizar órdenes mediante Excel
- `POST /order/editSupplierOrders` - Actualizar órdenes de proveedores mediante Excel
- `GET /order/editOrders/frame` - Descargar plantilla Excel para actualización
- `GET /order/editSupplierOrders/frame` - Descargar plantilla Excel para proveedores
- `POST /order/add/api` - Crear orden vía API
- `GET /order/supplierOrders/download` - Descargar Excel de órdenes de proveedores

#### Endpoints para Administradores (AdminStoreUsers)
- `POST /order/add` - Agregar orden mediante Excel
- `GET /order/add/frame` - Descargar plantilla para agregar órdenes
- `POST /order/addByRequest` - Crear orden mediante JSON
- `PATCH /order` - Editar órdenes

### 2. Módulo de Líneas de Orden (Order Lines)

#### Endpoints para Todos los Usuarios (AllStoreUsers)
- `GET /order/orderLines` - Obtener líneas de órdenes
- `GET /order/orderLines/labels` - Obtener etiquetas de líneas
- `PATCH /order/orderLines/add` - Agregar cantidad a líneas
- `PATCH /order/orderLines/remove` - Remover cantidad de líneas
- `PATCH /order/orderLines` - Editar líneas de órdenes

#### Asignación de Líneas
- `POST /order/orderLines/asignation` - Crear asignación de líneas
- `PATCH /order/orderLines/asignation` - Editar asignación de líneas

### 3. Módulo de Órdenes Padre (Father Orders)

#### Endpoints para Todos los Usuarios (AllStoreUsers)
- `GET /fatherOrder` - Obtener datos de órdenes padre
- `PATCH /fatherOrder` - Actualizar órdenes padre
- `GET /fatherOrder/orderLines` - Obtener líneas por ID de orden padre
- `GET /fatherOrder/orderLines/downloadPicking` - Descargar Excel de picking
- `GET /fatherOrder/amazonExcel` - Descargar Excel para Amazon
- `POST /fatherOrder/closePicking` - Cerrar orden de picking
- `PATCH /fatherOrder/close` - Cerrar líneas de orden
- `PATCH /fatherOrder/open` - Abrir líneas de orden

### 4. Módulo de Almacén (Store)

#### Endpoints Autenticados
- `PATCH /store/:order_code` - Actualizar almacén por código de orden
- `GET /store/:store_name` - Obtener stock de almacén
- `GET /store` - Obtener listado de almacenes
- `POST /store/excel` - Actualizar stock mediante Excel
- `GET /store/excel/frame` - Descargar plantilla Excel para actualización

### 5. Módulo de Déficit de Stock

#### Endpoints Autenticados
- `GET /stock_deficit` - Obtener déficit de stock
- `GET /stock_deficit/calc` - Calcular déficit por orden
- `GET /stock_deficit/download` - Descargar Excel de déficit

#### Endpoints para Administradores (AdminStoreUsers)
- `PATCH /stock_deficit/clean` - Limpiar déficit de stock

### 6. Módulo de Stock

#### Endpoints Autenticados
- `GET /stock/getExcel` - Obtener datos de stock en Excel
- `POST /stock/byAsins` - Obtener stock por ASINs

### 7. Módulo de Ubicaciones de Almacén (Store Locations)

#### Endpoints Autenticados
- `GET /store_location` - Obtener ubicaciones
- `POST /store_location` - Crear ubicación
- `PATCH /store_location` - Actualizar ubicación

### 8. Módulo de Pallets

#### Endpoints Autenticados
- `GET /pallet` - Obtener pallets
- `GET /pallet/crossDataByOrderId` - Obtener pallet por ID de orden
- `POST /pallet` - Crear pallet
- `PATCH /pallet` - Actualizar pallet

### 9. Módulo de Cajas (Boxes)

#### Endpoints Autenticados
- `GET /box` - Obtener cajas
- `POST /box` - Crear caja
- `PATCH /box` - Actualizar caja

#### Endpoints para Administradores (AdminStoreUsers)
- `DELETE /box` - Eliminar caja

### 10. Módulo de Líneas de Orden en Cajas

#### Endpoints Autenticados
- `GET /order_line_boxes` - Obtener líneas en cajas
- `POST /order_line_boxes` - Crear línea en caja
- `POST /order_line_boxes/withProcess` - Crear línea con proceso
- `PATCH /order_line_boxes` - Actualizar línea en caja

### 11. Módulo de Ubicación de Stock de Items

#### Endpoints Autenticados
- `GET /item_stock_location` - Obtener ubicaciones de stock
- `DELETE /item_stock_location` - Eliminar ubicación
- `POST /item_stock_location` - Crear ubicación
- `PATCH /item_stock_location` - Actualizar ubicación
- `PATCH /item_stock_location/stockChanges` - Cambios de stock
- `PATCH /item_stock_location/stockMovement` - Movimientos de stock

#### Endpoints para Administradores (AdminStoreUsers)
- `DELETE /item_stock_location/cleanZeroStock` - Limpiar stock en cero

### 12. Módulo de Proveedores (Suppliers)

#### Endpoints para Todos los Usuarios (AllStoreUsers)
- `GET /supplier` - Obtener listado de proveedores

### 13. Módulo de Estadísticas (Statistics)

#### Endpoints para Todos los Usuarios (AllStoreUsers)
- `GET /stadistics/olHisotricByFatherOrder` - Obtener histórico de órdenes
- `GET /stadistics/lines` - Obtener estadísticas de líneas
- `GET /stadistics/items` - Obtener estadísticas de items

## Configuración CORS

El sistema permite peticiones desde los siguientes orígenes:
- `http://localhost:3000`
- `http://127.0.0.1:3000`
- `http://localhost:3001`
- `http://127.0.0.1:3001`
- `https://testerp.zarivy.com`
- `https://erp.zarivy.com`

**Métodos permitidos:** GET, POST, PUT, PATCH, DELETE, OPTIONS

**Headers permitidos:** Origin, Content-Type, Accept, Authorization

**Credenciales:** Habilitadas

**Cache:** 12 horas

## Estructura de Base de Datos

### Tabla: users
- `id`: ID único del usuario
- `name`: Nombre del usuario
- `email`: Email único
- `email_verified_at`: Timestamp de verificación
- `password`: Contraseña hasheada (bcrypt)
- `remember_token`: Token de sesión persistente
- `created_at` / `updated_at`: Timestamps

### Tabla: user_token_erps
- `id`: ID único del token
- `user_id`: Referencia al usuario
- `token`: Token Bearer
- `valid`: Estado del token
- `created_at` / `updated_at`: Timestamps

### Tabla: model_has_roles
- `model_type`: Tipo de modelo ('App\\Models\\User')
- `model_id`: ID del usuario
- `role_id`: ID del rol asignado

## Flujo de Autenticación

```
1. Usuario → POST /login (email, password)
2. Sistema → Verifica credenciales en tabla users
3. Sistema → Compara password con bcrypt
4. Sistema → Genera token aleatorio (240 caracteres)
5. Sistema → Guarda token en user_token_erps
6. Sistema → Consulta rol en model_has_roles
7. Sistema ← Retorna {token, id, role}
8. Usuario → Incluye token en header Authorization: Bearer {token}
9. Middleware Auth → Verifica validez del token
10. Middleware CheckRoles → Verifica permisos según rol
11. Sistema → Procesa petición si todo es válido
```

## Seguridad

### Encriptación de Contraseñas
- Se utiliza bcrypt para hashear contraseñas
- Las contraseñas nunca se almacenan en texto plano

### Validación de Tokens
- Todos los endpoints (excepto /login) requieren token válido
- Los tokens se verifican en cada petición
- Los tokens pueden ser invalidados cambiando el campo `valid` a false

### Control de Acceso Basado en Roles (RBAC)
- Tres niveles de roles: Manager, Supervisor, Worker
- Permisos granulares por endpoint
- Verificación automática mediante middlewares

## Notas Importantes

1. **Tokens**: Los tokens no tienen expiración automática, deben ser invalidados manualmente
2. **Roles**: Los roles se gestionan mediante la tabla `model_has_roles` (sistema Laravel Spatie)
3. **CORS**: Configurado para entornos de desarrollo y producción específicos
4. **Base de Datos**: El sistema asume una base de datos MySQL/MariaDB compatible con GORM
5. **Archivos Excel**: Múltiples endpoints soportan carga y descarga de archivos Excel para operaciones masivas
