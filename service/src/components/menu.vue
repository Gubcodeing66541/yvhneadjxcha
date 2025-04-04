<template>
    <div class="menu-container bg-them">
        <!-- 个人信息 -->
        <div class="tc ptb20">
            <a-popover trigger="click" placement="rightTop">
                <template #content>
                    <div style="width:300px">
                        <div class="flex algin-center">
                            <a-upload 
                                accept=".png, .jpg, .jpeg" 
                                :show-upload-list="false" 
                                :before-upload="beforeUpload"
                                :customRequest="handleChange"
                                capture="camera,user">
                                <a-avatar :src="info.head" shape="square" :size="40" class="mr10 pointer" />
                            </a-upload>                            
                            <div class="flex align-center flex-warp">
                                <div class="layout">
                                    <div v-if='!editInput'>
                                        <span>{{ info.name }}</span>
                                        <EditOutlined class="pointer ml10" @click="editInput=true" />
                                    </div>
                                    <div class="f12" v-else>
                                        <a-input v-model:value="nickname" placeholder="请输入昵称" class="ww150 mr10" />
                                        <span class="bg-them br4 t-white plr6 ptb4 pointer" @click="editNickname">修改</span>
                                        <span class="bg-tipc br4 plr6 ptb4 pointer ml10" @click="editInput=false">取消</span>
                                    </div>
                                </div>
                                <p class=" t-tipc f12 mt2">到期时间：{{ info.time_out }}</p>
                            </div>
                        </div>
                        <div class="mt10 flex just-between plr50">
                            <div class="item strong tc pointer" @click="tabIndex = 'mailList', visible['mailList'] = true">
                                <p>{{ info.room_count }}</p>
                                <p>通讯录</p>
                            </div>
                            <div class="item strong tc pointer" @click="tabIndex = 'blacklist', visible['blacklist'] = true">
                                <p>{{ info.black_count }}</p>
                                <p>黑名单</p>
                            </div>
                        </div>
                    </div>
                </template>
                <a-avatar :src="info.head" shape="square" :size="50" style="background-color: #f56a00;border:2px solid #fff" class="pointer br6" @click="getInfo" />
            </a-popover>            
        </div>
        
        <!-- 菜单列表 -->
        <div class="list">
            <!-- 菜单列表 -->
            <div class="menu-item ptb8 tc pointer" :class="tabIndex == item.tabIndex ? 'active':''"
                v-for="(item, index) in menuList" :key="index" @click="tabIndex = item.tabIndex, visible[tabIndex] = true">
                <i class="iconfont f24" :class="item.ico"></i>
                <p class="tc">{{ item.title }}</p>
            </div>
            <!-- 菜单*更多 -->
            <a-popover placement="rightBottom">
                <template #content>
                    <div class="popover-pointer">
                        <div class="item" @click="tabIndex = 'userMenu', visible['userMenu'] = true">快捷菜单</div>
                        <div class="item" @click="tabIndex = 'userNotice', visible['userNotice'] = true">用户公告</div>
                    </div>
                </template>
                <div class="menu-item ptb8 tc pointer" :class="tabIndex == 'userMenu' || tabIndex == 'userNotice' ? 'active':''">
                    <i class="iconfont icon-gengduo1 f24"></i>
                    <p class="tc">更多</p>
                </div>
            </a-popover>
        </div>

        <!-- 二维码 -->
        <a-modal v-model:visible="visible.qrCode" :centered="true" :footer="null" :width="500" :bodyStyle="{padding:0}" @cancel="close">
            <template #title>
                <div class="relative">
                    二维码<span class="t-them f-normal pointer ml20 " style="position: absolute;right: 26%;" @click="copyHland(info.web)">[复制链接]</span>
                </div>
            </template>
            <div class="bt">
                <vue-qr :logoSrc="info.head" :text="info.web" :backgroundColor="info.code_background" :colorLight="info.code_background" :colorDark="info.code_color" :logoMargin="5" :size="500"></vue-qr>
            </div>
            <div class="plr20 ptb10">
                <!-- 二维码设置部分 -->
                <div class="flex just-between align-center mb10">
                    <div class="flex align-center">
                        <span class="mr6">二维码背景</span>
                        <input type="color" :value="info.code_background" @input="e => onColorChange('code_background', e)" />
                    </div>
                    <div class="flex align-center">
                        <span class="mr6">二维码颜色</span>
                        <input type="color" :value="info.code_color" @input="e => onColorChange('code_color', e)" />
                    </div>
                    <a-popconfirm title="重置后当前二维码会立即失效，是否重置？" @confirm="resetCode">
                        <div class="pointer" title="重置后当前二维码会立即失效">
                            <redo-outlined />
                            <span class="ml6">重置二维码</span>
                        </div>
                    </a-popconfirm>         
                    <a-popconfirm title="更换后当前二维码，旧码不会失效？" @confirm="updateCode">
                        <div class="pointer" title="更换当前二维码 旧码不会失效">
                            <redo-outlined />
                            <span class="ml6">更换二维码</span>
                        </div>
                    </a-popconfirm>
                </div>
                
                <!-- 域名绑定部分 -->
                <div class="flex just-between align-center">
                    <div class="flex align-center">
                        <a-button type="primary" size="small" class="mr10" @click="visible.bindDomainModal = true">绑定入口</a-button>
                        <a-button type="primary" size="small" class="mr10" @click="visible.bindActionModal = true">绑定落地</a-button>
                    </div>
                </div>
            </div>
        </a-modal>

        <!-- 绑定入口弹窗 -->
        <a-modal v-model:visible="visible.bindDomainModal" title="绑定入口域名" @ok="handleBindDomain" :width="500">
            <div class="mb10 t-tipc">当前入口域名：{{ info.bind_domain || '未设置' }}</div>
            <a-input v-model:value="bindDomain" placeholder="请输入入口域名" />
        </a-modal>

        <!-- 绑定落地弹窗 -->
        <a-modal v-model:visible="visible.bindActionModal" title="绑定落地域名" @ok="handleBindAction" :width="500">
            <div class="mb10 t-tipc">当前落地域名：{{ info.bind_action || '未设置' }}</div>
            <a-input v-model:value="bindAction" placeholder="请输入落地域名" />
        </a-modal>

        <!-- 通讯录 -->
        <a-drawer v-model:visible="visible.mailList" :width="drawer.width" :closable="false" class="custom-class" title="通讯录" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <mailList v-if="visible.mailList"></mailList>
        </a-drawer>

        <!-- 黑名单 -->
        <a-drawer v-model:visible="visible.blacklist" :width="drawer.width" :closable="false" class="custom-class" title="黑名单" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <blacklist v-if="visible.blacklist"></blacklist>
        </a-drawer>

        <!-- 群发 -->
        <a-drawer v-model:visible="visible.massDispatch" :width="drawer.width" :closable="false" class="custom-class" title="群发" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <template #extra>
                <a-space>
                    <a-button type="primary" size="small" @click="visible.addMassVisible = true"><plus-outlined />新建</a-button>
                </a-space>
            </template>
            <massDispatch :addVisible="visible.addMassVisible" @addClose="addClose($event,'addMassVisible')" v-if="visible.massDispatch"></massDispatch>
        </a-drawer>

        <!-- 打招呼 -->
        <a-drawer v-model:visible="visible.sayHello" :width="drawer.width" :closable="false" class="custom-class" title="打招呼" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <template #extra>
                <a-space>
                    <a-button type="primary" size="small" @click="visible.addSayHelloVisible = true">
                        <plus-outlined />新建
                    </a-button>
                </a-space>
            </template>
            <sayHello :addVisible="visible.addSayHelloVisible" @addClose="addClose($event,'addSayHelloVisible')" v-if="visible.sayHello"></sayHello>
        </a-drawer>

        <!-- 快捷回复 -->
        <a-drawer v-model:visible="visible.quickReply" :width="drawer.width" :closable="false" class="custom-class" title="快捷回复" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <template #extra>
                <a-space>
                    <a-button type="primary" size="small" @click="visible.addQuickReplyVisible = true">
                        <plus-outlined />新建
                    </a-button>
                </a-space>
            </template>
            <quickReply :addVisible="visible.addQuickReplyVisible" @addClose="addClose($event,'addQuickReplyVisible')" v-if="visible.quickReply"></quickReply>
        </a-drawer>

        <!-- 离线回复 -->
        <a-drawer v-model:visible="visible.offlineReply" :width="drawer.width" :closable="false" class="custom-class" title="离线回复" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <template #extra>
                <a-space>
                    <a-button type="primary" size="small" @click="visible.addOfflineReplyVisible = true">
                        <plus-outlined />新建
                    </a-button>
                </a-space>
            </template>
            <offlineReply :addVisible="visible.addOfflineReplyVisible" @addClose="addClose($event,'addOfflineReplyVisible')" v-if="visible.offlineReply"></offlineReply>
        </a-drawer>

        <!-- 统计 -->
        <a-drawer v-model:visible="visible.statistics" :width="drawer.width" :closable="false" class="custom-class" title="统计" placement="right" :bodyStyle="drawer.bodyStyle" @close="close">
            <template #extra>
                <a-space>
                    <a-button type="primary" size="small" @click="copyHland('share')"> 分享统计 </a-button>
                </a-space>
            </template>
            <statistics v-if="visible.statistics"></statistics>
        </a-drawer>

        <!-- 设置用户快捷菜单 -->
        <a-modal v-model:visible="visible.userMenu" :centered="true" :footer="null" title="设置用户快捷菜单" :width="1200" :bodyStyle="{padding:0}" @cancel="close">
           <userMenu></userMenu>
        </a-modal>

        <!-- 设置公告设置 -->
        <a-modal v-model:visible="visible.userNotice" :centered="true" :footer="null" title="设置用户公告" :width="1000" :bodyStyle="{padding:10}" @cancel="close">
            <userNotice></userNotice>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, computed } from "vue";
