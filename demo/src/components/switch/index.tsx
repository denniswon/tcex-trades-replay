import { UiText } from "@uireact/text";

import { ChangeEvent } from "react";
import styled from "styled-components";

const Label = styled.label`
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  margin: 0 0 0.5rem 0.7rem;
`;

const Switch = styled.div`
  position: relative;
  width: 28px;
  height: 8px;
  background: #b3b3b3;
  border-radius: 8px;
  padding: 4px;
  transition: 300ms all;

  &:before {
    transition: 300ms all;
    content: "";
    position: absolute;
    width: 10px;
    height: 10px;
    border-radius: 10px;
    top: 50%;
    left: 4px;
    background: white;
    transform: translate(0, -50%);
  }
`;

const Input = styled.input`
  opacity: 0;
  position: absolute;

  &:checked + ${Switch} {
    background: #577ae1;

    &:before {
      transform: translate(18px, -50%);
    }
  }
`;

const LabelText = styled(UiText)`
  margin-bottom: 2px;
`;

const ToggleSwitch = ({
  label,
  labelPosition = "start",
  checked,
  onToggle,
}: {
  label?: string;
  labelPosition?: "start" | "end";
  checked: boolean;
  onToggle?: (checked: boolean) => void;
}) => {
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    onToggle?.(e.target.checked);
  };

  return (
    <Label>
      {labelPosition === "start" && label && (
        <LabelText size="small" category="tertiary">
          {label}
        </LabelText>
      )}
      <Input checked={checked} type="checkbox" onChange={handleChange} />
      {labelPosition === "end" && label && (
        <LabelText size="small" category="tertiary">
          {label}
        </LabelText>
      )}
      <Switch />
    </Label>
  );
};

export default ToggleSwitch;
