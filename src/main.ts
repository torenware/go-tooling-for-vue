import { createApp } from 'vue';
import App from './App.vue';

type EPLookup = Record<string, typeof App>;

const keys: EPLookup = {};
keys['app'] = App;

const appMounts = document.querySelectorAll('[data-entryp]');
appMounts.forEach((mp) => {
  const ep = mp.getAttribute('data-entryp');
  if (ep && ep in keys) {
    createApp(keys[ep]).mount(mp);
  } else {
    console.log(`${ep}: key was not found.`);
  }
});
