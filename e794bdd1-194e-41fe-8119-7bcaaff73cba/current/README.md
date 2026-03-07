# Express REST API with TypeScript

A simple Node.js REST API built with Express and TypeScript, featuring Zod validation.

## Setup Instructions

1. Install dependencies:
```bash
npm install
```

2. Run the development server:
```bash
npm run dev
```

The server will start on `http://localhost:3000`.

## API Endpoints

### GET /api/health
Health check endpoint.

**Example:**
```bash
curl http://localhost:3000/api/health
```

**Response:**
```json
{"status":"ok"}
```

### GET /api/users
Retrieve all users.

**Example:**
```bash
curl http://localhost:3000/api/users
```

**Response:**
```json
[
  {"id":"1","name":"Alice Johnson","email":"alice@example.com"},
  {"id":"2","name":"Bob Smith","email":"bob@example.com"},
  {"id":"3","name":"Charlie Davis","email":"charlie@example.com"}
]
```

### POST /api/users
Create a new user with validation.

**Example:**
```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'
```

**Response:**
```json
{"id":"a1b2c3d4-e5f6-7890-abcd-ef1234567890","name":"John Doe","email":"john@example.com"}
```

**Validation Error Example:**
```bash
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"","email":"invalid-email"}'
```

**Response:**
```json
{
  "error":"Validation failed",
  "details":[
    {"code":"too_small","minimum":1,"type":"string","message":"Name is required","path":["name"]},
    {"validation":"email","code":"invalid_string","message":"Invalid email address","path":["email"]}
  ]
}
```

## Build for Production

```bash
npm run build
npm start
```
