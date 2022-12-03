const { fontFamily } = require('tailwindcss/defaultTheme');

/** @type {import('tailwindcss').Config} */
module.exports = {
	content: ['./src/**/*.svelte'],
	theme: {
		extend: {
			fontFamily: {
				sans: ['"Source Sans Pro"', ...fontFamily.sans]
			}
		}
	},
	plugins: []
};
