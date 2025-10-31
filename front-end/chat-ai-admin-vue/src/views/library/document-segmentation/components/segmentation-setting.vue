<style lang="less" scoped>
.segmentation-setting {
  .setting-item {
    margin-bottom: 16px;
    transition: all 0.2s;
    border: 1px solid #fff;

    .setting-item-header {
      display: flex;
      align-items: center;
      height: 72px;
      padding-left: 16px;
      border-radius: 2px;
      background-color: #f5f5f5;
      transition: all 0.2s;
      cursor: pointer;

      .setting-item-icon {
        width: 40px;
        height: 40px;
        background-repeat: no-repeat;
        background-size: cover;
      }

      .setting-item-info {
        padding-left: 8px;
      }

      .setting-item-name {
        line-height: 22px;
        font-size: 14px;
        font-weight: 600;
        color: #242933;
      }

      .setting-item-desc {
        line-height: 20px;
        padding-top: 2px;
        font-size: 12px;
        font-weight: 400;
        color: #3a4559;
      }
    }

    &:hover {
      border-radius: 4px;
      overflow: hidden;

      .setting-item-header {
        background: #e5efff;
      }
    }

    &.active {
      border-radius: 4px;
      overflow: hidden;

      .setting-item-header {
        background: #e5efff;
      }
    }

    .intelligent-icon {
      background-image: url('@/assets/img/library/document-segmentation/intelligent.svg');
    }

    .intelligent-icon.active {
      background-image: url('@/assets/img/library/document-segmentation/intelligent_active.svg');
    }

    .custom-icon {
      background-image: url('@/assets/img/library/document-segmentation/custom.svg');
    }

    .custom-icon.active {
      background-image: url('@/assets/img/library/document-segmentation/custom_active.svg');
    }
  }

  .sub-setting-item-name {
    height: 30px;
    line-height: 30px;
    padding-left: 8px;
    margin-bottom: 16px;
    font-size: 14px;
    font-weight: 600;
    color: #242933;
    background-color: #f2f4f7;
  }

  .document-types {
    display: flex;
    margin-bottom: 16px;

    .document-type {
      display: inline-flex;
      justify-content: center;
      align-items: center;
      border-radius: 2px;
      padding: 5px 16px;
      margin-right: 8px;
      font-size: 14px;
      font-weight: 400;
      color: #595959;
      background: #fff;
      border: 1px solid #d9d9d9;
      cursor: pointer;
      transition: all 0.2s;

      &:last-child {
        margin-right: 0;
      }

      &.active {
        color: #2475fc;
        background: #f5f9ff;
        border: 1px solid #2475fc;
      }
    }
  }

  .custom-setting-form {
    .form-item {
      margin-bottom: 16px;
      flex: 1;
    }
    .flex {
      display: flex;
      gap: 16px;

      .form-item-box {
        flex: 1;
      }
    }

    .form-item-label-box {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .ai-generate {
        cursor: pointer;
        display: inline-flex;
        align-items: center;
        padding: 1px 8px;
        justify-content: center;
        align-items: center;
        gap: 4px;
        border-radius: 6px;
        margin-bottom: 3px;

        &:hover {
          background: #E4E6EB;
        }

        .ai-generate-text {
          color: #6524fc;
          font-size: 14px;
          font-style: normal;
          font-weight: 400;
          line-height: 22px;
        }
      }
    }

    .form-item-label {
      line-height: 22px;
      margin-bottom: 4px;
      font-size: 14px;
      font-weight: 400;
      color: #262626;
      &.not-required {
        &::before {
          content: '';
        }
      }
      &::before {
        content: '*';
        font-size: 14px;
        line-height: 22px;
        padding-right: 2px;
        color: #fb363f;
      }
    }

    .unit-text {
      font-size: 14px;
      color: #595959;
    }
    .form-tip {
      color: #8c8c8c;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      margin-top: 2px;
    }

    .form-item-tip {
      color: #8c8c8c;
      font-size: 12px;
      font-weight: 400;
      line-height: 14px;
      margin: 2px 0 6px;
    }

    .form-item-title {
      color: #262626;
      font-size: 14px;
      font-style: normal;
      font-weight: 600;
      line-height: 22px;
      margin-bottom: 8px;
    }
  }
  .excel-qa-form {
    .form-item {
      flex-wrap: wrap;
      display: flex;
      align-items: center;
      gap: 8px 16px;
      .form-item-box {
        width: calc(50% - 8px);
      }
      .form-item-label {
        padding-top: 5px;
      }
    }
  }

  .indexing-methods-box {
    .list-item {
      margin-top: 8px;
      padding: 16px;
      border: 1px solid #d9d9d9;
      position: relative;
      cursor: pointer;
      &.active {
        border: 2px solid #2475fc;
        .list-title-block {
          color: #2475fc;
          .svg-action {
            font-size: 16px;
            color: #2475fc;
          }
        }
        .check-icon {
          display: block;
        }
      }
      .check-icon {
        position: absolute;
        right: 0;
        bottom: 0;
        font-size: 18px;
        color: #fff;
        display: none;
      }
      .list-title-block {
        display: flex;
        align-items: center;
        font-size: 14px;
        font-weight: 600;
        line-height: 22px;
        color: #262626;
        .svg-action {
          font-size: 16px;
          margin-right: 4px;
          color: #262626;
        }
      }
      .list-content {
        margin-top: 4px;
        color: #595959;
        font-size: 14px;
        line-height: 22px;
      }
    }
  }
  .btn-box-block {
    margin-top: 16px;
    display: flex;
    align-items: center;
    gap: 12px;
  }
  .select-card-box {
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 8px;
    width: 100%;
    .select-card-item {
      width: 100%;
      position: relative;
      padding: 16px;
      border-radius: 6px;
      border: 1px solid #d9d9d9;
      cursor: pointer;
      .check-arrow {
        position: absolute;
        display: block;
        right: -1px;
        bottom: -1px;
        width: 24px;
        height: 24px;
        font-size: 24px;
        color: #fff;
        opacity: 0;
        transition: all 0.2s cubic-bezier(0.8, 0, 1, 1);
      }
      .card-title {
        display: flex;
        align-items: center;
        line-height: 22px;
        margin-bottom: 4px;
        color: #262626;
        font-weight: 600;
        font-size: 14px;
      }
      .title-icon {
        margin-right: 4px;
        font-size: 16px;
      }
      .card-desc {
        min-height: 44px;
        line-height: 22px;
        font-size: 14px;
        color: #595959;
      }

      &.active {
        background: var(--01-, #e5efff);
        border: 2px solid #2475fc;
        .check-arrow {
          opacity: 1;
        }
        .card-title {
          color: #2475fc;
        }
      }
    }
  }
}

.form-item-inptu-numbner-box {
  display: flex;
  align-items: center;
}

.form-item-inptu-numbner {
  flex: 1;
  margin-right: 5px;
}
</style>

<template>
  <div class="segmentation-setting">
    <div class="setting-items">
      <div class="setting-item intelligent-setting active">
        <div class="setting-item-body">
          <!-- <div class="sub-setting-item-name">文档类型</div> -->

          <div class="document-types" v-if="false">
            <div
              class="document-type"
              :class="{ active: formState.is_qa_doc == 0 }"
              @click="changeDocumentType(0)"
            >
              普通文档
            </div>
            <div
              class="document-type"
              :class="{ active: formState.is_qa_doc == 1 }"
              @click="changeDocumentType(1)"
            >
              QA文档
            </div>
          </div>
          <!-- <a-space class="custom-setting-form" v-if="isHtmlOrDocx">
            <div class="form-item">
              <div class="form-item-label">提取图片</div>
              <div class="form-item-body">
                <a-switch
                  @change="onChange"
                  v-model:checked="formState.enable_extract_image"
                ></a-switch>
                <div class="form-tip">
                  开启提取图片功能，将自动从文档中提取图片，作为图片上方文本分段的附件。仅支持.docx文件和网页
                </div>
              </div>
            </div>
          </a-space> -->

          <template v-if="formState.is_qa_doc == 1">
            <div class="sub-setting-item-name">文件切分</div>
            <!-- 表格类型的QA文档 -->
            <div class="custom-setting-form excel-qa-form" v-if="props.mode == 1">
              <div class="form-item">
                <div class="form-item-box">
                  <div class="form-item-label">问题所在列：</div>
                  <div class="form-item-body">
                    <a-select
                      v-model:value="formState.question_column"
                      @change="onChagneFormInput"
                      placeholder="请选择列名"
                      style="width: 100%"
                    >
                      <a-select-option
                        v-for="item in props.excellQaLists"
                        :value="item.value"
                        :key="item.value"
                        >{{ item.lable }}</a-select-option
                      >
                    </a-select>
                  </div>
                </div>
                <div class="form-item-box">
                  <div class="form-item-label not-required">相似问法所在列：</div>
                  <div class="form-item-body">
                    <a-select
                      allowClear
                      v-model:value="formState.similar_column"
                      @change="onChagneFormInput"
                      placeholder="请选择相似问法列名"
                      style="width: 100%"
                    >
                      <a-select-option
                        v-for="item in props.excellQaLists"
                        :value="item.value"
                        :key="item.value"
                        >{{ item.lable }}</a-select-option
                      >
                    </a-select>
                  </div>
                </div>
                <div class="form-item-box">
                  <div class="form-item-label">答案所在列：</div>
                  <div class="form-item-body">
                    <a-select
                      v-model:value="formState.answer_column"
                      @change="onChagneFormInput"
                      placeholder="请选择答案所在列"
                      style="width: 100%"
                    >
                      <a-select-option
                        v-for="item in props.excellQaLists"
                        :value="item.value"
                        :key="item.value"
                        >{{ item.lable }}</a-select-option
                      >
                    </a-select>
                  </div>
                </div>
              </div>
              <div class="sub-setting-item-name">索引方式</div>
              <div class="indexing-methods-box">
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 1 }"
                  @click="handleChangeQaIndexType(1)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    问题与答案一起生成索引
                  </div>
                  <div class="list-content">
                    回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复
                  </div>
                </div>
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 2 }"
                  @click="handleChangeQaIndexType(2)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="comment-search"></svg-icon>
                    仅对问题生成索引
                  </div>
                  <div class="list-content">
                    回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复
                  </div>
                </div>
                <!-- 新增的操作按钮 -->
                <div class="btn-box-block">
                  <!-- <a-button
                    @click="reChange"
                    :loading="reLoading"
                    style="flex: 1"
                    type="primary"
                    ghost
                    >生成分段预览</a-button
                  > -->
                </div>
                <div class="btn-box-block">
                  <a-button type="primary" block @click="onSave">保存</a-button>
                </div>
              </div>
            </div>
            <div class="custom-setting-form excel-qa-form" v-else>
              <div class="form-item">
                <div class="form-item-box">
                  <div class="form-item-label">问题开始标识符：</div>
                  <div class="form-item-body">
                    <a-input
                      placeholder="请输入标识符"
                      v-model:value="formState.question_lable"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>
                <div class="form-item-box">
                  <div class="form-item-label not-required">相似问法开始标识符：</div>
                  <div class="form-item-body">
                    <a-input
                      placeholder="请输入标识符"
                      v-model:value="formState.similar_label"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>

                <div class="form-item-box">
                  <div class="form-item-label">答案开始标识符：</div>
                  <div class="form-item-body">
                    <a-input
                      placeholder="请输入标识符"
                      v-model:value="formState.answer_lable"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>
              </div>

              <div class="sub-setting-item-name">索引方式</div>
              <div class="indexing-methods-box">
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 1 }"
                  @click="handleChangeQaIndexType(1)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    问题与答案一起生成索引
                  </div>
                  <div class="list-content">
                    回答用户提问时，将用户提问与导入的问题和答案一起对比相似度，根据相似度高的问题和答案回复
                  </div>
                </div>
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 2 }"
                  @click="handleChangeQaIndexType(2)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="comment-search"></svg-icon>
                    仅对问题生成索引
                  </div>
                  <div class="list-content">
                    回答用户提问时，将用户提问与导入的问题一起对比相似度，再根据相似度高的问题和对应的答案来回复
                  </div>
                </div>
                <!-- 新增的操作按钮 -->
                <div class="btn-box-block">
                  <!-- <a-button
                    @click="reChange"
                    :loading="reLoading"
                    style="flex: 1"
                    type="primary"
                    ghost
                    >生成分段预览</a-button
                  > -->
                </div>
                <div class="btn-box-block">
                  <a-button type="primary" block @click="onSave">保存</a-button>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="custom-setting-form subsection-form" v-if="props.mode != 1">
              <div class="form-item">
                <div class="form-item-title" v-if="status === 'paragraphsSegmented'">分段及清洗规则</div>
                <div class="form-item-label">分段方式：</div>
                <div class="form-item-body">
                  <div class="form-item-tip">
                    提示：语义分段更适合没有排版过的文章，即没有明显换行符号的文本，否则更推荐使用普通分段
                  </div>
                  <div class="select-card-box">
                    <div
                      class="select-card-item"
                      @click="formState.chunk_type = 1"
                      :class="{ active: formState.chunk_type == 1 }"
                    >
                      <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                      <div class="card-title">
                        <svg-icon name="ordinary-segmentation" class="title-icon"></svg-icon>
                        普通分段
                      </div>
                      <div class="card-desc">
                        基于文章中句号、空行，或者自定义符号进行分段，不会消耗模型token
                      </div>
                    </div>
                    <div
                      class="select-card-item"
                      @click="formState.chunk_type = 2"
                      :class="{ active: formState.chunk_type == 2 }"
                    >
                      <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                      <div class="card-title">
                        <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                        语义分段
                      </div>
                      <div class="card-desc">
                        将文章拆分成句子后，通过语句向量相似度进行分段，会消耗模型token
                      </div>
                    </div>
                    <div
                      class="select-card-item"
                      @click="formState.chunk_type = 3"
                      :class="{ active: formState.chunk_type == 3 }"
                    >
                      <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                      <div class="card-title">
                        <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                        AI分段
                      </div>
                      <div class="card-desc">
                        将文章提交给大模型，大模型基于设定的提示词进行分段，会消耗大量模型token
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <template v-if="formState.chunk_type == 1">
                <div class="form-item" style="margin-bottom: 18px">
                  <div class="form-item-label">分段标识符：</div>
                  <div class="form-item-body">
                    <a-select
                      placeholder="请选择"
                      style="width: 100%"
                      mode="multiple"
                      v-model:value="formState.separators_no"
                    >
                      <a-select-option
                        :value="item.no"
                        v-for="item in separatorsOptions"
                        :key="item.no"
                        >{{ item.name }}</a-select-option
                      >
                    </a-select>
                  </div>
                </div>
                <a-flex :gap="16">
                  <div class="form-item">
                    <div class="form-item-label">分段最大长度：</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.chunk_size"
                          placeholder="分段最大长度"
                          :min="200"
                          :max="10000"
                          :precision="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">字符</span>
                      </a-flex>
                    </div>
                  </div>

                  <div class="form-item">
                    <div class="form-item-label not-required">分段重叠长度：</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.chunk_overlap"
                          placeholder="分段重叠长度"
                          :min="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">字符</span>
                      </a-flex>
                    </div>
                  </div>
                </a-flex>
                <div class="form-item">
                  <div class="form-item-body">
                    <a-flex align="center">
                      <div>自动合并较小分段
                        <a-tooltip title="开启后，如果分段长度不足设置的最大分段长度，会尝试与下一分段合并，直至合并后的分段字符数大于分段最大长度">
                           <QuestionCircleOutlined />
                        </a-tooltip>
                       ：</div>
                      <a-switch 
                        :checkedValue="false" 
                        :unCheckedValue="true" 
                        v-model:checked="formState.not_merged_text" 
                        checked-children="开" 
                        un-checked-children="关" 
                      />
                    </a-flex>
                  </div>
                </div>
              </template>
              <div v-show="formState.chunk_type == 2">
                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">嵌入模型：</div>
                  <div class="form-item-body">
                    <ModelSelect
                      modelType="TEXT EMBEDDING"
                      v-model:modeName="formState.semantic_chunk_use_model"
                      v-model:modeId="formState.semantic_chunk_model_config_id"
                      style="width: 100%"
                      @loaded="onVectorModelLoaded"
                    />
                  </div>
                </div>
                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">
                    分段阈值
                    <a-tooltip>
                      <template #title
                        >用于控制分段拆分的标准，数值0~100,数值越大，分段越少，数值越小，分段越多。</template
                      >
                      <QuestionCircleOutlined
                        style="cursor: pointer; margin-left: 2px"
                      /> </a-tooltip
                    >：
                  </div>
                  <div class="form-item-body">
                    <a-input-number
                      style="width: 100%"
                      v-model:value="formState.semantic_chunk_threshold"
                      placeholder="请输入分段阈值"
                      :precision="0"
                      :min="0"
                      :max="100"
                    />
                  </div>
                </div>
                <a-flex :gap="16">
                  <div class="form-item">
                    <div class="form-item-label">分段最大长度：</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.semantic_chunk_size"
                          placeholder="分段最大长度"
                          :min="200"
                          :max="10000"
                          :precision="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">字符</span>
                      </a-flex>
                    </div>
                  </div>

                  <div class="form-item">
                    <div class="form-item-label not-required">分段重叠长度：</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.semantic_chunk_overlap"
                          placeholder="分段重叠长度"
                          :min="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">字符</span>
                      </a-flex>
                    </div>
                  </div>
                </a-flex>
              </div>
              <!-- AI分段 -->
              <div v-show="formState.chunk_type == 3">
                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">AI大模型：</div>
                  <div class="form-item-body">
                    <ModelSelect
                      modelType="LLM"
                      placeholder="请选择AI大模型"
                      v-model:modeName="formState.ai_chunk_model"
                      v-model:modeId="formState.ai_chunk_model_config_id"
                      :modeName="formState.ai_chunk_model"
                      :modeId="formState.ai_chunk_model_config_id"
                      style="width: 100%"
                      @loaded="onVectorModelLoaded"
                    />
                  </div>
                </div>

                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label-box">
                    <div class="form-item-label">提示词设置：</div>
                    <!-- <div class="ai-generate" @click="handleAIGenerate">
                      <svg-icon name="ai-generate" style="font-size: 14px;"></svg-icon>
                      <div class="ai-generate-text">AI自动生成</div>
                    </div> -->
                  </div>
                  <div class="form-item-body">
                    <a-textarea
                      :maxLength="500"
                      style="height: 80px"
                      v-model:value="formState.ai_chunk_prumpt"
                      :placeholder="defaultAiChunkPrumpt"
                    />
                  </div>
                </div>

                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">
                    单次最大字符数
                    <a-tooltip>
                      <template #title
                        >由于大模型支持的上下文数量有限制，如果上传的文档较大，会按照最大字符数先拆分成多个分段，再提交给大模型进行分段。</template
                      >
                      <QuestionCircleOutlined
                        style="cursor: pointer; margin-left: 2px"
                      /> </a-tooltip
                    >：
                  </div>
                  <div class="form-item-body form-item-inptu-numbner-box">
                    <a-input-number
                      class="form-item-inptu-numbner"
                      v-model:value="formState.ai_chunk_size"
                      placeholder="请输入单次最大字符数"
                      :precision="0"
                      :min="0"
                      :formatter="(value) => parseInt(value)"
                      :parser="(value) => parseInt(value)"
                    />
                    字符
                  </div>
                </div>
              </div>
              <div class="btn-box-block">
                <a-button @click="handleReset" style="flex: 1">重置</a-button>
                <a-button
                  @click="reChange"
                  :loading="reLoading"
                  style="flex: 1"
                  type="primary"
                  ghost
                  >生成分段预览</a-button
                >
              </div>
              <div class="btn-box-block">
                <a-button type="primary" v-if="!props.hideSave" block :disabled="saveLoading" @click="onSave"
                  >保存</a-button
                >
              </div>
            </div>
          </template>
        </div>
      </div>
    </div>
    <AiGenerate ref="aiGenerateRef" @handleEdit="handleEdit" />
  </div>
