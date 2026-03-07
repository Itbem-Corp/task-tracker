import express, { Request, Response } from 'express';
import { z } from 'zod';
import { randomUUID } from 'crypto';

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

// Mock user data
interface User {
  id: string;
  name: string;
  email: string;
}

const users: User[] = [
  { id: '1', name: 'Alice Johnson', email: 'alice@example.com' },
  { id: '2', name: 'Bob Smith', email: 'bob@example.com' },
  { id: '3', name: 'Charlie Davis', email: 'charlie@example.com' }
];

// Zod schema for POST /api/users
const createUserSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  email: z.string().email('Invalid email address')
});

// GET /api/health
app.get('/api/health', (req: Request, res: Response) => {
  res.json({ status: 'ok' });
});

// GET /api/users
app.get('/api/users', (req: Request, res: Response) => {
  res.json(users);
});

// POST /api/users
app.post('/api/users', (req: Request, res: Response) => {
  try {
    const validatedData = createUserSchema.parse(req.body);
    
    const newUser: User = {
      id: randomUUID(),
      name: validatedData.name,
      email: validatedData.email
    };
    
    users.push(newUser);
    res.status(201).json(newUser);
  } catch (error) {
    if (error instanceof z.ZodError) {
      res.status(400).json({ error: 'Validation failed', details: error.errors });
    } else {
      res.status(500).json({ error: 'Internal server error' });
    }
  }
});

app.listen(PORT, () => {
  console.log(`Server running on http://localhost:${PORT}`);
});
