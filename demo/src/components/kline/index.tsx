import styled from "styled-components";

import { EOF } from "@/data/eof";
import { isEOF, isKline, tryParse } from "@/utils/data";
import { useMemo } from "react";
import { KlineChartView } from "../chart";
import { ErrorTextbox } from "../error";

const Container = styled.div`
  gap: 2rem;
  padding: 2rem;
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
  const klineData = useMemo(
    () =>
      messages
        .map((event) => {
          if ("code" in JSON.parse(event.data || null)) {
            console.log(event.data);
            return null;
          }

          const kline = tryParse(isKline, event.data);
          if (kline) {
            return kline;
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
        .filter((o) => o !== null)
        .reduce((r, o, i) => {
          if (r.length === 0) {
            r.push(o);
          } else {
            const last = r[r.length - 1];
            if (last.timestamp < o.timestamp || i === r.length - 1) {
              r.push(o);
            }
          }
          return r;
        }, []),
    [messages]
  );

  const _granularity = useMemo(() => {
    if (klineData.length > 0) {
      return klineData[0].granularity;
    }
    return 0;
  }, [klineData]);

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
      <KlineChartView data={klineData} />
    </Container>
  );
};

export default Kline;
