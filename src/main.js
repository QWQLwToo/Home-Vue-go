import { createApp } from 'vue';
import App from './App.vue';
import './style.less';
import { MotionPlugin } from '@vueuse/motion';
import router from './router';
import { loadAndApplyFrontendConfig } from './utils/frontendConfig';

// 加载前端配置
loadAndApplyFrontendConfig();

const app = createApp(App);

app.use(MotionPlugin);
app.use(router);

app.mount('#app');