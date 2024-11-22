import { createApp } from 'vue'
import 'ant-design-vue/dist/antd.less';
import './assets/icon/iconfont.css' //图标样式(用阿里矢量图,创建项目后覆盖)
import './assets/css/common.less' //全局公共样式
import './assets/css/mobile.less' //全局手机端样式
import App from './App.vue'
import Antd from 'ant-design-vue';
import router from './router';
import store from './store';

const app = createApp(App)
app.use(Antd)
app.use(router)
app.use(store)
app.mount('#app')


// ----自定义指令-----
// 拖拽指令
app.directive('drag', {
    mounted(el, binding, vnode) {        
        el.addEventListener('drop', (event) => {//文件落下            
            event.stopPropagation();
            event.preventDefault();
            if (event.dataTransfer && event.dataTransfer.items) { //判断是否有拖拽进来的文件                
                event = Array.from(event.dataTransfer.items); //将拖拽进来的items类数组对象，转成真正的数组
                let files = [];
                event.map((i) => {
                    files.push(i.getAsFile());
                });
                binding.value(files)
            }
        });
    }
})

// 复制指令
app.directive('paste', {
    mounted(el, binding, vnode) {
        el.addEventListener('paste', function (event) { //这里直接监听元素的粘贴事件
            if (event.clipboardData.files[0]) {                
                event.preventDefault(); // 阻止默认的复制事件
                binding.value([event.clipboardData.files[0]])
            }
        })
    }
})
