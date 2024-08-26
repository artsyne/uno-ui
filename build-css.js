const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');

// Step 1: Generate tailwind.css
execSync('npx tailwindcss -i ./static/css/src/tailwind.css -o ./static/css/dist/tailwind.css', { stdio: 'inherit' });

// Step 2: Combine CSS files
const isWindows = os.platform() === 'win32';
const combineCommand = isWindows
  ? 'type .\\static\\css\\src\\diary-theme.css .\\static\\css\\dist\\tailwind.css > .\\static\\css\\dist\\combined.css'
  : 'cat ./static/css/src/diary-theme.css ./static/css/dist/tailwind.css > ./static/css/dist/combined.css';

execSync(combineCommand, { stdio: 'inherit' });

// Step 3: Minify combined CSS file
execSync('npx postcss ./static/css/dist/combined.css -u cssnano -o ./static/css/dist/idiary.min.css --no-map', { stdio: 'inherit' });

// Clean up
fs.unlinkSync('./static/css/dist/combined.css');
fs.unlinkSync('./static/css/dist/tailwind.css');

console.log('CSS build process completed successfully.');