<template>
    <div class="menu">
        <div class="item lh22" v-for="item in menu" :key="item" @click="menuHland(item)">{{item.title}}</div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
import useClipboard from "vue-clipboard3";
import { Toast } from 'vant';
import axios from "@/utils/axios";
export default {
    name: "bottomMenu",
    props: ['menu'],
    setup() {
        const { toClipboard } = useClipboard();
        const state = reactive({

        });

        //点击
        const menuHland = async(item) => {
            switch (item.action) {
                case 'copy':
                    try {
                        await toClipboard(item.content);
                        Toast(item.tag || '复制成功')
                    } catch (e) {
                        Toast("复制失败")
                    }                    
                    break;
                case 'banswer': //自动回复
                    const res = await axios.post('/message/button', {type:'text', content: item.title, service_content: item.content }, false);
                    break;
                case 'http': //点击链接跳转
                    window.open(item.content,"_blank")
                    break;
                case 'call': //点击打电话
                    window.location.href = 'tel://' + item.content
                    break;
            }
        }

        return {
            ...toRefs(state),
            menuHland, //点击
        };
    },
};
</script>

<style lang="less" scoped>
.menu{
    white-space: nowrap;
    overflow-x: scroll;
    overflow-y: hidden;
    padding: 5px 10px;
    scrollbar-width: none;
    -ms-overflow-style: none;
    &::-webkit-scrollbar {
        display: none;
    }
    .item {
        white-space: nowrap;
        display: inline-block;
        background-color: #fff;
        border: 1px solid #eee;
        border-radius: 20px;
        font-size: 12px;
        padding: 0 8px;
        cursor: pointer;
        margin-right: 6px;
    }
}
</style>