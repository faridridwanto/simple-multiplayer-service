# Simple Multiplayer WebSocket Service

A Go-based WebSocket service that allows users to connect, communicate with each other, and disconnect.

## Features

- WebSocket connection handling
- User identification by connection ID
- Message routing between users
- Connection cleanup when users leave

## Requirements

- Go 1.24 or higher
- [gorilla/websocket](https://github.com/gorilla/websocket) package

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/simple-multiplayer-service.git
   cd simple-multiplayer-service
   ```

2. Install dependencies:
   ```
   go mod download
   ```

## Running the Server

Start the WebSocket server:

```
go run cmd/server/main.go
```

The server will start on port 8080 by default. You can connect to the WebSocket endpoint at `ws://localhost:8080/ws`.

## Using the Client

1. Open the `client.html` file in a web browser.
2. Click the "Connect" button to establish a WebSocket connection.
3. The server will assign you a connection ID, which will be displayed in the messages area.
4. To send a message to another user:
   - Enter their connection ID in the "Recipient Connection ID" field.
   - Type your message in the "Message" field.
   - Click "Send Message".
5. To disconnect, click the "Disconnect" button.

## Testing

### Running Unit Tests

Run the unit tests with:

```
go test -v
```

### Manual Testing

For manual testing, you can:

1. Open multiple instances of the `client.html` file in different browser windows.
2. Connect each client to the server.
3. Copy the connection ID from one client and use it as the recipient ID in another client.
4. Send messages between clients to verify the functionality.

## Project Structure

The project follows the standard Go module layout:

- `cmd/server/`: Contains the main application entry point
  - `main.go`: The main application that starts the WebSocket server
- `internal/`: Contains packages that are internal to the application
  - `websocket/`: Contains the WebSocket server implementation
    - `manager.go`: Manages WebSocket connections
    - `handler.go`: Handles WebSocket requests
    - `manager_test.go`: Tests for the connection manager
    - `handler_test.go`: Tests for the WebSocket handler
- `pkg/`: Contains packages that can be used by external applications
  - `client/`: Contains the client implementation
    - `client.go`: Defines the WebSocket client
  - `message/`: Contains the message implementation
    - `message.go`: Defines the message structure
- `client.html`: A simple HTML/JavaScript client for manual testing

## How It Works

1. **Connection**: When a user connects to the WebSocket endpoint, they are assigned a unique connection ID.
2. **Message Routing**: Users can send messages to other users by specifying the recipient's connection ID.
3. **Disconnection**: When a user disconnects, their connection is cleaned up and removed from the connection manager.

## Future Improvements

- User authentication
- Persistent message storage
- Room-based chat functionality
- Matchmaking for multiplayer games
