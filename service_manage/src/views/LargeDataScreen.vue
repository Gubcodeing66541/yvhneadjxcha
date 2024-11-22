<template>
    <div class="warp">
        <div class="container flex dir-column pd10 scrollbar">
            <div class="header-box flex just-between">
                <div class="item icon_k_lf1 pl60">
                    欢迎来到在线客服实时监控平台！
                </div>
                <div class="strong f30 pt10 title">在线客服实时监控平台</div>
                <div class="item icon_k_lr1 tr pr60">
                    {{ time }}
                </div>
            </div>
            <div class="mt14 mb7"><div style="width: 150px; height: 1px; background-color: #3F96A5" /></div>
            <div class="space-box flex flex1 just-between">
                <div class="box shrink">                    
                    <div class="flex just-center ptb8 star" style="border: 1px solid #0C77E6;">
                        <div class="w50">
                            <p class="f16 tc" style="color:#0C77E6">总访客数</p>
                            <p class="f16 tc" style="color:#34BFA9">{{ data.top.all_user_cnt }}个</p>
                        </div>
                        <div class="w50">
                            <p class="f16 tc" style="color:#0C77E6">总消息数</p>
                            <p class="f16 tc" style="color:#34BFA9">{{ data.top.all_message_cnt }}个</p>
                        </div>
                    </div>
                    <div class="tab-box layout icon_k_pm_b ptb6 plr16 mt6">
                        <div class="ptb6 mb10 t-title">
                            <i class="ico icon_k_user mr10"></i>实时访客列表
                        </div>
                        <div class="table flex dir-column">
                            <a-row type="flex" justify="center" class="item">
                                <a-col :span="5">头像</a-col>
                                <a-col :span="5">昵称</a-col>
                                <a-col :span="8">地区</a-col>
                                <a-col :span="6">客服名称</a-col>
                            </a-row>
                            <div class="tab-cont">
                                <a-row type="flex" justify="center" class="item" v-for="item in data.user_online" :key="item">
                                    <a-col :span="5"><a-avatar shape="square" :size="30" :src="item.user_head" style="background: #ffffff;" /></a-col>
                                    <a-col :span="5" class="text1">{{ item.user_name }}</a-col>
                                    <a-col :span="8" class="text1">{{ item.map }}</a-col>
                                    <a-col :span="6" class="text1 pl6">{{ item.service_name }}</a-col>
                                </a-row>
                            </div>
                        </div>
                    </div>
                    <div class="tab-box layout icon_k_pm_b ptb11 plr16 mt6">                        
                        <div class="ptb6 mb10 t-title">
                            <i class="ico icon_k_pm mr10"></i>客服访客数排名
                        </div>
                        <div class="table flex dir-column">
                            <a-row type="flex" justify="center" class="item">
                                <a-col :span="5">头像</a-col>
                                <a-col :span="8">昵称</a-col>
                                <a-col :span="5">访客数</a-col>
                                <a-col :span="6">排名</a-col>
                            </a-row>
                            <div class="tab-cont">
                                <a-row type="flex" justify="center" class="item" v-for="(item, index) in data.service_rank" :key="index">
                                    <a-col :span="5"><a-avatar shape="square" :size="30" :src="item.head" style="background: #ffffff;" /></a-col>
                                    <a-col :span="8" class="text1">{{ item.name }}</a-col>
                                    <a-col :span="5" class="text1 pl6">{{ item.cnt }}</a-col>
                                    <a-col :span="6" class="text1 pl6">{{ index+1 }}</a-col>
                                </a-row>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="box layout plr4">
                    <div class="layout icon_k_t hh67 plr6 flex align-center just-between">
                        <div class="ico icon_k_btn1 flex just-between align-center plr7">
                            <i class="ico icon_k_left"></i>
                            <span class="f12">今日访客：{{ data.top.today_user_cnt }}</span>
                            <i class="ico icon_k_right"></i>
                        </div>
                        <span class="f18" style="color:#34BFA9">总计{{ data.top.all_message_cnt }}条消息</span>
                        <div class="ico icon_k_btn1 flex just-between align-center plr7">
                            <i class="ico icon_k_left"></i>
                            <span class="f12">今日消息数：{{ data.top.today_message_cnt }}</span>
                            <i class="ico icon_k_right"></i>
                        </div>
                    </div>
                    <div class="icon_k_map relative mlr3 mt20">
                        <div class="head-China flex">
                            <div class="icon_fx left flex align-bottom">
                                <i class="item"></i>
                                <i class="item"></i>
                                <i class="item"></i>
                            </div>
                            <div class="ico icon_k_btn2 tc">中国访客人数分布</div>
                            <div class="icon_fx right flex align-bottom">
                                <i class="item"></i>
                                <i class="item"></i>
                                <i class="item"></i>
                            </div>
                        </div> 
                        <div class="map-box h100">
                            <div id="chinaMap" style="width: 100%; height: 100%;"></div>
                        </div>                        
                    </div>
                </div>
                <div class="box shrink">
                    <div class="flex ptb8 pl16 star" style="border: 1px solid #0C77E6;">
                        <div>
                            <p class="f16 tc" style="color:#0C77E6">在线客服数</p>
                            <p class="f16 tc" style="color:#34BFA9">{{ data.top.online_service }}个</p>
                        </div>
                        <div class="ml100">
                            <p class="f16 tc" style="color:#0C77E6">在线访客数</p>
                            <p class="f16 tc" style="color:#34BFA9">{{ data.top.online_user }}个</p>
                        </div>
                    </div>
                    <div class="tab-box layout icon_k_pm_b ptb6 plr16 mt6">
                        <div class="ptb6 mb10 t-title">
                            <i class="ico icon_k_user mr10"></i>实时在线客服
                        </div>
                        <div class="table flex dir-column">
                            <a-row type="flex" justify="center" class="item">
                                <a-col :span="5">头像</a-col>
                                <a-col :span="9">昵称</a-col>
                                <a-col :span="10">账号</a-col>
                            </a-row>
                            <div class="tab-cont">
                                <a-row type="flex" justify="center" class="item" v-for="item in data.service_online" :key="item">
                                    <a-col :span="5"><a-avatar shape="square" :size="30" :src="item.head" style="background: #ffffff;" /></a-col>
                                    <a-col :span="9" class="text1">{{ item.name }}</a-col>
                                    <a-col :span="10" class="text1">{{ item.username }}</a-col>
                                </a-row>
                            </div>
                        </div>
                    </div>
                    <div class="tab-box layout icon_k_pm_b ptb6 plr16 mt6">
                        <div class="ptb6 mb10 t-title">
                            <i class="ico icon_k_tab mr10"></i>最近7天访客数
                        </div>
                        <div class="table flex dir-column">
                            <a-row type="flex" justify="center" class="item">
                                <a-col :span="12">时间</a-col>
                                <a-col :span="12">访客数</a-col>
                            </a-row>
                            <div class="tab-cont">
                                <a-row type="flex" justify="center" class="item" v-for="item in data.count_time" :key="item">
                                    <a-col :span="12" class="text1">{{ item.time }}</a-col>
                                    <a-col :span="12" class="text1 pl10">{{ item.cnt }}</a-col>
                                </a-row>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { onMounted, reactive, toRefs, onBeforeUnmount } from "vue";
