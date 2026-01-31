# Texas Hold'em Poker - Hand Evaluator & Probability Calculator

A distributed Texas Hold'em poker application with hand evaluation and Monte Carlo probability calculation. Built with Flutter (frontend), Go (backend), containerized with Docker, and deployed on Google Kubernetes Engine (GKE) with AMD64 nodes.

## Tech Stack

| Layer     | Technology                    |
|----------|-------------------------------|
| Frontend | Flutter / Dart                |
| Backend  | Go (REST API)                 |
| Runtime  | Docker, GKE (AMD64)           |

## Card Format

Cards are 2-character strings:
- **First char**: Suit – `H` (Hearts), `D` (Diamonds), `C` (Clubs), `S` (Spades)
- **Second char**: Rank – `A` (Ace), `K` (King), `Q` (Queen), `J` (Jack), `T` (Ten), `9`–`2`

Examples: `HA` (Heart Ace), `S7` (Spade 7), `CT` (Club Ten)

## Project Structure

```
.
├── backend/          # Go REST API
│   ├── cmd/
│   ├── internal/
│   ├── Dockerfile
│   └── go.mod
├── frontend/         # Flutter app
│   └── ...
├── k8s/              # Kubernetes manifests for GKE
├── load-test/        # Load testing (k6)
└── docs/             # Step-by-step deployment guide
```

## Quick Start

### Prerequisites

- Go 1.21+
- Flutter 3.x
- Docker
- kubectl
- gcloud CLI

### Local Development

```bash
# Backend
cd backend && go run ./cmd/server

# Frontend (separate terminal)
cd frontend && flutter run -d chrome
```

### Docker Build (AMD64)

```bash
# Backend
docker build --platform linux/amd64 -t texas-holdem-backend:latest ./backend

# Frontend (web)
cd frontend && flutter build web
docker build --platform linux/amd64 -t texas-holdem-frontend:latest .
```

## API Endpoints

| Method | Endpoint            | Description                                           |
|--------|---------------------|-------------------------------------------------------|
| POST   | `/api/v1/evaluate`  | Best hand from 2 hole + 5 community cards             |
| POST   | `/api/v1/compare`   | Compare two hands, return winner                      |
| POST   | `/api/v1/probability` | Win probability via Monte Carlo simulation         |

## Step-by-Step Guide

See [docs/PROJECT_GUIDE.md](docs/PROJECT_GUIDE.md) for the complete walkthrough from development to GKE deployment.

## References

- [Texas Hold'em (Wikipedia)](https://en.wikipedia.org/wiki/Texas_hold_%27em)
- [Poker Hand Rankings](https://en.wikipedia.org/wiki/List_of_poker_hands)
- Peter Norvig's poker rules
