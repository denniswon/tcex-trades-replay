"use client";

import { Box, Container, Main } from "@/components";
import { Form, FormInputData } from "@/components/form";
import Kline from "@/components/kline";
import Orders from "@/components/orders";
import { UploadHeader, useFileUpload } from "@/hooks/useFileUpload";
import { useWs } from "@/hooks/useWs";
import { UiHeading, UiText } from "@uireact/text";
import { useEffect, useState } from "react";

export default function Home() {
  const [wsError, setWsError] = useState<any>();
  const [mode, setMode] = useState<"order" | "kline">("kline");
  const [formInput, setFormInput] = useState<FormInputData>();

  const { ws, open, messages, clearMessages } = useWs({
    url: "ws://localhost:8080/v1/ws",
    onError: (error) => {
      setWsError(error);
    },
  });

  const { header, isError, isLoading } = useFileUpload({
    file: formInput?.file,
  });

  const replayRequestHandler = async (data: FormInputData) => {
    if (!ws || !open) {
      return;
    }

    clearMessages();
    setWsError(undefined);

    setFormInput(data);
    return;
  };

  useEffect(() => {
    if (!ws || !open || !formInput || !header || isError || isLoading) {
      return;
    }

    subscribe(header, formInput);
  }, [ws, open, formInput, header, isError, isLoading]);

  const subscribe = async (header: UploadHeader, data: FormInputData) => {
    console.log(`Subscribing for request ${header.id}`, header, data);

    try {
      const body = {
        type: "subscribe",
        name: mode,
        id: header.id,
        filename: header.filepath,
        replay_rate: Number(data.replay_rate),
        granularity: data.granularity,
      };
      console.log(body);

      ws.send(JSON.stringify(body));
    } catch (error) {
      console.error(error);
      setWsError(error);
    } finally {
      setFormInput(undefined);
    }
  };

  const unsubscribe = async (request_id: string) => {
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
        <Form
          onSubmit={replayRequestHandler}
          disabled={!open || isError}
          uploading={isLoading}
          mode={mode}
          onSetMode={setMode}
        />
      </Container>
      <Container>
        <UiText fontStyle="bold" margin={{ all: "four" }}>
          {mode.replace(/^.{1}/g, mode[0].toUpperCase())}
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
