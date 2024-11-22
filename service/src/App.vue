<template>
  <a-config-provider :locale="locale">
    <router-view />
  </a-config-provider>
</template>

<script>
import { defineComponent, reactive, toRefs, computed,onMounted } from "vue";
import "dayjs/locale/zh-cn";
import zhCN from 'ant-design-vue/lib/locale-provider/zh_CN';
import { useStore } from "vuex";

export default defineComponent({
  name: "App",
  setup() {
    const store = useStore();
    const state = reactive({
      locale: zhCN,
      messageReminder: computed(() => store.state.messageReminder), //声音播放
    });

    onMounted(()=>{
      if (window.Notification) {  //浏览器通知
        if (Notification.permission == "default") {  //从未授权
          Notification.requestPermission()
        }
      }
    })

    return {
      ...toRefs(state),
    };
  },
});
</script>
<style lang="less">
/** drawer抽屉列表hover动画统一样式 **/
.custom-class {
  .ant-drawer-body {
    padding: 0;

    .drawer-list {
      height: 100%;
      padding: 20px;

      .item {
        background: white;
        border-radius: 20px;
        padding: 20px;
        margin-bottom: 20px;
        transition: all .3s;
        cursor: default;

        &:hover {
          box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);

          .content {
            overflow: auto;
            text-overflow: initial;
            white-space: normal;
          }
        }
      }

      .draggable {
        cursor: copy !important;
      }
    }
  }
}
/********** 气泡弹框 **********/
.popover-pointer {
  .item {
    cursor: pointer;
    padding: 10px 0;
    border-bottom: 1px solid #eee;

    &:hover,
    span:hover {
      color: @primary-color;
    }

    &:last-child {
      border: none;
      padding-bottom: 0;
    }
  }

  .item-auto {
    cursor: auto;

    &:hover,
    span:hover {
      color: #333;
    }
  }
}
/********** 弹窗统一样式 **********/
.ant-modal-content {
  border-radius: 20px;
  overflow: hidden;

  .ant-modal-header {
    border-bottom: none;

    .ant-modal-title {
      text-align: center;
      font-weight: bold;
    }
  }

  .ant-modal-footer {
    text-align: center;
    border: none;

    .ant-btn {
      border-radius: 12px;
    }
  }
}
.ant-modal-confirm .ant-modal-confirm-btns .ant-btn {
  border-radius: 6px;
}
.left-modal {
  .ant-modal-content {
    // left: -22% !important;
  }
}
/**** 输入框统一样式 *****/
.ant-input,.ant-input-affix-wrapper {
  border-radius: 8px;
  overflow: hidden;
}
.ant-radio-button-wrapper:first-child{
  border-radius: 8px 0 0 8px;
}
.ant-radio-button-wrapper:last-child{
  border-radius: 0 8px 8px 0;
}
/****** 表格统一样式 ******/
.tab-box {
  .ant-table-wrapper {
    border: 1px solid #EEEEEE;

    tr {
      &:last-child {
        td {
          border-bottom: none;
        }
      }
    }
  }
}
/******* 右键菜单 ******/
.contextmenu {
  margin: 0;
  background: #fff;
  z-index: 3000;
  position: fixed;
  border-radius: 4px;
  overflow: hidden;
  font-size: 12px;
  box-shadow: 2px 2px 3px 0 rgba(0, 0, 0, 0.3);
  border: 1px solid #d7d7d7;

  li {
    margin: 0;
    padding: 7px 16px;
    cursor: pointer;

    &:hover {
      background: #d7d7d7;
    }
  }
}
/******* 组件按钮 ******/
.btn-add {
  background: #31DAC3;
  border-color: #31DAC3;

  &:hover,
  &:focus {
    background: #70e2d3;
    border-color: #70e2d3;
  }
}

.btn-red {
  background: #FF6E6E;
  border-color: #FF6E6E;

  &:hover,
  &:focus {
    background: #fc8989;
    border-color: #fc8989;
  }
}

.btn-green {
  background: #039B84;
  border-color: #039B84;

  &:hover,
  &:focus {
    background: #039B84;
    border-color: #039B84;
  }
}
</style>