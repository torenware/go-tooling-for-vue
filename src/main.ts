import { createApp } from 'vue';
import App from './App.vue';

const mountPoint = '#app';
const mount = document.querySelector(mountPoint);
const ep = mount?.getAttribute('data-entryp');
console.log(ep ? `ep = ${ep}` : 'no ep');

createApp(App).mount(mountPoint);