</template>

<script setup>
import { getSeparatorsList } from '@/api/library/index'
import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { reactive, ref, toRaw, computed, watch, nextTick } from 'vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import { Form } from 'ant-design-vue'
import { message } from 'ant-design-vue'
import AiGenerate from './ai-generate-modal.vue'

const useForm = Form.useForm
const emit = defineEmits(['change', 'validate', 'save', 'changeChunkType'])

const props = defineProps({
  mode: {
    type: Number,
    default: 0
  },
  excellQaLists: {
    type: Array,
    default: () => []
  },
  libFileInfo: {
    type: Object,
    default: () => {}
  },
  library_type: {
    default: 0
  },
  hideSave: {
    type: Boolean,
    default: false
  },
  status: {
    type: String,
    default: ''
  }
})

const isHtmlOrDocx = computed(() => {
  return (
    props.libFileInfo.doc_type == '1' &&
    (props.libFileInfo.file_ext == 'docx' || props.libFileInfo.file_ext == 'html')
  )
})

const defaultAiChunkPrumpt = '你是一位文章分段助手，根据文章内容的语义进行合理分段，确保每个分段表述一个完整的语义，每个分段字数控制在500字左右，最大不超过1000字。请严格按照文章内容进行分段，不要对文章内容进行加工，分段完成后输出分段后的内容。'