import { getNowTime } from "@/utils";
import * as echarts from 'echarts';
import chinaJSON from '@/utils/china.json';
import axios from "@/utils/axios";
export default {
    name: "largeDataScreen",
    setup() {        
        const state = reactive({
            timer:null,
            dataTimer:null,
            time:'',
            data: {
                region:[], //地图
                top:{}, //顶上数据
                service_online: [], //实时在线客服
                user_online:[], //实时访客
                service_rank:[], //排名
                count_time:[], //最近7天访客数
            },
        });

        onMounted(()=>{
            document.title = "数据大屏"; //标题
            dataTime();
            getTime();            
        })

        onBeforeUnmount(()=>{
            clearInterval(state.timer);
            clearInterval(state.dataTimer);
            state.timer = null;
            state.dataTimer = null;
        })

        //定时请求数据
        const dataTime = () => {
            getData();
            state.dataTimer = setInterval(() => {
                getData();
            }, 5000);
        }

        //获取数据
        const getData = async() => {
            const res = await axios.post("/data", '', false);
            state.data.region = res.data.region;
            state.data.top = res.data.top;
            state.data.service_online = res.data.service_online;
            state.data.user_online = res.data.user_online;
            state.data.service_rank = res.data.service_rank;
            state.data.count_time = res.data.count_time;
            chinaMap(); //中国地图
        }

        //地图
        const chinaMap = () => {
            const myChart = echarts.init(document.getElementById("chinaMap"));
            echarts.registerMap('china', chinaJSON) //注册可用的地图
            const option = {
                 geo: {
                    map: 'china',
                    roam: true, //是否允许缩放，拖拽
                    zoom: 1.6, //初始化大小
                    center: [105, 36],
                    scaleLimit: { min: 1.7, max: 6 },
                    itemStyle: { areaColor: '#354EC4', color: '#ffffff', borderColor: '#A7AECD', borderWidth: 1 },
                    label: { show: true, fontSize: 10, color: '#EDEEF9' },               
                    emphasis: { itemStyle: { areaColor: '#1D297D', color: '#fff' },label:{color:"#fff"} }, //鼠标移入样式
                    // regions: [{ name: '黑龙江省', label: { show: true, fontSize: 10, color: '#EDEEF9' }, itemStyle: { areaColor: '#1D297D' } }, { name: '西藏自治区', label: { show: true, fontSize: 10, color: '#EDEEF9' }, itemStyle: { areaColor: '#1D297D' } }], //选中的地方
                },
                series: [
                    {
                        type: "scatter",
                        coordinateSystem: "geo",
                        symbol: "pin",
                        legendHoverLink: true,
                        symbolSize: [70, 70],
                        label: { 
                            show: true, 
                            formatter(value) { return value.data.value[2] }, 
                            color: "#fff",
                        },
                        // 标志的样式
                        itemStyle: {
                            color: "rgba(255,0,0,.7)",
                            shadowBlur: 2,
                            shadowColor: "D8BC37",
                        },
                        // 数据格式，其中name,value是必要的，value的前两个值是数据点的经纬度，其他的数据格式可以自定义
                        // 至于如何展示，完全是靠上面的formatter来自己定义的
                        data: convertData(),
                        showEffectOn: "render",
                        rippleEffect: { brushType: "stroke", },
                        zlevel: 1,
                    },
                ],
            }
            myChart.setOption(option)
            window.onresize = function () { // 自适应大小                
                myChart.resize();
            };
        }

        //位置经纬度
        const convertData = () => {
            let data = state.data.region;
            var res = [];
            var geoCoordMap = {
                北京: [116.4, 40.4],
                天津: [117.04, 39.52],
                河北: [115.21, 38.44],
                山西: [111.95, 37.65],
                内蒙古: [112.17, 42.81],
                辽宁: [123.42, 41.29],
                吉林: [126.32, 43.38],
                黑龙江: [128.34, 47.05],
                上海: [121.46, 31.28],
                江苏: [120.26, 32.54],
                浙江: [120.15, 29.28],
                安徽: [117.28, 31.86],
                福建: [118.31, 26.07],
                江西: [115.89, 27.97],
                山东: [118.01, 36.37],
                河南: [113.46, 34.25],
                湖北: [112.29, 30.98],
                湖南: [112.08, 27.79],
                广东: [113.98, 22.82],
                广西: [108.67, 23.68],
                海南: [110.03, 19.33],
                重庆: [107.51, 29.63],
                四川: [103.36, 30.65],
                贵州: [106.91, 26.67],
                云南: [101.71, 24.84],
                西藏: [89.13, 30.66],
                陕西: [108.94, 34.46],
                甘肃: [103.82, 36.05],
                青海: [97.07, 35.62],
                宁夏: [106.27, 36.76],
                新疆: [86.61, 40.79]
            };
            for (var i = 0; i < data.length; i++) {
                var geoCoord = geoCoordMap[data[i].map];
                if (geoCoord) {
                    res.push({
                        name: data[i].map,
                        value: [...geoCoord, data[i].cnt]
                    });
                }
            }
            return res;
        }

        //获取实时时间
        const getTime = async() => {
            clearInterval(state.timer)
            state.timer = null;
            state.time = await getNowTime();
            state.timer = setInterval(async() => {
                state.time = await getNowTime();
            }, 1000);
        }

        return {
            ...toRefs(state),
        };
    },
};
</script>

