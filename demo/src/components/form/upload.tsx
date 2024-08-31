import { useTheme } from "@/contexts/theme";
import { UiPrimaryButton } from "@uireact/button";
import { Theme } from "@uireact/foundation";
import { ChangeEvent, useRef, useState } from "react";
import styled from "styled-components";
import { Box } from "..";

const FileInput = styled.input`
  display: none;
`;

const FileUpload = ({
  onFileSelected,
  onError,
  disabled,
}: {
  onFileSelected?: (file: File) => void;
  onError?: (error: any) => void;
  disabled?: boolean;
}) => {
  const [selectedFile, setSelectedFile] = useState<File>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const { theme } = useTheme();

  const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
    if (fileInputRef.current && fileInputRef.current.files.length > 0) {
      console.log("File selected", event.target.files);
      const file = fileInputRef.current.files[0];
      if (file.type !== "text/plain") {
        const err = new Error(
          `unsupported file type ${file.type} for ${file.name} expected text/plain`
        );
        console.error(err);

        onError?.(err);
        setSelectedFile(null);

        return;
      }

      console.log("File selected", file);
      setSelectedFile(event.target.files[0]);
      onFileSelected?.(event.target.files[0]);
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current.click();
  };

  return (
    <Container>
      <UploadButton
        theme={theme}
        disabled={disabled}
        onClick={handleUploadClick}
        padding={{ right: "two", left: "two" }}
      >
        {selectedFile
          ? selectedFile.name.length > 32
            ? selectedFile.name.slice(0, 32) + "..."
            : selectedFile.name
          : "Choose File"}
      </UploadButton>
      <FileInput type="file" ref={fileInputRef} onChange={handleFileChange} />
    </Container>
  );
};

const UploadButton = styled(UiPrimaryButton)<{ theme?: Theme }>`
  white-space: nowrap;
  overflow: hidden;
  background-color: ${(props) =>
    props.theme.colors.primary.token_100 || "#000"};
  text-overflow: ellipsis;
  color: ${(props) => props.theme.colors.fonts.token_100 || "#fff"};
  font-size: small;
  font-weight: 500;
  max-width: 20vh;
  border-color: ${(props) => props.theme.colors.fonts.token_100 || "#fff"};
  border-radius: 4px;
  padding: 0.4rem 0.8rem;

  &:hover {
    background-color: ${(props) =>
      props.theme.colors.fonts.token_100 || "#fff"};
    color: ${(props) => props.theme.colors.primary.token_100 || "#000"};
  }
`;

const Container = styled(Box)`
  margin: 0;
  padding: 0.6rem;
  flex: 0;
`;

export default FileUpload;
