const esbuild = require('esbuild');
const fs = require('fs');
const path = require('path');

const isDev = process.argv.includes('--dev');

// Paths
const projectRoot = path.resolve(__dirname, '../../');
const srcDir = path.join(projectRoot, 'website/static/src/js');
const distDir = path.join(projectRoot, 'website/static/dist/js');
const unoSrc = path.join(projectRoot, 'uno/dist/uno.umd.js');
const unoDest = path.join(distDir, 'uno.umd.js');

// Ensure dist dir exists
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true });
}

// Copy Uno
function copyUno() {
  try {
    if (fs.existsSync(unoSrc)) {
      fs.copyFileSync(unoSrc, unoDest);
      console.log(`[Build] Copied Uno UI to ${unoDest}`);
    } else {
      console.warn(`[Build] Warning: Uno UI source not found at ${unoSrc}`);
      console.warn(`[Build] Please ensure 'uno/dist/uno.umd.js' exists.`);
    }
  } catch (e) {
    console.error(`[Build] Failed to copy Uno UI: ${e.message}`);
  }
}

// Build JS
async function build() {
  console.log(
    `[Build] Starting frontend build (Mode: ${isDev ? 'Dev' : 'Prod'})...`,
  );

  try {
    // 1. Bundle JS
    await esbuild.build({
      entryPoints: [path.join(srcDir, 'core/init.js')],
      bundle: true,
      outfile: path.join(distDir, 'bundle.js'),
      minify: true,
      sourcemap: isDev,
      platform: 'browser',
      target: ['es2020'],
      logLevel: 'info',
      define: {
        'process.env.NODE_ENV': isDev ? '"development"' : '"production"',
      },
    });
    console.log(
      `[Build] JS bundle built successfully at ${path.join(distDir, 'bundle.js')}`,
    );

    // 2. Copy Uno UI
    copyUno();
  } catch (e) {
    console.error('[Build] Build failed:', e);
    process.exit(1);
  }
}

build();
