import forms from "@tailwindcss/forms";
import plugin from "tailwindcss/plugin";

/** @type {import("tailwindcss").Config} */
const config = {
  content: { files: ["internal/web/**/*.templ"] },
  plugins: [
    forms({ strategy: "class" }),
    plugin(({ addVariant }) => {
      addVariant("hocus", ["&:hover", "&:focus"]);
      addVariant("list-outer-hocus", ["&:has(a:focus)", "&:has(a:hover)"]);
    })
  ],
  theme: {
    extend: {
      colors: {
        border: "#1c1c23",
        inactive: "#86868c"
      }
    }
  }
};

export default config;