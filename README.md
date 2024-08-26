# TCEX order replay server

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

## Demo

### Replay Order at x60 rate (demo video below trimmed in the middle due to github readme upload size limit)
https://github.com/user-attachments/assets/d3978f17-d6dd-4ac9-983e-83805b0af578

### Replay Order at x600 rate (demo video below trimmed in the middle due to github readme upload size limit)
https://github.com/user-attachments/assets/1b7cf986-9e96-458c-9a08-425ccb995e3e

### Replay Order at x0.1 rate (demo video below trimmed in the middle due to github readme upload size limit)
https://github.com/user-attachments/assets/c22d1b8f-9b17-45e0-9ee7-086154f52b38


