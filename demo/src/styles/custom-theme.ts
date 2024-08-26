import { Coloration, Theme } from "@uireact/foundation";

export const CustomTheme: Theme = {
  name: "theme",
  light: {
    error: {
      token_10: "#e77979",
      token_50: "#dc3d3d",
      token_100: "#c72424",
      token_150: "#851818",
      token_200: "#6f1414",
    },
    warning: {
      token_10: "#ebc184",
      token_50: "#e1a045",
      token_100: "#d38a22",
      token_150: "#8d5c17",
      token_200: "#754d13",
    },
    negative: {
      token_10: "#e47373",
      token_50: "#d93a3a",
      token_100: "#c02525",
      token_150: "#801919",
      token_200: "#6a1515",
    },
    positive: {
      token_10: "#5dd353",
      token_50: "#3abb30",
      token_100: "#319c28",
      token_150: "#20681b",
      token_200: "#1b5716",
    },
    fonts: {
      token_10: "#ffffff",
      token_50: "#f9f9f9",
      token_100: "#f2f2f2",
      token_150: "#bababa",
      token_200: "#adadad",
    },
    primary: {
      token_10: "#302b2b",
      token_50: "#282424",
      token_100: "#1b1818",
      token_150: "#131111",
      token_200: "#0f0d0d",
    },
    secondary: {
      token_10: "#ffffff",
      token_50: "#f9fddb",
      token_100: "#edf990",
      token_150: "#daf313",
      token_200: "#b9cf0b",
    },
    tertiary: {
      token_10: "#ffffff",
      token_50: "#e2e3d9",
      token_100: "#c1c3ae",
      token_150: "#8a8e68",
      token_200: "#737657",
    },
  },
  dark: {
    error: {
      token_10: "#f08989",
      token_50: "#e84646",
      token_100: "#e01c1c",
      token_150: "#961212",
      token_200: "#7d0f0f",
    },
    warning: {
      token_10: "#f9e5a5",
      token_50: "#f5cf57",
      token_100: "#f1c023",
      token_150: "#ad870b",
      token_200: "#917009",
    },
    negative: {
      token_10: "#f4b0b0",
      token_50: "#ea6666",
      token_100: "#e33535",
      token_150: "#a41717",
      token_200: "#891313",
    },
    positive: {
      token_10: "#98e791",
      token_50: "#5eda53",
      token_100: "#39cf2d",
      token_150: "#268a1e",
      token_200: "#207319",
    },
    fonts: {
      token_10: "#302b2b",
      token_50: "#282424",
      token_100: "#1b1818",
      token_150: "#131111",
      token_200: "#0f0d0d",
    },
    primary: {
      token_10: "#ffffff",
      token_50: "#f9f9f9",
      token_100: "#f2f2f2",
      token_150: "#bababa",
      token_200: "#adadad",
    },
    secondary: {
      token_10: "#f0f064",
      token_50: "#eaea26",
      token_100: "#cece14",
      token_150: "#8a8a0d",
      token_200: "#73730b",
    },
    tertiary: {
      token_10: "#ffffff",
      token_50: "#e1e2d8",
      token_100: "#c1c3ae",
      token_150: "#8a8d69",
      token_200: "#737657",
    },
  },
  spacing: {
    one: "0.1rem",
    two: "0.2rem",
    three: "0.6rem",
    four: "1rem",
    five: "1.2rem",
    six: "1.5rem",
    seven: "3rem",
  },
  texts: {
    font: "Sens, Helvetica, Arial, sans-serif",
  },
  sizes: {
    texts: {
      xsmall: "0.75rem",
      small: "1rem",
      regular: "1.2rem",
      large: "2.5rem",
      xlarge: "3.75rem",
    },
    headings: {
      level1: "4rem",
      level2: "3rem",
      level3: "2.5rem",
      level4: "2rem",
      level5: "1.5rem",
      level6: "1rem",
    },
  },
};

export type ContextTheme = Pick<Theme, "texts" | "sizes" | "spacing"> & {
  colors: Coloration;
};

export const LightTheme: ContextTheme = {
  colors: CustomTheme.light,
  spacing: {
    one: "0.1rem",
    two: "0.2rem",
    three: "0.6rem",
    four: "1rem",
    five: "1.2rem",
    six: "1.5rem",
    seven: "3rem",
  },
  texts: {
    font: "Sens, Helvetica, Arial, sans-serif",
  },
  sizes: {
    texts: {
      xsmall: "0.75rem",
      small: "1rem",
      regular: "1.2rem",
      large: "2.5rem",
      xlarge: "3.75rem",
    },
    headings: {
      level1: "4rem",
      level2: "3rem",
      level3: "2.5rem",
      level4: "2rem",
      level5: "1.5rem",
      level6: "1rem",
    },
  },
};

export const DarkTheme: ContextTheme = {
  colors: CustomTheme.dark,
  spacing: {
    one: "0.1rem",
    two: "0.2rem",
    three: "0.6rem",
    four: "1rem",
    five: "1.2rem",
    six: "1.5rem",
    seven: "3rem",
  },
  texts: {
    font: "Sens, Helvetica, Arial, sans-serif",
  },
  sizes: {
    texts: {
      xsmall: "0.75rem",
      small: "1rem",
      regular: "1.2rem",
      large: "2.5rem",
      xlarge: "3.75rem",
    },
    headings: {
      level1: "4rem",
      level2: "3rem",
      level3: "2.5rem",
      level4: "2rem",
      level5: "1.5rem",
      level6: "1rem",
    },
  },
};
