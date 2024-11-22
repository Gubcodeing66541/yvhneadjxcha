<template>
    <div class="report-warp f13 tc">
        <router-link to="/report" style="color:#333">举报</router-link>        
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
export default {
    name: "report",
    setup() {
        const state = reactive({

        });

        return {
            ...toRefs(state),
        };
    },
};
</script>

<style lang="less" scoped>
.report-warp{
    position: fixed;
    right: 10px;
    top: 10px;
    width: 40px;
    height: 40px;
    line-height: 40px;
    background-color: #d0cfcf;
    border-radius: 50%;
}
</style>