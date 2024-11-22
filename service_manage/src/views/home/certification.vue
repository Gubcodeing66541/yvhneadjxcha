<template>
    <div class="content-box">
        <div class="main bg-white relative">
            <div class="mb40">
                <span class="f20 t-them">实名认证</span>
            </div>
            <div>
                <div class="steps flex">
                    <p class="item ico icon_1">1.个人信息认证</p>
                    <p class="item ico icon_2">2.证件审核中</p>
                    <p class="item ico icon_3">3.审核通过</p>
                </div>
                <p class="f14 mt10">
                    <span class="t-tipc">填写个人信息</span> <span class="t-red">（身份证照片仅用于实名认证，请确认好个人信息，如识别错误请手动填写，确保信息通过）</span>
                </p>
            </div>
            <div class="mt30">
                <a-form :model="form" name="basic" autocomplete="off" @finish="onFinish"  :wrapper-col="{ span: 10 }">
                    <a-form-item label="真实姓名" name="real_name" :rules="[{ required: true, message: '请输入真实姓名号!' }]">
                        <a-input v-model:value="form.real_name" placeholder="请输入真实姓名" />
                    </a-form-item>
                
                    <a-form-item label="身份证号" name="id_card_number" :rules="[{ required: true, message: '请输入您的证件号码!' }]">
                        <a-input v-model:value="form.id_card_number" placeholder="请输入您的证件号码" />
                    </a-form-item>

                    <a-form-item label="上传证件" :rules="[{ required: true, message: '请上传身份证正面' }]">
                        <div class="flex">
                            <a-form-item name="id_card_a" >
                                <a-upload name="id_card_a" list-type="picture-card" class="avatar-uploader" accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChange">
                                    <img v-if="form.id_card_a" :src="form.id_card_a" class="img" />
                                    <div class="identity_zm" v-else>
                                        <div class="upload-bg flex dir-column align-center just-center">
                                            <i class="ico icon_sc"></i>
                                            <span class="t-white mt10 f12">上传正面</span>
                                        </div>
                                    </div>
                                </a-upload>
                            </a-form-item>
                            <a-form-item name="id_card_b" :rules="[{ required: true, message: '请上传身份证反面' }]">
                                <a-upload name="id_card_b" list-type="picture-card" class="avatar-uploader" accept=".png, .jpg, .jpeg"
                                    :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChangeF">
                                    <img v-if="form.id_card_b" :src="form.id_card_b" class="img" />
                                    <div class="identity_fm" v-else>
                                        <div class="upload-bg flex dir-column align-center just-center">
                                            <i class="ico icon_sc"></i>
                                            <span class="t-white mt10 f12">上传反面</span>
                                        </div>
                                    </div>
                                </a-upload>
                            </a-form-item>
                        </div>
                    </a-form-item>

                    <a-form-item label="图例及要求">
                        <div class="f16">
                            <p>1.格式要求JPG、JPEG、png，大小不超过2.0MB</p>
                            <p class="mt20">2.上传身份证正反面的原件照片，文字清晰可辨认</p>
                        </div>
                    </a-form-item>
                    <div class="mt50">
                        <div class="flex just-end">
                            <a-button type="primary" html-type="submit" class="br12">提交</a-button>
                            <a-button type="primary" class="br12 bg-tipc bn ml26" @click="reset()">重置</a-button>
                        </div>
                    </div>
                </a-form>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, computed } from "vue";
import { message } from 'ant-design-vue';
import { beforeUpload, uploadFile } from '@/utils';
import { useStore } from "vuex";
import axios from "@/utils/axios";
export default {
    name: "certification",
    setup() {
        const store = useStore();
        const state = reactive({
            form: computed(() => store.state.info), //系统信息
        });

        onMounted(()=>{
        })
        
        //上传图片
        const handleChange = async (files) => {
            const res = await uploadFile(files.file);
            state.form.id_card_a = res;
        };

        //上传图片
        const handleChangeF = async (files) => {
            const res = await uploadFile(files.file);
            state.form.id_card_b = res;
        };

        //提交
        const onFinish = async (values) => {
            const res = await axios.post("/update", values);
            store.dispatch("asyncSetInfo");
        };
        
        //重置
        const reset = () => {
            message.success("重置成功")
        }

        return {
            ...toRefs(state),
            onFinish, //提交
            reset, //重置
            beforeUpload,  //上传前
            handleChange, //上传正面
            handleChangeF, //上传发面
        };
    },
};
</script>

<style lang="less" scoped>
.steps{
    .item{
        width: 300px;
        height: 45px;
        line-height: 45px;
        text-align: center;
        color: white;
        font-size: 18px;
    }
}
.ant-input{
    border-radius: 12px;
}
:deep(.ant-form-item-label>label){
    font-size: 16px;
}
.ant-form-item{
    margin-bottom: 30px;
}
.avatar-uploader{
    :deep(.ant-upload-select-picture-card) {
        width: 258px;
        height: 152px;        
        border-radius: 12px;
        overflow: hidden;
    }    
    .identity_zm {
        width: 100%;
        height: 100%;
        background: url(../../assets/image/icon_sfzzm.png) no-repeat;
        background-size: 100% 100%;
    }
    .identity_fm {
        width: 100%;
        height: 100%;
        background: url(../../assets/image/icon_sfzfm.png) no-repeat;
        background-size: 100% 100%;
    }
    .upload-bg{
        width: 100%;
        height: 100%;
        background: rgba(51, 51, 51, 0.5);
    }
}
</style>