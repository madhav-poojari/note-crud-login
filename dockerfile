# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# Build the server binary
RUN go build -o myserver

# Expose the port that the server will listen on
EXPOSE 4068

# Define an environment variable (if needed)
# ENV MY_ENV_VARIABLE=value

# Command to run the server when the container starts
CMD ["./myserver"]
