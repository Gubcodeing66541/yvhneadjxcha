<template>
    <div class="notice bg-white mg10 br6 hidden">
        <van-image width="100%" height="300" fit="cover" position="center" :src="notice.image" class="bg-them" v-if="notice.image && notice.image!= ''" @click="imagesZoom(notice.image)" />
        <p class="lh30 plr10 pb4 t-red" v-if="notice.text && notice.text!= ''">{{ notice.text}}</p>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import { ImagePreview } from 'vant';
export default {
    name: "notice",
    props: ['notice'],
    setup() {
        const state = reactive({

        });

        // 预览图片
        const imagesZoom = (img) => {
            ImagePreview([img]);
        }

        return {
            ...toRefs(state),
            imagesZoom,
        };
    },
};
</script>

<style lang="less" scoped>
.item {
    padding: 1.6px;
}
</style>