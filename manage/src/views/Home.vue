<template>
    <div class="warp">
        <div class="menu scrollbar fl t-3s" :class="!menuIsShow ? 'tel-menu-hidden':''">
            <div class="title flex align-center just-between pl22 mb26  tel-title">
                <img :src="config.system_logo" width="40" height="40" class="round">
                <div class="layout"><span class="f20 ml12">{{ config.system_name }}总后台</span></div>
                <div class="mr40 tel-mr20">
                    <router-link to="/"><i class="ico icon_sy pc" @click="home"></i></router-link>
                    <menu-unfold-outlined v-if="menuIsShow" class="mobile f24" @click="menuIsShow = !menuIsShow" />
                    <MenuFoldOutlined class="mobile f24" v-else @click="menuIsShow = !menuIsShow" />
                </div>
            </div>
            <a-menu v-model:selectedKeys="selectedKeys" class="tel-menu" style="width: 256px" mode="inline" :open-keys="openKeys" @openChange="onOpenChange" @click="menuClick">
                <a-menu-item key="home" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><home-outlined /></div>
                    </template>
                    <span>首页</span>
                </a-menu-item>
                <a-menu-item key="accountManage" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><contacts-outlined /></div>
                    </template>
                    <span>账号管理</span>
                </a-menu-item>
                <a-sub-menu key="domainManage">
                    <template #icon>
                        <div class="ico-hover"><global-outlined /></div>
                    </template>
                    <template #title>域名管理</template>
                    <a-menu-item key="independentDomainName">独立域名</a-menu-item>
                    <a-menu-item key="transfer">中转域名</a-menu-item>
                    <a-menu-item key="LandedDomainName">落地域名</a-menu-item>
                    <template #expandIcon>
                        <i class="ico icon_sq t-3s"></i>
                    </template>
                </a-sub-menu>
                <a-menu-item key="authorizationManage" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><branches-outlined /></div>
                    </template>
                    <span>授权管理</span>
                </a-menu-item>
                <a-menu-item key="orderManage" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><file-text-outlined /></div>                        
                    </template>
                    <span>订单管理</span>
                </a-menu-item>
                <a-menu-item key="userManage" class="menu-row">
                    <template #icon>
                        <div class="ico-hover"><contacts-outlined /></div>
                    </template>
                    <span>用户管理</span>
                </a-menu-item>
                <a-sub-menu key="systemSettings">
                    <template #icon>
                        <div class="ico-hover"><setting-outlined /></div>
                    </template>
                    <template #title>系统配置</template>
                    <a-menu-item key="service">系统</a-menu-item>
                    <a-menu-item key="notice">公告</a-menu-item>
                    <!-- <a-menu-item key="sensitiveWords">敏感词</a-menu-item> -->
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
            <div class="title flex align-center just-between mb26 plr40">
                <a-breadcrumb>
                    <a-breadcrumb-item v-for="item in matched" :key="item"><span>{{ item.meta.title}}</span></a-breadcrumb-item>
                </a-breadcrumb>
                <div class="flex align-center">
                    <a-avatar style="background-color: #f56a00" :size="50">
                        <template #icon>
                            <UserOutlined />
                        </template>
                    </a-avatar>
                    <a-popover>
                        <template #content>
                            <div class="popover-pointer">
                                <div class="item">
                                    <router-link to="/ResetPassword" style="color:#333">重置密码</router-link>
                                </div>
                            </div>
                        </template>
                        <div class="pointer flex align-center just-end">
                            <span class="f16 ml10">admin</span>
                            <div class="lh26 ml10 t-tipc">
                                <caret-down-outlined />
                            </div>
                        </div>
                    </a-popover>
                </div>
            </div>
            <div class="container">
                <router-view />
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, watch, computed } from "vue";
import { useStore } from "vuex";
import { useRouter,useRoute } from "vue-router"
import { localEmpty } from '@/utils';
import { UserOutlined, HomeOutlined, ContactsOutlined, GlobalOutlined, BranchesOutlined, FileTextOutlined, SettingOutlined, CaretDownOutlined, MenuUnfoldOutlined, MenuFoldOutlined } from '@ant-design/icons-vue';
import axios from "@/utils/axios";
export default {
    name: "Home",
    components: { UserOutlined,HomeOutlined, ContactsOutlined, GlobalOutlined, BranchesOutlined, FileTextOutlined, SettingOutlined, CaretDownOutlined, MenuUnfoldOutlined, MenuFoldOutlined },
    setup() {
        const store = useStore();
        const router = useRouter();
        const route = useRoute();
        const state = reactive({
            config: computed(() => store.state.config), //系统配置信息
            menuIsShow: false, //手机端控制菜单显示隐藏
            matched:[], //面包屑
            rootSubmenuKeys: ['home', 'accountManage', 'domainManage', 'authorizationManage', 'orderManage', 'userManage','systemSettings'],
            openKeys: ['home'],//选中的一级菜单
            selectedKeys: ['home'],//选中的二级菜单
        });

        watch(() => route, (v1) => {
            state.matched = v1.matched;
        }, { deep: true })

        onMounted(() => {
            store.dispatch("asyncSetInfo"); //获取配置信息
            router.push({
                name: 'home'
            })
            state.matched = route.matched;
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

        //退出登录
        const loginOut = () => {
            localEmpty('token');
            window.location.reload();
        }

        //回到首页
        const home = () => {
            state.openKeys = ['home'];//选中的一级菜单
            state.selectedKeys = ['home'];//选中的二级菜单
        }

        return {
            ...toRefs(state),
            onOpenChange,
            menuClick, //菜单点击获取key并跳转
            loginOut, //退出登录
            home, //回到首页
        };
    },
};
</script>

<style lang="less" scoped>
// 菜单宽度
@menuW : 300px;
.warp {
    background-color: #F5F7FF;
    width: 100vw;
    height: 100vh;
    min-width: 900px;
    padding: 38px 48px;
    .title {
        height: 75px;
        background: #FFFFFF;
        box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.28);
        border-radius: 18px;
    }
    .menu {
        width: @menuW;
        height: calc(100% - 70px);
        position: relative;        
        .ant-menu {
            background: transparent;
            border: none;
            width: 100% !important;
            padding-right: 14px;
            :deep(.ant-menu-title-content) {
                    font-size: 15px;
                }
            :deep(.ant-menu-sub.ant-menu-inline) {
                background: transparent;
            }
            :deep(.ant-menu-item-selected) {
                background: transparent;
            }            
            .ant-menu-submenu-open {
                .icon_sq {
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
                .anticon{
                    font-size: 20px !important;
                }
            }
            .ant-menu-submenu-selected, .ant-menu-item-selected {
                .ico-hover {
                    background-color: white;
                }
                .icon_sq {
                    background-image: url(@/assets/image/icon_zk.png);
                }                
            }
            :deep(.ant-menu-submenu-title) {
                height: 68px;
            }
            :deep(.menu-row) {
                height: 68px;
                &.ant-menu-item-selected {
                    background: url(@/assets/image/menu_bg.png) no-repeat !important;
                    background-size: 100% 100% !important;
                    .ant-menu-title-content {
                        font-weight: bold;
                    }
                }
            }            
            :deep(.ant-menu-submenu-selected .ant-menu-submenu-title) {
                background: url(@/assets/image/menu_bg.png) no-repeat;
                background-size: 100% 100%;
                height: 68px;
                color: #6971F8 !important;
                font-weight: bold;
            }            
            :deep(.ant-menu-item:active, .ant-menu-submenu-title:active) {
                background: transparent;
            }
            :deep(.ant-menu-item::after) {
                background: none !important;
            }
            :deep(.ant-menu-item-selected::after) {
                border: none !important;
            }
        }
        .sign-out {
            position: fixed;
            bottom: 0;
            height: 70px;
            width: @menuW;
            padding-bottom: 24px;
            background-color: #F5F7FF;
        }
    }
    
    .left-conent {
        width: calc(100% - @menuW);
        padding-left: 28px;
    }
}
</style>