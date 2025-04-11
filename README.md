# My Project

## Overview

This project is a .NET Core web application that utilizes OpenTelemetry for distributed tracing. The application captures HTTP request and response data, including sensitive information, and sends it to an OTL Collector for processing. The traces can be visualized using Jaeger.

## Features

- **Automatic Instrumentation**: The application uses OpenTelemetry to automatically instrument HTTP requests and responses.
- **Data Capture**: Captures the body of requests and responses, along with HTTP method, path, and status code.
- **Sensitive Data Handling**: Configured to remove sensitive fields from traces before sending them to the collector.
- **Database Integration**: Uses SQLite for data storage and initialization.

## Getting Started

### Prerequisites

- [.NET SDK](https://dotnet.microsoft.com/download) (version 6.0 or later)
- [Docker](https://www.docker.com/get-started) (for running the OTL Collector and Jaeger)

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <project-directory>
   ```

2. Build the applications and OTL Collector:
   ```bash
   docker-compose up --build
   ```

3. Access the application at `http://localhost:8888`.

5. Access Jaeger at `http://localhost:16686` to visualize traces.

## Configuration

The OTL Collector is configured to receive traces from the application and export them to Jaeger. The configuration file can be found in the `otlcollector/config.yaml`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
