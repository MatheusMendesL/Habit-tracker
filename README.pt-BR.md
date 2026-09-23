# Habit Tracker

[🇺🇸 Read in English](README.md)

Uma plataforma escalável de gerenciamento de hábitos e rotinas, construída em Go com arquitetura de microsserviços.

O projeto começa como um habit tracker pessoal e é pensado para evoluir para uma plataforma de produtividade social e empresarial, onde usuários podem gerenciar hábitos e rotinas individualmente ou dentro de grupos.

## Visão Geral

O Habit Tracker permite que os usuários:

- Criem e gerenciem hábitos
- Criem e gerenciem rotinas
- Adicionem hábitos a rotinas
- Marquem hábitos e rotinas como concluídos
- Acompanhem o histórico de conclusões
- Acompanhem estatísticas e streaks
- Gerenciem perfis de usuário
- Se comuniquem através de chat em tempo real

O principal diferencial do projeto é o **sistema de Grupos**, que permite que a plataforma suporte tanto casos de uso sociais quanto empresariais.

## Grupos

Usuários podem criar ou participar de **grupos públicos ou privados**.

Ao criar um grupo, o tipo dele define como ele funciona.

### Grupos Sociais

Pensados para amigos, comunidades e atividades colaborativas.

Membros podem:

- Compartilhar hábitos e rotinas
- Participar de desafios
- Competir através de leaderboards
- Desbloquear conquistas
- Interagir através do chat do grupo

### Grupos Empresariais

Pensados para times e organizações.

Administradores podem:

- Gerenciar membros
- Definir hábitos obrigatórios
- Definir hábitos opcionais
- Criar rotinas compartilhadas
- Criar desafios
- Criar conquistas personalizadas
- Acompanhar a atividade do time através de dashboards
- Gerenciar cargos e permissões

Isso permite que a mesma plataforma seja usada tanto para:

> "Eu e meus amigos queremos criar um desafio."

quanto para:

> "Minha empresa quer gerenciar rotinas e objetivos de um time."

---

# Arquitetura

O projeto segue uma **arquitetura de microsserviços**, com cada serviço responsável por um domínio específico.

## Serviços Atuais

- Auth Service
- User Service
- Habit Service
- Stats Service
- Social Service

O Routine Service eventualmente será separado do Habit Service conforme o projeto cresce.

## Serviços Planejados

- Group Service (com subgrupos)
- Challenge Service
- Achievement Service
- Notification Service
- Dashboard Service
- Workspace Service
- Feed Service
- File Service
- Routine Service

---

# Tecnologias

## Backend

- Go
- NodeJS
- gRPC
- Protocol Buffers
- PostgreSQL
- Redis

## Infraestrutura

- Docker
- Docker Compose
- RabbitMQ
- Circuit Breaker

## Ferramentas de Desenvolvimento

- sqlc
- Goose
- protoc

---

# Responsabilidades dos Serviços

### Auth Service

Responsável pela autenticação e autorização.

### User Service

Responsável pelos dados e perfis dos usuários.

### Habit Service

Responsável por:

- Hábitos
- Rotinas
- Logs de hábitos
- Logs de rotinas
- Relacionamento entre hábitos e rotinas

### Stats Service

Responsável pelas estatísticas agregadas do usuário:

- Hábitos concluídos
- Rotinas concluídas
- Streak atual de hábitos
- Maior streak de hábitos
- Streak atual de rotinas
- Maior streak de rotinas

O Stats Service não é dono dos dados de hábitos ou rotinas. Esses dados pertencem ao Habit Service.

### Social Service

Responsável pela funcionalidade social e pelas interações entre usuários.

---

# Comunicação

Atualmente os serviços se comunicam principalmente através de **gRPC**.

A arquitetura também é pensada para suportar comunicação assíncrona através do **RabbitMQ**.

Por exemplo:

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

---

# Como Rodar

> ⚠️ **Status:** Os serviços principais (Auth, User, Habit, Stats, Social) funcionam individualmente. Ainda não existe um API Gateway unificado — cada serviço precisa ser executado e acessado separadamente.

## Pré-requisitos

- Go 1.2x+
- Docker & Docker Compose
- PostgreSQL
- Node.js (necessário apenas para o Auth Service, que é feito em Node.js)

## Rodando a partir do Código-Fonte

Clone o repositório:

```bash
git clone https://github.com/MatheusMendesL/Habit-tracker.git
cd Habit-tracker
```

Rode um serviço em Go localmente:

```bash
cd services/habit-service
go run cmd/main.go
```

Repita para `user-service`, `stats-service` e `social-service`.

Rode o Auth Service (Node.js):

```bash
cd backend/auth-service-node
npm install
npm start
```

## Variáveis de Ambiente

Cada serviço espera um arquivo `.env` — veja o `.env.example` na pasta de cada serviço para as variáveis necessárias (conexão com o banco de dados, portas gRPC, etc.).

## Em Breve

- Imagens hospedadas no OCI para todos os serviços (atualmente, apenas o **Auth Service** está publicado no OCI)
- API Gateway (ponto de entrada unificado)
- Orquestração completa via docker-compose para rodar toda a stack de uma vez
