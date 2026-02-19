import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ToastPlugin from 'vue-toast-notification';
import { Api } from './api'
import Vue3Lottie from 'vue3-lottie'
import "./main.css"
import 'vue-toast-notification/dist/theme-bootstrap.css';

const app = createApp(App)

app.use(router)
app.use(Vue3Lottie)

const api = new Api()
app.use(api)
app.use(ToastPlugin)

app.mount('#app')
