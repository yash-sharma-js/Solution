# Prime Trade Client

A React-based frontend application for the Prime Trade microservices backend.

## Features

- User authentication (login/register)
- Task management (CRUD operations)
- Role-based access control
- Responsive design with Tailwind CSS
- Real-time data fetching with React Query

## Tech Stack

- React 18
- TypeScript
- Vite
- React Router v6
- TanStack React Query
- Tailwind CSS
- Axios

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Environment Variables

Create a `.env` file in the root directory:

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## Project Structure

```
src/
├── components/        # Reusable UI components
│   ├── Layout.tsx
│   └── TaskModal.tsx
├── contexts/          # React contexts
│   └── AuthContext.tsx
├── pages/             # Page components
│   ├── Login.tsx
│   ├── Register.tsx
│   └── Dashboard.tsx
├── services/          # API services
│   └── api.ts
├── App.tsx            # Main app component
├── main.tsx           # Entry point
└── index.css          # Global styles
```

## API Integration

The client connects to the API Gateway at `http://localhost:8080` which routes requests to:

- `/api/v1/auth/*` - Auth Service
- `/api/v1/tasks/*` - Task Service

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run lint` - Run ESLint
- `npm run preview` - Preview production build
