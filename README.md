# TCEX order replay server

The purpose of this server is to replay historical trade messages as if it were simulating production data flow at requested replay rate (1x, 5x, 60x, etc). This server currently loads a static set of trades from a file, and stream them on-demand via WebSocket.

The server hosts the subscription and unsubscription API for clients, and the demo webpage displays the simple data visualization as the client receives the replayed order messages over subscribed webSocket.

The replay server runs on <http://localhost:8080/>.

## Installation

```bash
git clone git@github.com:denniswon/tcex.git

cd tcex

cp .env.example .env
```

## Commands

```bash
# runs the initial setup script
make setup

# generate protobufs
make proto_gen

# installs backend go mod dependencies
go mod tidy
# builds backend
make build
# run backend
./tcex

# or to build & run together
make run

# run frontend
cd demo; yarn dev;

# builds and runs both back and frontend locally on separate processes
make demo
```

## High-Level Architecture

The webserver has been designed with concurrency and scalability in mind. That is, it tries to minimize the number and durations of files being opened on the server for processing. Also, it utilizes processing queue _kappa architecture_ for processing order replays in a streaming manner. In particular, the system uses **Redis** cache and pubsub based processing queues, one for processing trades input file and another for handling actual order replays.

This way, the server is able to handle multiple requests concurrently without blocking any of them while minimizing the resources used and modularity among the subcomponents for easier incremental optimization in the future.

### Demo

#### Replay Order at x60 rate (demo video below trimmed in the middle due to github readme upload size limit)

<https://github.com/user-attachments/assets/d3978f17-d6dd-4ac9-983e-83805b0af578>

#### Replay Order at x600 rate

<https://github.com/user-attachments/assets/1b7cf986-9e96-458c-9a08-425ccb995e3e>

#### Replay Order at x0.1 rate

<https://github.com/user-attachments/assets/c22d1b8f-9b17-45e0-9ee7-086154f52b38>

### Subscribing with Order Replay Requests

For requesting and listening to orders being replayed, connect to `/v1/ws` endpoint using websocket client library & once connected, send **subscription** request with payload _( JSON encoded )_

```json
{
  "type": "subscribe",
  "name": "order", // "order" or "kline"
  "filename": "trades.txt", // optional, defaults to 'trades.txt' if not supplied.
  "replay_rate": 60, // optional, defaults to 60 for x60 replay rate.
  "granularity": 60 // optional, defaults to 60. used only for "kline" requests. in seconds.
}
```

Subscription confirmeation response _( JSON encoded )_

```json
{
  "code": 1,
  "message": "Subscribed to <subscription_id>"
}
```

Real-time notification about orders being replayed at requested rate:

```json
{
  "price": "1347.41",
  "quantity": 200,
  "aggressor": "ask",
  "timestamp": 1722527801638
}
```

Cancel subscription:

```json
{
  "type": "unsubscribe",
  "id": "<subscription_id>" // subscription id returned from subscription request above
}
```

Unsubscription confirmation response:

```json
{
  "code": 1,
  "message": "Unsubscribed from <subscription_id>"
}
```
