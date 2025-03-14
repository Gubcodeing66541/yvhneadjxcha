<template>
    <div class="userInfo-content">
        <div class="tab-box flex">
            <div class="item" :class="tabIndex==index ? 'active':''" v-for="(item, index) in tab" :key="index" @click="tabHland(index)">{{item}}</div>
        </div>
        <div class="cont-box">
            <div class="userinfo-box ptb20 plr10" v-if="tabIndex == 0">
                <div class="flex just-between algin-center">
                    <span>基本资料</span>
                    <a-popover placement="bottomRight">
                        <template #content>
                            <div class="popover-pointer">
                                <div class="item" @click="addBlack('user')">用户拉黑</div>
                                <div class="item" @click="addBlack('ip')">IP 拉黑</div>
                            </div>
                        </template>
                        <a-button type="dashed" size="small">加入黑名单</a-button>
                    </a-popover>                    
                </div>
                <div class="info-cont mt20">
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">备注</span>
                        </a-avatar>
                        <div class="input-box">
                            <a-input v-model:value="detail.rename" :maxlength="30" placeholder="请输入备注" :bordered="false" @blur="changeDetail(true)" />
                        </div>
                    </div>
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">位置</span>
                        </a-avatar>
                        <div class="input-box">
                            <a-input v-model:value="detail.map" :bordered="false" readOnly />
                        </div>
                    </div>
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">IP</span>
                        </a-avatar>
                        <div class="input-box">
                            <a-input v-model:value="detail.ip" readOnly :bordered="false" />
                        </div>
                    </div>
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">设备</span>
                        </a-avatar>
                        <div class="input-box">
                            <a-input v-model:value="detail.drive" readOnly :bordered="false" />
                        </div>
                    </div>
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">电话</span>
                        </a-avatar>
                        <div class="input-box">
                            <a-input v-model:value="detail.mobile" placeholder="请输入电话" :bordered="false" @blur="changeDetail" />
                        </div>
                    </div>
                    <div class="item flex">
                        <a-avatar :size="46" class="shrink" :style="bgColorFor()">
                            <span class="f12">标签</span>
                        </a-avatar>
                        <div class="input-box layout">
                            <a-textarea v-model:value="detail.tag" placeholder="" :bordered="false" :rows="6" style="background:white" @blur="changeDetail" />
                        </div>
                    </div>
                </div>
            </div>
            <div class="record-cont plr20" v-else>
                <p class="tc mt20 strong" v-if="tracks">总共扫码 {{ tracks.length }} 次</p>
                <div class="list-box scrollbar pt20">
                    <a-timeline>
                        <a-timeline-item v-for="item in tracks" :key="item">
                            <p>扫码进入 {{ item.create_time }}</p>
                            <p class="t-tipc">IP : {{ item.ip }}</p>
                        </a-timeline-item>
                    </a-timeline>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, watch } from "vue";
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
import bus from "@/utils/bus";
export default {
    name: "userInfo",
    props: ['userInfo','login_log'],
    setup(props, contact) {
        const state = reactive({           
            detail:{}, //用户详情
            tab: ["用户资料","用户轨迹"],
            tabIndex:0, //菜单索引
            listForm: { page: 1, size:20 }, //分页
            pageCount: 2, //总页数
            tracks: [], //用户轨迹
        });

        watch(() => props.userInfo, (v1) => {
            init();
        }, { deep: true })

        onMounted(()=>{
            init();
        })

        //初始化
        const init = () => {
            tabHland(0);
            state.detail = JSON.parse(JSON.stringify(props.userInfo));
            state.tracks = props.login_log;
        }

        //菜单切换
        const tabHland = (index) => {
            state.tabIndex = index
        }
        
        //拉黑
        const addBlack = async(type) => {
            const res = await axios.post("/rooms/black", { user_id: state.detail['user_id'], ip: state.detail['ip'], is_black: 1, type }, false);
            bus.emit("blackUser", state.detail['user_id']);
        }

        //修改详细信息
        const changeDetail = async(refresh=false) => {
            const res = await axios.post("/rooms/update", state.detail, false);
            if (refresh){ 
                contact.emit("changeDetail","")
                bus.emit("onRefresh", '');
            }
        }

        //随机渐变色
        const bgColorFor = () => {
            let R = Math.floor(Math.random() * 130 + 110);
            let G = Math.floor(Math.random() * 130 + 110);
            let B = Math.floor(Math.random() * 130 + 110);
            return {
                background: 'rgba(' + R + ',' + B + ',' + G + ',0.5)'
            };
        }

        //监听滚动(暂时没有分页)
        const listScroll = (e) => {
            if (e.srcElement.clientHeight + e.srcElement.scrollTop >= e.srcElement.scrollHeight) { //滚动到底部
                state.listForm.page += 1;
                message.success("滚动到底部");
            }
        }

        return {
            ...toRefs(state),
            tabHland, //菜单切换
            bgColorFor, //颜色随机
            addBlack, //拉黑
            changeDetail, //修改详细信息
        };
    },
};
</script>

<style lang="less" scoped>
.userInfo-content{
    min-width: 200px;
    height: 100%;
    .tab-box {
        height: 50px;
        border-bottom: 1px solid #e1e1e1;
        .item {
            line-height: 50px;
            border: 1px solid #e1e1e1;
            padding: 0 16px;
            cursor: pointer;
            border-radius: 8px 8px 0 0;
            color: #999999;
            margin-bottom: -1px;
            &:hover {
                font-weight: bold;
            }
        }
        .active {
            border-bottom: none;
            font-weight: bold;
            background: #f5f5f5;            
            color: #333333;
        }
    }
    .cont-box{
        height: calc(100% - 50px);
        .info-cont {
            height: 100%;
            padding: 0 10px;
            .item {
                margin-bottom: 30px;
    
                .input-box {
                    border-bottom: 1px solid white;
                }
    
                .shrink {
                    margin-right: 20px;
                    color: #fff;
                }
            }
        }
        .record-cont{
            height: 100%;
            .list-box{
                height: calc(100% - 50px);
            }
            
        }
    }
}
</style>