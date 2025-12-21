<!--
Sync Impact Report:
- Version change: 1.0.0 → 2.0.0 (MAJOR: Complete architectural realignment from Python CLI to Go+Vue web service)
- Modified principles: All 5 principles completely redefined for web service architecture
  - "Python Standards" → "Go Backend Standards"
  - "Modular Package Design" → "Service-Oriented Architecture"
  - "Test-First Development" → "Web Service Testing"
  - "Robust Error Handling" → "Asynchronous Processing & Error Handling"
  - "CLI Interface First" → "Web API First"
- Added sections: Frontend Standards, External Integrations
- Removed sections: CLI-specific testing requirements
- Templates requiring updates: ⚠ plan-template.md (update technical context), ⚠ tasks-template.md (update paths and structure)
- Follow-up TODOs: Update project structure references in templates
-->

# Podcast Reader Constitution

## Core Principles

### I. Go Backend Standards
Go 1.21+ with strict formatting and best practices. All code MUST use gofmt, follow Go conventions, include comprehensive error handling, and use Go modules for dependency management. Concurrent processing with goroutines for long-running audio processing tasks.

### II. Service-Oriented Architecture
Clear separation between frontend (Vue.js) and backend (Go server). Backend exposes RESTful APIs, frontend handles UI/UX. File-based storage for processing state and results (no database). Each major function (download, transcription, LLM processing) MUST be implemented as independent services.

### III. Asynchronous Processing First
All audio processing tasks MUST be asynchronous. Return process ID immediately, provide status endpoints for tracking. Use job queues for managing concurrent processing. Implement timeout handling and progress tracking for long-running operations.

### IV. External Integration Resilience
All external service calls (Tencent Cloud, LLM APIs) MUST implement circuit breakers, retry logic with exponential backoff, and graceful degradation. Store raw responses for debugging. Implement cost controls and usage monitoring for third-party services.

### V. Web API First Design
All functionality MUST be accessible via REST APIs with OpenAPI documentation. Support JSON request/response format. Implement proper HTTP status codes, CORS policies for frontend communication, and API versioning for future compatibility.

## Frontend Standards

### Vue.js Architecture
Vue 3 with Composition API, TypeScript for type safety. Vite for build tooling with hot reload. Component-based architecture with reusable UI elements. Responsive design for mobile and desktop access.

### State Management
Use Pinia for complex state, local state for simple components. Implement proper loading states, error handling, and user feedback for all async operations. Store process IDs and results in component state.

### Build Quality
ESLint + Prettier for code formatting. TypeScript strict mode enabled. Automated testing with Vitest for unit tests and Cypress for e2e tests. Bundle size optimization and lazy loading for performance.

## Backend Standards

### Go Code Quality
gofmt + golint + go vet for code quality. Use context for request-scoped values and cancellation. Structured logging with proper log levels. Graceful shutdown handling for in-flight processing.

### API Design
RESTful endpoints with consistent patterns: GET /status/{id}, POST /process, GET /result/{id}. Request validation with proper error responses. Rate limiting to prevent abuse. Health check endpoints for monitoring.

### File Management
Organized directory structure for processing artifacts: downloads/, transcripts/, briefings/. Cleanup policies for old files. Atomic file operations to prevent corruption. Unique filenames with process IDs.

## Testing Requirements

### Frontend Testing
- Unit tests with Vitest for components and utilities
- Integration tests for API interactions
- E2E tests with Cypress for complete user workflows
- Visual regression tests for UI consistency

### Backend Testing
- Unit tests for service functions with table-driven tests
- Integration tests for API endpoints
- Mock external services (Tencent Cloud, LLM)
- Load testing for concurrent processing scenarios

### System Testing
- Full workflow tests: URL submission → processing → result retrieval
- Error scenario testing: network failures, invalid URLs, service outages
- Performance testing: concurrent user load, large file processing
- Accessibility testing for web interface

## External Integrations

### Tencent Cloud Integration
Secure API key management with environment variables. Request/response logging for debugging. Implement audio format validation and size limits. Cost monitoring and usage quotas.

### LLM Integration
Abstract LLM interface to support multiple providers. Prompt template management for consistent briefing generation. Response caching to reduce costs. Content filtering and safety checks.

### Error Recovery
Automatic retry with exponential backoff for transient failures. Manual retry options for permanent failures. Fallback processing paths when external services are unavailable. User notification system for integration issues.

## Governance

### Amendment Process
- Proposed changes MUST be documented with impact analysis
- Changes affecting core architecture require consensus approval
- Frontend/backend coordination required for API changes
- All changes update version according to semantic versioning

### Compliance Review
- PR reviews MUST check constitutional compliance for both frontend and backend
- Automated checks verify code quality, test coverage, and API contract compliance
- Monthly reviews of external service usage and costs
- Security reviews for API keys and external integrations

### Versioning Policy
- MAJOR: Architectural changes, breaking API changes, technology stack changes
- MINOR: New features, new external integrations, additional endpoints
- PATCH: Bug fixes, UI improvements, documentation updates, performance optimizations

**Version**: 2.0.0 | **Ratified**: 2025-12-21 | **Last Amended**: 2025-12-21