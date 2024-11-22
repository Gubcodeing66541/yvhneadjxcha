import { createStore } from 'vuex'
import axios from "@/utils/axios";
import { localSet, localGet } from "@/utils";
const store = createStore({
    state: {
        token: '登录成功', //登录token
        info: localGet('info') || {}, //账号信息
        config: localGet('config') || {}, //后台配置信息
    },
    // 获取状态(变量的值)
    getters: {
        getToken: state => state.token,
        getInfo: state => state.info,
        getConfig: state => state.config,
    },
    //修改state里面的变量(修改变量的值)
    mutations: {
        setToken(state, newValue) {
            state.token = newValue
        },
        setInfo(state, newValue) {
            state.info = newValue
        },
        setConfig(state, newValue) {
            state.config = newValue
        },
    },
    // Action 触发 mutation 函数,从而修改状态，通过dispatch()调用
    actions: {
        async asyncSetInfo (context) {
            const res = await axios.post("/info",'',false);
            res.data.info.account = res.data.info.account.toFixed(2);
            localSet('info', res.data.info);
            localSet('config', res.data.config);
            context.commit("setInfo", res.data.info);
            context.commit("setConfig", res.data.config);
        }
    },
    // Module 当状态很多时,把状态分开来管理
    modules: {

    }
})

export default store