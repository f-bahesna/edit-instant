FROM node:18-alpine

WORKDIR /app

# Copy package.json to install dependencies during the Docker build phase
COPY package.json ./

# Install npm dependencies inside the container
RUN npm install

# Expose Vite dev server port
EXPOSE 5173

# Run development server with host configuration for access outside container
CMD ["npm", "run", "dev", "--", "--host"]
