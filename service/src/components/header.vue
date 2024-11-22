<template>
    <div class="header-container flex algin-center just-between plr10">
        <!-- 公告 -->
        <div class="notice">
            <div class=" flex algin-center" v-if="noticeFirst && noticeFirst!= ''">
                <div class="lh40 bg-white"><i class="iconfont icon-notice mr10"></i></div>
                <div class="text-cont">
                    <div class="text-box lh40 pointer" @click="visible.notice = true">
                        {{ noticeFirst }}
                    </div>
                </div>
            </div>            
        </div>

        <!-- 设置 -->
        <div class="setting flex algin-center">
            <a-popover>
                <template #content>
                    <div class="popover-pointer" style="width:220px">
                        <div class="item item-auto flex just-between">
                            <span>消息提醒</span>
                            <a-switch v-model:checked="reminder" checked-children="开" un-checked-children="关" />
                        </div>
                    </div>
                </template>
                <div class="item">
                    <sound-outlined />
                </div>
            </a-popover>
            
            <a-tooltip placement="bottom">
                <template #title>
                    <span>系统刷新</span>
                </template>
                <div class="item ml20" @click="reload">
                    <undo-outlined />
                </div>
            </a-tooltip>
            
            <a-tooltip placement="bottom">
                <template #title>
                    <span>下载应用</span>
                </template>
                <div class="item ml20">
                    <arrow-down-outlined @click="downloadApp" />
                </div>
            </a-tooltip>

            <a-tooltip placement="bottom">
                <template #title>
                    <span>操作视频</span>
                </template>
                <div class="item ml20" @click="visible.video = true">
                    <laptop-outlined />
                </div>
            </a-tooltip>

            <a-tooltip placement="bottom">
                <template #title>
                    <span>版本号{{ config.SystemVersion }}</span>
                </template>
                <div class="item ml20" @click="visible.video = true">
                    <windows-outlined />
                </div>
            </a-tooltip>

            <!-- 用户信息 -->
            <a-popover placement="bottomRight">
                <template #content>
                    <div class="popover-pointer" style="width:250px">
                        <div class="item item-auto flex just-between">
                            <span>授权码</span>
                            <span>{{ info.username }}</span>
                        </div>
                        <div class="item item-auto flex just-between">
                            <span>到期时间</span>
                            <span>{{ info.time_out }}</span>
                        </div>
                        <div class="item t-tipc" @click="loginOut">
                            <PoweroffOutlined />
                            <span class="ml10">退出登录</span>
                        </div>
                    </div>
                </template>
                <div class="name-box ml50 flex algin-center pointer">
                    <p class="text1 mr10 lh40" style="max-width:200px">{{ info.name }}</p>
                    <caret-down-outlined class="t-tipc lh48 t-3s" />
                </div>
            </a-popover>
        </div>

        <!-- 操作视频 -->
        <a-modal title="操作视频" :width="1200" v-model:visible="visible.video" :footer="null" :centered="true">
            <video controls="controls" style="width: 1150px;" :src="config.SystemVideo" v-if="visible.video"></video>
        </a-modal>

        <!-- 公告 -->
        <a-modal title="系统公告" :width="600" v-model:visible="visible.notice" :footer="null" :centered="true">
            <div style="min-height: 360px">
                <div class="notice-list" style="min-height: 330px">
                    <div class="item pointer br8" v-for="item in noticeData" :key="item">
                        <div class="text2">{{ item.value }}</div>
                        <p class="tr f12 mt4" style="color: #bab0e7;">{{ item.create_time }}</p>
                    </div>
                </div>
                <div class="tc">
                    <a-pagination size="small" :total="total" v-model:current="pagination.page" v-model:page-size="pagination.offset"
                        show-quick-jumper @change="change" />
                </div>
            </div>            
        </a-modal>
    </div>
</template>

<script>
import { reactive, toRefs, computed,onMounted } from "vue";
import { SoundOutlined, UndoOutlined, ArrowDownOutlined, LaptopOutlined, WindowsOutlined, CaretDownOutlined, LockOutlined, PoweroffOutlined } from '@ant-design/icons-vue';
import { localEmpty, localSet } from '../utils';
import { useStore } from "vuex";
import axios from "@/utils/axios";
import { websocketClose } from "@/socket";
export default {
    name: "myHeader",
    components: { SoundOutlined, UndoOutlined, ArrowDownOutlined, LaptopOutlined, WindowsOutlined, CaretDownOutlined, LockOutlined, PoweroffOutlined },
    setup() {
        const store = useStore();
        const reminder = computed({ //消息提醒开关
            get() {
                return store.state.messageReminder;
            },
            set(val) {
                localSet('messageReminder',val)
                store.commit('setMessageReminder', val);
            },
        });
        const state = reactive({
            info: computed(() => store.state.info), //账号信息
            config: computed(() => store.state.config), //系统配置信息
            visible:{}, //弹窗
            pagination: { page: 1, offset: 4, type:'notice' }, //公告分页
            total: 0, //总条数
            noticeData:[], //公告数据
            noticeFirst:"", //第一条公告
        });

        onMounted(()=>{
            getSetting(); //获取系统配置
            getNotice(); //获取公告
        })

        //系统配置
        const getSetting = async() => {
            const res = await axios.post("/config", '', false);
            localSet('config', res.data.config);
            store.commit("setConfig", res.data.config);
        }

        //获取公告
        const getNotice = async () => {
            const res = await axios.post("/setting", state.pagination, false);
            if (res.data.list && res.data.list!=''){
                state.noticeData = res.data.list;
                state.noticeFirst = res.data.list[0]['value']
                state.total = res.data.count;
            }            
        }

        //分页
        const change = (e) => {
            state.pagination.page = e;
            getNotice(); //获取列表
        }

        //下载应用
        const downloadApp = () => {
            window.open(state.config['PcDownloadLink']);
        }

        //系统刷新
        const reload = () => {
            window.location.reload();
        }

        //退出登录
        const loginOut = () => {
            websocketClose();
            localEmpty('token');
            window.location.reload();
        }

        return {
            ...toRefs(state),
            reminder, //消息提醒开关
            reload, //系统刷新
            downloadApp, //下载应用
            loginOut, //退出登录
            change, //页码监听
        };
    },
};
</script>

<style lang="less" scoped>
.header-container{
    height: 100%;
    .notice{        
        .icon-notice{
            color: #ffc931;
        }
        .text-cont{
            white-space: nowrap;
            overflow: hidden;
            .text-box {
                width: 50%;
                position: relative;
                width: fit-content;
                animation: move 20s linear infinite;
        
                &::after {
                    position: absolute;
                    right: -100%;
                    content: attr(text);
                }
            }
        
            @keyframes move {
                0% {
                    transform: translateX(5%);
                }
        
                100% {
                    transform: translateX(-100%);
                }
            }
        }        
    }    
    .setting{
        .item{
            display: flex;
            align-items: center;
            cursor: pointer;
            font-size: 16px;
            &:hover{
                color: #666;
            }
        }
    }
}
.notice-list {
    .item {
        border-left: 4px solid #d4ceed;
        background: #efeef5;
        padding: 10px 6px;
        margin-bottom: 10px;
        &:hover{
            background: #e3e2e9;
            .text2{
                -webkit-line-clamp: inherit;
            }
        }
    }
}
</style>