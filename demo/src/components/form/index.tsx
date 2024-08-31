import { useTheme } from "@/contexts/theme";
import { useFormValidate } from "@/hooks/useFormValidate";
import { Theme } from "@uireact/foundation";
import { UiReactViewRotating } from "@uireact/framer-animations";
import { UiIcon } from "@uireact/icons";
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
import FileUpload from "./upload";

export type FormInputData = InferType<typeof schema> & {
  file: File;
  granularity: number;
};

const FormInput = styled.input<{ theme: Theme; disabled: boolean }>`
  height: ${(props: any) => (props.size === "large" ? "36px" : "26px")};
  margin-top: 10px;
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
  margin: 0.4rem 0 0 0.6rem;
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

export const granularityOptions = [
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
  uploading,
  mode,
  onSetMode,
}: {
  onSubmit: (data: FormInputData) => void;
  onError?: (error: any) => void;
  disabled?: boolean;
  uploading?: boolean;
  mode?: "order" | "kline";
  onSetMode?: (mode: "order" | "kline") => void;
}) => {
  const { theme } = useTheme();
  const resolver = useFormValidate(schema);
  const {
    register,
    handleSubmit,
    formState: { isSubmitting, isLoading, errors },
  } = useForm<FormInputData>({
    resolver: resolver,
  });

  const _disabled = useMemo(
    () => !!disabled || isSubmitting || isLoading || uploading,
    [disabled, isSubmitting, isLoading, uploading]
  );

  const [file, setFile] = useState<File | undefined>();
  const [fileError, setFileError] = useState<string | undefined>();
  const [replayRateError, setReplayRateError] = useState<string | undefined>(
    errors?.replay_rate?.message?.toString()
  );

  const [granularity, setGranularity] = useState<number>(
    granularityOptions[0].value
  );

  const clearError = () => {
    setFileError(undefined);
    setReplayRateError(undefined);
  };

  const _onSubmit = async (data: FormInputData) => {
    try {
      clearError();
      onSubmit({ ...data, file, granularity });
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
          <Row style={{ flex: 1, marginTop: "0.6rem" }}>
            <FileUpload
              disabled={_disabled}
              onError={(err) => setFileError(err.message)}
              onFileSelected={setFile}
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
                onOptionSelected={setGranularity}
                defaultValue={{ label: "1m", value: 60 }}
              />
            </Row>
          )}
        </Row>

        <FormButton type="submit" disabled={_disabled}>
          Submit
          {uploading && (
            <UiIcon
              icon="LoadingSpinner"
              category="tertiary"
              size="small"
              motion={UiReactViewRotating}
            />
          )}
        </FormButton>
      </StyledForm>
      {fileError && <ErrorTextbox message={fileError} />}
      {replayRateError && <ErrorTextbox message={replayRateError} />}
    </Container>
  );
};

export { Form, FormInput };
