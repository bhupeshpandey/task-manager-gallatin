FROM alpine:3.18

# Set environment variables
ENV PORT=50051
ENV CONFIG_PATH=/app/config.yaml

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the locally built Linux binary
COPY ./build/task-manager-gallatin_linux_amd64 ./gallatin

# Copy the config.yaml file into the container
COPY config.yaml .

# Expose the port the app runs on
EXPOSE 50051

# Run the binary
CMD ["./gallatin", "-config", "/app/config.yaml"]
