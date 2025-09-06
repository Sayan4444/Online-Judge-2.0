# Docker Setup for Online Judge 2.0

This guide will help you set up the entire Online Judge application using Docker and Docker Compose.

## Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 2.0+)
- At least 4GB of RAM available for containers

## Architecture

The application consists of the following services:

- **Frontend**: Next.js application (Port 3000)
- **Backend**: Go REST API (Port 8080)
- **Worker**: Go worker service for code execution
- **PostgreSQL**: Database (Port 5432)
- **RabbitMQ**: Message broker (Port 5672, Management UI: 15672)

## Quick Start

1. Clone the repository:
```bash
git clone <repository-url>
cd Online-Judge-2.0
```

2. Build and start all services:
```bash
docker-compose up --build
```

3. Wait for all services to be healthy (this may take a few minutes on first run).

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - RabbitMQ Management: http://localhost:15672 (guest/guest)

## Environment Configuration

### Backend Environment Variables

Create a `.env` file in the backend directory if you need custom configuration:

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=online_judge
RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/
JWT_SECRET=your-jwt-secret-here
SERVER_PORT=8080
```

### Frontend Environment Variables

Create a `.env.local` file in the frontend directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-nextauth-secret-here
```

### Worker Environment Variables

The worker uses the same database and RabbitMQ configuration as the backend.

## Development Mode

For development with hot reloading:

```bash
# Start only database and message broker
docker-compose up postgres rabbitmq

# Run backend locally
cd backend
go run server.go

# Run frontend locally (in another terminal)
cd frontend
npm run dev

# Run worker locally (in another terminal)
cd worker
go run worker.go
```

## Production Deployment

1. Update environment variables with production values
2. Use production-ready secrets for JWT and database passwords
3. Consider using external managed database and message broker services
4. Set up proper logging and monitoring

```bash
# Production build
docker-compose -f docker-compose.yml up --build -d
```

## Useful Commands

```bash
# View logs
docker-compose logs -f [service-name]

# Stop all services
docker-compose down

# Stop and remove volumes (⚠️ This will delete all data)
docker-compose down -v

# Rebuild specific service
docker-compose build [service-name]
docker-compose up -d [service-name]

# Execute commands in containers
docker-compose exec backend sh
docker-compose exec postgres psql -U postgres -d online_judge

# Scale worker instances
docker-compose up -d --scale worker=3
```

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 3000, 5432, 5672, 8080, 15672 are available
2. **Permission issues with worker**: The worker container runs with privileged mode for isolate
3. **Memory issues**: Ensure Docker has enough memory allocated (4GB recommended)

### Health Checks

The compose file includes health checks for PostgreSQL and RabbitMQ. Services will wait for dependencies to be healthy before starting.

### Database Initialization

If you need to initialize the database with specific data, modify the `init.sql` file.

### Worker Isolation

The worker service uses `isolate` for secure code execution. This requires:
- Privileged container mode
- Proper volume mounting for isolate directories

## Security Considerations

- Change default passwords in production
- Use proper secrets management
- Consider network security and firewall rules
- Regular security updates for base images
- Limit worker container permissions where possible

## Monitoring

For production deployments, consider adding:
- Prometheus for metrics
- Grafana for dashboards
- ELK stack for centralized logging
- Health check endpoints

## Backup

Regular backups should include:
- PostgreSQL database
- RabbitMQ configuration and queues (if persistent)
- Application configuration files
