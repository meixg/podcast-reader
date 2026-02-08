# Podcast Reader - Web Application

A web-based podcast management application for browsing downloaded episodes and managing download tasks.

## Features

- **Browse Episodes**: View all downloaded podcast episodes with pagination
- **Episode Details**: View show notes in a modal dialog
- **Download Management**: Create and monitor download tasks
- **Real-time Updates**: Task status updates automatically every 2-3 seconds

## Tech Stack

### Backend
- Go 1.21+
- Standard library `net/http`
- In-memory task queue

### Frontend
- Vue 3 with Composition API
- TypeScript (strict mode)
- Vite 5.x
- Tailwind CSS 3.x
- Vue Router 4.x

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- npm or yarn

## Installation

### Backend Setup

1. Navigate to the backend directory:
```bash
cd backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the server:
```bash
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`.

### Frontend Setup

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Run the development server:
```bash
npm run dev
```

The frontend will start on `http://localhost:5173`.

## Usage

1. Open your browser and navigate to `http://localhost:5173`
2. Use the navigation bar to switch between "Podcasts" and "Tasks" pages
3. On the Podcasts page, click any episode to view its show notes
4. On the Tasks page, click "New Download" to create a download task

## License

MIT
