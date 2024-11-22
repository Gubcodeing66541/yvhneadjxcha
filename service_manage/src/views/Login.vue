<template>
    <div class="warp">
        <div class="container flex align-center">
            <div class="w50 pc"><img src="../assets/image/login_1.png" class="layout" /></div>
            <div class="w50 pl100 pr50 tel-layout tel-pd10">
                <h1 class="f32 mb40">Hi~ <br/> 欢迎来到{{ config.SystemName }}管理后台</h1>
                <a-form :model="form" name="basic" autocomplete="off" @finish="onFinish" :wrapper-col="{ span: 20 }">
                    <a-form-item name="username" :rules="[{ required: true, message: '请输入账号!' }]">
                        <a-input v-model:value="form.username" :bordered="false" placeholder="请输入账号" />
                    </a-form-item>
                
                    <a-form-item name="password" :rules="[{ required: true, message: '请输入密码!' }]">
                        <a-input-password v-model:value="form.password" :bordered="false" placeholder="请输入密码" />
                    </a-form-item>
                
                    <a-form-item>
                        <a-button type="primary" html-type="submit" class="layout hh75 br18 f20">登录</a-button>
                    </a-form-item>
                </a-form>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import axios from "@/utils/axios";
import { useRouter } from "vue-router";
import { localSet, localGet } from "@/utils";
export default {
    name: "login",
    setup() {
        const router = useRouter();
        const state = reactive({
            config: {},
            form:{},
        });

        onMounted(()=>{
            getSetting();
            state.form = localGet('loginForm') || {};
        })

        //系统配置
        const getSetting = async () => {
            const res = await axios.post("/config", '', false);
            state.config = res.data.config;
            document.title = res.data.config.SystemName + '管理后台'
        }

        const onFinish = async (values) => {
            const res = await axios.post("/auth/login", values);
            localSet("token", res.data.token);
            localSet("loginForm", values); //记住账号密码
            router.push("/");
        };

        return {
            ...toRefs(state),
            onFinish,
        };
    },
};
</script>

<style lang="less" scoped>
.warp {
    background: url(../assets/image/login_bg.png) no-repeat;
    background-size: 100% 100%;
    width: 100vw;
    height: 100vh;
    padding: 4% 7%;
    .container{
        width: 100%;
        height: 100%;
        background: rgba(255, 255, 255, 0.24);
        border-radius: 62px;
        padding: 6%;
        .ant-input,.ant-input-password{
            background: white;
            height: 75px;
            box-shadow: 0px 4px 6px 1px rgba(219, 219, 219, 0.46);
            border-radius: 18px 18px 18px 18px;
        }
    }
}
</style>