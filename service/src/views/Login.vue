<template>
    <div class="warp flex align-center just-center">
        <div class="login-content bg-white">
            <div class="logo-box bg-them flex align-center just-center t-white">
                <i class="iconfont icon-xiaolian f50"></i>
            </div>
            <a-form :model="form" autocomplete="off" @finish="onFinish" class="mt40">
                <a-form-item name="username" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="form.username" placeholder="请输入授权码" :bordered="false" class="hh50 br18 code-input" />
                </a-form-item>            
                <a-form-item>
                    <a-checkbox v-model:checked="agreement">请阅读并同意</a-checkbox>
                    <span class="t-them pointer" @click="visible.text=true">《用户协议&隐私政策》</span>
                </a-form-item>            
                <a-form-item>
                    <a-button type="primary" html-type="submit" class="btn layout hh50 br18 f16">登录</a-button>
                </a-form-item>
            </a-form>
        </div>
    </div>

    <!-- 协议 -->
    <a-modal title="用户协议&隐私政策" :width="800" v-model:visible="visible.text" :footer="null" :centered="true">
        <div v-html="config.AdDefault" class="scrollbar" style="height: 500px;"></div>
    </a-modal>
</template>

<script>
import { reactive, toRefs, onMounted, computed } from "vue";
import { localSet, localGet } from "@/utils";
import { useRouter } from "vue-router";
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { message } from 'ant-design-vue';
export default {
    name: "Login",
    setup() {
        const router = useRouter();
        const store = useStore();
        const state = reactive({
            form:{
                username:'',                
            },
            agreement: true,
            config:"", //系统配置
            visible:{}, //弹窗
            config: computed(() => store.state.config), //系统配置信息
        });
        onMounted(() => {
            getSetting();
            state.form.username = localGet('username') || '';
        })

        //系统配置
        const getSetting = async () => {
            const res = await axios.post("/config", '', false);
            localSet('config', res.data.config);
            store.commit("setConfig", res.data.config);
            document.title = store.state.config.SystemName + '云客服'
        }

        //登录
        const onFinish = async (values) => {
            if ( !state.agreement ) return message.error("请先阅读并同意用户协议")
            const res = await axios.post("/auth/login", values);
            localSet("token", res.data.token);
            localSet("username", values.username);
            store.commit('setToken', res.data.token)
            router.push("/");
        }
        return {
            ...toRefs(state),
            onFinish, //登录
        };
    },
};
</script>

<style lang="less" scoped>
.warp {
    width: 100vw;
    height: 100vh;
    background: url(../assets/image/login_bg.png) no-repeat;
    background-size: 100% 100%;
    .login-content{
        width: 460px;
        border-radius: 16px;
        padding: 40px;
        box-shadow: 0px 2px 11px 1px rgba(203, 189, 240, 0.46);
        position: relative;
        .logo-box{
            width: 94px;
            height: 94px;
            border-radius: 50%;
            overflow: hidden;
            position: absolute;
            top: -47px;
            left: calc(50% - 47px);
            border: 4px solid white;
        }
        .code-input{
            box-shadow: 0px 1px 8px 5px rgb(219 219 219 / 46%);
        }
        .btn{
            box-shadow: 0px 1px 8px 5px rgb(219 219 219 / 46%);
        }
    }
}
</style>