<template>
    <div class="content-box flex">
        <div class="left">
            <div class="sma-box bg-white pointer mb40 t-3s" @click="dayBind(1)" title="点击查看客服访客数">
                <p><i class="ico icon_jrfks"></i></p>
                <p class="f20 mt40">今日访客数</p>
                <p class="f30 mt10">{{ data.count.room_cnt }}</p>
            </div>
            <div class="sma-box bg-white pointer t-3s" @click="dayBind(2)" title="点击查看客服消息数">
                <p><i class="ico icon_jrxxs"></i></p>
                <p class="f20 mt40">今日消息数</p>
                <p class="f30 mt10">{{ data.count.message_cnt }}</p>
            </div>
        </div>
        <div class="right bg-white">
            <div id="myChartnum" class="layout" style="height: 100%"></div>
        </div>
    </div>

    <!-- 今日访客弹窗 -->
    <a-modal v-model:visible="visible.jrfk_visible" :centered="true" :closable="false" :footer="null">
        <div>
            <div class="hh300">
                <a-table :dataSource="tabDetail.dataSource" :columns="tabDetail.columns" :scroll="{ y: 300 }" :pagination=false />
            </div>
            <div class="tc pt20">
                <a-pagination size="small" :total="tabDetail.total" v-model:current="pagination.page" v-model:page-size="pagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="change" />
            </div>
        </div>
    </a-modal>
    <!-- 今日消息弹窗 -->
    <a-modal v-model:visible="visible.jrxx_visible" :centered="true" :closable="false" :footer="null">
        <div>
            <div class="hh300">
                <a-table :dataSource="tabDetail.xx_dataSource" :columns="tabDetail.xx_columns" :scroll="{ y: 300 }" :pagination=false />
            </div>
            <div class="tc pt20">
                <a-pagination size="small" :total="tabDetail.xxTotal" v-model:current="xx_pagination.page" v-model:page-size="xx_pagination.offset" :showTotal="(total) => `共有 ${total} 条数据`" show-quick-jumper @change="xxChange" />
            </div>
        </div>
    </a-modal>
</template>

