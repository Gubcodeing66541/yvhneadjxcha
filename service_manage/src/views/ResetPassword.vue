<template>
    <div class="warp">
        <div class="container flex align-center">
            <div class="w50 pc"><img src="../assets/image/login_1.png" class="layout" /></div>
            <div class="w50 pl100 pr50 tel-layout tel-pd10">
                <h1 class="f30">重置密码</h1>
                <a-form :model="form" autocomplete="off" @finish="onFinish" :wrapper-col="{ span: 20 }">
                    <a-form-item name="password" :rules="[{ required: true, message: '请输入原密码!' }]">
                        <a-input-password v-model:value="form.password" autocomplete="off" :bordered="false" placeholder="请输入原密码" class="input" />
                    </a-form-item>                    
                    <a-form-item name="new_password" :rules="[{ required: true, pattern: /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])[\da-zA-Z]{6,14}$/, message: '密码6至14位包含大写和小写字母' }]">
                        <a-input-password v-model:value="form.new_password" autocomplete="off" :bordered="false" placeholder="请输入新密码" class="input" />
                    </a-form-item>
                    <a-form-item name="two_new_password" :rules="[{ required: true, validator:validatepasssure }]">
                        <a-input-password v-model:value="form.two_new_password" autocomplete="off" :bordered="false" placeholder="请确认新密码" class="input" />
                    </a-form-item>
                    <a-form-item>
                        <a-button type="primary" html-type="submit" class="layout hh75 br18 f20">确认重置密码</a-button>
                    </a-form-item>
                    <a-form-item class="tc">
                        <span class="t-white">取消重置？</span>
                        <router-link to="/" class="pointer t-them">返回首页</router-link>
                    </a-form-item>
                </a-form>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, } from "vue";
import axios from "@/utils/axios";
import { localEmpty } from "@/utils";
export default {
    name: "Login",
    setup() {
        const state = reactive({
            form: {
                password: '',
                new_password: '',
                two_new_password: '',
            },
        });

        const validatepasssure = (rule, value, callback) => {
            return new Promise((resolve, reject) => {
                if (value == '') reject('请再次输入新密码')
                else if (value != state.form.new_password) reject('两次密码不一致')
                else resolve()
            })            
        }

        //提交
        const onFinish = async (values) => {
            let form = JSON.parse(JSON.stringify(values))
            delete form.two_new_password
            const res = await axios.post("/auth/reset_password", form);
            localEmpty('token');
            window.location.reload();
        };

        return {
            ...toRefs(state),
            validatepasssure, //密码验证
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

    .container {
        width: 100%;
        height: 100%;
        background: rgba(255, 255, 255, 0.24);
        border-radius: 62px;
        padding: 6%;

        .input {
            background: white;
            height: 75px;
            box-shadow: 0px 4px 6px 1px rgba(219, 219, 219, 0.46);
            border-radius: 18px;
        }

        .code-cont {
            box-shadow: 0px 4px 6px 1px rgba(219, 219, 219, 0.46);
            border-radius: 18px;
            overflow: hidden;

            .code-box {
                width: 179px;
                height: 75px;
                background: #29A19C;
            }
        }
    }
}
</style>