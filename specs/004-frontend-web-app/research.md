# Research: Podcast Management Web Application

**Feature**: 004-frontend-web-app
**Date**: 2026-02-08

## Overview

This document consolidates research findings for building a Vue 3 + TypeScript frontend with Go backend API for podcast management.

## Frontend Technology Stack

### Decision: Vue 3 + Composition API + TypeScript

**Rationale**:
- Aligns with project constitution (Vue 3 with Composition API required)
- TypeScript provides type safety for API contracts and data models
- Composition API enables better code reuse through composables
- Reactive state management built-in, no external library needed

**Alternatives Considered**:
- React: Rejected - constitution mandates Vue 3
- Vue 2 Options API: Rejected - constitution requires Composition API
- JavaScript: Rejected - TypeScript strict mode required by constitution

### Decision: Vite 5.x for Build Tooling

**Rationale**:
- Fast HMR (Hot Module Replacement) for development
- Optimized production builds with code splitting
- Native ESM support
- Official Vue 3 recommendation

**Alternatives Considered**:
- Webpack: Rejected - slower build times, more complex configuration
- Rollup: Rejected - Vite provides better DX with similar output

### Decision: Tailwind CSS 3.x for Styling

**Rationale**:
- Required by project constitution
- Utility-first approach speeds development
- Responsive design built-in
- Small production bundle with PurgeCSS

**Alternatives Considered**:
- CSS Modules: Rejected - constitution mandates Tailwind
- Styled Components: Rejected - not aligned with constitution

## Backend API Design

### Decision: RESTful API with Go net/http

**Rationale**:
- Aligns with constitution (Web API First Design)
- Go standard library sufficient for simple CRUD operations
- No need for heavy frameworks (Gin, Echo) for this scope
- Existing codebase uses standard library patterns

**Alternatives Considered**:
- Gin/Echo frameworks: Rejected - unnecessary complexity for 4-6 endpoints
- GraphQL: Rejected - REST sufficient for simple data fetching
- gRPC: Rejected - web browser client requires REST/HTTP

### Decision: In-Memory Task Queue

**Rationale**:
- Single-user local deployment (no persistence needed across restarts)
- Simple implementation with Go channels
- Fast access for status polling (2-3 second intervals)
- Aligns with existing CLI downloader architecture

**Alternatives Considered**:
- Redis: Rejected - overkill for single-user local deployment
- Database: Rejected - constitution prefers file-based storage
- Persistent queue: Rejected - task history not critical across restarts

## State Management Pattern

### Decision: Vue 3 Composition API with Composables

**Rationale**:
- Constitution explicitly prohibits external state management libraries
- Composables provide reusable reactive state logic
- Sufficient for single-page application scope
- Built-in reactivity system handles updates efficiently

**Best Practices**:
- Create composables for each domain (useEpisodes, useTasks, usePagination)
- Use `ref()` for primitive values, `reactive()` for objects
- Implement loading/error states in composables
- Poll backend using `setInterval` in composables with cleanup

**Alternatives Considered**:
- Pinia: Rejected - constitution removed Pinia, use Composition API only
- Vuex: Rejected - deprecated in favor of Composition API
- Global reactive objects: Rejected - composables provide better encapsulation

## Pagination Strategy

### Decision: Server-Side Pagination with Client-Side Page Size Control

**Rationale**:
- Supports 1000+ episodes efficiently (SC-006)
- Reduces initial load time (SC-001: < 2 seconds)
- Backend can optimize file system scanning
- Client controls page size (20/50/100) for user preference

**Implementation**:
- Backend: Accept `page` and `pageSize` query parameters
- Frontend: usePagination composable manages state
- Cache current page in memory, fetch on page change
- Show loading state during page transitions

**Alternatives Considered**:
- Client-side pagination: Rejected - loading 1000 episodes upfront violates SC-001
- Infinite scroll: Rejected - clarification specified pagination
- Virtual scrolling: Rejected - pagination provides better UX for large lists

## Polling Strategy for Task Updates

### Decision: Client-Side Polling Every 2-3 Seconds

**Rationale**:
- Meets SC-005 requirement (updates visible within 5 seconds)
- Simple implementation with setInterval
- No WebSocket infrastructure needed for single-user app
- Balances responsiveness with server load

**Implementation**:
- Frontend: setInterval in useTasks composable
- Poll only when tasks page is active (cleanup on unmount)
- Backend: Lightweight GET /api/tasks endpoint
- Return only changed tasks (compare timestamps)

**Alternatives Considered**:
- WebSockets: Rejected - overkill for single-user local deployment
- Server-Sent Events: Rejected - polling simpler for this scope
- Long polling: Rejected - standard polling sufficient

## Error Handling Strategy

### Decision: Centralized Error Handling with User-Friendly Messages

**Rationale**:
- FR-027 to FR-030 require backend unavailability detection
- Consistent error UX across all API calls
- Retry mechanism for transient failures

**Implementation**:
- Frontend: API client wrapper with try-catch and error mapping
- Display toast/banner for errors with retry button
- Composables expose error state for component-level handling
- Backend: Consistent error response format (JSON with error field)

**Best Practices**:
- Network errors: Show "Server unavailable" with retry button
- Validation errors: Show field-specific messages
- 404 errors: Show "Not found" with navigation back
- 500 errors: Show "Something went wrong" with retry

## Testing Strategy

### Decision: Vitest for Unit Tests, Cypress for E2E

**Rationale**:
- Constitution requires Vitest and Cypress
- Vitest: Fast, Vite-native, compatible with Vue Test Utils
- Cypress: Reliable e2e testing with good DX

**Test Coverage**:
- Unit tests: Composables, utility functions, API client
- Component tests: Vue components with Vue Test Utils
- E2E tests: Complete user workflows (view episodes, create tasks)
- Target: 80%+ coverage for business logic

**Best Practices**:
- Mock API calls in unit tests
- Test loading/error states in components
- E2E tests cover all acceptance scenarios from spec
- Use data-testid attributes for stable selectors

## Summary

All technical decisions align with project constitution and feature requirements:

1. **Frontend**: Vue 3 + Composition API + TypeScript + Vite + Tailwind CSS
2. **Backend**: Go with net/http standard library, RESTful API
3. **State**: Composables pattern (no external state library)
4. **Storage**: File-based for episodes, in-memory for tasks
5. **Pagination**: Server-side with configurable page size
6. **Updates**: Client polling every 2-3 seconds
7. **Testing**: Vitest + Cypress

No unresolved clarifications remain. Ready to proceed to Phase 1 (Design & Contracts).
