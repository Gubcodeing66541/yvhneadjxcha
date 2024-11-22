<template>
    <div class="warp flex dir-column">
        <!-- 头部 -->
        <div class="header">            
            <MyHeader></MyHeader>
        </div>
        <div class="container flex">
            <!-- 菜单 -->
            <div class="menu">
                <MyMenu></MyMenu>
            </div>
            <!-- 内容 -->
            <div class="main flex dir-column">
                <!-- 已读未读状态 -->
                <div class="state">
                    <Mystate></Mystate>
                </div>
                <div class="content flex">
                    <!-- 用户列表 -->
                    <div class="user">
                        <MyUser @user_id="user_id"></MyUser>
                    </div>
                    <!-- 聊天界面 -->
                    <div class="chat">
                        <MyChat :user_id = checkedUserID v-if="checkedUserID && checkedUserID!= ''"></MyChat>   
                        <MyEmpty v-else></MyEmpty>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, onBeforeUnmount } from "vue";
import MyHeader from "@/components/header.vue";
import MyMenu from "@/components/menu.vue";
import Mystate from "@/components/state.vue";
import MyUser from "@/components/user.vue";
import MyChat from "@/components/chat.vue";
import MyEmpty from "@/components/empty.vue";
import { initWebSocket } from "@/socket";
import { websocketClose } from "@/socket";
import { useStore } from "vuex";
export default {
    name: "Home",
    components: { MyHeader, MyMenu, Mystate, MyUser, MyChat, MyEmpty },
    setup() {
        const store = useStore();
        const state = reactive({
            checkedUserID:'', //当前聊天的用户ID
        });

        onMounted(()=>{
            document.title = store.state.config.SystemName + '云客服'
            store.dispatch("asyncSetInfo");
            initWebSocket(); //连接websocket
        })

        onBeforeUnmount(()=>{
            websocketClose(); //关闭websocket
        })

        //当前聊天的用户ID
        const user_id = (id) => {
            state.checkedUserID = id;
        }

        return {
            ...toRefs(state),
            user_id, //点击聊天列表的用户ID
        };
    },
};
</script>

<style lang="less" scoped>
.warp{
    width: 100vw;
    height: 100vh;
    overflow-y: hidden;
    display: flex;
    background: #f5f5f5;
    .header{
        width: 100vw;
        height: 40px;
        flex-shrink: 0;
        border-bottom: 1px solid #f1f1f1;
        background: white;        
    }
    .container{
        height: calc(100% - 40px);
        overflow-y: hidden;
        .menu {
            width: 70px;            
            flex-shrink: 0;
        }
        .main{
            width: 100%;
            height: 100%;
            .state{
                height: 40px;
                flex-shrink: 0;
                border-bottom: 1px solid #f1f1f1;
                background: #ffffff;
            }
            .content{
                height: calc(100% - 40px);
                overflow-x: overlay;
                overflow-y: hidden;
                .user {
                    width: 16%;
                    min-width: 160px;
                    flex-shrink: 0;
                    border-right: 1px solid #eee;
                    background: #ebe9e8;
                }            
                .chat {
                    width: 100%;
                }
            }            
        }
    }    
}
</style>