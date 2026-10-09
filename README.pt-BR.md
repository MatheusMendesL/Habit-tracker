# Habit Tracker

[🇺🇸 Read in English](README.md)

Uma plataforma escalável de gerenciamento de hábitos e rotinas, construída em Go com arquitetura de microsserviços.

O projeto começa como um habit tracker pessoal e é pensado para evoluir para uma plataforma de produtividade social e empresarial, onde usuários podem gerenciar hábitos e rotinas individualmente ou dentro de grupos.

## Visão Geral

O Habit Tracker permite que os usuários:

* Criem e gerenciem hábitos
* Criem e gerenciem rotinas
* Adicionem hábitos a rotinas
* Marquem hábitos e rotinas como concluídos
* Acompanhem o histórico de conclusões
* Acompanhem estatísticas e streaks
* Gerenciem perfis de usuário
* Se comuniquem através de chat em tempo real

O principal diferencial do projeto é o **sistema de Grupos**, que permite que a plataforma suporte tanto casos de uso sociais quanto empresariais.

## Grupos

Usuários podem criar ou participar de **grupos públicos ou privados**.

Ao criar um grupo, o tipo dele define como ele funciona.

### Grupos Sociais

Pensados para amigos, comunidades e atividades colaborativas.

Membros podem:

* Compartilhar hábitos e rotinas
* Participar de desafios
* Competir através de leaderboards
* Desbloquear conquistas
* Interagir através do chat do grupo

### Grupos Empresariais

Pensados para times e organizações.

Administradores podem:

* Gerenciar membros
* Definir hábitos obrigatórios
* Definir hábitos opcionais
* Criar rotinas compartilhadas
* Criar desafios
* Criar conquistas personalizadas
* Acompanhar a atividade do time através de dashboards
* Gerenciar cargos e permissões

Isso permite que a mesma plataforma seja usada tanto para:

> "Eu e meus amigos queremos criar um desafio."

quanto para:

> "Minha empresa quer gerenciar rotinas e objetivos de um time."

---

# Arquitetura

O projeto segue uma **arquitetura de microsserviços**, com cada serviço responsável por um domínio específico.

## Serviços Atuais

* API Gateway
* Auth Service
* User Service
* Habit Service
* Stats Service
* Social Service

O Routine Service eventualmente será separado do Habit Service conforme o projeto cresce.

## Serviços Planejados

* Group Service (com subgrupos)
* Challenge Service
* Achievement Service
* Notification Service
* Dashboard Service
* Workspace Service
* Feed Service
* File Service
* Routine Service

---

# Tecnologias

## Backend

* Go
* Node.js
* gRPC
* Protocol Buffers
* PostgreSQL
* Redis

## Infraestrutura

* Docker
* Docker Compose
* Nginx

## Ferramentas de Desenvolvimento

* sqlc
* Goose
* protoc

---

# Responsabilidades dos Serviços

### Auth Service

Responsável pela autenticação e autorização.

### User Service

Responsável pelos dados e perfis dos usuários.

### Habit Service

Responsável por:

* Hábitos
* Rotinas
* Logs de hábitos
* Logs de rotinas
* Relacionamento entre hábitos e rotinas

### Stats Service

Responsável pelas estatísticas agregadas do usuário:

* Hábitos concluídos
* Rotinas concluídas
* Streak atual de hábitos
* Maior streak de hábitos
* Streak atual de rotinas
* Maior streak de rotinas

O Stats Service não é dono dos dados de hábitos ou rotinas. Esses dados pertencem ao Habit Service.

### Social Service

Responsável pela funcionalidade social e pelas interações entre usuários.

### API Gateway

Atua como ponto de entrada unificado para a plataforma.

O Gateway é responsável por:

* Expor endpoints HTTP/REST
* Encaminhar requisições para os serviços internos
* Comunicar-se com os serviços através de gRPC
* Gerenciar autenticação na entrada da API
* Padronizar as respostas da API
* Fornecer metadados das requisições e informações de tempo de execução

---

# Comunicação

Atualmente, os serviços se comunicam principalmente através de **gRPC**.

O API Gateway utiliza HTTP/REST como interface externa e gRPC para comunicação com os serviços internos.

Exemplo:

```text
Client
   |
   | HTTP/REST
   v
API Gateway
   |
   | gRPC
   v
User Service
```

---

# Respostas da API

O API Gateway utiliza uma estrutura padronizada para as respostas.

Exemplo:

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

A estrutura fornece metadados consistentes entre os endpoints, incluindo versionamento da API, tempo de execução, informações HTTP, status da operação e dados retornados.

---

# Autenticação

Os endpoints protegidos da API utilizam **autenticação baseada em JWT**.

O cliente deve fornecer um JWT válido através do header `Authorization`:

```http
Authorization: Bearer <token>
```

O API Gateway valida as requisições autenticadas antes de encaminhá-las para o serviço interno apropriado.

---

# Como Rodar

> ⚠️ **Status:** Os principais serviços (Auth, User, Habit, Stats, Social) funcionam individualmente, e o API Gateway está atualmente em desenvolvimento. A plataforma completa ainda não está orquestrada como uma única stack.

## Pré-requisitos

* Go 1.2x+
* Docker & Docker Compose
* PostgreSQL
* Node.js (necessário apenas para o Auth Service, que é desenvolvido em Node.js)

## Rodando a partir do Código-Fonte

Clone o repositório:

```bash
git clone https://github.com/MatheusMendesL/Habit-tracker.git
cd Habit-tracker
```

Execute um serviço Go localmente:

```bash
cd services/habit-service
go run cmd/main.go
```

Repita o processo para `user-service`, `stats-service` e `social-service`.

Execute o Auth Service:

```bash
cd backend/auth-service-node
npm install
npm start
```

O API Gateway pode ser executado localmente a partir do diretório correspondente ao serviço.

## Variáveis de Ambiente

Cada serviço espera seu próprio arquivo `.env`.

Consulte o arquivo `.env.example` dentro de cada diretório de serviço para verificar as variáveis de ambiente necessárias, incluindo conexões com o banco de dados, portas gRPC e outras configurações específicas de cada serviço.

---

# Infraestrutura

O projeto foi desenvolvido para utilizar serviços containerizados através de Docker e Docker Compose.

A infraestrutura também utiliza:

* PostgreSQL
* Redis
* Nginx
* Oracle Cloud Infrastructure (OCI)

Atualmente, o Auth Service está hospedado no Oracle Cloud Infrastructure (OCI).

---

# Roadmap

As próximas etapas planejadas incluem:

* Hospedar todos os serviços no OCI
* Finalizar o API Gateway
* Completar a orquestração de toda a stack através do Docker Compose
* Implementar o Group Service
* Implementar o Challenge Service
* Implementar o Achievement Service
* Implementar o Notification Service
* Implementar o Dashboard Service
* Implementar o Workspace Service
* Implementar o Feed Service
* Implementar o File Service
* Separar o Routine Service do Habit Service conforme o domínio crescer

---

# Status do Projeto

O projeto está em desenvolvimento ativo.

Os microsserviços principais estão sendo desenvolvidos de forma independente, enquanto o API Gateway, o fluxo de autenticação, a comunicação entre serviços, a infraestrutura e os futuros serviços de domínio estão sendo progressivamente integrados à plataforma.
