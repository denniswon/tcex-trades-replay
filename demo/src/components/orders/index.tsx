import styled from "styled-components";

import type { Order } from "@/data/order";
import { ErrorTextbox } from "../error";
import GridList from "../list";
import OrderCard from "./card";

const Container = styled.div`
  gap: 2rem;
  padding: 0.5rem;
  margin-top: 1vh;
  background: #f0f0f0;
  border-radius: 10px;
  flex: 1;
`;

const Orders = ({ orders, error }: { orders: Order[]; error?: any }) => {
  return (
    <Container>
      {error && (
        <ErrorTextbox
          inverse
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

export default Orders;
