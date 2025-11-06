# Development setup for Vite/React frontend (assumes npm; adjust for yarn if needed)
FROM node:22.13.1-alpine

# Install dependencies
RUN apk add --no-cache tzdata && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

ENV TZ=Asia/Almaty

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production=false  # Install dev deps for dev server

# Copy source code
COPY . .

# Expose port (Vite default)
EXPOSE 3000

# Run dev server, binding to all interfaces for container access
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"]