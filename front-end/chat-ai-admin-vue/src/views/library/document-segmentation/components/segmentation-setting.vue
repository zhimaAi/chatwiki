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

.main-title-block{
  margin: 16px 0;
  padding-bottom: 8px;
  font-size: 14px;
  font-weight: 600;
  border-bottom: 1px solid #d9d9d9;
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
              {{ t('document_type_normal') }}
            </div>
            <div
              class="document-type"
              :class="{ active: formState.is_qa_doc == 1 }"
              @click="changeDocumentType(1)"
            >
              {{ t('document_type_qa') }}
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
            <div class="sub-setting-item-name">{{ t('file_split') }}</div>
            <!-- 表格类型的QA文档 -->
            <div class="custom-setting-form excel-qa-form" v-if="props.mode == 1">
              <div class="form-item">
                <div class="form-item-box">
                  <div class="form-item-label">{{ t('question_column') }}</div>
                  <div class="form-item-body">
                    <a-select
                      v-model:value="formState.question_column"
                      @change="onChagneFormInput"
                      :placeholder="t('ph_select_column')"
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
                  <div class="form-item-label not-required">{{ t('similar_question_column') }}</div>
                  <div class="form-item-body">
                    <a-select
                      allowClear
                      v-model:value="formState.similar_column"
                      @change="onChagneFormInput"
                      :placeholder="t('ph_select_similar_column')"
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
                  <div class="form-item-label">{{ t('answer_column') }}</div>
                  <div class="form-item-body">
                    <a-select
                      v-model:value="formState.answer_column"
                      @change="onChagneFormInput"
                      :placeholder="t('ph_select_answer_column')"
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
              <div class="sub-setting-item-name">{{ t('indexing_method') }}</div>
              <div class="indexing-methods-box">
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 1 }"
                  @click="handleChangeQaIndexType(1)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    {{ t('qa_index_with_answer') }}
                  </div>
                  <div class="list-content">
                    {{ t('qa_index_with_answer_desc') }}
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
                    {{ t('qa_index_question_only') }}
                  </div>
                  <div class="list-content">
                    {{ t('qa_index_question_only_desc') }}
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
                    >{{ t('btn_generate_segment_preview') }}</a-button
                  > -->
                </div>
                <div class="btn-box-block">
                  <a-button type="primary" block @click="onSave">{{ t('btn_save') }}</a-button>
                </div>
              </div>
            </div>
            <div class="custom-setting-form excel-qa-form" v-else>
              <div class="form-item">
                <div class="form-item-box">
                  <div class="form-item-label">{{ t('question_start_label') }}</div>
                  <div class="form-item-body">
                    <a-input
                      :placeholder="t('ph_input_label')"
                      v-model:value="formState.question_lable"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>
                <div class="form-item-box">
                  <div class="form-item-label not-required">{{ t('similar_question_start_label') }}</div>
                  <div class="form-item-body">
                    <a-input
                      :placeholder="t('ph_input_label')"
                      v-model:value="formState.similar_label"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>

                <div class="form-item-box">
                  <div class="form-item-label">{{ t('answer_start_label') }}</div>
                  <div class="form-item-body">
                    <a-input
                      :placeholder="t('ph_input_label')"
                      v-model:value="formState.answer_lable"
                      @change="onChagneFormInput"
                    />
                  </div>
                </div>
              </div>

              <div class="sub-setting-item-name">{{ t('indexing_method') }}</div>
              <div class="indexing-methods-box">
                <div
                  class="list-item"
                  :class="{ active: formState.qa_index_type == 1 }"
                  @click="handleChangeQaIndexType(1)"
                >
                  <svg-icon class="check-icon" name="check-arrow-filled"></svg-icon>
                  <div class="list-title-block">
                    <svg-icon name="file-search"></svg-icon>
                    {{ t('qa_index_with_answer') }}
                  </div>
                  <div class="list-content">
                    {{ t('qa_index_with_answer_desc') }}
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
                    {{ t('qa_index_question_only') }}
                  </div>
                  <div class="list-content">
                    {{ t('qa_index_question_only_desc') }}
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
                    >{{ t('btn_generate_segment_preview') }}</a-button
                  > -->
                </div>
                <div class="btn-box-block">
                  <a-button type="primary" block @click="onSave">{{ t('btn_save') }}</a-button>
                </div>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="custom-setting-form subsection-form" v-if="props.mode != 1">
              <div class="form-item">
                <div class="form-item-title" v-if="status === 'paragraphsSegmented'">{{ t('segmentation_cleaning_rules') }}</div>
                <div class="form-item-label">{{ t('segmentation_method') }}</div>
                <div class="form-item-body">
                  <div class="form-item-tip">
                    {{ t('segmentation_method_tip') }}
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
                        {{ t('ordinary_segmentation') }}
                      </div>
                      <div class="card-desc">
                        {{ t('ordinary_segmentation_desc') }}
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
                        {{ t('semantic_segmentation') }}
                      </div>
                      <div class="card-desc">
                        {{ t('semantic_segmentation_desc') }}
                      </div>
                    </div>
                    <div
                      v-if="props.library_type == 0"
                      class="select-card-item"
                      @click="formState.chunk_type = 4"
                      :class="{ active: formState.chunk_type == 4 }"
                    >
                      <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
                      <div class="card-title">
                        <svg-icon name="semantic-segmentation" class="title-icon"></svg-icon>
                        {{ t('parent_child_segmentation') }}
                      </div>
                      <div class="card-desc">
                        {{ t('parent_child_segmentation_desc') }}
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
                        {{ t('ai_segmentation') }}
                      </div>
                      <div class="card-desc">
                        {{ t('ai_segmentation_desc') }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <template v-if="formState.chunk_type == 1">
                <div class="form-item" style="margin-bottom: 18px">
                  <div class="form-item-label">{{ t('segmentation_separator') }}</div>
                  <div class="form-item-body">
                    <a-select
                      :placeholder="t('ph_select')"
                      style="width: 100%"
                      mode="tags"
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
                    <div class="form-item-label">{{ t('segmentation_max_length') }}</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.chunk_size"
                          :placeholder="t('ph_segmentation_max_length')"
                          :min="200"
                          :max="10000"
                          :precision="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">{{ t('unit_characters') }}</span>
                      </a-flex>
                    </div>
                  </div>

                  <div class="form-item">
                    <div class="form-item-label not-required">{{ t('segmentation_overlap_length') }}</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.chunk_overlap"
                          :placeholder="t('ph_segmentation_overlap_length')"
                          :min="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">{{ t('unit_characters') }}</span>
                      </a-flex>
                    </div>
                  </div>
                </a-flex>
                <div class="form-item">
                  <div class="form-item-body">
                    <a-flex align="center">
                      <div>{{ t('auto_merge_small_segments') }}
                        <a-tooltip :title="t('auto_merge_small_segments_tooltip')">
                           <QuestionCircleOutlined />
                        </a-tooltip>
                       ：</div>
                      <a-switch
                        :checkedValue="false"
                        :unCheckedValue="true"
                        v-model:checked="formState.not_merged_text"
                        :checked-children="t('switch_on')"
                        :un-checked-children="t('switch_off')"
                      />
                    </a-flex>
                  </div>
                </div>
              </template>
              <div v-show="formState.chunk_type == 2">
                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">{{ t('embedding_model') }}</div>
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
                    {{ t('segmentation_threshold') }}
                    <a-tooltip>
                      <template #title
                        >{{ t('segmentation_threshold_tooltip') }}</template
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
                      :placeholder="t('ph_input_segmentation_threshold')"
                      :precision="0"
                      :min="0"
                      :max="100"
                    />
                  </div>
                </div>
                <a-flex :gap="16">
                  <div class="form-item">
                    <div class="form-item-label">{{ t('segmentation_max_length') }}</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.semantic_chunk_size"
                          :placeholder="t('ph_segmentation_max_length')"
                          :min="200"
                          :max="10000"
                          :precision="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">{{ t('unit_characters') }}</span>
                      </a-flex>
                    </div>
                  </div>

                  <div class="form-item">
                    <div class="form-item-label not-required">{{ t('segmentation_overlap_length') }}</div>
                    <div class="form-item-body">
                      <a-flex align="center" :gap="8">
                        <a-input-number
                          style="flex: 1"
                          v-model:value="formState.semantic_chunk_overlap"
                          :placeholder="t('ph_segmentation_overlap_length')"
                          :min="0"
                          :formatter="(value) => parseInt(value)"
                          :parser="(value) => parseInt(value)"
                        /><span class="unit-text">{{ t('unit_characters') }}</span>
                      </a-flex>
                    </div>
                  </div>
                </a-flex>
              </div>
              <!-- AI分段 -->
              <div v-show="formState.chunk_type == 3">
                <div class="form-item" style="margin-bottom: 16px">
                  <div class="form-item-label">{{ t('ai_llm') }}</div>
                  <div class="form-item-body">
                    <ModelSelect
                      modelType="LLM"
                      :placeholder="t('ph_select_ai_llm')"
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
                    <div class="form-item-label">{{ t('prompt_setting') }}</div>
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
                    {{ t('max_characters_per_time') }}
                    <a-tooltip>
                      <template #title
                        >{{ t('max_characters_per_time_tooltip') }}</template
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
                      :placeholder="t('ph_input_max_characters')"
                      :precision="0"
                      :min="0"
                      :formatter="(value) => parseInt(value)"
                      :parser="(value) => parseInt(value)"
                    />
                    {{ t('unit_characters') }}
                  </div>
                </div>
              </div>
              <template v-if="formState.chunk_type == 4 && props.library_type == 0">
                <div class="main-title-block">{{ t('parent_block_context') }}</div>
                <div class="form-item">
                  <div class="form-item-label">{{ t('segmentation_type') }}</div>
                  <div class="form-item-body">
                    <a-radio-group v-model:value="formState.father_chunk_paragraph_type">
                      <a-radio :value="1">{{ t('segmentation_type_full_text') }}
                        <a-tooltip :title="t('segmentation_type_full_text_tooltip')">
                          <QuestionCircleOutlined />
                        </a-tooltip>
                      </a-radio>
                      <a-radio :value="2">{{ t('segmentation_type_paragraph') }}
                        <a-tooltip :title="t('segmentation_type_paragraph_tooltip')">
                          <QuestionCircleOutlined />
                        </a-tooltip>
                      </a-radio>
                    </a-radio-group>
                  </div>
                </div>
                <div class="form-item" v-if="formState.father_chunk_paragraph_type == 2">
                  <div class="form-item-label">{{ t('segmentation_separator') }}</div>
                  <div class="form-item-body">
                    <a-select
                      :placeholder="t('ph_select')"
                      style="width: 100%"
                      mode="tags"
                      v-model:value="formState.father_chunk_separators_no"
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
                <div class="form-item" v-if="formState.father_chunk_paragraph_type == 2">
                  <div class="form-item-label">{{ t('segmentation_max_length') }}</div>
                  <div class="form-item-body">
                    <a-flex align="center" :gap="8">
                      <a-input-number
                        style="flex: 1"
                        v-model:value="formState.father_chunk_chunk_size"
                        :placeholder="t('ph_segmentation_max_length')"
                        :min="200"
                        :max="10000"
                        :precision="0"
                        :formatter="(value) => parseInt(value)"
                        :parser="(value) => parseInt(value)"
                      /><span class="unit-text">{{ t('unit_characters') }}</span>
                    </a-flex>
                  </div>
                </div>
                <div class="main-title-block">{{ t('child_block_retrieval') }}</div>
                <div class="form-item">
                  <div class="form-item-label">{{ t('segmentation_separator') }}</div>
                  <div class="form-item-body">
                    <a-select
                      :placeholder="t('ph_select')"
                      style="width: 100%"
                      mode="tags"
                      v-model:value="formState.son_chunk_separators_no"
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
                <div class="form-item">
                  <div class="form-item-label">{{ t('segmentation_max_length') }}</div>
                  <div class="form-item-body">
                    <a-flex align="center" :gap="8">
                      <a-input-number
                        style="flex: 1"
                        v-model:value="formState.son_chunk_chunk_size"
                        :placeholder="t('ph_segmentation_max_length')"
                        :min="200"
                        :max="10000"
                        :precision="0"
                        :formatter="(value) => parseInt(value)"
                        :parser="(value) => parseInt(value)"
                      /><span class="unit-text">{{ t('unit_characters') }}</span>
                    </a-flex>
                  </div>
                </div>
              </template>

              <div class="btn-box-block">
                <a-button @click="handleReset" style="flex: 1">{{ t('btn_reset') }}</a-button>
                <a-button
                  @click="reChange"
                  :loading="reLoading"
                  style="flex: 1"
                  type="primary"
                  ghost
                  >{{ t('btn_generate_segment_preview') }}</a-button
                >
              </div>
              <div class="btn-box-block">
                <a-button type="primary" v-if="!props.hideSave" block :disabled="saveLoading" @click="onSave"
                  >{{ t('btn_save') }}</a-button
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
import { formatSeparatorsNo } from '@/utils/index'
import { useI18n } from '@/hooks/web/useI18n'

const { t } = useI18n('views.library.document-segmentation.components.segmentation-setting')

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

const defaultAiChunkPrumpt = t('default_ai_chunk_prompt')

watch(props, (val) => {
  let libFileInfo = val.libFileInfo
  formState.separators_no = formatSeparatorsNo(libFileInfo.separators_no, [12, 11])
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

  formState.father_chunk_paragraph_type = +libFileInfo.father_chunk_paragraph_type || 2
  formState.father_chunk_separators_no = formatSeparatorsNo(libFileInfo.father_chunk_separators_no, [12, 11])
  formState.father_chunk_chunk_size = +libFileInfo.father_chunk_chunk_size || 1024
  formState.son_chunk_separators_no = formatSeparatorsNo(libFileInfo.son_chunk_separators_no, [8, 10])
  formState.son_chunk_chunk_size = +libFileInfo.son_chunk_chunk_size || 512

  setTimeout(() => {
    formState.semantic_chunk_use_model = libFileInfo.semantic_chunk_use_model || ''
    formState.semantic_chunk_model_config_id =
      libFileInfo.semantic_chunk_model_config_id > 0
        ? libFileInfo.semantic_chunk_model_config_id
        : ''
  }, 200)
  if (libFileInfo.chunk_type == 0) {
    formState.separators_no = formatSeparatorsNo(libFileInfo.normal_chunk_default_separators_no, [12, 11])
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
  ai_chunk_prumpt: defaultAiChunkPrumpt, // ai大模型分段提示词设置

  father_chunk_paragraph_type: 2,
  father_chunk_separators_no: [],
  father_chunk_chunk_size: 1024,
  son_chunk_separators_no: [],
  son_chunk_chunk_size: 512,
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
      message: t('validation_error_question_label'),
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject(t('validation_error_question_label'))
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  semantic_chunk_model_config_id: [
    {
      message: t('validation_error_embedding_model'),
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 0 && formState.chunk_type == 2) {
          console.log(value, '==')
          if (!value) {
            return Promise.reject(t('validation_error_embedding_model'))
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  ai_chunk_model_config_id: [
    {
      message: t('validation_error_ai_llm'),
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 0 && formState.chunk_type == 3) {
          console.log(value, '==')
          if (!value) {
            return Promise.reject(t('validation_error_ai_llm'))
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_lable: [
    {
      message: t('validation_error_answer_label'),
      validator: async (rule, value) => {
        if (formState.is_qa_doc == 1 && props.mode == 0) {
          if (!value) {
            return Promise.reject(t('validation_error_answer_label'))
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  separators_no: [
    {
      message: t('validation_error_separator'),
      validator: async (rule, value) => {
        if (
          props.mode != 1 &&
          formState.is_qa_doc == 0 &&
          value.length == 0 &&
          formState.chunk_type == 1
        ) {
          return Promise.reject(t('validation_error_separator'))
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
            return Promise.reject(t('validation_error_max_length'))
          } else if (value > 10000) {
            return Promise.reject(t('validation_error_max_length_exceeded'))
          }
        } else if (formState.chunk_type == 3) {
          if (!value) {
            return Promise.reject(t('validation_error_max_characters'))
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
            return Promise.reject(t('validation_error_max_characters'))
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
            return Promise.reject(t('validation_error_prompt'))
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
            return Promise.reject(t('validation_error_overlap_exceeded'))
          }
        }
        return Promise.resolve()
      }
    }
  ],
  question_column: [
    {
      message: t('validation_error_question_column'),
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject(t('validation_error_question_column'))
          }

          return Promise.resolve()
        }

        return Promise.resolve()
      }
    }
  ],
  answer_column: [
    {
      message: t('validation_error_answer_column'),
      validator: async (rule, value) => {
        if (props.mode == 1 && formState.is_qa_doc == 1) {
          if (!value) {
            return Promise.reject(t('validation_error_answer_column'))
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
