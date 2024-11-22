import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index'
import store from './store'
import vant from 'vant';
import 'vant/lib/index.css';
import './assets/icon/iconfont.css' //图标样式(用阿里矢量图,创建项目后覆盖)
import './assets/css/common.less' //全局公共样式

const app = createApp(App)
app.use(router)
app.use(store)
app.use(vant)
app.mount('#app')
