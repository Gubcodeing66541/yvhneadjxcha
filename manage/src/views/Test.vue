<template>
    <div class="warp">
        <div>
            <a-button type="primary">Primary Button</a-button>
            <a-button>Default Button</a-button>
            <a-button type="dashed">Dashed Button</a-button>
            <a-button type="text">Text Button</a-button>
            <a-button type="link">{{token}}</a-button>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, computed, onMounted } from "vue";
import axios from "@/utils/axios";
import { useStore } from "vuex";
export default {
    name: "Index",
    setup() {
        const store = useStore();
        const state = reactive({
            token: computed(() => store.state.token), //vuex调用
        });

        onMounted(() => {
            request();
        });

        const request = (async()=>{
            try {
                const res = await axios.get("/json/API_kdd531mytfdzm06i?sdAS1dsnuUa3sd=190001&Jsdh4bajs99dii=sohpuisypf4nfaei");
                console.log(res);
            } catch (err) {}
        })

        return {
            ...toRefs(state),
        };
    },
};
</script>

<style lang="less" scoped>
.warp {
    background-color: aquamarine;
}
</style>