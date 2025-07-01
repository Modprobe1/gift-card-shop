# Build stage
FROM node:20-alpine AS builder

# Accept build args
ARG REACT_APP_API_URL=http://45.12.111.208:8080/api/v1

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Set environment variable for build
ENV REACT_APP_API_URL=$REACT_APP_API_URL

# Build the app
RUN npm run build

# Production stage with Caddy
FROM caddy:2-alpine

# Copy build files
COPY --from=builder /app/build /usr/share/caddy

# Create simple Caddyfile
RUN echo ":3000" > /etc/caddy/Caddyfile && \
    echo "root * /usr/share/caddy" >> /etc/caddy/Caddyfile && \
    echo "file_server" >> /etc/caddy/Caddyfile && \
    echo "try_files {path} /index.html" >> /etc/caddy/Caddyfile

# Expose port
EXPOSE 3000

# Start Caddy
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile"]