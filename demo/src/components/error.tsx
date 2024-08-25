import { useTheme } from "@/contexts/theme";
import { ContextTheme } from "@/styles/custom-theme";
import { UiText, UiTextProps } from "@uireact/text";
import styled, { CSSProperties } from "styled-components";

const ErrorText = styled(UiText)<{ theme: ContextTheme; color?: string }>`
  color: ${({ theme, color }) => color ?? theme.colors.error.token_100};
`;

const ErrorTextbox = ({
  message,
  textStyle,
  inverse,
  containerStyle,
}: {
  message: string;
  textStyle?: UiTextProps;
  // if inverse true, text color will be regular text color, but border will be red
  inverse?: boolean;
  containerStyle?: CSSProperties;
}) => {
  const { theme } = useTheme();
  return (
    <StyledContainer
      theme={theme}
      style={{
        ...containerStyle,
        ...{
          borderRadius: 8,
          border: inverse
            ? `1px solid ${theme.colors.error.token_100}`
            : undefined,
        },
      }}
    >
      <ErrorText
        theme={theme}
        category="error"
        size="xsmall"
        color={inverse ? theme.colors.tertiary.token_200 : undefined}
        {...textStyle}
      >
        {message}
      </ErrorText>
    </StyledContainer>
  );
};

const StyledContainer = styled.div`
  padding: 1rem;
  margin-bottom: 1rem;
`;

export { ErrorText, ErrorTextbox };
