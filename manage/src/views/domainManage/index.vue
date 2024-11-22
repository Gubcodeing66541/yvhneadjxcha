<template>
    <div class="content">
        <router-view></router-view>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
export default {
    name: "domainManage",
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

</style>