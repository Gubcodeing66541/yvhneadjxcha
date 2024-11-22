<template>
    <div class="editBox-content">
        <!-- 发送选项 -->
        <div class="option-box t-tipc f18">
            <!-- 表情 -->
            <a-popover trigger="click" placement="topLeft">
                <template #content>
                    <div class="emojis-box flex align-center flex-warp ww300">
                        <div v-for="(elImg, text) in emoji.emojis" v-html="elImg" :key="elImg" class="item br6" @click="emojisHland(text)"></div>
                    </div>
                </template>
                <SmileOutlined class="mr14 pointer" title="表情" />
            </a-popover>
            <!-- 发送图片 -->
            <a-upload accept=".png, .jpg, .jpeg" :show-upload-list="false" :before-upload="beforeUpload" :customRequest="uploadImgHland">
                <PictureOutlined class="mr14 t-tipc f18 pointer" title="发送图片" />
            </a-upload>
            <!-- 发送视频 -->
            <a-upload accept=".mp4" :show-upload-list="false" :before-upload="videoBeforeUpload" :customRequest="uploadHland">
                <VideoCameraOutlined class="mr14 t-tipc f18 pointer" title="发送视频" />
            </a-upload>
            <!-- 快捷回复 -->
            <ReadOutlined class="mr14 pointer" title="快捷回复" @click="visible.quickReply = true" />
            <!-- 一键撤回 -->
            <a-popconfirm title="撤回后将无法恢复，是否确定撤回所有记录？" ok-text="确定" cancel-text="取消" @confirm="withdraw()">
                <RotateLeftOutlined class="mr14 pointer" title="一键撤回" />
            </a-popconfirm>
            <qrcode-outlined class="mr14 pointer" title="向用户发送二维码" @click="sendCode" />
            <!-- <alert-outlined class="pointer" title="向用户发起提醒" @click="sendRemind" /> -->
        </div>
        <!-- 文本编辑 -->
        <div class="mt10">
            <a-textarea v-model:value="sendText" :rows="4" :bordered="false" @keydown.enter="sendMessage" v-drag="handleFiles" v-paste="handleFiles" />
        </div>
        <!-- 发送文本 -->
        <div class="tr mt10">
            <a-button type="primary" class="send-btn" title="按Enter键发送，按shift+Enter回车" @click="sendMessage">发送</a-button>
        </div>
    </div>

    <!-- 快捷回复 -->
    <a-drawer v-model:visible="visible.quickReply" :width="drawer.width" :closable="false" class="custom-class" title="快捷回复" placement="right" :bodyStyle="drawer.bodyStyle">
        <quickReply :userId="user_id" @sendQuick="sendQuick" v-if="visible.quickReply"></quickReply>
    </a-drawer>
</template>

<script>
import { reactive, toRefs, createVNode, computed } from "vue";
import { beforeUpload, videoBeforeUpload, uploadFile, getTime } from '@/utils';
import { emojiList as emoji } from "@/utils/emojis";
import { SmileOutlined, PictureOutlined, VideoCameraOutlined, ReadOutlined, RotateLeftOutlined, ExclamationCircleOutlined, QrcodeOutlined, AlertOutlined } from '@ant-design/icons-vue';
import { Modal, message } from 'ant-design-vue';
import quickReply from '@/components/menu/quickReply.vue';
import axios from "@/utils/axios";
import { useStore } from "vuex";
export default {
    name: "editBox",
    components: { SmileOutlined, PictureOutlined, VideoCameraOutlined, ReadOutlined, RotateLeftOutlined, ExclamationCircleOutlined, quickReply, QrcodeOutlined, AlertOutlined },
    props: ['user_id'],
    emits: ['funHland'],
    setup(props, contact) {
        const store = useStore();
        const state = reactive({
            info: computed(() => store.state.info), //系统信息
            sendText: "", //文本编辑
            visible: {}, //弹窗显示
            drawer: { bodyStyle: { backgroundColor: '#f5f5f5' }, width: 500 }
        });

        //表情选中
        const emojisHland = (text) => {
            state.sendText += text;
        }

        //上传图片
        const uploadImgHland = async (files) => {
            const res = await uploadFile(files.file);
            messagePush('image', res);
        }

        //上传文件
        const uploadHland = async (files) => {
            const res = await uploadFile(files.file);
            messagePush('video', res);
        };

        //撤回所有聊天记录
        const withdraw = async() => {
            contact.emit('funHland', { type:'withdraw'});
            const res = await axios.post("/message/clear_message", { user_id: props.user_id },false);
        }

        //发送文本信息
        const sendMessage = (e) => {
            if (e.shiftKey === true) return //shift+回车
            e.preventDefault(); // 阻止浏览器默认换行操作
            if (state.sendText != '' && state.sendText.split(" ").join("").length) {
                messagePush('text', state.sendText);
                state.sendText = ''
            } else {
                message.warn("不能发送空白消息")
            }
        }

        //快捷回复发送
        const sendQuick = (e) => {
            messagePush(e.msg_type, e.msg_info);
        }

        //文件监听上传
        const handleFiles = (files) => {
            files.map((i) => {
                let file = i;
                Modal.confirm({
                    title: '是否发送该文件？',
                    icon: createVNode(ExclamationCircleOutlined),
                    content: file.name,
                    centered:true,
                    onOk() {
                        switch (file.type) {
                            case 'video/mp4':
                                uploadHland({ file: files[0] })
                                break;
                        
                            case 'image/png' || 'image/jpg' || 'image/jpeg':
                                uploadImgHland({ file: files[0] });
                                break;
                        }
                    }
                });
            });
        }

        //发送二维码
        const sendCode = () => {
            messagePush('code', state.info.web);
        }

        //向用户发起消息提醒
        const sendRemind = async () => {
            let time = await getTime(new Date());
            messagePush('remind', time + " 发送一次消息提醒");
        }

        //消息推送
        const messagePush = async (type, content) => {
            const res = await axios.post("/message/send_to_user", { user_id: props.user_id, type, content },false);
        }

        return {
            ...toRefs(state),
            emoji, //表情列表
            emojisHland, //表情选中
            sendMessage, //发送文本信息
            beforeUpload, //上传图片前
            videoBeforeUpload, //上传视频前
            uploadImgHland, //上传图片
            uploadHland, //上传文件
            withdraw, //撤回所有聊天记录
            sendQuick, //快捷回复发送
            handleFiles, //文件监听上传
            sendCode, //发送二维码
            sendRemind, //向用户发起消息提醒
        };
    },
};
</script>

<style lang="less" scoped>
.editBox-content {
    textarea.ant-input {
        resize: none;
        padding: 0;
        overflow: overlay;
        border-radius: 0;
    }
    .send-btn {
        background: #eef0f9;
        border-color: #eef0f9;
        color: @primary-color;
        padding: 0 30px;

        &:hover,
        &:focus {
            background: #d2d6e6;
            border-color: #d2d6e6;
        }
    }    
}
.emojis-box {
    .item {
        padding: 3.7px;

        &:hover {
            background: #eee;
        }
    }
}
</style>