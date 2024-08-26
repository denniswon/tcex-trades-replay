"use client";

import { Box, Container, Main } from "@/components";
import { Form, FormInputSchema } from "@/components/form";
import Kline from "@/components/kline";
import Orders from "@/components/orders";
import { useWs } from "@/hooks/useWs";
import { UiHeading, UiText } from "@uireact/text";
import { useState } from "react";

export default function Home() {
  const [wsError, setWsError] = useState<any>();
  const [mode, _setMode] = useState<"order" | "kline">("order");

  const { ws, open, messages, clearMessages } = useWs({
    url: "ws://localhost:8080/v1/ws",
    onError: (error) => {
      setWsError(error);
    },
  });

  const replayRequestHandler = async (data: FormInputSchema) => {
    if (!ws || !open) {
      return;
    }

    clearMessages();
    setWsError(undefined);

    try {
      const body = {
        type: "subscribe",
        name: mode,
        filename: data.filename,
        replay_rate: Number(data.replay_rate),
      };
      ws.send(JSON.stringify(body));
    } catch (error) {
      console.error(error);
      setWsError(error);
    }
  };

  const unsubscribe = async (request_id: string) => {
    console.log(`Unsubscribing from ${request_id}`, !ws, open);
    if (!ws) {
      return;
    }

    console.log(`Unsubscribing from ${request_id}`);
    try {
      const body = {
        type: "unsubscribe",
        id: request_id,
      };
      ws.send(JSON.stringify(body));
    } catch (error) {
      console.error(error);
      setWsError(error);
    }
  };

  return (
    <Main>
      <Box align="center">
        <UiHeading level={5}>TCEX order replay server demo</UiHeading>
      </Box>
      <Container>
        <Form onSubmit={replayRequestHandler} disabled={!open} />
      </Container>
      <Container>
        <UiText fontStyle="bold" margin={{ all: "four" }}>
          Orders
        </UiText>
        {mode === "order" ? (
          <Orders
            messages={messages}
            error={wsError || (open ? undefined : "Server disconnected.")}
            onEOF={(eof) => {
              unsubscribe(eof.request_id);
            }}
          />
        ) : (
          <Kline
            messages={messages}
            error={wsError || (open ? undefined : "Server disconnected.")}
          />
        )}
      </Container>
    </Main>
  );
}
