<template>
    <div class="content-box">
        <div class="main bg-white">
            <a-form :model="form" layout="vertical" name="robot" autocomplete="off" @finish="onFinish" :wrapper-col="{ span: 10 }" class="h100 relative">
                <div class="mb60  tel-mb20">
                    <label>机器人状态</label>
                    <a-switch v-model:checked="form.status" class="ml20" />
                </div>
            
                <a-form-item name="head" label="机器人头像" class="mb56 tel-mb20">
                    <div class="flex align-bottom">
                        <a-upload name="avatar" list-type="picture-card" class="avatar-uploader" accept=".png, .jpg, .jpeg"
                            :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChange">
                            <img v-if="form.head" :src="form.head" class="img" />
                            <div v-else>
                                <plus-outlined />
                            </div>
                        </a-upload>
                        <!-- <a-button type="primary" ghost size="small" class="br8 mb10">默认</a-button> -->
                    </div>
                </a-form-item>

                <!-- <a-form-item name="hello" label="智能机器人欢迎语：" :rules="[{ required: true, message: ' ' }]">
                    <a-textarea v-model:value="form.hello" placeholder="例如：客服正在忙碌哦，本次服务由智能小U为您服务哦~" :rows="8" />
                </a-form-item>
                <p style="color:#D6D6D6">客服全部不在线或者全部忙碌，自动由智能机器人接待</p> -->
            
                <div class="absolute-btn layout mt80" style="position:absolute;bottom:0">
                    <div class="tr">
                        <a-button type="primary" html-type="submit" class="br12 plr30">保存</a-button>
                        <a-button type="primary" class="br12 bg-tipc bn ml26 plr30" @click="reset()">重置</a-button>
                    </div>
                </div>
            </a-form>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import { PlusOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import { beforeUpload, uploadFile } from '@/utils'
import axios from "@/utils/axios";
export default {
    name: "knowledgeBase",
    components: {
        PlusOutlined,
    },
    setup() {
        const state = reactive({
            form: {
                status: false, //状态 run启动 stop停止
                head:"",
                hello: " ",
            },
        });

        onMounted(() => {
            getInfo(); //获取列表
        })

        const getInfo = async () => {
            const res = await axios.post("/bot/info",'',false);
            let info = res.data.info;
            state.form.status = info.status == 'run' ? true : false;
            state.form.head = info.head;
        }

        //保存
        const onFinish = async () => {
            let info = JSON.parse(JSON.stringify(state.form));            
            info.status = info.status ? 'run' : 'stop';
            const res = await axios.post("/bot/update_info", info);
        }

        //重置
        const reset = () => {
            state.form = {status: false, avatar:"", content: ""},
            message.success("重置成功")
        }

        //上传图片
        const handleChange = async(files) => {
            const res = await uploadFile(files.file);
            state.form.head = res;
        };

        return {
            ...toRefs(state),
            onFinish, //保存
            reset, //重置
            beforeUpload,  //上传前
            handleChange, //上传
        };
    },
};
</script>

<style lang="less" scoped>
.ant-upload-picture-card-wrapper{
    width: auto;
}
:deep(.ant-form label) {
    font-size: 16px;
}
:deep(.ant-upload.ant-upload-select-picture-card){
    width: 84px;
    height: 84px;
    border-radius: 12px;
    overflow: hidden;
}
</style>