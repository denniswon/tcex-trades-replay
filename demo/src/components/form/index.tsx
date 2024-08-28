import { useTheme } from "@/contexts/theme";
import { useFormValidate } from "@/hooks/useFormValidate";
import { Theme } from "@uireact/foundation";
import { UiText } from "@uireact/text";
import { normalize, transitions } from "polished";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import styled from "styled-components";
import { InferType } from "yup";
import { Button, Row } from "..";
import DropDown from "../dropdown";
import { ErrorTextbox } from "../error";
import Switch from "../switch";
import schema from "./schema";

export type FormInputSchema = InferType<typeof schema> & {
  granularity?: number;
};

const FormInput = styled.input<{ theme: Theme; disabled: boolean }>`
  height: ${(props: any) => (props.size === "large" ? "36px" : "24px")};
  margin: 8px;
  padding: 2px;
  width: 30%;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: ${(props) =>
    props.disabled ? "#ccc" : props.theme.colors.primary.token_100};
`;

const FormButton = styled(Button)`
  ${normalize()}
  ${transitions(["background"], "0.3s")}
  margin: 0.4rem 0 0 0.5rem;
`;

const StyledForm = styled.form`
  padding: 12px;
  border-radius: 10px;
`;

const Container = styled.div`
  gap: 2rem;
  padding: 0.5rem;
  border: 1px solid #ccc;
  border-radius: 10px;
  flex: 1;
`;

const granularityOptions = [
  { label: "1m", value: 60 },
  { label: "5m", value: 300 },
  { label: "15m", value: 900 },
  { label: "1h", value: 3600 },
  { label: "6h", value: 21600 },
  { label: "1d", value: 86400 },
];

const Form = ({
  onSubmit,
  onError,
  disabled,
  mode,
  onSetMode,
}: {
  onSubmit: (data: FormInputSchema) => void;
  onError?: (error: any) => void;
  disabled?: boolean;
  mode?: "order" | "kline";
  onSetMode?: (mode: "order" | "kline") => void;
}) => {
  const { theme } = useTheme();
  const resolver = useFormValidate(schema);
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, isLoading, errors },
  } = useForm({
    resolver: resolver,
  });

  const _disabled = useMemo(
    () => !!disabled || isSubmitting || isLoading,
    [disabled, isSubmitting, isLoading]
  );

  const [filenameError, setFilenameError] = useState<string | undefined>(
    errors?.filename?.message?.toString()
  );
  const [replayRateError, setReplayRateError] = useState<string | undefined>(
    errors?.replay_rate?.message?.toString()
  );

  const [granularity, setGranularity] = useState<number>();

  const clearError = () => {
    setFilenameError(undefined);
    setReplayRateError(undefined);
  };

  const _onSubmit = async (data: FormInputSchema) => {
    try {
      clearError();
      onSubmit({ ...data, granularity });
    } catch (error) {
      console.error(error);
      onError?.(error);
    }
  };

  return (
    <Container>
      <StyledForm onSubmit={handleSubmit(_onSubmit)}>
        <Switch
          label={`Mode: ${mode === "kline" ? "kline" : "order"}`}
          onToggle={(checked) => onSetMode(checked ? "kline" : "order")}
          checked={mode === "kline"}
        />
        <Row>
          <Row style={{ flex: 1 }}>
            <FormInput
              theme={theme}
              disabled={_disabled}
              style={{
                border:
                  !_disabled && filenameError ? "1px solid red" : undefined,
              }}
              type="text"
              placeholder="Filename (trades.txt)"
              {...register("filename")}
            />
            <FormInput
              theme={theme}
              disabled={_disabled}
              style={{
                border:
                  !_disabled && replayRateError ? "1px solid red" : undefined,
              }}
              type="text"
              placeholder="Replay Rate (60)"
              {...register("replay_rate")}
            />
          </Row>
          {mode === "kline" && (
            <Row
              style={{
                gap: "0.5rem",
                alignItems: "center",
                marginRight: "1.5rem",
              }}
            >
              <UiText margin={{ top: "two" }} size="small">
                Granularity
              </UiText>
              <DropDown<number>
                disabled={_disabled}
                options={granularityOptions}
                onOptionSelected={(value) => setGranularity(value)}
                defaultValue={{ label: "1m", value: 60 }}
              />
            </Row>
          )}
        </Row>

        <FormButton type="submit" disabled={_disabled}>
          Submit
        </FormButton>
      </StyledForm>
      {filenameError && <ErrorTextbox message={filenameError} />}
      {replayRateError && <ErrorTextbox message={replayRateError} />}
    </Container>
  );
};

export { Form, FormInput };
