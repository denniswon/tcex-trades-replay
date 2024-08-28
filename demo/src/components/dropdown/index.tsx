import { useTheme } from "@/contexts/theme";
import { Theme } from "@uireact/foundation";
import { useState } from "react";
import styled from "styled-components";

const DropDownContainer = styled.div`
  width: 2rem;
`;

const DropDownHeader = styled.div<{ theme: Theme; disabled?: boolean }>`
  margin-bottom: 0.2em;
  margin-top: 0.6em;
  width: 2rem;
  padding: 0.4rem 0.6em 0.4rem 0.6em;
  border-radius: 4px;
  text-align: center;
  font-size: 0.8rem;
  cursor: ${(props) => (props.disabled ? "not-allowed" : "pointer")};
  border: 1px solid #ccc;
  background: ${(props) =>
    props.disabled ? "#ccc" : props.theme.colors.primary.token_100};
  color: ${(props) => props.theme.colors.fonts.token_100};
`;

const DropDownListContainer = styled.div<{ theme: Theme; disabled?: boolean }>`
  position: absolute;
  z-index: 100;
  background: ${(props) => props.theme.colors.primary.token_100};
  width: 3rem;
`;

const DropDownList = styled.ul`
  padding: 0;
  margin: 0;
  border-radius: 4px;
  border: 1px solid #ccc;
  box-sizing: border-box;
  &:first-child {
    padding-top: 0.2em;
  }
`;

const ListItem = styled.li`
  list-style: none;
  margin-bottom: 0.4em;
  padding: 0.4em 0 0.4em 0;
  font-size: 0.8rem;
  text-align: center;
  &:hover {
    color: #577ae1;
  }
`;

export type DropdownOption<T> = { label: string; value: T };

export default function DropDown<T>({
  defaultValue,
  options,
  onOptionSelected,
  disabled = false,
}: {
  defaultValue?: DropdownOption<T>;
  options: DropdownOption<T>[];
  onOptionSelected?: (selected: T) => void;
  disabled?: boolean;
}) {
  const { theme } = useTheme();
  const [isOpen, setIsOpen] = useState(false);
  const [selectedOption, setSelectedOption] =
    useState<DropdownOption<T> | null>(defaultValue ?? null);

  const toggling = () => {
    if (disabled) {
      return;
    }
    setIsOpen(!isOpen);
  };

  const _onOptionSelected = (value: DropdownOption<T>) => () => {
    setSelectedOption(value);
    setIsOpen(false);
    onOptionSelected?.(value.value);
  };

  return (
    <DropDownContainer>
      <DropDownHeader onClick={toggling} disabled={disabled} theme={theme}>
        {selectedOption.label || "-"}
      </DropDownHeader>
      {isOpen && (
        <DropDownListContainer theme={theme}>
          <DropDownList>
            {options.map((option, i) => (
              <ListItem
                onClick={_onOptionSelected(option)}
                key={`${option}:${i}:${Math.random()}`}
              >
                {option.label}
              </ListItem>
            ))}
          </DropDownList>
        </DropDownListContainer>
      )}
    </DropDownContainer>
  );
}
