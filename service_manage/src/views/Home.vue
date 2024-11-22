<template>
    <div class="warp">
        <div class="menu scrollbar fl t-3s" :class="!menuIsShow ? 'tel-menu-hidden':''">
            <div class="left-title flex align-center just-between pl22 mb26 tel-title">
                <img :src="config.SystemLogo" width="40" height="40" class="round">
                <div class="layout"><span class="f20 ml12">{{ config.SystemName }}管理后台</span></div>
                <div class="mr40 tel-mr20">
                    <router-link to="/"><i class="ico icon_sy pc" @click="home"></i></router-link>
                    <menu-unfold-outlined v-if="menuIsShow" class="mobile f24" @click="menuIsShow = !menuIsShow" />
                    <MenuFoldOutlined class="mobile f24" v-else @click="menuIsShow = !menuIsShow" />
                </div>
            </div>
            <a-menu v-model:selectedKeys="selectedKeys" class="tel-menu" style="width: 256px" mode="inline" :open-keys="openKeys" @openChange="onOpenChange" @click="menuClick">
                <a-menu-item key="record" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_home"></i></div>
                    </template>
                    <span>首页</span>
                </a-menu-item>
                <a-sub-menu key="account">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_zh"></i></div>
                    </template>
                    <template #title>账户</template>
                    <a-menu-item key="certification" v-if="config.RealName==1">实名认证</a-menu-item>
                    <a-menu-item key="recharge">充值</a-menu-item>
                    <a-menu-item key="rechargeRecord">订单列表</a-menu-item>
                    <template #expandIcon>
                        <i class="ico icon_sq t-3s"></i>
                    </template>
                </a-sub-menu>
                <a-menu-item key="service" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_kfgl"></i></div>
                    </template>
                    <span>客服管理</span>
                </a-menu-item>
                <a-menu-item key="quickReply" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_kjhf"></i></div>
                    </template>
                    <span>快捷回复</span>
                </a-menu-item>
                <a-sub-menu key="robotManagement">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_jqrgl"></i></div>
                    </template>
                    <template #title>机器人管理</template>
                    <a-menu-item key="robot">机器人</a-menu-item>
                    <a-menu-item key="knowledgeBase">机器人知识库</a-menu-item>
                    <template #expandIcon>
                        <i class="ico icon_sq t-3s"></i>
                    </template>
                </a-sub-menu>                
                <a-sub-menu key="userManage">
                    <template #icon>
                        <div class="ico-hover"><i class="ico icon_yhgl"></i></div>
                    </template>
                    <template #title>用户管理</template>
                    <a-menu-item key="visitorRecords">访客记录</a-menu-item>
                    <a-menu-item key="blacklist">黑名单</a-menu-item>
                    <template #expandIcon>
                        <i class="ico icon_sq t-3s"></i>
                    </template>
                </a-sub-menu>
            </a-menu>
            <div class="sign-out flex align-center pl30 pointer" @click="loginOut">
                <i class="ico icon_tc"></i>
                <span class="ml22 f16 layout">退出</span>
            </div>
        </div>
        <div class="left-conent fl">
            <div class="header">
                <router-link target="_blank" to="/largeDataScreen" class="sjdp-box pointer bg-them br17 t-white f16 fl t-3s pc">
                    <div class="flex align-center just-center"><i class="ico icon_sjdp mr8"></i>数据大屏</div>
                </router-link>
                <div class="fr flex align-center">
                    <div class="money-box flex align-center">
                        <i class="ico icon_money"></i>
                        <span class="ml8 f20">{{ info.account }}</span>
                    </div>
                    <div class="ml26">
                        <a-popover>
                            <template #content>
                                <div class="popover-pointer">
                                    <div class="item">
                                        <a-upload class="avatar-uploader" accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="handleChange">
                                            <span>修改头像</span>
                                        </a-upload>
                                    </div>
                                    <div class="item" @click="editVisible=true">修改昵称</div>
                                    <div class="item">
                                        <router-link to="/ResetPassword" style="color:#333">重置密码</router-link>
                                    </div>
                                </div>
                            </template>
                            <div>
                                <a-avatar :src="info.head" style="background-color: #f56a00" :size="50"></a-avatar>
                                <span class="f16 ml10">{{ info.name }}</span>
                            </div>
                        </a-popover>
                    </div>
                </div>
            </div>
            <div class="container">
                <router-view />
            </div>
        </div>

        <!-- 编辑用户信息 -->
        <a-modal v-model:visible="editVisible" :centered="true" :closable="false" :footer="null" title="修改昵称">
            <a-form :model="editVisibleForm" layout="vertical" autocomplete="off" @finish="editUser">
                <a-form-item name="name" label="用户昵称" :rules="[{ required: true, message: ' ' }]">
                    <a-input v-model:value="editVisibleForm.name" placeholder="请输入昵称" />
                </a-form-item>
                <a-form-item class="mb0">
                    <div class="tc">
                        <a-button class="br12" @click="editVisible = false">取消</a-button>
                        <a-button type="primary" html-type="submit" class="br12 ml26">保存</a-button>
                    </div>
                </a-form-item>
            </a-form>
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, computed } from "vue";
import { beforeUpload, localEmpty, uploadFile } from '@/utils';
import { useStore } from "vuex";
import { useRouter } from "vue-router";
import { MenuUnfoldOutlined, MenuFoldOutlined } from '@ant-design/icons-vue';
import axios from "@/utils/axios";
export default {
    name: "home",
    components: { MenuUnfoldOutlined, MenuFoldOutlined },
    setup() {
        const router = useRouter();
        const store = useStore();
        const state = reactive({
            info: computed(() => store.state.info), //账号信息
            config: computed(() => store.state.config), //系统配置信息
            rootSubmenuKeys: ['home', 'account', 'quickReply', 'robotManagement', 'serviceManage', 'userManage'],
            openKeys: ['home'],//选中的一级菜单
            selectedKeys: ['record'],//选中的二级菜单
            editVisible:false,//修改昵称弹窗
            editVisibleForm:{},
            menuIsShow:false, //手机端控制菜单显示隐藏
        });

        onMounted (()=>{
            store.dispatch("asyncSetInfo");
            router.push({
                name: 'record'
            })
        })        

        //菜单展开变化监听
        const onOpenChange = (openKeys) => {
            // 将当前打开的父级菜单存入缓存中
            window.sessionStorage.setItem('openKeys', JSON.stringify(openKeys))
            //  控制只打开一个
            const latestOpenKey = openKeys.find(key => state.openKeys.indexOf(key) === -1);
            if (state.rootSubmenuKeys.indexOf(latestOpenKey) === -1) {
                state.openKeys = openKeys;
            } else {
                state.openKeys = latestOpenKey ? [latestOpenKey] : [];
            }
        };

        // 菜单点击获取key并跳转
        const menuClick = (key) => {
            router.push({
                name: key.key
            })
        }

        //修改用户信息
        const editUser = async (val) => {
            const res = await axios.post("/update", val);
            store.dispatch("asyncSetInfo");
            state.editVisible = false;
        }

        //上传图片
        const handleChange = async (files) => {
            const head = await uploadFile(files.file);
            const res = await axios.post("/update", { head });
            store.dispatch("asyncSetInfo");
        };

        //退出登录
        const loginOut = () => {
            localEmpty('token');
            window.location.reload();
        }

        //回到首页
        const home = () => {
            state.openKeys = ['home'];//选中的一级菜单
            state.selectedKeys = ['record'];//选中的二级菜单
        }

        return {
            ...toRefs(state),
            onOpenChange,
            menuClick, //菜单点击获取key并跳转
            editUser, //修改用户信息
            beforeUpload,  //上传前
            handleChange, //上传
            loginOut, //退出登录
            home, //回到首页
        };
    },
};
</script>

