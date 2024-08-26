import { EOF } from "@/data/eof";
import { Kline } from "@/data/kline";
import { Order } from "@/data/order";

export const tryParse = <T>(
  guard: (o: any) => o is T,
  json: string
): T | null => {
  try {
    const parsed = JSON.parse(json);
    if (guard(parsed)) {
      return parsed;
    }
  } catch (_e) {}
  return null;
};

export const isOrder = (o: any): o is Order => {
  return (
    "price" in o && "quantity" in o && "aggressor" in o && "timestamp" in o
  );
};

export const isKline = (k: any): k is Kline => {
  return (
    "low" in k &&
    "high" in k &&
    "open" in k &&
    "close" in k &&
    "volume" in k &&
    "granularity" in k &&
    "timestamp" in k
  );
};

export const isEOF = (o: any): o is EOF => {
  return "request_id" in o;
};
