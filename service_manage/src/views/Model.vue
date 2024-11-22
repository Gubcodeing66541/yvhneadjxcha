<template>
    <div class="warp">
        模板
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
export default {
    name: "",
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