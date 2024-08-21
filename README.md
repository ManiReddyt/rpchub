# RPC Hub

RPC Hub is a lightweight Go server that manages and forwards RPC requests to healthy endpoints. It periodically checks the health of RPC endpoints and dynamically routes requests based on the health status, ensuring reliable and efficient communication with blockchain nodes.

## Features

- **Dynamic Health Checks:** Periodically checks the health of RPC endpoints and updates their status.
- [WIP]**Load Balancing:** Routes requests to available healthy RPC endpoints.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/ManiReddyt/rpchub
   cd rpchub
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Configure the RPC endpoints:**

   - Edit the `config.json` file to include the RPC URLs categorized by chain ID:
   - An example is given in `config.json`

4. **Run the server:**

   ```bash
   go run main.go
   ```

5. **Access the RPC endpoint:**
   - You can now send requests to the server:
   ```
   POST http://localhost:8080/?chain_id=1
   ```

## Health Check

The server periodically checks the health of all RPC endpoints defined in `config.json`. Healthy endpoints are stored in an in-memory map, and only these endpoints are used to handle incoming requests.

The health check runs every 10 seconds by default. You can adjust this interval in the code by modifying the time interval in the cron job within `main.go`.

## Example

There is an example Thunder Client collection provided in `thunder-collection-rpchub.json`. Please check it out to see how to structure your requests. The response will be the same as what the RPC endpoints return.

## Contributing

Contributions are welcome! Feel free to submit a pull request or open an issue to discuss potential improvements or new features.
