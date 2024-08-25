"use client";

import {
  type ContextTheme,
  DarkTheme,
  LightTheme,
} from "@/styles/custom-theme";
import { ThemeColor } from "@uireact/foundation";
import React, { createContext, useContext, useState } from "react";

interface ThemeContextType {
  theme: ContextTheme;
  type: ThemeColor;
  toggleTheme: () => void;
}

// for placeholder function
const nil = () => {};
const ThemeContext = createContext<ThemeContextType>({
  theme: LightTheme,
  type: ThemeColor.light,
  toggleTheme: nil,
});

const ThemeProvider = ({
  children,
  theme: _theme = ThemeColor.light,
}: {
  children: React.ReactNode;
  theme?: ThemeColor;
}) => {
  const [type, setType] = useState<ThemeColor>(_theme);
  const [theme, setTheme] = useState<ContextTheme>(
    _theme === ThemeColor.light ? LightTheme : DarkTheme
  );

  const toggleTheme = () => {
    setType((prevType) =>
      prevType === ThemeColor.light ? ThemeColor.dark : ThemeColor.light
    );
    setTheme((prevTheme) =>
      prevTheme === LightTheme ? DarkTheme : LightTheme
    );
  };

  return (
    <ThemeContext.Provider value={{ theme, type, toggleTheme }}>
      {children}
    </ThemeContext.Provider>
  );
};

const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
};

export { ThemeProvider, useTheme };
