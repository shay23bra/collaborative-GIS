# Collaborative GIS Application

This application enables real-time collaboration on drawing and analyzing areas on a map. Users can draw, view, and calculate the area of polygons in real-time, switch between map layers, and view other users’ actions as they draw.

## Features
- Real-time polygon drawing on a map
- Satellite and regular map layer switching
- WebSocket-based updates for real-time collaboration
- Area calculation for polygons

## Tech Stack
- Backend: Golang with Gorilla WebSocket and PostgreSQL
- Frontend: Leaflet.js, HTML, JavaScript
- Database: PostgreSQL with PostGIS for spatial data storage
- Docker for containerized setup

## Getting Started

### Prerequisites
- Docker
- Docker Compose

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/collaborative-gis-app.git
   cd collaborative-gis-app

### Running the Application
1. Start the application:
   ```bash
   docker-compose up
   ```
2. Open `http://localhost:8080` in your browser.
3. Connect with one of 2 default users (shay1, shay2) with the password `password`.
4. You can draw polygons on the map and see the changes in real-time.
