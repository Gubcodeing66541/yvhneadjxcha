<template>
    <div class="user-notice layout flex">
        <div class="left shrink pc">
            <examplePhone :type="'userNotice'"></examplePhone>
        </div>
        <div class="right layout ml40 mt2">
            <a-form :model="form" layout="vertical" :wrapper-col="{ span: 18 }" autocomplete="off" @finish="onFinish">
                <a-form-item label="展示公告" name="is_show">
                    <a-radio-group v-model:value="form.is_show">
                        <a-radio :value="1"> 展示 </a-radio>
                        <a-radio :value="0"> 不展示 </a-radio>
                    </a-radio-group>
                </a-form-item>
                <a-form-item label="公告图片" name="image">
                    <a-upload accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="uploadImgHland">
                        <div class="img-box relative" v-if="form.image">
                            <img :src="form.image" width="200" height="100" class="image" />
                            <div class="remove-img flex align-center just-center"><span class="t-them strong pointer" @click.stop="form.image=''">删除图片</span></div>
                        </div>
                        <div class="uploader-box flex align-center just-center pointer" v-else>
                            <plus-outlined></plus-outlined>
                            <div class="ant-upload-text ml6">上传图片</div>
                        </div>
                    </a-upload>
                </a-form-item>
                <a-form-item label="公告文字" name="text">
                    <a-input v-model:value="form.text" />
                </a-form-item>
                <a-form-item :wrapper-col="{ span: 19 }" class="mt50">
                    <div class="layout">
                        <a-button class="br12" @click="form = {}">重置</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted } from "vue";
import examplePhone from "./examplePhone.vue";
import axios from "@/utils/axios";
import { beforeUpload, uploadFile, } from '@/utils';
import { PlusOutlined } from '@ant-design/icons-vue';
export default {
    name: "userNotice",
    components: { examplePhone, PlusOutlined },
    setup() {
        const state = reactive({
            form: { image:"http://43.155.7.124/static/upload/20221116/uplodservice422111606044444359.png"},
        });

        onMounted(() => {
            getDetail();
        })

        //获取详请
        const getDetail = async () => {
            const res = await axios.post("/notice_setting/info", '', false);
            state.form = res.data.notice;
        }

        //上传图片
        const uploadImgHland = async (files) => {
            const res = await uploadFile(files.file);
            state.form['image'] = res;
        }

        const onFinish = async () => {
            const res = await axios.post("/notice_setting/update", state.form);
        }

        return {
            ...toRefs(state),
            beforeUpload, //上传图片前
            uploadImgHland, //上传图片
            onFinish, //提交表单
        };
    },
};
</script>

<style lang="less" scoped>
.left {
    width: 370px;
}
:deep(.ant-form-vertical .ant-form-item-label>label){
    font-size: 15px;
    font-weight: bold;
}
.uploader-box,.image{
    border: 1px dashed #999;
    border-radius: 8px;
    width: 200px;
    height: 100px;
    overflow: hidden;
}
.remove-img{
    opacity: 0;
    position: absolute;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    z-index: 999;
}
.img-box {
    border-radius: 8px;
    &:hover {
        background: rgba(0, 0, 0, 0.5);
        .remove-img{
            opacity: 1;
        }        
    }
}
</style>