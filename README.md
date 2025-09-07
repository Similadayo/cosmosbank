# CosmosBank

CosmosBank is a blockchain application built using the Cosmos SDK, providing a custom bank module for managing digital assets and transactions on a decentralized network.

## Features

- **Send Transactions**: Transfer coins between accounts using the `send` CLI command.
- **Balance Queries**: Query the balance of any account using the `balance` CLI command.
- **CLI Interface**: Command-line tools for transactions and queries under the `bank` module.
- **Modular Architecture**: Built on Cosmos SDK for extensibility and interoperability.

## Prerequisites

- Go 1.25.0 or later
- Make sure you have the necessary dependencies installed as specified in `go.mod`

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/similadayo/cosmosbank.git
   cd cosmosbank
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   make build
   ```

## Usage

### Running the Node

Start a single-node testnet:
```bash
make run
```

### CLI Commands

The CLI commands are organized under the `bank` module:

#### Send Coins
```bash
./build/cosmosbankd tx bank send [from_address] [to_address] [amount]
```

#### Query Balance
```bash
./build/cosmosbankd query bank balance [address]
```

### API Endpoints

- **Send**: gRPC endpoint for sending coins (`MsgSend`).
- **Balance**: gRPC endpoint for querying balances (`QueryBalance`).

## Testing

Run the test suite:
```bash
make test
```

## Project Structure

- `app/`: Main application setup
- `cmd/`: CLI entry point
- `x/bank/`: Custom bank module including CLI commands and keeper logic
- `proto/`: Protocol buffer definitions for messages and queries

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.
