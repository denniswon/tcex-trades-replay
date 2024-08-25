import { UiButton } from "@uireact/button";
import styled from "styled-components";

interface BoxProps {
  align?: "center" | "left" | "right";
}
const Box = styled.div<BoxProps>`
  display: flex;
  justify-content: ${({ align }) =>
    align === "right"
      ? "flex-end"
      : align === "center"
      ? "center"
      : "flex-start"};
  align-items: ${({ align }) =>
    align === "right"
      ? "flex-end"
      : align === "center"
      ? "center"
      : "flex-start"};
  gap: 1rem;
  padding: 1rem;
  border-radius: 10px;
  flex: 1;
`;

const Row = styled.div`
  display: flex;
  flex-direction: row;
  gap: 1rem;
`;

const Col = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const Main = styled.div`
  padding: 2vh;
`;

const Container = styled.div`
  max-width: 80%;
  margin: 0 auto;
  gap: 1rem;
  padding-top: 1rem;
`;

const Text = styled.p`
  padding: 4px;
  margin: 4px;
  color: #fff;
  font-size: 0.8rem;
`;

const Input = styled.input`
  border: 1px solid #ccc;
  border-radius: 4px;
  padding: 12px;

  &:focus {
    outline: 2px solid blue;
  }

  &:disabled {
    background: #f0f0f0;
    cursor: not-allowed;
  }
`;

const Button = styled(UiButton)`
  border-radius: 4px;
  border: none;
  padding: 0.5rem 1rem;
  margin: 0rem 0.3rem;
  background: #577ae1;
  color: white;
  &:hover {
    background: #3762e6;
    color: white;
    cursor: pointer;
  }

  &:active {
    background: #10329a;
    color: white;
  }

  &:disabled {
    background: #a0a5b3;
    color: white;
    pointer-events: none;
  }
`;

export { Box, Button, Col, Container, Input, Main, Row, Text };
