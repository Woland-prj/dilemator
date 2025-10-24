/** @type {import('tailwindcss').Config} */
export default {
  content: ["../internal/view/**/*.{html,js,templ}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark", "cupcake"],
  },
};