<script>
import { reactive, toRefs, onMounted, onBeforeUnmount } from "vue";
import * as echarts from 'echarts';
import axios from "@/utils/axios";
const columns = [ //今日访客表格配置
    { title: '客服名称', dataIndex: 'name' },
    { title: '今日访客数量', dataIndex: 'room_cnt' }
]
const xx_columns = [ //今日消息表格配置
    { title: '客服名称', dataIndex: 'name' },
    { title: '今日消息数', dataIndex: 'message_cnt' }
]
export default {
    name: "record",
    setup() {
        var myChart = 1;
        var option = null;
        const state = reactive({
            timer: null, //定时器
            data: { count: {}, }, //首页数据
            tabDetail: {
                total: 0, //总条数
                xxTotal: 0, //总条数
                columns: columns, //今日消息表格配置
                xx_columns: xx_columns, //今日消息表格配置
                dataSource: [], //今日访客表格数据
                xx_dataSource: [], //今日访客表格数据
            },
            pagination: { page: 1, offset: 9 }, //访客分页
            xx_pagination: { page: 1, offset: 9 }, //消息分页
            visible: { //弹窗显示隐藏
                jrfk_visible: false, //今日访客弹窗
                jrxx_visible: false, //今日消息弹窗
            },
        });

        onMounted(() => {
            init();
        });

        onBeforeUnmount(() => {
            clearInterval(state.timer);
            state.timer = null;
        })

        //统计数据
        const init = async () => {
            const res = await axios.post("/count", '', false);
            state.data = res.data;
            chartPayInit();            
            state.timer = setInterval(async () => {
                const res = await axios.post("/count", '', false);
                if (JSON.stringify(state.data.time_key) == JSON.stringify(res.data.time_key) && JSON.stringify(state.data.user_online) == JSON.stringify(res.data.user_online) && JSON.stringify(state.data.service_online) == JSON.stringify(res.data.service_online)) return false
                state.data = res.data;
                chartPayInit();
            }, 5000);
        };

        //统计图初始化
        const chartPayInit = () => {
            myChart = echarts.init(document.getElementById("myChartnum"));
            option = {
                // title: { text: " " },
                tooltip: {
                    trigger: "axis",
                    axisPointer: {
                        type: "cross",
                        label: {
                            backgroundColor: "#6a7985",
                        },
                    },
                },
                legend: {
                    data: ["当前在线访客", "当前在线客服"],
                },
                toolbox: {
                    show: false,
                },
                grid: {
                    left: "3%",
                    right: "4%",
                    bottom: "3%",
                    containLabel: true,
                },
                xAxis: [
                    {
                        type: "category",
                        boundaryGap: false,
                        data: state.data.time_key,
                    },
                ],
                yAxis: [
                    {
                        type: "value",
                        splitLine: { show: false },
                    },
                ],
                series: [
                    {
                        name: "当前在线访客",
                        type: "line",
                        showSymbol: false,
                        smooth: true,
                        areaStyle: {
                            color: {
                                type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
                                colorStops: [
                                    {
                                        offset: 0, color: 'rgba(101, 109, 236, 0.8)' // 0% 处的颜色
                                    },
                                    {
                                        offset: 1, color: 'rgba(232, 255, 254, 0)' // 100% 处的颜色
                                    }
                                ],
                                global: false // 缺省为 false
                            },
                        },
                        color: "#656DEC", //改变折线点的颜色
                        // emphasis: {
                        //     focus: "series",
                        // },
                        data: state.data.user_online,
                    }, {
                        name: "当前在线客服",
                        type: "line",
                        showSymbol: false,
                        smooth: true,
                        areaStyle: {
                            color: {
                                type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
                                colorStops: [
                                    {
                                        offset: 0, color: 'rgba(253, 174, 99, 0.8)' // 0% 处的颜色
                                    },
                                    {
                                        offset: 1, color: 'rgba(232, 255, 254, 0)' // 100% 处的颜色
                                    }
                                ],
                                global: false // 缺省为 false
                            },
                        },
                        color: "#FDAB61", //改变折线点的颜色
                        // emphasis: {
                        //     focus: "series",
                        // },
                        data: state.data.service_online,
                    },
                ],
                animationDuration: 0,
                animationDurationUpdate: 3000,
                animationEasing: 'linear',
                animationEasingUpdate: 'linear'
            }
            myChart.setOption(option);
            window.onresize = function () { // 自适应大小                
                myChart.resize();
            };
        };

        //今日数据
        const dayBind = (type) => {
            if (type == 1) {
                state.visible.jrfk_visible = true
                state.pagination.page = 1;
                getUserList();
            } else {
                state.visible.jrxx_visible = true
                state.pagination.xx_pagination = 1;
                getMesgList();
            }
        }

        //获取今日客服访客数量
        const getUserList = async () => {
            const res = await axios.post("/count_room_detail", state.pagination, false);
            state.tabDetail.total = res.data.count;
            state.tabDetail.dataSource = res.data.list;
        }

        //获取今日客服消息数量
        const getMesgList = async () => {
            const res = await axios.post("/count_message_detail", state.pagination, false);
            state.tabDetail.xxTotal = res.data.count;
            state.tabDetail.xx_dataSource = res.data.list;
        }

        //访客分页
        const change = (e) => {
            state.pagination.page = e;
            getUserList(); //获取列表
        }

        //消息分页
        const xxChange = (e) => {
            state.xx_pagination.page = e;
            getMesgList(); //获取列表
        }

        return {
            ...toRefs(state),
            dayBind, //今日数据
            change,
            xxChange,
        };
    },
};
</script>

<style lang="less" scoped>
.content-box {
    padding-bottom: 90px;

    .sma-box {
        text-align: center;
        width: 368px;
        height: calc(50% - 20px);
        border-radius: 42px;
        box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.2);
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;

        &:hover {
            box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.5);
        }
    }

    .right {
        width: 100%;
        height: 100%;
        margin-left: 40px;
        box-shadow: 0px 3px 6px 1px rgba(211, 229, 233, 0.2);
        border-radius: 42px;
        padding: 20px 0 0 0;
    }
}
</style>