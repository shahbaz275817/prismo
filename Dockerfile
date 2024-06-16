# Use the official Golang image as the base image for building the application
FROM golang:1.19-bullseye AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the project files to the working directory
COPY . .

# Install dependencies
RUN apt-get update && apt-get install -y \
    git \
    && apt-get clean

# Copy and set up the configuration files
RUN make copy-config

# Install Go tools needed for the setup
RUN make setup

# Build the application
RUN make build

## Run migrations
#RUN make migrate

# Specify the command to run your application
#CMD ["make run"]

CMD ["sh", "-c", "cd out/ && ./prismo server"]
