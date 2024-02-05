/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
      'features/**/*.templ',
      'static/test/*.html'
    ],
    darkMode: 'class',
    theme: {
      extend: {
        fontFamily: {
          mono: ['Courier Prime', 'monospace'],
        },
        animation: {
          'grow': 'grow 250ms ease-in-out forwards',
        },
      },
    },
    plugins: [
    ],
    corePlugins: {
      preflight: true,
    }
  }