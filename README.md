# Notification Worker

A Go-based notification worker service that processes and sends notifications.

This service is part of the Travel Booking System ecosystem, specifically designed to handle asynchronous notification processing for the [Travel Booking API](https://github.com/vmdt/travel-booking-api).

## Prerequisites
- Go 1.22.5 or higher
- Docker and Docker Compose (for containerized deployment)
- MongoDB
- RabbitMQ

### Local Development

1. Clone the repository:
   ```bash
   git clone https://github.com/vmdt/notification-worker.git
   cd notification-worker
   ```

2. Copy the example environment file and configure it:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run the application:
   ```bash
   go run cmd/main.go
   ```

## Configuration

The application can be configured using environment variables or configuration files. See `.env.example` for available configuration options.

## License

This project is licensed under the MIT License - see the LICENSE file for details. 