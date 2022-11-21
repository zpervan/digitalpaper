/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        'main-black': '#171717',
        'main-gray': '#444444',
        'main-red': '#DA0037',
        'main-white': '#EDEDED',
      },
      fontFamily: {
        sourceSerifPro: ['Source Serif Pro'],
      },
    },
  },
  plugins: [],
};