import { EditOutlined, RedoOutlined, ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue';
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { beforeUpload, uploadFile } from '@/utils';
import { message } from 'ant-design-vue';
import vueQr from 'vue-qr/src/packages/vue-qr.vue';
import useClipboard from "vue-clipboard3";
import mailList from '@/components/menu/mailList.vue';
import blacklist from '@/components/menu/blacklist.vue';
import massDispatch from '@/components/menu/massDispatch.vue';
import sayHello from '@/components/menu/sayHello.vue';
import quickReply from '@/components/menu/quickReply.vue';
import offlineReply from '@/components/menu/offlineReply.vue';
import statistics from '@/components/menu/statistics.vue';
import userMenu from '@/components/menu/userMenu.vue';
import userNotice from '@/components/menu/userNotice.vue';
const menuList = [
    { title: "二维码", ico: "icon-erweima", tabIndex: "qrCode" },
    { title: "消息", ico: "icon-xiaoxi", tabIndex: "news" },
    { title: "通讯录", ico: "icon-tongxunlu", tabIndex: "mailList" },
    { title: "黑名单", ico: "icon-heimingdan", tabIndex: "blacklist" },
    { title: "群发", ico: "icon-fasong", tabIndex: "massDispatch" },
    { title: "打招呼", ico: "icon-dazhaohu", tabIndex: "sayHello" },
    { title: "快捷回复", ico: "icon-quickReply", tabIndex: "quickReply" },
    { title: "离线回复", ico: "icon-quickReply", tabIndex: "offlineReply" },
    { title: "统计", ico: "icon-tongji", tabIndex: "statistics" },
];
export default {
    name: "MyMenu",
    components: { EditOutlined, vueQr, RedoOutlined, ExclamationCircleOutlined, mailList, blacklist, massDispatch, PlusOutlined, sayHello, quickReply, offlineReply, statistics, userMenu, userNotice },
    setup() {
        const { toClipboard } = useClipboard();
        const store = useStore();
        const state = reactive({
            info: computed(() => store.state.info), //系统信息
            nickname:"",//昵称
            editInput:false, //修改昵称
            menuList, //菜单列表
            tabIndex:'news', //菜单索引
            visible:{}, //弹窗显示
            code:{},
            drawer:{ bodyStyle: { backgroundColor: '#f5f5f5' }, width: 500 },
            bindDomain: "", // 绑定入口域名
            bindAction: "", // 绑定落地域名
        });

        onMounted(()=>{
            
        })
        //获取账号信息
        const getInfo = () => {
            store.dispatch("asyncSetInfo");
        }

        //修改昵称
        const editNickname = async () => {
            const res = await axios.post("/update", { name: state.nickname }, false);
            store.dispatch("asyncSetInfo");
            state.editInput = false;
            state.nickname = '';
        }

        //修改头像
        const handleChange = async (files) => {
            const head = await uploadFile(files.file);
            const res = await axios.post("/update", { head: head },false);
            store.dispatch("asyncSetInfo");
        };

        //二维码样式调整
        const onColorChange = async (type,e) => {
            state.code[type] = e.target.value;
            const res = await axios.post("/update", { [type]: e.target.value }, false);
            store.dispatch("asyncSetInfo");
        }

        //刷新二维码
        const resetCode = async () => {
            const res = await axios.post("/reset_qrcode");
            store.dispatch("asyncSetInfo");
        }

        const updateCode = async () => {
            const res = await axios.post("/update_qrcode");
            store.dispatch("asyncSetInfo");
        }

        //菜单关闭
        const close = () => {
            state.tabIndex = 'news'
        }

        //控制新建显示隐藏
        const addClose = (e,type) => {
            state.visible[type] = e;
        }

        //复制
        const copyHland = async(content) => {
            if (content == 'share') content = `${import.meta.env.VITE_BASE_URL}/service/#/Statistics?username=${state.info.username}`
            try {
                await toClipboard(content);                
                message.success("复制成功")
            } catch (e) {
                message.error("复制失败")
            }
        }

        // 绑定入口域名
        const handleBindDomain = async () => {
       
            try {
                const res = await axios.post("/domain/bind_domain", {
                    domain: state.bindDomain
                });
                message.success("绑定入口域名成功");
                state.bindDomain = "";
            } catch (error) {
                message.error("绑定失败");
            }
        };

        // 绑定落地域名
        const handleBindAction = async () => {
 
            try {
                const res = await axios.post("/domain/bind_action", {
                    action: state.bindAction
                });
                message.success("绑定落地域名成功");
                state.bindAction = "";
            } catch (error) {
                message.error("绑定失败");
            }
        };

        return {
            ...toRefs(state),     
            getInfo, //获取用户信息       
            editNickname, //修改昵称
            beforeUpload, //图片上传前
            handleChange, //修改头像
            close, //关闭弹窗
            onColorChange, //二维码样式调整
            resetCode, //重置二维码
            updateCode,//更新二维码
            addClose, //控制新建显示隐藏
            copyHland, //分享统计
            handleBindDomain,
            handleBindAction,
        };
    },
};
</script>

<style lang="less" scoped>
.menu-container{
    height: 100%;    
    height: 100%;
    .list{
        overflow-y: scroll;
        height: calc(100% - 95px);
        &::-webkit-scrollbar {
            width: 0px;
            height: 0px;
        }
        .menu-item {
            color: #dddddd;
            border-radius: 18px 0 0 18px;
        }
        .active {
            background: #ebe9e8;
            color: @primary-color;
        }
    }    
}
.w100 {
    width: 100%;
}
.w80 {
    width: 80px;
}
.flex1 {
    flex: 1;
}
.flex-column {
    flex-direction: column;
}
.mt10 {
    margin-top: 10px;
}
.mb10 {
    margin-bottom: 10px;
}
</style>