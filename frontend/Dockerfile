# Build del frontend
FROM node:22-alpine AS builder

WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./

# Instalar dependencias
# Usar npm install con --legacy-peer-deps para resolver conflictos de peer dependencies con React 19
# npm install es más flexible que npm ci y puede manejar actualizaciones de dependencias
RUN npm install --legacy-peer-deps

# Copiar código fuente
COPY . .

# Build de producción
RUN npm run build

# Stage 2: Copiar archivos al volumen compartido
FROM alpine:latest

# Instalar herramientas necesarias
RUN apk add --no-cache bash

# Crear directorio de destino
RUN mkdir -p /usr/share/nginx/html

# Copiar el build a un directorio temporal en la imagen
COPY --from=builder /app/dist /tmp/dist

# Script que copia los archivos al volumen compartido al iniciar
RUN echo '#!/bin/sh' > /entrypoint.sh && \
    echo 'echo "Copiando archivos al volumen compartido..."' >> /entrypoint.sh && \
    echo 'cp -r /tmp/dist/* /usr/share/nginx/html/' >> /entrypoint.sh && \
    echo 'echo "Archivos copiados. Listando contenido:"' >> /entrypoint.sh && \
    echo 'ls -la /usr/share/nginx/html/assets/ | head -5' >> /entrypoint.sh && \
    echo 'echo "Iniciando contenedor..."' >> /entrypoint.sh && \
    echo 'while true; do sleep 3600; done' >> /entrypoint.sh && \
    chmod +x /entrypoint.sh

# Ejecutar el script de entrada
CMD ["/entrypoint.sh"]

