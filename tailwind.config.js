/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./ui/**/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [require('daisyui'),],
  daisyui: {
    themes: ["lofi","cyberpunk","retro","wireframe","luxury","night"],
  },
}

