# Get GO
FROM golang:1.21.0

# Create work dir
WORKDIR /

# Copy source code
COPY ./ /

# Build
RUN go build

# Open port
EXPOSE 4444

# Start
CMD ["/noti-share"]
