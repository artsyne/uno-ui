/**
 * Uno - A lightweight, modern UI component library.
 * Built on Alpine.js and Tailwind CSS.
 *
 * @module uno
 * @version 0.1.0
 * @license MIT
 */

import Alpine from 'alpinejs';
import htmx from 'htmx.org';
import { installAlpine } from './core/i18n';
import tooltip from './components/tooltip/tooltip';
import notify from './components/notify/notify';
import box from './components/box/box';
import valid from './components/valid/valid';
import tab from './components/tab/tab';

window.Alpine = Alpine;
window.htmx = htmx;

const register = (alpine) => {
  if (!alpine) return;
  installAlpine(alpine);
  tooltip(alpine);
  notify(alpine);
  box(alpine);
  valid(alpine);
  tab(alpine);
}

if (Alpine) {
  register(Alpine);
} else {
  document.addEventListener('alpine:init', () => {
    register(window.Alpine);
  });
}

// Alpine.start(); // Let the user or Alpine CDN start it
