# RSS Feed Aggregator Server

## Overview

This project is an efficient **RSS feed aggregator server** in **Go**. It enables users to manage their favorite RSS feeds seamlessly by providing a robust backend for collecting, storing, and serving aggregated posts. Built with performance and scalability in mind, this server leverages concurrent goroutines and a PostgreSQL database for reliable data storage and retrieval.

## Features
1. **Add Feeds**: Users can add RSS feed URLs via HTTP requests.
2. **Concurrent Feed Fetching**: The server uses a configurable number of goroutines to fetch and process feeds concurrently.
3. **Post Storage**: Fetched posts are parsed and stored in a PostgreSQL database.
4. **Follow/Unfollow Feeds**: Users can subscribe to specific feeds they are interested in.
5. **Aggregated View**: Users can view a summary of aggregated posts for their followed feeds, including links to the full articles.

## Getting Started

### Prerequisites

- **Go**: Version 1.19 or higher.
- **PostgreSQL**: Ensure you have a running PostgreSQL instance.
- **Make**: For using the provided Makefile commands.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jessie-sr/feed-aggregator.git
   cd feed-aggregator
   ```

2. Configure Environment Variables:
   - Set up the server's port number and database connection string by creating the .env file in the project root directory:
  
     ```ini
     # Server port number
     PORT=<your_server_port>
      
     # PostgreSQL database connection string
     DB_URL=<your_postgresql_connection_string>
      
     # Database migration directory (optional)
     MIGRATION_URL=<your_postgresql_connection_string_for_migration>
     ```

3. Run database migrations:
   ```bash
   make migrate
   ```

4. Build the application:
   ```bash
   make build
   ```

5. Start the server:
   ```bash
   make run
   ```
