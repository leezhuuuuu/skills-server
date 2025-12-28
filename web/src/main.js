import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import 'highlight.js/styles/github-dark.css' // 代码高亮样式

createApp(App).use(router).mount('#app')
