<style lang="less" scoped>
.message-item {
  display: flex;
  margin: 45px 12px;

  .markdown-content{
    white-space: normal !important;
  }

  .message-item-body {
    overflow: hidden;
    flex: 1;
    padding-left: 8px;
  }

  .avatar {
    display: block;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background-color: #fff;
    box-shadow: 0 2px 4px 0 rgba(0, 0, 0, 0.08);
  }

  .msg-img {
    width: 100%;
    height: 100%;
    max-width: 260px;
    max-height: 260px;
  }

  .message-content {
    margin-right: 40px;
    &:hover {
      .hover-copy-tool-block {
        opacity: 1;
      }
    }
  }

  .text-message {
    position: relative;
    display: inline-block;
    padding: 12px;
    font-size: 14px;
    min-height: 44px;
    line-height: 20px;
    text-align: left;
    font-weight: 400;
    white-space: pre-wrap;
    word-break: break-all;
    max-width: 100%;
  }

  .multiple-message-item{
    display: inline-block;
    padding: 12px;
    border-radius: 16px 4px 16px 16px;
    color: #f5f9ff;
    background-color: #2475fc;
  }

  &.robot-message-item {
    .text-message {
      color: #1a1a1a;
      background-color: #edeff2;
      border-radius: 4px 16px 16px 16px;
    }
    .hover-copy-tool-block {
      bottom: -26px;
      left: 8px;
    }
    .reply-item {
      background-color: #edeff2;
    }
  }

  &.user-message-item {
    flex-direction: row-reverse;

    .message-item-body {
      text-align: right;
      padding-left: 0;
      padding-right: 8px;
    }

    .message-content {
      margin-left: 40px;
      margin-right: 0;
    }

    .text-message {
      border-radius: 16px 4px 16px 16px;
      color: #f5f9ff;
      background-color: #2475fc;
    }
    .hover-copy-tool-block {
      left: -30px;
      width: 36px;
      bottom: 0px;
      height: 46px;
      display: flex;
    }
  }

  &.welcome-message-item {
    .text-message {
      width: 100%;
      border-radius: 4px 16px 16px 16px;
    }
  }

  .question-list {
    border: 1px solid #edeff2;
    border-radius: 0 0 16px 16px;
    background-color: #fff;
    .question-item {
      line-height: 20px;
      padding: 6px 12px;
      font-size: 14px;
      font-weight: 400;
      border-bottom: 1px solid #edeff2;
      color: #2475fc;
      cursor: pointer;

      &:last-child {
        border-bottom: 0;
      }
    }

    .question-item:hover {
      color: #4d94ff;
    }
  }
  .copy-block {
    display: flex;
    align-items: center;
    color: #7a8699;
    font-size: 14px;
    cursor: pointer;
    width: fit-content;
    transition: all 0.5s ease;
    padding: 0 8px;
    height: 24px;
    border-radius: 6px;
    span {
      display: flex;
      height: 100%;
      line-height: 25px;
      margin-left: 2px;
    }
    .copy-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/copy.png);
      background-size: 16px;
      width: 16px;
      height: 16px;
    }
    &:hover {
      background: #f2f4f7;
      text-align: center;
      color: #3a4559;
      .copy-icon {
        background-image: url(@/assets/img/copy-hover.png);
      }
    }
  }

  .operation {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 12px;
    .operation-right {
      display: flex;
      gap: 5px;
      align-items: center;
      justify-content: space-between;
    }
  }

  .praise-block {
    text-align: center;
    width: 26px;
    cursor: pointer;
    border-radius: 6px;

    .praise-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/praise.png);
      background-size: 26px;
      width: 26px;
      height: 26px;
      margin-right: 2px;
    }
    &:hover {
      background: #f2f4f7;
    }
  }

  .praise-active {
    .praise-icon {
      background-image: url(@/assets/img/praise-hover.png);
    }
  }

  .trample-block {
    text-align: center;
    width: 26px;
    cursor: pointer;
    border-radius: 6px;

    .trample-icon {
      transition: background-image 0.5s ease;
      background-image: url(@/assets/img/trample.png);
      background-size: 26px;
      width: 26px;
      height: 26px;
      margin-right: 2px;
    }
    &:hover {
      background: #f2f4f7;
    }
  }

  .trample-active {
    .trample-icon {
      background-image: url(@/assets/img/trample-hover.png);
    }
  }

  .hover-copy-tool-block {
    padding: 0;
    opacity: 0;
    position: absolute;
    max-width: 90px;
    align-items: center;
    justify-content: center;
    transition: all 0.5s ease;

    .copy-block {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 0;
      width: 26px;
    }

    .operation-right {
      display: flex;
      gap: 5px;
      align-items: center;
      justify-content: space-between;
    }
  }
}

