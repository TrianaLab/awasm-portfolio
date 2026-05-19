import { mount } from 'svelte';
import App from './components/App.svelte';
import './app.css';

const target = document.getElementById('app');
if (!target) {
  throw new Error('missing #app mount point');
}

mount(App, { target });
