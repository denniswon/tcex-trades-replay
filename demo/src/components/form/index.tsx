import { useForm } from "react-hook-form";
import styled from "styled-components";

import { useTheme } from "@/contexts/theme";
import { useFormValidate } from "@/hooks/useFormValidate";
import { Theme } from "@uireact/foundation";
import { normalize, transitions } from "polished";
import { useMemo, useState } from "react";
import { InferType } from "yup";
import { Button } from "..";
import { ErrorTextbox } from "../error";
import schema from "./schema";

export type FormInputSchema = InferType<typeof schema>;

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

const Form = ({
  onSubmit,
  onError,
  disabled,
}: {
  onSubmit: (data: FormInputSchema) => void;
  onError?: (error: any) => void;
  disabled?: boolean;
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

  const clearError = () => {
    setFilenameError(undefined);
    setReplayRateError(undefined);
  };

  const _onSubmit = async (data: FormInputSchema) => {
    try {
      clearError();
      onSubmit(data);
    } catch (error) {
      console.error(error);
      onError?.(error);
    }
  };

  return (
    <Container>
      <StyledForm onSubmit={handleSubmit(_onSubmit)}>
        <FormInput
          theme={theme}
          disabled={_disabled}
          style={{
            border: !_disabled && filenameError ? "1px solid red" : undefined,
          }}
          type="text"
          placeholder="Filename (trades.txt)"
          {...register("filename")}
        />
        <FormInput
          theme={theme}
          disabled={_disabled}
          style={{
            border: !_disabled && replayRateError ? "1px solid red" : undefined,
          }}
          type="text"
          placeholder="Replay Rate (60)"
          {...register("replay_rate")}
        />
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
