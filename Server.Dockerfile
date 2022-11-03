# Image
FROM golang:1.19.1

# Define the work directory
WORKDIR ./Backend

# Copy the project folders into the container's working directory
COPY . .

# Build it
RUN go build -o server digitalpaper/backend

# Give executable privileges
RUN chmod +x server

# Run the backend server
CMD [ "./server" ]
