# Ticket Booking System

A modern, full-stack web application for event management and real-time concurrent ticket booking. Built with a robust Go backend and a highly responsive, aesthetic SvelteKit frontend.

## Features

- **Conversational Booking Agent**: Integrates LLM-driven event discovery with guided conversational booking, powered by a controlled vocabulary and real-time availability filters.
- **Optimistic Concurrency Control**: Prevents double-booking and race conditions when multiple users attempt to book the exact same seat simultaneously.
- **Premium User Interface**: Built with Svelte, featuring beautiful micro-animations, skeleton loaders, and dynamic rendering.
- **Virtual Tickets**: Dynamically generated digital tickets with CSS-based aesthetic barcodes and precise purchase timestamps.
- **Super Admin Dashboard**: A secure control panel restricted to specific users, allowing for on-the-fly event creation and management.
- **Security & Protection**: Fully protected API utilizing JSON Web Tokens (JWT) for authentication and Global IP-based rate limiting to prevent spam and DDoS attacks.
- **Automated Seat Generation**: Generating an event automatically triggers a database transaction to map and generate all associated seating capacity.
- **Vercel Serverless Ready**: Backend structured to deploy seamlessly as a Vercel serverless function (`/api/index.go`), connecting directly to Supabase via IPv4 transaction poolers.

## Tech Stack

### Frontend
- **Framework**: Svelte / Vite
- **Styling**: Tailwind CSS
- **State Management**: Built-in Svelte reactivity

### Backend
- **Language**: Go
- **Framework**: Fiber (Fast HTTP Engine)
- **Database**: SQLite (Development) / PostgreSQL (Production ready via GORM & Supabase)
- **ORM**: GORM
- **Authentication**: JWT & bcrypt

## Quick Start (Development)

### Prerequisites
- Node.js (v18+)
- Go (v1.20+)

### 1. Start the Backend
```bash
cd backend
go mod tidy
go run ./cmd/api/main.go
```
The backend server will run on `http://localhost:8080`. It will automatically initialize the local SQLite database and apply schema migrations.

### 2. Seed the Database
To populate the database with realistic spread and consistent event data (spanning multiple categories, cities, and dates):
```bash
cd backend
go run ./cmd/seed/main.go
```
*Note: The seed script is idempotent and safely guarded with a marker event.*

### 3. Start the Frontend
```bash
cd frontend
npm install
npm run dev
```
The frontend application will be available at `http://localhost:5173`.

## Authentication & Admin Access
You can register a new account from the web interface. To access the **Admin Dashboard** and create new events, you must register an account with the exact email: `tmarked4l@gmail.com`.
