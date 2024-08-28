import { useTheme } from "@/contexts/theme";
import { Chart, dispose, init, KLineData } from "klinecharts";
import { useEffect, useRef, useState } from "react";
import styled from "styled-components";
import {
  KLineChart,
  KlineChartMenuButton,
  KlineChartMenuContainer,
  KlineContainer,
} from "./base";

const mainIndicators = ["MA", "EMA", "SAR"];
const subIndicators = ["VOL", "MACD", "KDJ"];

export const KlineChartView = ({ data }: { data: KLineData[] }) => {
  const chart = useRef<Chart | null>();
  const paneId = useRef<string>("");
  const { theme } = useTheme();

  const [mainIndicator, setMainIndicator] = useState<string>();
  const [subIndicator, setSubIndicator] = useState<string>();

  useEffect(() => {
    chart.current = init("indicator-k-line");
    paneId.current = chart.current?.createIndicator("VOL", false) as string;
    chart.current?.applyNewData(data);
    return () => {
      dispose("indicator-k-line");
    };
  }, [data]);

  return (
    <Container>
      <KlineContainer>
        <KLineChart id="indicator-k-line" />
        <KlineChartMenuContainer theme={theme}>
          <span style={{ paddingRight: 10, marginRight: 6 }}>
            Main Indicators
          </span>
          {mainIndicators.map((type) => {
            return (
              <KlineChartMenuButton
                theme={theme}
                isactive={type === mainIndicator ? "true" : "false"}
                key={type}
                onClick={(_) => {
                  setMainIndicator(type);
                  chart.current?.createIndicator(type, false, {
                    id: "candle_pane",
                  });
                }}
              >
                {type}
              </KlineChartMenuButton>
            );
          })}
          <span style={{ paddingRight: 10, paddingLeft: 12, marginRight: 6 }}>
            Sub-indicators
          </span>
          {subIndicators.map((type) => {
            return (
              <KlineChartMenuButton
                isactive={type === subIndicator ? "true" : "false"}
                theme={theme}
                key={type}
                onClick={(_) => {
                  setSubIndicator(type);
                  chart.current?.createIndicator(type, false, {
                    id: paneId.current,
                  });
                }}
              >
                {type}
              </KlineChartMenuButton>
            );
          })}
        </KlineChartMenuContainer>
      </KlineContainer>
    </Container>
  );
};

const Container = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
  flex-wrap: wrap;
`;
