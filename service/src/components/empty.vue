<template>
    <div class="empty-content hidden">
        <div class="h100 flex align-center just-center mobile-hidden flex align-center just-center">
            <a-empty :image="image" :image-style="{ height: '600px', opacity:'0.4' }">
                <template #description>
                    <span>
                        Customize <a href="#">Description</a>
                    </span>
                </template>
            </a-empty>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
export default {
    name: "empty",
    setup() {
        const state = reactive({
            image: "http://"+import.meta.env.VITE_BASE_URL + 'common/service/ico_empty.png'
        });

        return {
            ...toRefs(state),
        };
    },
};
</script>

<style lang="less" scoped>
.empty-content{
    height: 100%;
}
</style>