watch(props, (val) => {
  let libFileInfo = val.libFileInfo
  let separators_no = libFileInfo.separators_no ? libFileInfo.separators_no.split(',') : [11, 12]
  formState.separators_no = separators_no.map((item) => +item)
  formState.chunk_size = +libFileInfo.chunk_size || 512
  formState.not_merged_text = libFileInfo.not_merged_text == 'true'
  formState.ai_chunk_size = +libFileInfo.ai_chunk_size || 5000
  formState.chunk_overlap = +libFileInfo.chunk_overlap || 50
  // formState.is_qa_doc = +libFileInfo.is_qa_doc
  formState.question_lable = libFileInfo.question_lable
  formState.similar_label = libFileInfo.similar_label
  formState.answer_lable = libFileInfo.answer_lable
  formState.question_column = libFileInfo.question_column || void 0
  formState.answer_column = libFileInfo.answer_column || void 0
  formState.similar_column = libFileInfo.similar_column || void 0
  formState.qa_index_type = +libFileInfo.qa_index_type
  formState.enable_extract_image = libFileInfo.enable_extract_image == 'true'

  formState.chunk_type = +libFileInfo.chunk_type
  formState.semantic_chunk_size = +libFileInfo.semantic_chunk_size || 512
  formState.semantic_chunk_overlap = +libFileInfo.semantic_chunk_overlap || 50
  formState.semantic_chunk_threshold = +libFileInfo.semantic_chunk_threshold || 90
  formState.semantic_chunk_use_model = libFileInfo.semantic_chunk_use_model || ''
  formState.ai_chunk_prumpt = libFileInfo.ai_chunk_prumpt || defaultAiChunkPrumpt
  formState.ai_chunk_model = libFileInfo.ai_chunk_model || ''
  formState.semantic_chunk_model_config_id =
    libFileInfo.semantic_chunk_model_config_id > 0 ? libFileInfo.semantic_chunk_model_config_id : ''
  formState.ai_chunk_model_config_id =
    libFileInfo.ai_chunk_model_config_id > 0 ? libFileInfo.ai_chunk_model_config_id : ''

  setTimeout(() => {
    formState.semantic_chunk_use_model = libFileInfo.semantic_chunk_use_model || ''
    formState.semantic_chunk_model_config_id =
      libFileInfo.semantic_chunk_model_config_id > 0
        ? libFileInfo.semantic_chunk_model_config_id
        : ''
  }, 200)
  if (libFileInfo.chunk_type == 0) {
    separators_no = libFileInfo.normal_chunk_default_separators_no
      ? libFileInfo.normal_chunk_default_separators_no.split(',')
      : [11, 12]
    formState.separators_no = separators_no.map((item) => +item)
    formState.chunk_size = +libFileInfo.normal_chunk_default_chunk_size || 512
    formState.not_merged_text = libFileInfo.normal_chunk_default_not_merged_text == 'true'
    formState.chunk_overlap = +libFileInfo.normal_chunk_default_chunk_overlap || 50
    formState.chunk_type = +libFileInfo.default_chunk_type
    formState.semantic_chunk_size = +libFileInfo.semantic_chunk_default_chunk_size || 512
    formState.semantic_chunk_overlap = +libFileInfo.semantic_chunk_default_chunk_overlap || 50
    formState.semantic_chunk_threshold = +libFileInfo.semantic_chunk_default_threshold || 90
    setTimeout(() => {
      formState.semantic_chunk_use_model = libFileInfo.default_use_model
      formState.semantic_chunk_model_config_id =
        libFileInfo.default_model_config_id > 0 ? libFileInfo.default_model_config_id : ''
    }, 100)
  }

  formState.is_qa_doc = props.library_type == 2 ? 1 : 0
})

