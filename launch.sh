#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if Docker is running
check_docker() {
    if ! docker info >/dev/null 2>&1; then
        echo -e "${YELLOW}Docker is not running. Please start Docker and try again.${NC}"
        exit 1
    fi
}

# Function to clean up containers and volumes
cleanup() {
    echo -e "${YELLOW}Cleaning up existing containers and volumes...${NC}"
    docker-compose down -v
}

# Function to build and start services
start_services() {
    echo -e "${GREEN}Building and starting services...${NC}"
    docker-compose up --build -d

    echo -e "${GREEN}Waiting for services to be ready...${NC}"
    sleep 10

    echo -e "${GREEN}Services are running:${NC}"
    echo "- Frontend: http://localhost:3000"
    echo "- Backend: http://localhost:5001"
    echo "- PostgreSQL: localhost:5432"
    echo "- Redis: localhost:6379"
    
    echo -e "\n${GREEN}Showing logs (Ctrl+C to exit logs, services will continue running)${NC}"
    docker-compose logs -f
}

# Main execution
check_docker
cleanup
start_services