const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');

// Step 1: Generate tailwind.css
execSync(
  'npx tailwindcss -i ./website/static/src/css/tailwind.css -o ./website/static/dist/css/tailwind.css',
  { stdio: 'inherit' },
);

// Step 2: Combine CSS files
const isWindows = os.platform() === 'win32';
const combineCommand = isWindows
  ? 'type .\\website\\static\\src\\css\\diary-theme.css .\\website\\static\\dist\\css\\tailwind.css > .\\website\\static\\dist\\css\\combined.css'
  : 'cat ./website/static/src/css/diary-theme.css ./website/static/dist/css/tailwind.css > ./website/static/dist/css/combined.css';

execSync(combineCommand, { stdio: 'inherit' });

// Step 3: Minify combined CSS file
execSync(
  'npx postcss ./website/static/dist/css/combined.css -u cssnano -o ./website/static/dist/css/idiary.min.css --no-map',
  { stdio: 'inherit' },
);

// Clean up
fs.unlinkSync('./website/static/dist/css/combined.css');
fs.unlinkSync('./website/static/dist/css/tailwind.css');

console.log('CSS build process completed successfully.');