const formState = reactive({
  separators_no: [], // 自定义分段-分隔符序号集
  not_merged_text: false,
  chunk_size: 512, // 自定义分段-分段最大长度 默认512，最大值不得超过2000
  chunk_overlap: 50, // 自定义分段-分段重叠长度 默认为50，最小不得低于10，最大不得超过最大分段长度的50%
  is_qa_doc: 0, // 0 普通文档 1 QA文档
  question_lable: '', // QA文档-问题开始标识符
  similar_label: '', // QA文档-相似度标识符
  answer_lable: '', // QA文档-答案开始标识符
  question_column: void 0, // excel QA文档 问题所在列
  answer_column: void 0, // excel QA文档 答案所在列
  similar_column: void 0, // excel QA文档 相似度所在列
  qa_index_type: 1, // excel QA文档 索引方式
  enable_extract_image: true,
  chunk_type: 1,
  semantic_chunk_size: 512, // 语义分段最大长度
  semantic_chunk_overlap: 50, // 语义分段重叠长度
  semantic_chunk_threshold: 90, // 语义分段阈值
  semantic_chunk_use_model: '',
  semantic_chunk_model_config_id: '',
  ai_chunk_size: 5000, // ai大模型分段最大字符数
  ai_chunk_model:'', // ai大模型分段模型名称
  ai_chunk_model_config_id: '', // ai大模型分段模型配置id
  ai_chunk_prumpt: defaultAiChunkPrumpt // ai大模型分段提示词设置
})

