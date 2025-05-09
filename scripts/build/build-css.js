const { execSync } = require('child_process');
const fs = require('fs');
const os = require('os');

console.log('Starting CSS build process...');

// Step 1: Generate tailwind.css with the new CLI
console.log('Step 1: Generating tailwind.css...');
execSync(
  'npx @tailwindcss/cli -i ./website/static/src/css/tailwind.css -o ./website/static/dist/css/tailwind.css',
  { stdio: 'inherit' },
);

// Step 2: Combine CSS files
console.log('Step 2: Combining CSS files...');
const isWindows = os.platform() === 'win32';
const combineCommand = isWindows
  ? 'type .\\website\\static\\src\\css\\diary-theme.css .\\website\\static\\dist\\css\\tailwind.css > .\\website\\static\\dist\\css\\combined.css'
  : 'cat ./website/static/src/css/diary-theme.css ./website/static/dist/css/tailwind.css > ./website/static/dist/css/combined.css';

execSync(combineCommand, { stdio: 'inherit' });

// Step 3: Minify combined CSS file
console.log('Step 3: Minifying CSS...');
execSync(
  'npx postcss ./website/static/dist/css/combined.css -o ./website/static/dist/css/idiary.min.css --no-map',
  { stdio: 'inherit' },
);

// Clean up
console.log('Cleaning up temporary files...');
fs.unlinkSync('./website/static/dist/css/combined.css');
fs.unlinkSync('./website/static/dist/css/tailwind.css');

console.log('CSS build process completed successfully.');
