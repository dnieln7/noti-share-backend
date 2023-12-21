# Get GO
FROM golang:1.21.0

# Create work dir
WORKDIR /

# Copy source code
COPY ./ /

# Configure GOOGLE_APPLICATION_CREDENTIALS

COPY notishare.json /
ENV GOOGLE_APPLICATION_CREDENTIALS="/notishare.json"

# Build
RUN go build

# Open port
EXPOSE 4444

# Start
CMD ["/noti-share"]
