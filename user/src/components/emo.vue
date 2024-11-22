<template>
    <div class="emo-content flex align-center flex-warp">
        <div v-for="(elImg, text) in emoji.emojis" v-html="elImg" :key="elImg" class="item" @click.stop="emojisHland(text)"></div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import { emojiList as emoji } from "@/utils/emojis";
export default {
    name: "myEmo",
    setup(props, contact) {
        const state = reactive({

        });

        //表情选中
        const emojisHland = (text) => {
            contact.emit('emoHland',text);
        }

        return {
            ...toRefs(state),
            emoji, //表情列表
            emojisHland, //表情选中
        };
    },
};
</script>

<style lang="less" scoped>
.item{
    padding: 1.6px;
}
</style>