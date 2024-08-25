import { tryParse } from "@/utils/data";
import { useCallback, useEffect, useState } from "react";

export const useWs = <T,>({
  url,
  guard,
  onError,
}: {
  url: string;
  guard: (o: any) => o is T;
  onError?: (error: any) => void;
}) => {
  const [ws, setWs] = useState<WebSocket | undefined>(undefined);
  const [messages, setMessages] = useState<T[]>([]);
  const [open, setOpen] = useState(false);
  const [firstConnect, setFirstConnect] = useState(true);

  useEffect(() => {
    if (!url) {
      return;
    }
    const _ws = new WebSocket(url);
    setWs(_ws);
    setFirstConnect(false);
  }, [url]);

  // Keep trying to (re)connect if disconnected after first connect
  useEffect(() => {
    if (!ws || firstConnect || !url || open) {
      return;
    }

    const interval = setInterval(() => {
      if (!ws || open) {
        clearInterval(interval);
        return;
      }

      try {
        const _ws = new WebSocket(url);
        setWs(_ws);
      } catch (error) {
        console.error("WebSocket connection failed", error);
      }
    }, 3000);
    return () => {
      clearInterval(interval);
    };
  });

  useEffect(() => {
    if (!ws || open) {
      return;
    }

    ws.onopen = () => {
      console.log("Connected to server");
      setOpen(true);
    };

    ws.onclose = (event: CloseEvent) => {
      console.log("Disconnected from server", event);
      setOpen(false);
    };

    ws.onerror = (event: any) => {
      if (!ws || !open) {
        return;
      }
      console.error(event);
      onError?.(event);
    };

    ws.onmessage = (event: MessageEvent) => {
      const msg = tryParse<T>(guard, event.data);
      if (msg) {
        setMessages((messages) => [...messages, msg]);
      }
    };
  }, [ws, open, onError]);

  const clearMessages = useCallback(() => {
    setMessages([]);
  }, []);

  return {
    ws,
    open,
    messages,
    clearMessages,
  };
};