.feedback-modal {
  padding: 16px;
  flex-shrink: 0;
  border-radius: 12px;
  background: var(--Conditional-pop-over, #fff);
  box-shadow:
    0 6px 30px 5px #0000000d,
    0 16px 24px 2px #0000000a,
    0 8px 10px -5px #00000014;

  .modal-title {
    color: #262626;
    text-align: left;
    font-family: 'PingFang SC';
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    margin-bottom: 4px;
  }

  .btn-box {
    margin-top: 12px;
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 8px;
  }

  .model-btn {
    width: 74px;
    height: 32px;
    display: flex;
    padding: 5px 16px;
    justify-content: center;
    align-items: center;
    border-radius: 6px;
  }
}

.text-input {
  width: 268px;
  height: 98px;
  display: flex;
  padding: 5px 12px;
  color: #666;
  flex-direction: column;
  align-items: flex-start;
  align-self: stretch;
  border-radius: 2px;
  border: 1px solid var(--Neutral-5, #d9d9d9);
  background: var(--Neutral-1, #fff);
}

.answer-reference-box {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 24px;
  line-height: 22px;
  font-weight: 400;
  color: #164799;
  font-size: 14px;
  padding-bottom: 12px;
  border-bottom: 1px solid #edeff2;
  margin-bottom: 12px;
  .list-item {
    cursor: pointer;
    display: flex;
    align-items: baseline;
    gap: 4px;
    color: #164799;
    position: relative;
    .svg-action {
      position: relative;
      top: 2px;
    }
  }
}

.label-flex-block {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.thinking-label-wrapper {
  display: flex;
  align-items: center;
  margin-bottom: 8px;

  .thinking-label {
    display: flex;
    align-items: center;
    height: auto;
    padding: 8px 16px;
    border-radius: 8px;
    background: #e4e6eb;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #d8dde6;
    }
    .think-icon,
    .loading {
      margin-right: 8px;
      font-size: 16px;
      color: #262626;
    }
    .label-text {
      font-size: 14px;
      font-weight: 400;
      color: #262626;
    }
    .arrow-down {
      margin-left: 8px;
      font-size: 16px;
      color: #262626;
      cursor: pointer;
    }
  }
  .tip {
    margin-left: 8px;
    font-size: 16px;
    color: #8c8c8c;
    cursor: pointer;
  }

  &.reasoning_open {
    .down-arrow {
      transform: rotate(180deg);
    }
  }
}
.thinking-content {
  position: relative;
  line-height: 22px;
  padding-bottom: 16px;
  padding-left: 16px;
  margin-bottom: 12px;
  font-size: 14px;
  font-weight: 400;
  color: #8c8c8c;
  border-bottom: 1px solid #edeff2;
  // 4px竖线
  &::before {
    display: block;
    position: absolute;
    content: '';
    left: 0;
    top: 4px;
    bottom: 20px;
    width: 4px;
    background-color: #d9d9d9;
  }
}

.reply-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 8px;
}

.reply-url .url-link {
  color: #2475fc;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 20px;
  text-decoration: none;
}

.reply-url .url-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.reply-card .card-row {
  display: inline-flex;
  align-items: center;
  gap: 12px;
}

.reply-card .card-title-box {
  display: flex;
  align-items: center;
  gap: 4px;
}

.reply-card .card-thumb {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.reply-card .card-title {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.reply-imageText .imageText-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.reply-imageText .imageText-thumb {
  width: 40px;
  height: 40px;
  border-radius: 6px;
  object-fit: cover;
}

.reply-imageText .imageText-text {
  display: flex;
  flex-direction: column;
}

.reply-imageText .imageText-title {
  color: #595959;
  font-size: 14px;
  font-style: normal;
  font-weight: 400;
  line-height: 22px;
}

.reply-imageText .imageText-desc {
  color: #8c8c8c;
  font-size: 12px;
  font-style: normal;
  font-weight: 400;
  line-height: 16px;
}

.reply-item {
  margin-right: 50px;
  border-radius: 4px 16px 16px 16px;
  padding: 16px 12px;
  background: #fff;
}

.reply-image .msg-img {
  max-width: 500px;
}

.reply-smartMenu {
  color: #1a1a1a;
  .smart-menu-box {
    .card-title {
      white-space: pre-wrap;
      align-self: stretch;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
      margin-bottom: 14px;
    }
    .card-text {
      display: flex;
      flex-direction: column;
      gap: 10px;
      .reply-line {
        line-height: 22px;
        .line-text {
          color: #3a4559;
        }
        .empty-line {
          height: 22px;
        }
        .link {
          display: flex;
          align-items: center;
          gap: 4px;
        }
      }
    }
  }
}
</style>

<template>
  <div class="message-item" :class="messageItemClasses" :id="'msg-' + msg.uid">
    <div class="message-item-left">
      <!-- {{ props.msg }} -->
      <img class="avatar" :src="props.msg.avatar" />
    </div>
    <div class="message-item-body">
      <template v-if="parseReplyList(props.msg.reply_content_list).length">
        <div class="reply-list">
          <template v-for="(rc, idx) in parseReplyList(props.msg.reply_content_list)" :key="idx">
            <div v-if="(rc.reply_type || rc.type) === 'text'" class="reply-item reply-text">
              <div class="message-content" v-html="rc.description"></div>
            </div>
            <div v-else-if="(rc.reply_type || rc.type) === 'image'" class="reply-item reply-image">
              <div class="message-content">
                <img v-viewer class="msg-img" :src="rc.pic || rc.thumb_url" />
              </div>
            </div>
            <div v-else-if="(rc.reply_type || rc.type) === 'url'" class="reply-item reply-url">
              <div class="url-row">
                <a class="url-link" :href="rc.url" target="_blank">{{ rc.url }}</a>
              </div>
            </div>
            <div v-else-if="(rc.reply_type || rc.type) === 'card'" class="reply-item reply-card">
              <div class="card-row">
                <img v-if="rc.thumb_url" :src="rc.thumb_url" class="card-thumb" />
                <div class="card-title-box">
                  <svg-icon class="think-icon" name="applet"></svg-icon>
                  <span class="card-title">{{ rc.title }}</span>
                </div>
              </div>
            </div>
            <div v-else-if="(rc.reply_type || rc.type) === 'imageText'" class="reply-item reply-imageText">
              <a class="imageText-row" :href="rc.url" target="_blank">
                <img v-if="rc.thumb_url" :src="rc.thumb_url" class="imageText-thumb" />
                <div class="imageText-text">
                  <div class="imageText-title">{{ rc.title }}</div>
                  <div class="imageText-desc">{{ rc.description }}</div>
                </div>
              </a>
            </div>
            <div v-else-if="(rc.reply_type || rc.type) === 'smartMenu'" class="reply-item reply-smartMenu">
              <div class="smart-menu-box">
                <div class="card-title" v-if="rc.smart_menu && rc.smart_menu.menu_description">{{ rc.smart_menu.menu_description }}</div>
                <div class="card-text">
                  <div v-for="(line, li) in buildMenuLines(rc.smart_menu?.menu_content || [])" :key="li">
                    <div class="reply-line" @click="onSmartReplyLineClick($event)">
                      <span class="line-text" v-if="line.kind === 'text'">{{ line.text }}</span>
                      <div v-else-if="line.kind === 'newline'" class="empty-line"></div>
                      <span v-else-if="line.kind === 'html'" v-html="line.html"></span>
                      <a v-else-if="line.kind === 'keyword'" href="javascript:;" class="link" @click.prevent="onClickSmartMenuKeyword(line.text)">
                        <div v-if="line.serial_no">{{ line.serial_no }}</div>
                        {{ line.text }}
                      </a>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </template>

      <div class="label-flex-block">
        <!-- 检索知识库 -->
        <div
          class="thinking-label-wrapper"
          :class="{ reasoning_open: props.msg.show_quote_file }"
          v-if="props.msg.is_customer == 0 && props.msg.msg_type == 1 && isShowQuoteFileProgress"
        >
          <div class="thinking-label" @click="toggleQuoteFiel()">
            <template v-if="props.msg.quote_loading">
              <van-loading class="loading" color="#262626" size="16px" type="spinner" />
              <span class="label-text">{{ translate('正在检索知识库...') }}</span>
            </template>
            <template v-else>
              <svg-icon class="think-icon" name="quote-file"></svg-icon>
              <span class="label-text" v-if="externalConfigPC.lang == 'en-US'"
                >Found {{ props.msg.quote_file.length }}
                {{ props.msg.quote_file.length > 1 ? 'documents' : 'document' }} in knowledge
                base</span
              >
              <span class="label-text" v-else
                >检索到{{ props.msg.quote_file.length }}个知识库文档</span
              >
            </template>
            <svg-icon
              name="down-arrow"
              class="down-arrow"
              v-if="props.msg.quote_file.length"
            ></svg-icon>
          </div>
        </div>
        <!-- 思考过程label -->
        <div
          class="thinking-label-wrapper"
          :class="{ reasoning_open: props.msg.show_reasoning }"
          v-if="props.msg.reasoning_content && props.msg.reasoning_content.length > 0"
        >
          <div class="thinking-label" @click="toggleReasonProcess()">
            <van-loading
              class="loading"
              color="#262626"
              size="16px"
              type="spinner"
              v-if="props.msg.reasoning_status"
            />
            <svg-icon class="think-icon" name="think" v-else></svg-icon>
            <span class="label-text">{{
              props.msg.reasoning_status ? translate('深度思考中...') : translate('已完成深度思考')
            }}</span>

            <svg-icon
              name="down-arrow"
              class="down-arrow"
              v-if="!props.msg.reasoning_status"
            ></svg-icon>
          </div>
        </div>
      </div>

      <div class="message-content" :data-msg-type="props.msg.msg_type" v-if="!(props.msg.msg_type == 1 && props.msg.content == '')">
        <div class="thinking-content" v-if="props.msg.show_reasoning">
          <cherry-markdown :content="props.msg.reasoning_content" />
        </div>
        <template v-if="props.msg.msg_type == 1">
          <div class="text-message" v-if="props.msg.content !== ''" v-viewer>
            <div
              class="answer-reference-box"
              v-if="
                isShowQuoteFileProgress &&
                props.msg.show_quote_file &&
                props.msg.is_customer != 1 &&
                props.msg.quote_file &&
                props.msg.quote_file.length
              "
            >
              <div class="list-item" v-for="(item, index) in props.msg.quote_file" :key="index">
                <svg-icon name="quote-file" />
                <span @click="handleToLink(item)">
                  <span v-if="item.file_name">{{ item.file_name }}</span>
                  <span v-else>{{ item.library_name }}-精选</span>
                </span>
              </div>
            </div>
            <template v-if="props.msg.is_customer == 1">
              <template v-if="props.msg.received_message_type == 'image' && props.msg.media_id_to_oss_url">
                <img v-viewer class="msg-img" :src="props.msg.media_id_to_oss_url" />
              </template>
              <div v-else v-html="props.msg.content"></div>
            </template>
            <cherry-markdown :content="props.msg.content" v-else />
            <div class="operation" v-if="isShowCopy">
              <div @click="handleCopy" class="copy-block">
                <div class="copy-icon"></div>
                <span>{{ translate('复制') }}</span>
              </div>
              <div
                ref="operationRef"
                class="operation-right"
                v-if="!isCustomerMessage && !isWelcomeMessage && robot.feedback_switch"
              >
                <div
                  v-tooltip="translate('点赞')"
                  @click="handlePraise(props.msg)"
                  class="praise-block"
                  :class="{ 'praise-active': !isTrampleClick && isPraiseActive }"
                >
                  <div class="praise-icon"></div>
                </div>
                <van-popover placement="top-end" v-if="!isTrampleActive" @close="onCancel">
                  <div ref="feedbackModal" class="feedback-modal">
                    <div class="modal-title">{{ translate('结果反馈') }}</div>
                    <div class="content">
                      <textarea
                        class="text-input"
                        v-model="feedbackContent"
                        :placeholder="translate('请反馈你觉得回答不满意的地方')"
                      ></textarea>
                    </div>
                    <div class="btn-box">
                      <van-button class="model-btn cancel-btn" type="default" @click="onCancel"
                        >{{ translate('取消') }}</van-button
                      >
                      <van-button class="model-btn submit-btn" type="primary" @click="onSubmit"
                        >{{ translate('提交') }}</van-button
                      >
                    </div>
                  </div>
                  <template #reference>
                    <div
                      v-tooltip="translate('点踩')"
                      @click="handleTrample(props.msg, 4)"
                      class="trample-block"
                      :class="{ 'trample-active': isTrampleClick }"
                    >
                      <div class="trample-icon"></div>
                    </div>
                  </template>
                </van-popover>
                <div
                  v-tooltip="translate('点踩')"
                  v-else
                  @click="handleTrample(props.msg, 3)"
                  class="trample-block"
                  :class="{ 'trample-active': isTrampleActive }"
                >
                  <div class="trample-icon"></div>
                </div>
              </div>
            </div>
            <div class="hover-copy-tool-block" v-if="isShowHoverCopy">
              <div ref="operationRef" class="operation-right">
                <div v-tooltip="translate('复制')" @click="handleCopy" class="copy-block">
                  <div class="copy-icon"></div>
                </div>
                <template v-if="!isCustomerMessage && !isWelcomeMessage && robot.feedback_switch">
                  <div
                    v-tooltip="translate('点赞')"
                    @click="handlePraise(props.msg)"
                    class="praise-block"
                    :class="{ 'praise-active': !isTrampleClick && isPraiseActive }"
                  >
                    <div class="praise-icon"></div>
                  </div>
                  <van-popover placement="top-end" v-if="!isTrampleActive" @close="onCancel">
                    <div ref="feedbackModal" class="feedback-modal">
                      <div class="modal-title">{{ translate('结果反馈') }}</div>
                      <div class="content">
                        <textarea
                          class="text-input"
                          v-model="feedbackContent"
                          :placeholder="translate('请反馈你觉得回答不满意的地方')"
                        ></textarea>
                      </div>
                      <div class="btn-box">
                        <van-button class="model-btn cancel-btn" type="default" @click="onCancel"
                          >{{ translate('取消') }}</van-button
                        >
                        <van-button class="model-btn submit-btn" type="primary" @click="onSubmit"
                          >{{ translate('提交') }}</van-button
                        >
                      </div>
                    </div>
                    <template #reference>
                      <div
                        v-tooltip="translate('点踩')"
                        @click="handleTrample(props.msg, 2)"
                        class="trample-block"
                        :class="{ 'trample-active': isTrampleClick }"
                      >
                        <div class="trample-icon"></div>
                      </div>
                    </template>
                  </van-popover>
                  <div
                    v-else
                    v-tooltip="translate('点踩')"
                    @click="handleTrample(props.msg, 1)"
                    class="trample-block"
                    :class="{ 'trample-active': isTrampleActive }"
                  >
                    <div class="trample-icon"></div>
                  </div>
                </template>
              </div>
            </div>
          </div>
          <div v-else class="text-message">{{ textMessage }}</div>
          <div
            class="question-list"
            v-if="props.msg.menu_json && props.msg.menu_json.question.length"
          >
            <div
              class="question-item"
              @click="sendTextMessage(item)"
              v-for="(item, index) in props.msg.menu_json.question"
              :key="index"
            >
              {{ item }}
            </div>
          </div>
        </template>

        <template v-else-if="props.msg.msg_type == 2">
          <cherry-markdown class="text-message markdown-content" :content="props.msg.menu_json.content" v-if="props.msg.isWelcome" />
          <div class="text-message" v-html="escapeHTML(props.msg.menu_json.content)" v-else></div>
          <div
            class="question-list"
            v-if="props.msg.menu_json && props.msg.menu_json.question.length"
          >
            <div
              class="question-item"
              v-for="item in props.msg.menu_json.question"
              :key="item"
              @click="sendTextMessage(item)"
            >
              <span>{{ item }}</span>
            </div>
          </div>
        </template>

        <template v-else-if="props.msg.msg_type == 3">
          <img v-viewer class="msg-img" :src="props.msg.content" />
        </template>

        <template v-else-if="props.msg.msg_type == 99">
          <div class="multiple-message-item">
            <MultipleMessage :message="props.msg.content" />
          </div>
        </template>
      </div>
        <!-- 语音消息 -->
      <div :style="{'margin-top': props.msg.content == '' ? '-24px' : ''}" v-if="props.msg.voice_content && props.msg.voice_content.length">
        <VoiceMessage :voice_content="props.msg.voice_content" />
      </div>
    </div>
    <quote-modal ref="quoteModalRef"></quote-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import CherryMarkdown from '@/components/cherry-markdown/index.vue'
import type { Message } from '@/stores/modules/chat'
import { useChatStore } from '@/stores/modules/chat'
import { escapeHTML } from '@/utils/index'
import { showToast } from 'vant'
import useClipboard from 'vue-clipboard3'
import QuoteModal from '../quote-modal/index.vue'
import { translate } from '@/utils/translate.js'
import VoiceMessage from './voice-message.vue'
import MultipleMessage from './multiple-message.vue'


interface praiseParams {
  ai_message_id: string
  type: number
  customer_message_id?: string
}

const { toClipboard } = useClipboard()

const emit = defineEmits(['sendTextMessage'])
const chatStore = useChatStore()
const { robot, onAddFeedback, onDelFeedback } = chatStore
const externalConfigPC = computed(() => chatStore.externalConfigPC)

const textMessage = ref('.')
const feedbackContent = ref('')
const ai_message_id = ref('')
const customer_message_id = ref<any>('')
const operationRef = ref<HTMLElement | null>(null)
let interval: number

const props = defineProps({
  msg: {
    type: Object as () => Message,
    required: true
  },
  index: {
    type: [Number, String]
  },
  messageLength: {
    type: Number,
    default: 0
  },
  prevMsg: {
    type: Object as () => Message
  }
})

const isTrampleClick = ref(false)
const isPraiseActive = ref(props.msg.feedback_type == '1' ? true : false)
const isTrampleActive = ref(props.msg.feedback_type == '2' ? true : false)

const isShowCopy = computed(() => {
  // 最后一条消息 机器人的消息 消息类型为1 不是正在发送
  return (
    props.index === props.messageLength - 1 &&
    props.msg.msg_type == 1 &&
    !robot.is_sending &&
    !isCustomerMessage.value
  )
})

const isShowHoverCopy = computed(() => {
  return !isShowCopy.value && props.index !== props.messageLength - 1
})

const handleCopy = async () => {
  isTrampleClick.value = false
  await toClipboard(props.msg.content)
  showToast('复制成功')
}

const handlePraise = async (item) => {
  isTrampleClick.value = false
  if (props.prevMsg && props.prevMsg.is_customer == 1) {
    customer_message_id.value = props.prevMsg.id
  }
  ai_message_id.value = item.id
  const params: praiseParams = {
    ai_message_id: ai_message_id.value,
    type: 1
  }
  if (!isPraiseActive.value) {
    params.customer_message_id = customer_message_id.value
    await onAddFeedback(params)
    isTrampleClick.value = false
    isTrampleActive.value = false
    isPraiseActive.value = true
  } else {
    await onDelFeedback(params)
    isPraiseActive.value = false
  }
  // showToast('操作成功')
}

const handleTrample = async (item, type) => {
  isTrampleClick.value = true
  if (props.prevMsg && props.prevMsg.is_customer == 1) {
    customer_message_id.value = props.prevMsg.id
  }
  ai_message_id.value = item.id
  if (isTrampleActive.value) {
    // 取消踩
    const params = {
      ai_message_id: ai_message_id.value
    }
    await onDelFeedback(params)
    isTrampleActive.value = false
    isTrampleClick.value = false
    showToast('操作成功')
  }
}

const onCancel = async () => {
  const params = {
    ai_message_id: ai_message_id.value,
    customer_message_id: customer_message_id.value,
    type: 2
  }
  await onAddFeedback(params)
  isTrampleActive.value = true
  // operationRef.value!.click()
  feedbackContent.value = ''
  isPraiseActive.value = false
  // showToast('操作成功')
}

const onSubmit = async () => {
  const params = {
    ai_message_id: ai_message_id.value,
    customer_message_id: customer_message_id.value,
    type: 2,
    content: feedbackContent.value
  }
  await onAddFeedback(params)
  isTrampleActive.value = true
  // operationRef.value!.click()
  feedbackContent.value = ''
  isPraiseActive.value = false
  showToast('感谢反馈')
}

// 检查是否为用户消息
const isCustomerMessage = computed(() => props.msg.is_customer == 1)

// 是否显示引用
const isShowQuoteFileProgress = computed(() => {
  return (robot.chat_type == 1 || robot.chat_type == 3) && robot.answer_source_switch && robot.application_type == '0'
})

// 是否为欢迎语
const isWelcomeMessage = computed(() => props.msg.msg_type == 2)

// 计算消息项的类
const messageItemClasses = computed(() => ({
  'user-message-item': isCustomerMessage.value === true,
  'robot-message-item': isCustomerMessage.value === false,
  'welcome-message-item': props.msg.menu_json && props.msg.menu_json.question
}))

// 等待机器人回复增加动态...
const startLoadingAnimation = () => {
  const dots = ['.', '..', '...']
  let dotIndex = 0
  interval = window.setInterval(() => {
    dotIndex = (dotIndex + 1) % dots.length
    textMessage.value = dots[dotIndex]
  }, 500)
}

const sendTextMessage = (text: string) => {
  emit('sendTextMessage', text)
}
const quoteModalRef = ref<any>(null)
const handleToLink = (item: any) => {
  quoteModalRef.value &&
    quoteModalRef.value.showPopup({
      message_id: item.message_id || props.msg.id,
      file_id: item.id,
      robot_key: robot.robot_key,
      openid: robot.openid,
      file_name: item.file_name || item.library_name + '-精选',
      answer_source_data: item.answer_source_data ? JSON.parse(item.answer_source_data) : null
    })
}

const toggleReasonProcess = () => {
  props.msg.show_reasoning = !props.msg.show_reasoning
}
const toggleQuoteFiel = () => {
  props.msg.show_quote_file = !props.msg.show_quote_file
}

function parseReplyList(val) {
  try {
    if (!val) return []
    if (Array.isArray(val)) return val
    if (typeof val === 'string') return JSON.parse(val || '[]')
    return []
  } catch (_e) {
    return []
  }
}

type SmartMenuLine =
  | { kind: 'text'; text: string }
  | { kind: 'newline' }
  | { kind: 'html'; html: string }
  | { kind: 'keyword'; text: string; serial_no?: string }

function buildMenuLines(menu_content: Array<Record<string, any>> | undefined | null): SmartMenuLine[] {
  const out: SmartMenuLine[] = []
  ;(Array.isArray(menu_content) ? menu_content : []).forEach((mc) => {
    const t = String(mc?.menu_type || '')
    const txt = String(mc?.content || '')
    if (t === '0') {
      if (txt === '') { out.push({ kind: 'newline' }) }
      else if (/<a[\s\S]*?<\/a>/.test(txt)) {
        const sanitized = /href\s*=\s*['"]\s*#\s*['"]/i.test(txt)
          ? txt.replace(/href\s*=\s*['"]\s*#\s*['"]/ig, 'href="javascript:;"')
          : txt.replace(/href=/ig, 'target="_blank" href=')
        out.push({ kind: 'html', html: sanitized })
      } else { out.push({ kind: 'text', text: txt }) }
    } else if (t === '1') { out.push({ kind: 'keyword', text: txt, serial_no: mc?.serial_no || '' }) }
  })
  return out.slice(0, 20)
}

function onSmartReplyLineClick(e: any) {
  const a = (e?.target && (e.target.closest && e.target.closest('a'))) || null
  if (!a) return
  const href = String(a.getAttribute('href') || '')
  if (href === '#' || href === 'javascript:;') { e.preventDefault() }
}

function onClickSmartMenuKeyword(text: string) {
  const t = String(text || '').trim()
  if (!t) return
  emit('sendTextMessage', t)
}
onMounted(() => {
  startLoadingAnimation()
})

onUnmounted(() => {
  clearInterval(interval)
})
</script>
