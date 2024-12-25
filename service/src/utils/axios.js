import axios from 'axios'
import { localGet, localRemove } from './index'
import { message } from 'ant-design-vue';
import { websocketClose } from "@/socket";

axios.defaults.baseURL = 'http://'+import.meta.env.VITE_BASE_URL + 'service';
axios.defaults.timeout = 60000;

//请求拦截
axios.interceptors.request.use(
    config => {
        // 在发送请求之前做些什么
        config.headers['token'] = localGet('token') || '';
        return config;
    },
    error => {
        // 对请求错误做些什么
        return Promise.reject(error);
    }
);

// 添加响应拦截器
axios.interceptors.response.use(
    response => {
        // 这个状态码是和后端约定的code
        const code = response.data.code
        switch (code) {
            case 200:
                return Promise.resolve(response);
                break;
            case 301:
                websocketClose();
                localRemove();
                window.location.reload();
                message.error('登录过期');
                return Promise.reject(response);
                break;

            default:
                message.error(response.data.msg);
                return Promise.reject(response);
                break;
        }
    },
    error => {
        return Promise.reject(error);
    }
);
export default {
    post(url, data, msg = true) {
        return new Promise((resolve, reject) => {
            axios({
                method: 'post',
                url,
                data: data,
            })
                .then(res => {
                    if (msg) message.success(res.data.msg);
                    resolve(res.data)
                })
                .catch(err => {
                    // reject(err)
                });
        })
    },
    get(url, data, msg = true) {
        return new Promise((resolve, reject) => {
            axios({
                method: 'get',
                url,
                params: data,
            })
                .then(res => {
                    if (msg) message.success(res.data.msg);
                    resolve(res.data)
                })
                .catch(err => {
                    // reject(err)
                })
        })
    }
}


// axios.get("/api/index", {}).then((res) => {
//   console.log(res);
// });