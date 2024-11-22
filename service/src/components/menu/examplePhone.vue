<template>
    <div class="example-content relative">
        <i class="ico icon_tel" />
        <div class="tc bb-tipc f16 pb10">用户{{ type == 'userMenu'? '快捷菜单':'公告' }}示例</div>
        <div class="notice plr10" v-show="type == 'userNotice'">
            <div class="flex dir-column align-center bg-white br6 hidden">
                <div class="layout tc hh150 lh150 bg-red">我是系统公告图片的位置</div>
                <p class="lh30 plr10 pb4 t-red">我是公告文字的位置</p>
            </div>
        </div>
        <div class="bot-content layout">            
            <div class="menu layout flex" v-show="type == 'userMenu'">
                <div class="item">点击链接跳转</div>
                <div class="item">我是快捷菜单</div>
                <div class="item">电话咨询</div>
                <div class="item">可退货吗</div>
            </div>
            <div class="input-box plr12 ptb6 flex algin-center">
                <div class="hh40 layout bg-white br6"></div>
                <div class="flex align-bottom lh36">
                    <i class="iconfont icon-biaoqing f30 ml10"></i>
                    <i class="iconfont icon-jia f30 ml14"></i>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs } from "vue";
export default {
    name: "examplePhone",
    props: ['type'],
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
.example-content{
    border: 4px solid #ddd;
    border-radius: 16px;
    overflow: hidden;
    height: 650px;
    background: #ededed;
    .bot-content{
        position: absolute;
        bottom: 0;
        left: 0;
        .input-box{
            background-color: #f6f6f6;
        }
    }
    .menu{
        padding: 10px;
        .item{
            background-color: #fff;
            border: 1px solid #eee;
            border-radius: 20px;
            font-size: 12px;
            padding: 0 8px;
            cursor: pointer;    
            margin-right: 6px;
            &:hover{
                background-color: #eee;
            }
        }
    }
}
</style>