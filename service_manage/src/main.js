import { createApp } from 'vue'
import 'ant-design-vue/dist/antd.less';
import './assets/css/common.less' //全局公共样式
import './assets/css/mobile.less' //全局手机端样式
// import 'default-passive-events'
import App from './App.vue'
import Antd from 'ant-design-vue';
import router from './router/index'
import store from './store'

const app = createApp(App)
app.use(Antd)
app.use(router)
app.use(store)
app.mount('#app')
