# Quickstart Guide: Podcast Management Web Application

**Feature**: 004-frontend-web-app
**Date**: 2026-02-08

## Overview

This guide provides step-by-step instructions for setting up and running the podcast management web application locally.

## Prerequisites

- Go 1.21+ installed
- Node.js 18+ and npm installed
- Git installed
- Existing podcast downloads in the downloads directory

## Backend Setup

### 1. Navigate to Backend Directory

```bash
cd backend
```

### 2. Install Go Dependencies

```bash
go mod tidy
```

### 3. Run Backend Server

```bash
go run cmd/server/main.go
```

The backend server will start on `http://localhost:8080`.

## Frontend Setup

### 1. Navigate to Frontend Directory

```bash
cd frontend
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Run Development Server

```bash
npm run dev
```

The frontend will start on `http://localhost:5173` (default Vite port).

## Accessing the Application

1. Open your browser and navigate to `http://localhost:5173`
2. You should see the navigation bar with "Podcasts" and "Tasks" links
3. Click "Podcasts" to view your downloaded episodes
4. Click "Tasks" to create new download tasks

## Testing the Application

### View Downloaded Podcasts

1. Navigate to the Podcasts page
2. You should see a paginated list of downloaded episodes
3. Click on any episode to view its show notes in a modal
4. Use pagination controls to navigate through pages

### Create a Download Task

1. Navigate to the Tasks page
2. Click the "New Download" button
3. Enter a valid Xiaoyuzhou FM episode URL
4. Click "Submit"
5. The task will appear in the list with "pending" status
6. Watch as the status updates automatically (polls every 2-3 seconds)

## Running Tests

### Backend Tests

```bash
cd backend
go test ./...
```

### Frontend Tests

```bash
cd frontend

# Unit tests
npm run test

# E2E tests
npm run test:e2e
```

## Troubleshooting

### Backend server won't start

- Check if port 8080 is already in use
- Verify Go version: `go version` (should be 1.21+)
- Check downloads directory exists and is accessible

### Frontend won't connect to backend

- Verify backend is running on `http://localhost:8080`
- Check browser console for CORS errors
- Ensure API base URL is correctly configured in frontend

### No episodes showing

- Verify podcast files exist in the downloads directory
- Check backend logs for file scanning errors
- Ensure file permissions allow backend to read downloads

## Next Steps

- Review the [API documentation](./contracts/api.yaml) for endpoint details
- Check [data-model.md](./data-model.md) for entity definitions
- See [research.md](./research.md) for technology decisions
- Run `/speckit.tasks` to generate implementation tasks
