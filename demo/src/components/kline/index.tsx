import styled from "styled-components";

import { EOF } from "@/data/eof";
import { isEOF, isOrder, tryParse } from "@/utils/data";
import { useMemo } from "react";
import { ErrorTextbox } from "../error";
import GridList from "../list";
import OrderCard from "./card";

const Container = styled.div`
  gap: 2rem;
  padding: 0.5rem;
  margin-top: 1vh;
  border: 1px solid #ccc;
  border-radius: 10px;
  flex: 1;
`;

const Kline = ({
  messages,
  error,
  onEOF,
}: {
  messages: MessageEvent[];
  error?: any;
  onEOF?: (eof: EOF) => void;
}) => {
  const orders = useMemo(
    () =>
      messages
        .map((event) => {
          const order = tryParse(isOrder, event.data);
          if (order) {
            return order;
          } else {
            const eof = tryParse(isEOF, event.data);
            if (eof) {
              console.log("EOF", eof);
              onEOF?.(eof);
              return null;
            }

            console.error("Failed to parse message", event.data);
            return null;
          }
        })
        .filter((o) => o !== null),
    [messages]
  );

  return (
    <Container>
      {error && (
        <ErrorTextbox
          message={
            typeof error === "string"
              ? error
              : "message" in error
              ? error.message
              : JSON.stringify(error)
          }
        />
      )}
      <GridList
        items={orders}
        cols={["Aggressor", "Price", "Quantity", "Timestamp"]}
        itemRender={(order) => <OrderCard order={order} />}
      />
    </Container>
  );
};

export default Kline;
