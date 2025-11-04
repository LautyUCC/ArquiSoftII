FROM node:18-alpine

WORKDIR /app

# Archivo temporal para que no falle
RUN echo '{"name":"spotly-frontend","version":"1.0.0"}' > package.json
RUN echo 'console.log("Frontend - Coming soon"); setInterval(() => {}, 1000);' > index.js

EXPOSE 3000

CMD ["node", "index.js"]