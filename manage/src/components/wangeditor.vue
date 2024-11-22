<template>
    <div class="wangeditor-container">
        <div ref="editor" class="h100"></div>
    </div>
</template>
<script>
import { reactive, ref, watchEffect, toRefs, onMounted, nextTick } from "vue";
import WangEditor from "wangeditor";
export default {
    name:"wangeditor",
    props: ['detail'],
    setup(props, context) {
        const editor = ref(null);
        const state = reactive({
            content: "",
        });
        watchEffect(() => {
            state.content = props.detail;
        });
        onMounted(() => {
            nextTick(() => {
                //富文本
                richTxt();
            });
        });
        //富文本
        let instance;
        const richTxt = () => {
            instance = new WangEditor(editor.value);
            instance.config.showLinkImg = true;
            instance.config.showLinkImgAlt = false;
            instance.config.showLinkImgHref = false;
            // instance.config.uploadImgMaxSize = 2 * 1024 * 1024; // 2M
            //上传并插入图片
            // instance.config.customUploadImg = (resultFiles, insertImgFn) => {
            //     resultFiles.forEach((item) => {
            //         qiniuUploadImg({ file: item }, (files) => {
            //             let url = config.imgSrc.url + files.key;
            //             insertImgFn(url);
            //         });
            //     });
            // };
            //富文本内容发生变化
            Object.assign(instance.config, {
                onchange() {
                    context.emit("editorHtml", instance.txt.html()); //触发父组件
                },
            });
            instance.create();
            if (state.content != "") {
                if (instance) {
                    instance.txt.html(state.content);
                }
            }
        };
        return {
            ...toRefs(state),
            editor,
        };
    },
};
</script>
<style lang="less" scoped>
.wangeditor-container {
    position: relative;
    z-index: 1;
}
</style>