# Build stage
FROM node:20-alpine AS builder

# Accept build args
ARG REACT_APP_API_URL

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

# Production stage
FROM nginx:alpine

# Copy build files
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Start nginx
CMD ["nginx", "-g", "daemon off;"]