let baseFormState = {}
setTimeout(() => {
  baseFormState = JSON.parse(JSON.stringify(formState))
}, 500)

const vectorModelList = ref([])
const aiGenerateRef = ref(null)
// 处理AI生成按钮点击
const handleAIGenerate = () => {
  if (aiGenerateRef.value) {
    aiGenerateRef.value.showModal()
  }
}
const onVectorModelLoaded = (list) => {
  vectorModelList.value = list

  nextTick(() => {
    if (!formState.ai_chunk_model || !Number(formState.ai_chunk_model_config_id)) {
      setDefaultModel()
    }
  })
}

const setDefaultModel = () => {
  if (vectorModelList.value.length > 0) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    // 云版默认选中qwen-max
    for (let item of vectorModelList.value) {
      if (item.model_define === 'tongyi') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'qwen-max') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig && model) {
      formState.ai_chunk_model = model.name
      formState.ai_chunk_model_config_id = model.model_config_id
    }
  }
}

const handleReset = () => {
  Object.assign(formState, baseFormState)
}

const onSave = () => {
  emit('save')
}

const formRules = reactive({
  question_lable: [
    {
      message: '请输入问题开始标识符',
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject('请输入问题开始标识符')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  semantic_chunk_model_config_id: [
    {
      message: '请选择嵌入模型',
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 0 && formState.chunk_type == 2) {
          console.log(value, '==')
          if (!value) {
            return Promise.reject('请选择嵌入模型')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  ai_chunk_model_config_id: [
    {
      message: '请选择AI大模型',
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 0 && formState.chunk_type == 3) {
          console.log(value, '==')
          if (!value) {
            return Promise.reject('请选择AI大模型')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_lable: [
    {
      message: '请输入答案开始标识符',
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject('请输入答案开始标识符')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  separators_no: [
    {
      message: '请选择分段标识符',
      validator: async (rule, value) => {
        if (
          props.mode != 1 &&
          formState.is_qa_doc == 0 &&
          value.length == 0 &&
          formState.chunk_type == 1
        ) {
          return Promise.reject('请选择分段标识符')
        }

        return Promise.resolve()
      }
    }
  ],
  chunk_size: [
    {
      validator: async (rule, value) => {
        if (props.mode != 1 && formState.chunk_type == 1) {
          if (!value) {
            return Promise.reject('请输入分段最大长度')
          } else if (value > 10000) {
            return Promise.reject('最大分段长最大值不得超过10000')
          }
        } else if (formState.chunk_type == 3) {
          if (!value) {
            return Promise.reject('请输入分段最大长度')
          }
        }

        return Promise.resolve()
      }
    }
  ],
  ai_chunk_size: [
    {
      validator: async (rule, value) => {
        if (formState.chunk_type == 3) {
          if (!value) {
            return Promise.reject('请输入分段最大长度')
          }
        }

        return Promise.resolve()
      }
    }
  ],
  ai_chunk_prumpt: [
    {
      validator: async (rule, value) => {
        if (!formState.ai_chunk_prumpt) {
          if (!value) {
            return Promise.reject('请输入提示词')
          }
        }

        return Promise.resolve()
      }
    }
  ],
  chunk_overlap: [
    {
      validator: async (rule, value) => {
        if (props.mode != 1 && formState.chunk_type == 1) {
          if (value > parseInt(formState.chunk_size / 2)) {
            return Promise.reject('分段重叠长度最大不得超过最大分段长度的50%')
          }
        }
        return Promise.resolve()
      }
    }
  ],
  question_column: [
    {
      message: '请选择问题所在列',
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject('请选择问题所在列')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_column: [
    {
      message: '请选择答案所在列',
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject('请选择答案所在列')
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ]
})

const { resetFields, validate, validateInfos } = useForm(formState, formRules)

// 切换文档类型
const changeDocumentType = (val) => {
  formState.is_qa_doc = val
  onChange()
}

const reLoading = ref(false)
const saveLoading = ref(false)
const reChange = () => {
  reLoading.value = true
  if (formState.chunk_type != 3) {
    saveLoading.value = true
  }
  onChange()
}

const onChange = (showErrorAlert = true) => {
  let form = toRaw(formState)
  validate()
    .then(() => {
      emit('change', { ...form })
      emit('validate', '')
    })
    .catch((err) => {
      const { errorFields } = err
      if (errorFields.length > 0) {
        // console.log('error', err)
        if (showErrorAlert) {
          let errorName = errorFields[0].name
          // if (errorName == 'chunk_overlap') {
          //   message.error(errorFields[0]['errors'][0])
          // }
        }
        message.error(errorFields[0]['errors'][0])
        emit('validate', errorFields[0]['errors'][0])
        reLoading.value = false
        saveLoading.value = false
        return
      }
      emit('validate', '')
      emit('change', { ...form })
    })
}

// 查询分段节流
let formInEditTimer = null
const onChagneFormInput = () => {
  if (formInEditTimer) {
    clearTimeout(formInEditTimer)
    formInEditTimer = null
  }
  formInEditTimer = setTimeout(() => {
    onChange()
  }, 500)
}

// 分段标识符列表
const separatorsOptions = ref([])

const getSeparatorsOptions = () => {
  getSeparatorsList().then((res) => {
    separatorsOptions.value = res.data || []
  })
}

getSeparatorsOptions()

const handleChangeQaIndexType = (type) => {
  if (type == formState.qa_index_type) {
    return
  }
  formState.qa_index_type = type
  onChange()
}

const handleEdit = (val) => {
  formState.ai_chunk_prumpt = val
  // if (!formState.ai_chunk_prumpt) {
  //   formState.ai_chunk_prumpt += val
  // } else {
  //   formState.ai_chunk_prumpt += `\n${val}`
  // }
}

watch(() => formState.chunk_type, (val) => {
  emit('changeChunkType', val)
})

watch(() => reLoading.value, (val) => {
})

defineExpose({
  formState,
  reLoading,
  saveLoading,
  reChange
})
</script>
