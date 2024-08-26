const { execSync } = require('child_process');
const fs = require('fs');

// Step 1: Generate tailwind.css
execSync('npx tailwindcss -i ./static/css/src/tailwind.css -o ./static/css/dist/tailwind.css', { stdio: 'inherit' });

// Step 2: Combine and minify CSS files
execSync('cat ./static/css/src/diary-theme.css ./static/css/dist/tailwind.css > ./static/css/dist/combined.css', { stdio: 'inherit' });
execSync('npx postcss ./static/css/dist/combined.css -u cssnano -o ./static/css/dist/idiary.min.css --no-map', { stdio: 'inherit' });

// Clean up
fs.unlinkSync('./static/css/dist/combined.css');
fs.unlinkSync('./static/css/dist/tailwind.css');