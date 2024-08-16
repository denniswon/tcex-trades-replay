const { client } = require("websocket");

const _client = new client();

let state = true;

_client.on("connectFailed", (e) => {
  console.error(`[!] Failed to connect : ${e}`);
  process.exit(1);
});

// connect for listening to any order being mined
// event & any transaction being mined in any of those orders
// & any event being emitted from contract interaction transactions
_client.on("connect", (c) => {
  c.on("close", (d) => {
    console.log(`[!] Closed connection : ${d}`);
    process.exit(0);
  });

  // receiving json encoded message
  c.on("message", (d) => {
    console.log(JSON.parse(d.utf8Data));
  });

  // periodic subscription & unsubscription request performed
  handler = (_) => {
    c.send(
      JSON.stringify({
        name: "transaction/*/*",
        type: state ? "subscribe" : "unsubscribe",
      })
    );
    c.send(
      JSON.stringify({
        name: "transaction/0x78566ED47127e2F08EB4DD03F89a03e996e6Fcca/*",
        type: state ? "subscribe" : "unsubscribe",
      })
    );
    c.send(
      JSON.stringify({
        name: "transaction/*/0xa5E0829CaCEd8fFDD4De3c43696c57F7D7A678ff",
        type: state ? "subscribe" : "unsubscribe",
      })
    );

    state = !state;
  };

  setInterval(handler, 10000);
  handler();
});

_client.connect("ws://localhost:7000/v1/ws", null);
