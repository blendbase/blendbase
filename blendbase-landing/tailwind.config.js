module.exports = {
  content: ["./pages/**/*.{js,ts,jsx,tsx}", "./components/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        brand: {
          50: "#dee5fc",
          100: "#b6c4f7",
          200: "#8da3f1",
          300: "#6583ec",
          400: "#3c62e6",
          500: "#1441e1",
          600: "#1035b8",
          700: "#0d2a8f",
          800: "#091e66",
          900: "#05123d"
        },
        blue: {
          50: "#dee5fc",
          100: "#b6c4f7",
          200: "#8da3f1",
          300: "#6583ec",
          400: "#3c62e6",
          500: "#1441e1",
          600: "#1035b8",
          700: "#0d2a8f",
          800: "#091e66",
          900: "#05123d"
        },
        pink: {
          50: "#ffffff",
          100: "#fbd6f1",
          200: "#f8ade3",
          300: "#f485d5",
          400: "#f15cc7",
          500: "#ed33b9",
          600: "#cb299e",
          700: "#a91e82",
          800: "#871467",
          900: "#65094b"
        },
        teal: {
          100: "#E6FFFA",
          200: "#B2F5EA",
          300: "#81E6D9",
          400: "#4FD1C5",
          500: "#3ABAB4",
          600: "#319795",
          700: "#2C7A7B",
          800: "#285E61",
          900: "#234E52"
        }
      },
      boxShadow: {
        xs: "0 0 0 1px rgba(0, 0, 0, 0.16)",
        sm: "0 1px 2px 0 rgba(0, 0, 0, 0.16)",
        default: "0 1px 3px 0 rgba(0, 0, 0, 0.12), 0 1px 2px 0 rgba(0, 0, 0, 0.03)",
        md: "0 4px 6px -1px rgba(0, 0, 0, 0.04), 0 2px 4px -1px rgba(0, 0, 0, 0.03)",
        lg: "0 10px 15px -3px rgba(0, 0, 0, 0.04), 0 4px 6px -2px rgba(0, 0, 0, 0.02)",
        xl: "0 20px 25px -5px rgba(0, 0, 0, 0.12), 0 10px 10px -5px rgba(0, 0, 0, 0.02)",
        "2xl": "0 25px 50px -12px rgba(0, 0, 0, 0.15)",
        inner: "inset 0 2px 4px 0 rgba(0, 0, 0, 0.04)",
        outline: "0 0 0 3px rgba(66, 153, 225, 0.5)",
        none: "none"
      },
      spacing: {
        "9/16": "56.25%",
        "3/4": "75%",
        "1/1": "100%"
      },
      fontFamily: {
        inter: ["Inter", "sans-serif"]
      },
      fontSize: {
        xs: "0.75rem",
        sm: "0.875rem",
        base: "1rem",
        lg: "1.125rem",
        xl: "1.25rem",
        "2xl": "1.5rem",
        "3xl": "2rem",
        "4xl": "2.625rem",
        "5xl": "3.25rem",
        "6xl": "5.5rem"
      },
      inset: {
        "1/2": "50%",
        full: "100%"
      },
      letterSpacing: {
        tighter: "-0.02em",
        tight: "-0.01em",
        normal: "0",
        wide: "0.01em",
        wider: "0.02em",
        widest: "0.4em"
      },
      lineHeight: {
        none: "1",
        tighter: "1.125",
        tight: "1.25",
        snug: "1.375",
        normal: "1.5",
        relaxed: "1.625",
        loose: "2",
        3: ".75rem",
        4: "1rem",
        5: "1.2rem",
        6: "1.5rem",
        7: "1.75rem",
        8: "2rem",
        9: "2.25rem",
        10: "2.5rem"
      },
      minWidth: {
        10: "2.5rem",
        48: "12rem"
      },
      opacity: {
        90: "0.9"
      },
      scale: {
        98: ".98"
      },
      animation: {
        float: "float 3s ease-in-out infinite"
      },
      keyframes: {
        float: {
          "0%, 100%": { transform: "translateY(0)" },
          "50%": { transform: "translateY(-5%)" }
        }
      },
      zIndex: {
        "-1": "-1"
      }
    }
  },
  plugins: []
};
