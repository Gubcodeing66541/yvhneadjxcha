<template>
    <div class="report pd30">
        <div class="flex just-center mt40"><p class="gou tc f28"><van-icon name="success" /></p></div>
        <p class="tc f18 strong mt40">投诉已提交</p>
        <p class="tc mt20 lh28">微信团队会尽快核实，并通过“微信团队”通知你审核结果。感谢你的支持</p>
        <p class="tc button" style="posi"><van-button type="success" class="w60 strong" @click="close">关闭</van-button></p>
    </div>
</template>

<script>
import { reactive, toRefs,onMounted } from "vue";
export default {
    name: "Report",
    setup() {
        const state = reactive({

        });

        onMounted(()=>{
            document.title = "举报成功"; //标题
            if (typeof WeixinJSBridge != "undefined") {
                WeixinJSBridge.call('hideToolbar');
            }
        })

        const close = () => {
            if (typeof WeixinJSBridge != "undefined") {
                WeixinJSBridge.call("closeWindow");
            }
        }

        return {
            ...toRefs(state),
            close,
        };
    },
};
</script>

<style lang="less" scoped>
.gou{
    width: 50px;
    height: 50px;
    line-height: 50px;
    background-color: #07c160;
    border-radius: 50%;
    color: white;
}
.button{
    position: fixed;
    bottom: 20%;
    left: 0;
    width: 100%;
}
</style>