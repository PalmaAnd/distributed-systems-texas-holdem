# Stage 1: Build Flutter web (run manually: flutter build web)
# The build output in build/web/ must exist before building this image.
# Build: cd frontend && flutter build web && docker build -t texas-holdem-frontend .
FROM nginx:alpine AS runtime
# AMD64 is default for nginx:alpine

# Serve built Flutter web files
COPY build/web /usr/share/nginx/html

# SPA routing: serve index.html for all routes
RUN echo 'server { listen 80; root /usr/share/nginx/html; index index.html; location / { try_files $uri $uri/ /index.html; } }' > /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
