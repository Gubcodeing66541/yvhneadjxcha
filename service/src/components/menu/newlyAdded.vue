<template>
    <div class="add-msg_info">
        <div class="bd br8">
            <a-form-item class="ml11 mr11 pt11 pb4 bb-tipc mb0">
                <div class="flex align-center just-between pl6 pr6">
                    <a-upload accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChange">
                        <PictureOutlined class="t-tipc f18 pointer" title="上传图片" />
                    </a-upload>
                    <span class="f12 t-red pointer" @click="visibleForm.msg_type = 'text', visibleForm.msg_info=''" v-if="visibleForm.msg_type == 'image'">删除图片</span>
                </div>
            </a-form-item>
            <a-textarea v-model:value="visibleForm.msg_info" :bordered="false" placeholder="请输入内容" :rows="8"  v-if="visibleForm.msg_type=='text'" />
            <div class="mg10">
                <a-image :src="visibleForm.msg_info" v-if="visibleForm.msg_type=='image'" />
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, watch, onMounted } from "vue";
import { PictureOutlined } from '@ant-design/icons-vue';
import { beforeUpload, uploadFile } from '@/utils';
export default {
    name: "newlyAdded",
    components: { PictureOutlined },
    props: ['detailForm'],
    setup(props,contact) {
        const state = reactive({
            visibleForm: {
                msg_type: "text",
                msg_info: "",
                
            },
        });

        watch(() => state.visibleForm, (v1) => {
            contact.emit('visibleForm',v1)
        }, { deep: true })

        onMounted(()=>{
            if (props.detailForm.id && props.detailForm.id != ''){
                state.visibleForm.msg_info = props.detailForm.msg_info;
                state.visibleForm.msg_type = props.detailForm.msg_type;
                state.visibleForm['id'] = props.detailForm.id;
            }
        })

        //上传图片
        const handleChange = async (files) => {
            const res = await uploadFile(files.file);
            state.visibleForm.msg_info = res;
            state.visibleForm.msg_type = 'image';          
        };

        return {
            ...toRefs(state),
            beforeUpload, //上传图片前
            handleChange, //上传图片
        };
    },
};
</script>

<style lang="less" scoped>

</style>