import WebSocket from "ws";

const c = new WebSocket("ws://localhost:8080/v1/ws");

let state = true;

c.on("open", () => {
  console.log("Connected to server");

  // periodic subscription & unsubscription request performed
  const handler = () => {
    c.send(
      JSON.stringify({
        type: state ? "subscribe" : "unsubscribe",
        replay_rate: 0.01, // [optional]
        filename: "trades.txt", // [optional]
      })
    );

    state = !state;
  };

  setInterval(handler, 10000);
  handler();
});

c.on("message", (message: string) => {
  console.log(JSON.parse(message));
});

c.on("close", () => {
  console.log(`[!] Closed connection`);
  process.exit(0);
});
