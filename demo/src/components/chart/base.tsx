import { Theme } from "@uireact/foundation";
import styled from "styled-components";

export const KlineContainer = styled.div`
  display: flex;
  flex-direction: column;
  border-radius: 2px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  padding: 2rem;
  width: 100%;
  height: 65vh;
`;

export const KLineChart = styled.div`
  display: flex;
  flex: 1;
  width: 100%;
  height: 100%;
`;

export const KlineChartMenuContainer = styled.div<{ theme?: Theme }>`
  display: flex;
  flex-direction: row;
  align-items: center;
  font-size: 12px;
  margin-top: 36px;
`;

export const KlineChartMenuButton = styled.button<{
  theme?: Theme;
  isactive?: "true" | "false";
}>`
  cursor: pointer;
  background-color: ${({ isactive, theme }) =>
    isactive === "true" ? theme.colors.primary.token_100 : undefined};
  border: 1px solid #ccc;
  justify-content: center;
  border-radius: 2px;
  margin-right: 12px;
  padding: 6px 12px;
  font-size: 12px;
  color: ${({ isactive, theme }) =>
    isactive === "true" ? "white" : theme.colors.primary.token_100};
`;
