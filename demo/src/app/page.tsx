"use client";

import { Box, Container, Main } from "@/components";
import { Form, FormInputSchema } from "@/components/form";
import Orders from "@/components/orders";
import { Order } from "@/data/order";
import { useWs } from "@/hooks/useWs";
import { isOrder } from "@/utils/data";
import { UiHeading, UiText } from "@uireact/text";
import { useState } from "react";

export default function Home() {
  const [wsError, setWsError] = useState<any>();

  const { ws, open, messages, clearMessages } = useWs<Order>({
    url: "ws://localhost:8080/v1/ws",
    guard: isOrder,
    onError: (error) => {
      setWsError(error);
    },
    onEOF: (eof) => {
      unsubscribe(eof.request_id);
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

    setWsError(undefined);

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
        <Orders
          orders={messages}
          error={wsError || (open ? undefined : "Server disconnected.")}
        />
      </Container>
    </Main>
  );
}
