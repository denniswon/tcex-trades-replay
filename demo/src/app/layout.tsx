import type { Metadata } from "next";
import { Sen } from "next/font/google";

import { ThemeProvider } from "@/contexts/theme";
import { GlobalStyles, StyledComponentsRegistry, ViewWrapper } from "@/styles";
import "./global.css";

const sen = Sen({
  style: "normal",
  subsets: ["latin"],
  variable: "--font-family",
});

export const metadata: Metadata = {
  title: "TCEX order replay server demo",
  description: "TCEX order replay server demo",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={sen.variable}>
        <StyledComponentsRegistry>
          <GlobalStyles />
          <ViewWrapper>
            <ThemeProvider>{children}</ThemeProvider>
          </ViewWrapper>
        </StyledComponentsRegistry>
      </body>
    </html>
  );
}
