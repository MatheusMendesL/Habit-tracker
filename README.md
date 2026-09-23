# Habit Tracker

[🇧🇷 Leia em Português](README.pt-BR.md)

A scalable habit and routine management platform built with Go and a microservices architecture.

The project starts as a personal habit tracker and is designed to evolve into a social and enterprise productivity platform, where users can manage habits and routines individually or within groups.

## Overview

Habit Tracker allows users to:

* Create and manage habits
* Create and manage routines
* Add habits to routines
* Mark habits and routines as completed
* Track completion history
* Track statistics and streaks
* Manage user profiles
* Communicate through real-time chat

The main differentiator of the project is the **Group system**, which allows the platform to support both social and enterprise use cases.

## Groups

Users can create or join **public or private groups**.

The group's type determines how it operates.

### Social Groups

Designed for friends, communities, and collaborative activities.

Members can:

* Share habits and routines
* Participate in challenges
* Compete through leaderboards
* Unlock achievements
* Interact through group chat

### Enterprise Groups

Designed for teams and organizations.

Administrators can:

* Manage members
* Define mandatory habits
* Define optional habits
* Create shared routines
* Create challenges
* Create custom achievements
* Monitor team activity through dashboards
* Manage roles and permissions

This allows the same platform to support use cases ranging from:

> "My friends and I want to create a challenge."

to:

> "Our company wants to manage routines and goals for a team."

---

# Architecture

The project follows a **microservices architecture**, with each service responsible for a specific domain.

## Current Services

* API Gateway
* Auth Service
* User Service
* Habit Service
* Stats Service
* Social Service

The Routine Service may eventually be separated from the Habit Service as the project grows.

## Planned Services

* Group Service (with subgroups)
* Challenge Service
* Achievement Service
* Notification Service
* Dashboard Service
* Workspace Service
* Feed Service
* File Service
* Routine Service

---

# Technologies

## Backend

* Go
* Node.js
* gRPC
* Protocol Buffers
* PostgreSQL
* Redis

## Infrastructure

* Docker
* Docker Compose
* Nginx

## Development Tools

* sqlc
* Goose
* protoc

---

# Service Responsibilities

### Auth Service

Responsible for authentication and authorization.

### User Service

Responsible for user data and profiles.

### Habit Service

Responsible for:

* Habits
* Routines
* Habit logs
* Routine logs
* Relationships between habits and routines

### Stats Service

Responsible for aggregated user statistics:

* Completed habits
* Completed routines
* Current habit streak
* Longest habit streak
* Current routine streak
* Longest routine streak

The Stats Service does not own habit or routine data. That data belongs to the Habit Service.

### Social Service

Responsible for social functionality and interactions between users.

### API Gateway

Acts as the unified entry point for the platform.

The Gateway is responsible for:

* Exposing HTTP/REST endpoints
* Routing requests to internal services
* Communicating with backend services through gRPC
* Handling authentication at the API entry point
* Standardizing API responses
* Providing request metadata and response timing information

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
```

This allows services such as the Stats Service to process events without requiring synchronous communication with the Habit Service.

---

# API Responses

The API Gateway uses a standardized response structure.

Example:

```json
{
  "timestamp": "2026-09-23T02:55:07Z",
  "api_version": "v1",
  "duration_ms": 35,
  "message": "user found",
  "data": {
    "id": "9979b255-9042-45ad-b321-8f1875d8d562",
    "name": "Matheus",
    "email": "matheusteste4@gmail.com"
  },
  "status": "success",
  "http": {
    "method": "GET",
    "url": "/api/v1/user/GetUserByID/9979b255-9042-45ad-b321-8f1875d8d562"
  }
}
```

The response format provides consistent metadata across API endpoints, including API versioning, execution time, HTTP information, operation status, and response data.

---

# Authentication

Protected API endpoints use **JWT-based authentication**.

Clients must provide a valid JWT through the `Authorization` header:

```http
Authorization: Bearer <token>
```

The API Gateway validates authenticated requests before forwarding them to the appropriate internal service.

---

# How to Run

> ⚠️ **Status:** The main services (Auth, User, Habit, Stats, Social) work individually, and the API Gateway is currently under development. The complete platform is not yet orchestrated as a single stack.

## Prerequisites

* Go 1.2x+
* Docker & Docker Compose
* PostgreSQL
* Node.js (required only for the Auth Service, which is built with Node.js)

## Running from Source

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

Repeat the process for `user-service`, `stats-service`, and `social-service`.

Run the Auth Service:

```bash
cd backend/auth-service-node
npm install
npm start
```

The API Gateway can be run locally from its respective service directory.

## Environment Variables

Each service expects its own `.env` file.

See the `.env.example` file inside each service directory for the required environment variables, including database connections, gRPC ports, and other service-specific configuration.

---

# Infrastructure

The project is designed to run using containerized services through Docker and Docker Compose.

The infrastructure also includes:

* PostgreSQL
* Redis
* Nginx
* OCI infrastructure

The Auth Service is currently deployed to Oracle Cloud Infrastructure (OCI).

---

# Roadmap

Planned improvements include:

* Deploy all services to OCI
* Complete the API Gateway
* Complete the Docker Compose orchestration for the entire stack
* Introduce asynchronous event-driven communication through RabbitMQ
* Implement Group Service
* Implement Challenge Service
* Implement Achievement Service
* Implement Notification Service
* Implement Dashboard Service
* Implement Workspace Service
* Implement Feed Service
* Implement File Service
* Separate the Routine Service from the Habit Service as the domain grows

---

# Project Status

The project is actively under development.

The core microservices are being developed independently, while the API Gateway, authentication flow, inter-service communication, infrastructure, and future domain services are being progressively integrated into the platform.
