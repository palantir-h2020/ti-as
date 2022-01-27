# Base golang image with alpine3.13 as base os
FROM golang:1.17.0-alpine3.13

# Create project directory
WORKDIR /PalantirIpAnonymizationService

# Copy project dependent files
COPY go.mod ./
COPY go.sum ./

# Download all needed libraries
RUN go mod download

# Copy .go files and project directories
COPY backend/ backend/
COPY cryptopan/ cryptopan/
COPY logger/ logger/
COPY utils/ utils/
COPY config/ config/
COPY *.go ./
COPY start_anonymization_service.sh ./

# Build Go application
RUN go build -o /PalantirIpAnonymizationService

# Set ENV variables
# Port, where IP Anonymization service will be running (Default: 8100)
ENV PORT=8100
# Database instance
ENV DB_INSTANCE=redis
# Database host
ENV DB_HOST=localhost
# Database port
ENV DB_PORT=6379
# Database username
ENV DB_USERNAME=
# Database password
ENV DB_PASSWORD=
# Database to use
ENV DB_NAME=0

# Expose service port
EXPOSE ${PORT}

# Execute Go service when starting container
CMD ["/bin/sh", "start_anonymization_service.sh"]