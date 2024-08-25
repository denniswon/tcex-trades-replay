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

// Validate this value with a custom type guard (extend to your needs)
export const isOrder = (o: any): o is Order => {
  return (
    "price" in o && "quantity" in o && "aggressor" in o && "timestamp" in o
  );
};
