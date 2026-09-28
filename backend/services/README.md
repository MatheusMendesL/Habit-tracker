# Habit Tracker

## Future Features

- Groups
- Shared routines
- Challenges
- Achievements
- Internal progression system
- Feed
- Leaderboards
- Group chat
- Business workspaces
- Team dashboard
- Custom achievements
- Role & permission system
- Required & optional habits
- Calendar view
- Recurring challenges
- File attachments
- AI habit suggestions (future)

## Future Architecture

- Event-driven communication (RabbitMQ)
- Circuit Breaker
- Observability (Prometheus & Grafana)
- Centralized logging
- Rate limiting

## Future Microservices

- Challenge Service
- Achievement Service
- Notification Service
- Group Service
- Dashboard Service
- Workspace Service
- Feed Service
- File Service
- Routine Service (split from Habit Service)

## TLS certificates

Generate the local CA and one server certificate for each Go service before
starting the services:

```bash
./backend/services/shared/gen-certs.sh
```

The script creates the shared CA in `backend/services/certs` and copies the
service certificates to each `cert` directory. Private keys and generated
certificates are ignored by Git.