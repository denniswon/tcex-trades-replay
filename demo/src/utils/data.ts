import { EOF } from "@/data/eof";
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

export const isEOF = (o: any): o is EOF => {
  return "request_id" in o;
};
