/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./web/views/**/*.{html,js}"],
  plugins: [require("@tailwindcss/forms"), require("@tailwindcss/typography")],
  darkMode: "class",
};
// standalone cli: https://github.com/tailwindlabs/tailwindcss/releases
