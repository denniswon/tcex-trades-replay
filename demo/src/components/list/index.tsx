import { UiGrid, UiGridItem } from "@uireact/grid";

import { UiText } from "@uireact/text";
import styled from "styled-components";
import { Box, Col, Row } from "..";

interface GridListProps<T> {
  items: T[];
  cols: string[];
  itemRender: (item: T) => React.ReactNode;
}

const TitleRow = styled(Row)`
  flex: 1;
  width: 100%;
  min-width: 80vh;
  overflow-x: scroll;
  margin-bottom: 1rem;
`;

function GridList<T>({ items, cols, itemRender }: GridListProps<T>) {
  return (
    <StyledContainer>
      <StyledGrid cols={cols.length} rows={items.length} justifyItems="center">
        <UiGridItem cols={4}>
          <TitleRow>
            {...cols.map((col) => (
              <Col key={col} style={{ width: "25%" }}>
                <UiText size="small" category="tertiary">
                  {col}
                </UiText>
              </Col>
            ))}
          </TitleRow>
        </UiGridItem>
        {...items.map((item) => (
          <StyledRow cols={cols.length}>{itemRender(item)}</StyledRow>
        ))}
      </StyledGrid>
    </StyledContainer>
  );
}

export default GridList;

const StyledContainer = styled(Box)`
  max-height: 65vh;
  overflow-y: scroll;
  justify-content: center;
  overflow-x: scroll;
`;

const StyledGrid = styled(UiGrid)`
  overflow-x: scroll;
`;

const StyledRow = ({
  children,
  cols,
}: {
  children: React.ReactNode;
  cols: number;
}) => <UiGridItem cols={cols}>{children}</UiGridItem>;