<style lang="less" scoped>
.warp{
    color: #FFFFFF;
    position: fixed;
    width: 100vw;
    height: 100vh;
    background-color: #060916;
    background-image: url(../assets/image/dp_bg.png);
    background-size: 100% 100%;
    .container{
        width: 100%;
        height: 100%;
        .header-box{
            .item{
                width: 30%;
                height: 35px;
                line-height: 35px;
            }
        }
        .box{
            display: flex;
            flex-direction: column;
        }
        .space-box{
            overflow: hidden;
            .shrink{
                width: 370px;                
            }
            .tab-box{
                height: calc(50% - 40px);
            }
            .table{
                height: calc(100% - 55px);
                .item{
                    line-height: 36px;
                    &:nth-of-type(odd) {
                        background: #0F1325;
                    }
                
                    &:nth-of-type(even) {
                        background: #171C33;
                    }
                }
                .tab-cont{
                    height: 100%;
                    overflow-y: scroll;
                }
            }
            .ico-tipc{
                i{
                    width: 8px;
                    height: 2px;
                    background: #0C77E6;
                    border-radius: 1px 1px 1px 1px;
                    margin-right: 3px;
                }
            }
        }
        .icon_k_map {
            height: calc(100% - 80px);
            .head-China {
                position: absolute;
                top: -15px;
                left: calc(50% - 194px);
                .icon_fx {
                    .item {
                        margin-right: 10px;
                        width: 20px;
                        height: 8px;
            
                        &:nth-child(2) {
                            background-color: #3D4F79;
                        }
                    }
                }
            
                .left {
                    transform: skew(30deg);
                    margin-right: 10px;
            
                    .item {
                        &:nth-child(1) {
                            background-color: #232E4B;
                        }
            
                        &:nth-child(3) {
                            background-color: #647DBB;
                        }
                    }
                }
            
                .right {
                    transform: skew(-30deg);
                    margin-left: 20px;
            
                    .item {
                        &:nth-child(1) {
                            background-color: #647DBB;
                        }
            
                        &:nth-child(3) {
                            background-color: #232E4B;
                        }
                    }
                }
            }
            
        }
                  
    }
}
.flex{
    display: flex !important;
}
:deep(.ant-progress-inner){
    border-radius: 0;
    border: 2px solid #fff;
    border-image: linear-gradient(to right, #5CD9E8, #BB4CF6) 1;
}
:deep(.ant-progress-text){
    color: #D6D6D6;
    font-size: 10px;
}
.icon_k_lf1 {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_lf1.png);
}
.icon_k_lr1 {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_lr1.png);
}
.icon_k_user_b{
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_user_b.png);
}
.icon_k_pm_b {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_pm_b.png);
}
.icon_k_t {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_t.png);
}
.icon_k_map {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_map.png);
}
.icon_k_b {
    background-size: 100% 100%;
    background-image: url(../assets/image/icon_k_b.png);
}
.ico.icon_k_yu {
    width: 16px;
    height: 15px;
    background-image: url(../assets/image/icon_k_yu.png);
}
.ico.icon_k_user {
    width: 10px;
    height: 11px;
    background-image: url(../assets/image/icon_k_user.png);
}
.ico.icon_k_pm {
    width: 10px;
    height: 10px;
    background-image: url(../assets/image/icon_k_pm.png);
}
.ico.icon_k_hq {
    width: 10px;
    height: 11px;
    background-image: url(../assets/image/icon_k_hq.png);
}
.ico.icon_k_btn1 {
    width: 217px;
    height: 40px;
    background-image: url(../assets/image/icon_k_btn1.png);
}
.ico.icon_k_left {
    width: 6px;
    height: 20px;
    background-image: url(../assets/image/icon_k_left.png);
}
.ico.icon_k_right {
    width: 6px;
    height: 20px;
    background-image: url(../assets/image/icon_k_right.png);
}
.ico.icon_k_btn2 {
    width: 189px;
    height: 32px;
    line-height: 32px;
    background-image: url(../assets/image/icon_k_btn2.png);
}
.ico.icon_k_pm {
    width: 10px;
    height: 8px;
    background-image: url(../assets/image/icon_k_pm.png);
}
.ico.icon_k_user {
    width: 10px;
    height: 11px;
    background-image: url(../assets/image/icon_k_user.png);
}
.ico.icon_k_tab {
    width: 10px;
    height: 10px;
    background-image: url(../assets/image/icon_k_tab.png);
}
.title {
    -webkit-animation: shining 1s alternate infinite;
    animation: shining 1s alternate infinite;
}
@keyframes shining {
    from {
        text-shadow: 0 0 6px #056c86, 0 0 9px #056c86, 0 0 12px #056c86, 0 0 15px #056c86, 0 0 18px #056c86, 0 0 21px #056c86;
    }

    to {
        text-shadow: 0 0 1px #056c86, 0 0 2px #056c86, 0 0 3px #056c86, 0 0 4px #056c86, 0 0 5px #056c86, 0 0 6px #056c86;
    }
}
.star {
    position: relative;
}
.star::before {
    content: "";
    position: absolute;
    display: block;
    height: 6px;
    width: 6px;
    background: #507ca9;
    box-shadow: 0px 0px 10px 3px #eff4f9;
    left: -3px;
    top: -3px;
    transition: all 0.1s linear;
    border-radius: 100%;
    animation: star 10s linear infinite;
}
@keyframes star {
    0% {
        left: -3px;
        top: -3px;
    }

    25% {
        left: calc(100% - 3px);
        top: -3px;
    }

    50% {
        left: calc(100% - 3px);
        top: calc(100% - 3px);
    }

    75% {
        left: -3px;
        top: calc(100% - 3px);
    }

    100% {
        left: -3px;
        top: -3px;
    }
}
.t-title{
    color: #b1fcfc;
}
</style>