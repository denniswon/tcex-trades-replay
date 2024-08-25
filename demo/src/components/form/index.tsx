import { useForm } from "react-hook-form";
import styled from "styled-components";

import { useFormValidate } from "@/hooks/useFormValidate";
import { normalize, transitions } from "polished";
import { useMemo, useState } from "react";
import { Button } from "..";
import { ErrorTextbox } from "../error";
import schema from "./schema";

export type FormInputSchema = typeof schema.shape;

const FormInput = styled.input`
  height: ${(props: any) => (props.size === "large" ? "36px" : "24px")};
  margin: 8px;
  padding: 2px;
  width: 30%;
`;

const FormButton = styled(Button)`
  ${normalize()}
  ${transitions(["background"], "0.3s")}
`;

const StyledForm = styled.form`
  background: white;
  padding: 12px;
  border-radius: 10px;
`;

const Container = styled.div`
  gap: 2rem;
  padding: 0.5rem;
  background: #f0f0f0;
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

  const _onSubmit = async (data) => {
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
          disabled={_disabled}
          style={{
            border: !_disabled && filenameError ? "1px solid red" : undefined,
          }}
          type="text"
          placeholder="Filename"
          {...register("filename")}
        />
        <FormInput
          disabled={_disabled}
          style={{
            border: !_disabled && replayRateError ? "1px solid red" : undefined,
          }}
          type="text"
          placeholder="Replay Rate"
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
