import './app.css'
import App from './App.svelte'
import 'dotenv'

const app = new App({
  target: document.getElementById('app'),
});

export default app
