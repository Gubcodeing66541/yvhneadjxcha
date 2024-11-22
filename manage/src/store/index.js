import { createStore } from 'vuex'
import axios from "@/utils/axios";
import { localSet, localGet } from "@/utils";

const store = createStore({
    state: {
        config: localGet('config') || {}, //后台配置信息
    },
    // 获取状态(变量的值)
    getters: {
        getConfig: state => state.config,
    },
    //修改state里面的变量(修改变量的值)
    mutations: {
        setConfig(state, newValue) {
            state.config = newValue
        },
    },
    // Action 触发 mutation 函数,从而修改状态
    actions: {
        async asyncSetInfo(context) {
            const res = await axios.post("/config/get", '', false);
            localSet('config', res.data.config);
            context.commit("setConfig", res.data.config);
        }
    },
    // Module 当状态很多时,把状态分开来管理
    modules: {

    }
})

export default store