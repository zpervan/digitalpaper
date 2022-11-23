# Images
FROM node:18.12.1-alpine

# Define work directory
WORKDIR ./frontend

# Copy the project folders into the container's working directory
COPY . .

# Front application workflow
RUN yarn --cwd frontend install
CMD ["yarn", "--cwd", "frontend", "run", "dev"]