<template>
    <div class="com-content">
        <div class="top f18 pt12 pb10 pl20 strong bb-tipc flex align-center just-between">
            <span>{{ formMsg.name }} 与 {{ formMsg.user_name}}的聊天记录</span>
        </div>
        <div class="mesg-box flex layout">
            <div id="chatBox" class="chat-box left pt12 pb12 scrollbar" @scroll="divScroll">
                <div class="item" :class="item.send_role == 'service' ?'item-right':''" v-for="item in list" :key="item">
                    <div v-if="item.type != 'time'">
                        <p class="t-time">{{ item.create_time }}</p>
                        <div class="message-row mt10">
                            <div class="message-box">
                                <div class="pt10 pl16 pb10 pr16 ico-arrow" v-if="item.type=='text'">
                                    <div style="white-space:pre-wrap" v-html="textReplaceEmoji(item.content)"></div>
                                </div>
                                <div class="img-cont" v-else-if="item.type == 'image'">
                                    <a-image :src="item.content" :height="200" :placeholder="true" />
                                </div>
                                <div class="video-cont" v-else-if="item.type == 'video'">
                                    <video :src="item.content" :height="200" controls style="object-fit:fill;width: 100%;"></video>
                                </div>
                                <div class="code-cont pd10" v-else>
                                    <div class="tc" style="height:240px">
                                        <vue-qr :logoSrc="formMsg.service_head" :text="item.content" :logoMargin="4" :size="240"
                                            :margin="10"></vue-qr>
                                    </div>
                                    <p class="tc mt10 f14">长按保存专属二维码，聊天记录不丢失</p>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div v-else class="flex just-center">
                        <p class="bg-tipc tc t-tipc" style="width:200px;color: #91d382;">{{ item.content }}</p>
                    </div>
                </div>
            </div>
            <div class="right pl14 pr14 pc">
                <div class="tc mb20 mt20">
                    <a-avatar :src="formMsg.user_head" style="background-color: #f56a00" :size="90"></a-avatar>
                </div>
                <div class="bb-tipc pb30">
                    <p>
                        <span class="lable">昵称：</span>
                        <span class="t-them">{{ formMsg.user_name }}</span>
                    </p>
                    <p>
                        <span class="lable">IP：</span>
                        <span class="t-them">{{ formMsg.ip }}</span>
                    </p>
                </div>
                <div class="pt30">
                    <p>
                        <span class="lable">备注：</span>
                        <span>{{ formMsg.rename }}</span>
                    </p>
                    <p>
                        <span class="lable">位置：</span>
                        <span>{{ formMsg.map }}</span>
                    </p>
                    <p>
                        <span class="lable">设备：</span>
                        <span>{{ formMsg.drive }}</span>
                    </p>
                    <p>
                        <span class="lable">电话：</span>
                        <span>{{ formMsg.mobile }}</span>
                    </p>
                    <p>
                        <span class="lable">标签：</span>
                        <span>{{ formMsg.tag }}</span>
                    </p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { reactive, toRefs, onMounted, nextTick } from "vue";
import { message } from 'ant-design-vue';
import axios from "@/utils/axios";
import { textReplaceEmoji } from "@/utils/emojis";
import vueQr from 'vue-qr/src/packages/vue-qr.vue';
export default {
    name: "chatRecord",
    props: ['formMsg'],
    components: { vueQr },
    setup(props, context) {
        const state = reactive({
            scrollHeight: 0, //滚动条总高度
            form: { page: 1, offset: 20 },
            pageCount: 0, //聊天记录总页数
            list:[]
        });

        onMounted (()=>{
            getMessageList()       
        })

        //监听滚动
        const divScroll = (e) => {
            nextTick(() => {
                if (e.srcElement.scrollTop == 0 && state.pageCount > state.form.page) {
                    state.form.page += 1;
                    getMessageList(); //获取聊天记录
                    state.scrollHeight = e.srcElement.scrollHeight; //滚动条高度
                }
            });
        }

        //获取聊天记录
        const getMessageList = async () => {
            state.form['service_id'] = props.formMsg.service_id;
            state.form['user_id'] = props.formMsg.user_id;
            const res = await axios.post('/users/message', state.form, false);
            state.list = state.form.page == 1 ? res.data.list : res.data.list.concat(state.list);
            state.pageCount = res.data.page;
            nextTick(() => {
                if (state.form.page == 1){
                    document.getElementById('chatBox').scrollTop = document.getElementById('chatBox').scrollHeight; //滚动到底部
                }else{
                    document.getElementById('chatBox').scrollTop = document.getElementById('chatBox').scrollHeight - state.scrollHeight; //重置滚动高度
                }
            });
        }

        return {
            ...toRefs(state),
            divScroll,
            textReplaceEmoji, //文字表情处理
        };
    },
};
</script>

<style lang="less" scoped>
.com-content{
    .mesg-box{
        height: 538px;
        .left {
            width: calc(100% - 310px);
        }
        .chat-box{
            font-size: 16px;
            .item{
                overflow: hidden;
                margin-bottom: 20px;
                padding: 0 20px;
                .message-row{
                    overflow: hidden;
                }
                .t-tipc{
                    color: #AEAEAE;
                    height: 29px;
                    line-height: 29px;
                    background: #EEEEEE;
                    font-size: 14px;
                    padding: 0 20px;
                }
                .t-time{
                    color: #AEAEAE;
                }
                .message-box{
                    background: #F5F5FF;
                    border-radius: 5px;
                    float: left;
                    max-width: 50%;
                    overflow: hidden;
                    .ico-arrow{
                        position: relative;
                        &::after {
                            content: '';
                            position: absolute;
                            left: -5px;
                            top: calc(50% - 5px);
                            width: 10px;
                            height: 10px;
                            background: #F5F5FF;
                            transform: rotate(45deg);
                        }
                    }                   
                    .img-cont {
                        :deep(.ant-image-img) {
                            width: auto;
                            height: 100%;
                        }
                        :deep(.ant-image-img-placeholder) {
                            background: none;
                        }
                    }
                    .video-cont{
                        height: 200px;
                        overflow: hidden;
                        border-radius: 5px;
                    }
                    .code-cont{
                        height: 296px;
                    }
                }                
            }
            .item-right{
                .t-time{
                    text-align: right;
                }
                .message-box{
                    float: right;
                    .ico-arrow{
                        &::after {
                            left: auto;
                            right: -5px;
                        }
                    }                    
                }
            }
        }
        .right {
            width: 310px;
            background: #F8F8F8;
            font-size: 16px;
            p{
                padding: 6px 0;
            }
        }
    }    
}
</style>