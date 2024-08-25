import type { Order } from "@/data/order";
import { UiCard } from "@uireact/card";
import { ColorCategory } from "@uireact/foundation";
import { UiText, UiTextProps } from "@uireact/text";
import styled from "styled-components";
import { Col, Row } from "..";

const OrderRow = styled(Row)`
  flex: 1;
  width: 100%;
  min-width: 80vh;
  overflow-x: scroll;
`;

const OrderCard = ({ order }: { order: Order }) => {
  return (
    <UiCard category="fonts" padding={{ all: "two" }}>
      <OrderRow>
        <OrderInfo k="aggressor" v={order.aggressor} />
        <OrderInfo k="price" v={order.price} />
        <OrderInfo k="quantity" v={order.quantity} />
        <OrderInfo k="timestamp" v={order.timestamp} />
      </OrderRow>
    </UiCard>
  );
};

export default OrderCard;

const OrderInfo = ({ k, v }: { k: string; v: any } & UiTextProps) => {
  let category: ColorCategory = "fonts";
  switch (k) {
    case "aggressor":
      switch (v) {
        case "ask":
          category = "positive";
          break;
        case "bid":
          category = "negative";
          break;
        default:
          break;
      }
      break;
    default:
      break;
  }
  return (
    <Col style={{ width: "25%" }}>
      <UiText size="small" inverseColoration category={category}>
        {v}
      </UiText>
    </Col>
  );
};