<style lang="less" scoped>
.warp{
    background-color: #F5F7FF;
    width: 100vw;
    height: 100vh;
    min-width: 900px;
    padding: 38px 48px;
    .menu{
        width: 326px;
        height: calc(100% - 70px);
        position: relative;
        .left-title{
            height: 75px;
            background: #FFFFFF;
            box-shadow: 0px 3px 6px 1px rgba(211,229,233,0.28);
            border-radius: 18px;
        }
        .ant-menu{
            background: transparent;
            border: none;
            width: 100% !important;
            padding-right: 14px;
            :deep(.ant-menu-sub.ant-menu-inline){
                background: transparent;
            }
            :deep(.ant-menu-submenu-title:hover){
                color: #333;
            }
            :deep(.ant-menu-submenu:hover){
                color: #333;
            }
            :deep(.ant-menu-submenu-selected){
                color: #6971F8 !important;
            }
            :deep(.ant-menu-item-selected){
                background:transparent;
            }
            .ant-menu-submenu-open{
                .icon_sq{
                    transform: rotate(90deg);
                }
            }
            .ico-hover {
                width: 34px;
                height: 34px;                
                border-radius: 12px;
                display: flex;
                align-items: center;
                justify-content: center;
                background: transparent;
            }
            .ant-menu-submenu-selected,.ant-menu-item-selected{
                .ico-hover {
                    background-color: white;            
                    .icon_home {
                        background-image: url(@/assets/image/icon_home2.png);
                    }
                    .icon_zh {
                        background-image: url(@/assets/image/icon_zh2.png);
                    }
                    .icon_kjhf {
                        background-image: url(@/assets/image/icon_kjhf2.png);
                    }
                    .icon_kfgl {
                        background-image: url(@/assets/image/icon_kfgl2.png);
                    }
                    .icon_jqrgl {
                        background-image: url(@/assets/image/icon_jqrgl2.png);
                    }
                    .icon_kfgl {
                        background-image: url(@/assets/image/icon_kfgl2.png);
                    }
                    .icon_yhgl {
                        background-image: url(@/assets/image/icon_yhgl2.png);
                    }
                }
                .icon_sq{
                    background-image: url(@/assets/image/icon_zk.png);
                }
            }
            :deep(.ant-menu-title-content){
                font-size: 15px;
            }
            :deep(.ant-menu-submenu-title){
                height: 68px;
            }
            :deep(.ant-menu-submenu-selected .ant-menu-submenu-title){
                background: url(@/assets/image/menu_bg.png) no-repeat;
                background-size: 100% 100%;
                height: 68px;
                color: #6971F8 !important;
                font-weight: bold;
            }
            :deep(.menu-row){
                height: 68px;
                &.ant-menu-item-selected{
                    background: url(@/assets/image/menu_bg.png) no-repeat !important;
                    background-size: 100% 100% !important;
                }
            }
            
            :deep(.ant-menu-item:active, .ant-menu-submenu-title:active){
                background: transparent;
            }
            :deep(.ant-menu-item::after){
                background: none !important;
            }
            :deep(.ant-menu-item-selected::after){
                border: none !important;
            }            
        }
        .sign-out{
            position: fixed;            
            bottom: 0;
            height: 70px;
            width: 326px;
            padding-bottom: 24px;
            background-color: #F5F7FF;
        }
    }
    .left-conent{
        width: calc(100% - 326px);
        padding-left: 28px;
        .header {
            overflow: hidden;
            padding-top: 16px;
            padding-bottom: 36px;
            .sjdp-box {
                width: 160px;
                height: 48px;
                line-height: 48px;
                text-align: center;
                background: #6971F8;
                &:hover{
                    color: white;
                    box-shadow: 0px 3px 6px 1px rgba(137, 141, 216, 0.5);
                }
            }
        }
    }    
}
</style>