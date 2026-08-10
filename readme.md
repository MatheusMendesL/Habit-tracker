# Habit Tracker

[🇧🇷 Leia em português](README.pt-BR.md)

A scalable habit and routine management platform built with Go and a microservices architecture.

The project starts as a personal habit tracker and is designed to evolve into a social and business productivity platform, where users can manage habits and routines individually or within groups.

## Overview

Habit Tracker allows users to:

- Create and manage habits
- Create and manage routines
- Add habits to routines
- Mark habits and routines as completed
- Track completion history
- Track statistics and streaks
- Manage user profiles
- Communicate through real-time chat

The main differentiator of the project is the **Groups system**, which allows the platform to support both social and business use cases.

## Groups

Users can create or participate in **public or private groups**.

When creating a group, its type can define how it works.

### Social Groups

Designed for friends, communities and collaborative activities.

Members can:

- Share habits and routines
- Participate in challenges
- Compete through leaderboards
- Unlock achievements
- Interact through group chat

### Business Groups

Designed for teams and organizations.

Administrators can:

- Manage members
- Define required habits
- Define optional habits
- Create shared routines
- Create challenges
- Create custom achievements
- Monitor team activity through dashboards
- Manage roles and permissions

This allows the same platform to be used for both:

> "My friends and I want to create a challenge."

and:

> "My company wants to manage routines and goals for a team."

---

# Architecture

The project follows a **microservices architecture**, with each service responsible for a specific domain.

## Current Services

- Auth Service
- User Service
- Habit Service
- Stats Service
- Social Service

The Routine Service will eventually be separated from the Habit Service as the project grows.

## Planned Services

- Group Service (with subgroups)
- Challenge Service
- Achievement Service
- Notification Service
- Dashboard Service
- Workspace Service
- Feed Service
- File Service
- Routine Service

---

# Technologies

## Backend

- Go
- NodeJS
- gRPC
- Protocol Buffers
- PostgreSQL
- Redis

## Infrastructure

- Docker
- Docker Compose
- RabbitMQ
- Circuit Breaker

## Development Tools

- sqlc
- Goose
- protoc

---

# Service Responsibilities

### Auth Service

Responsible for authentication and authorization.

### User Service

Responsible for user data and profiles.

### Habit Service

Responsible for:

- Habits
- Routines
- Habit logs
- Routine logs
- Habit/routine relationships

### Stats Service

Responsible for aggregated user statistics:

- Completed habits
- Completed routines
- Current habit streak
- Longest habit streak
- Current routine streak
- Longest routine streak

The Stats Service does not own habit or routine data. That data belongs to the Habit Service.

### Social Service

Responsible for social functionality and interactions between users.

---

# Communication

Services currently communicate primarily through **gRPC**.

The architecture is also designed to support asynchronous communication through **RabbitMQ**.

For example:

```text
Habit Service
      |
      | HabitCompleted
      v
   RabbitMQ
      |
      v
Stats Service

---

# Getting Started

> ⚠️ **Status:** Core services (Auth, User, Habit, Stats, Social) are functional individually. There is no unified API Gateway yet — each service must be run and accessed separately.

## Prerequisites

- Go 1.2x+
- Docker & Docker Compose
- PostgreSQL
- Node.js (only required for Auth Service, which is built in Node.js)

## Building from Source

Clone the repository:

```bash
git clone https://github.com/MatheusMendesL/Habit-tracker.git
cd Habit-tracker
```

Run a Go service locally:

```bash
cd services/habit-service
go run cmd/main.go
```

Repeat for `user-service`, `stats-service`, and `social-service`.

Run the Auth Service (Node.js):

```bash
cd backend/auth-service-node
npm install
npm start
```

## Environment Variables

Each service expects a `.env` file — see `.env.example` in each service folder for required variables (database connection, gRPC ports, etc.).

## Coming Soon

- OCI-hosted images for all services (currently, only the **Auth Service** is published to OCI)
- API Gateway (unified entrypoint)
- Full docker-compose orchestration for running the whole stack at once