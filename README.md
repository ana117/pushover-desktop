# Pushover Desktop
This Go application acts as a lightweight HTTP server that receives notification requests and displays them as desktop notifications on the local machine using `beeep`.


## Features
- Receives notification requests via HTTP POST.
- Displays desktop notifications/alerts (title, message, optional icon).
- Configurable port and application name via environment variables.


## Getting Started
### Prerequisites
- [Go](https://golang.org/doc/install)

### Installation
1.  **Clone the repository:**
    ```bash
    git clone https://github.com/shaala/pushover-desktop.git
    cd pushover-desktop
    ```
2.  **Download dependencies:**
    ```bash
    go mod tidy
    ```
3.  **Build the executable:**
    ```bash
    go build -o pushover-desktop.exe .
    ```

## Configuration
The application uses environment variables for custom configuration (optional). You can create a `.env` file in the same directory as your executable to set these:
- `PORT`: The port on which the HTTP server will listen (default: `3000`).
- `APP_NAME`: The name that will appear on desktop notifications (default: "Pushover Desktop Notification Service").

**Example `.env` file:**
```dotenv
PORT=8080
APP_NAME="My Custom Notifier"
```


## Running the Application
After building, you can run the executable directly. The application will start an HTTP server and print log messages to the console.

**On Windows:**
```powershell
.\pushover-notifier.exe
```

The server will start listening on the configured port (default `3000`). You will see log messages in your terminal. To stop the server, press `Ctrl+C` in the terminal.


## API Endpoints
The server exposes the following HTTP endpoints:
### `GET /`
A simple health check endpoint.
- **Response (HTTP 200)**

    ```json
    {
      "message": "Pushover Desktop Notification Service is running"
    }
    ```

### `POST /notification`
Sends a standard desktop notification.
- **Request Body (JSON):**

    ```json
    {
      "title": "New Message",
      "message": "You have a new message from Jane Doe.",
    }
    ```
    - `title` (string, required): The title of the notification.
    - `message` (string, required): The main message body.
    - `icon` (string, optional): Not yet implemented, but can be used to specify a path to an icon file for the notification.

- **Success Response (HTTP 200 OK):**
    ```json
    {
      "message": "Notification sent successfully"
    }
    ```

*   **Error Responses:**
    - HTTP 400 Bad Request: Invalid JSON payload or missing `title`/`message`.
    - HTTP 500 Internal Server Error: Failed to display notification on the desktop (e.g., `beeep` error or "Access Denied" if running in a non-interactive session).

### `POST /alert`
Sends a desktop alert, which often has a more prominent display or sound (depending on OS settings).
- **Request Body (JSON):**
    ```json
    {
      "title": "Urgent Alert!",
      "message": "Server CPU usage is critical!"
    }
    ```

- **Success Response (HTTP 200 OK):**
    ```json
    {
      "message": "Alert sent successfully"
    }
    ```

- **Error Responses:**
    - HTTP 400 Bad Request: Invalid JSON payload or missing `title`/`message`.
    - HTTP 500 Internal Server Error: Failed to display alert on the desktop.


## Example Usage (using `curl`)
Make sure your server is running on `http://localhost:3000` (or your configured port).

**Send a notification:**
```bash
curl -X POST \
  http://localhost:3000/notification \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "Reminder",
    "message": "Notification message!"
  }'
```

**Send an alert with an icon (replace with an actual path on your system):**
```bash
curl -X POST \
  http://localhost:3000/alert \
  -H 'Content-Type: application/json' \
  -d '{
    "title": "System Warning",
    "message": "Beep boop! Something needs your attention.",
  }'
```