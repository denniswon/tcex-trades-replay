import { useQuery } from "@tanstack/react-query";
import axios from "axios";

export type UploadHeader = {
  id: string;
  filepath: string;
  size: number;
};

const baseUrl = "http://localhost:8080/v1/upload";

export function useFileUpload({ file }: { file: File | undefined }): {
  header: UploadHeader | undefined;
  isError: boolean;
  isLoading: boolean;
} {
  const { data, isLoading, isError } = useQuery({
    queryFn: async () => {
      const formData = new FormData();
      formData.append("file", file);

      const config = {
        headers: {
          "content-type": "multipart/form-data",
        },
      };

      const response = await axios.post<UploadHeader>(
        baseUrl,
        formData,
        config
      );
      return response.data;
    },
    queryKey: [file],
    enabled: !!file,
  });

  return {
    header: data,
    isError,
    isLoading,
  };
}
