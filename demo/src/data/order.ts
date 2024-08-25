export type Order = {
  price: string;
  quantity: number;
  aggressor: "ask" | "bid";
  timestamp: number;
};
