import { useCallback, useEffect, useState } from "react";

export function useWs<T>({
  url,
  onError,
}: {
  url: string;
  onError?: (error: any) => void;
}) {
  const [ws, setWs] = useState<WebSocket | undefined>(undefined);
  const [messages, setMessages] = useState<MessageEvent[]>([]);
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
      setMessages((messages) => [...messages, event]);
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
}
