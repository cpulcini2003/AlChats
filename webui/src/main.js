// src/main.js

import { createApp } from 'vue';
import App from './App.vue';
import router from './router/index';
// import store from './store';

// import axios from './services/axios.js';

// import './assets/dashboard.css'
// import './assets/main.css'


const app = createApp(App);
// app.config.globalProperties.$axios = axios;
// app.component("ErrorMsg", ErrorMsg);
// app.component("LoadingSpinner", LoadingSpinner);

// app.use(store);

app.use(router); // Use the router

app.mount('#app');


