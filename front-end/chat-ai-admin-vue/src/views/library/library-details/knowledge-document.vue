<template>
  <div class="details-library-page">
    <div
      class="group-content-box"
      :class="[{ collapsed: isHiddenGroup }, { 'no-transition': isDragging }]"
      :style="
        isHiddenGroup
          ? {}
          : {
              width: groupBoxWidth + 'px',
              minWidth: minGroupBoxWidth + 'px',
              maxWidth: maxGroupBoxWidth + 'px'
            }
      "
    >
      <div class="group-head-box">
        <div class="head-title">
          <a-tooltip title="收起问答分组">
            <div class="hover-btn-wrap" @click="handleChangeHideStatus">
              <svg-icon name="put-away"></svg-icon>
            </div>
          </a-tooltip>
          <div class="flex-between-box">
            <div>知识分组</div>
            <a-tooltip title="新建分组">
              <div class="hover-btn-wrap" @click="openGroupModal({})"><PlusOutlined /></div>
            </a-tooltip>
          </div>
        </div>
        <div class="search-box" v-if="false">
          <a-input
            v-model:value="groupSearchKey"
            allowClear
            placeholder="搜索分组"
            style="width: 100%"
          >
            <template #suffix>
              <SearchOutlined @click.stop="" />
            </template>
          </a-input>
        </div>
      </div>
      <div class="classify-box" style="margin-top: 16px">
        <cu-scroll style="padding-right: 24px">
          <Draggable
            v-model="groupLists"
            item-key="id"
            :disabled="isDragging"
            @end="handleDragEnd"
            :move="checkMove"
            handle=".drag-btn"
          >
            <template #item="{ element: item }">
              <div
                class="classify-item"
                @click="handleChangeGroup(item)"
                :class="{ active: item.id == groupId }"
              >
                <!-- 原有内容保持不变 -->
                <span v-if="item.id > 0" class="drag-btn"><svg-icon name="drag" /></span>
                <span v-else style="width: 20px"></span>
                <div class="classify-title">{{ item.group_name }}</div>
                <div class="right-content">
                  <div class="num" :class="{ 'num-block': item.id <= 0 }">{{ item.total }}</div>
                  <div class="btn-box" v-if="item.id > 0">
                    <a-dropdown placement="bottomRight">
                      <div class="hover-btn-wrap">
                        <EllipsisOutlined />
                      </div>
                      <template #overlay>
                        <a-menu>
                          <a-menu-item>
                            <div @click.stop="openGroupModal(item)">重命名</div>
                          </a-menu-item>
                          <a-menu-item>
                            <div @click.stop="handleDelGroup(item)">删 除</div>
                          </a-menu-item>
                        </a-menu>
                      </template>
                    </a-dropdown>
                  </div>
                </div>
              </div>
            </template>
          </Draggable>
        </cu-scroll>
      </div>

      <div v-if="!isHiddenGroup" class="resize-bar" @mousedown="handleResizeMouseDown"></div>
    </div>
    <div class="main-content-box">
      <cu-scroll :scrollbar="false">
        <div class="group-header-title">
          <a-popover :title="null" placement="bottom" v-if="isHiddenGroup">
            <template #content>
              <div class="pover-group-content-box">
                <div class="group-head-box">
                  <div class="head-title">
                    <div class="flex-between-box">
                      <div>问答分组</div>
                      <a-tooltip title="新建分组">
                        <div class="hover-btn-wrap" @click="openGroupModal({})">
                          <PlusOutlined />
                        </div>
                      </a-tooltip>
                    </div>
                  </div>
                  <div class="search-box">
                    <a-input
                      allowClear
                      v-model:value="groupSearchKey"
                      placeholder="搜索分组"
                      style="width: 100%"
                    >
                      <template #suffix>
                        <SearchOutlined @click.stop="" />
                      </template>
                    </a-input>
                  </div>
                </div>
                <div class="classify-box classify-scroll-box">
                  <div
                    class="classify-item"
                    @click="handleChangeGroup(item)"
                    :class="{ active: item.id == groupId }"
                    v-for="item in filterGroupLists"
                    :key="item.id"
                  >
                    <div class="classify-title">{{ item.group_name }}</div>
                    <div class="right-content">
                      <div class="num" :class="{ 'num-block': item.id <= 0 }">
                        {{ item.total }}
                      </div>
                      <div class="btn-box" v-if="item.id > 0">
                        <a-dropdown placement="bottomRight">
                          <div class="hover-btn-wrap">
                            <EllipsisOutlined />
                          </div>
                          <template #overlay>
                            <a-menu>
                              <a-menu-item>
                                <div @click.stop="openGroupModal(item)">重命名</div>
                              </a-menu-item>
                              <a-menu-item>
                                <div @click.stop="handleDelGroup(item)">删 除</div>
                              </a-menu-item>
                            </a-menu>
                          </template>
                        </a-dropdown>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
            <div class="hover-btn-wrap" v-if="viewMode !== 'folder' || viewMode === 'list'" @click="handleChangeHideStatus">
              <svg-icon name="expand"></svg-icon>
            </div>
          </a-popover>
          <div class="title">
            <template v-if="inFolderDetail">
              <a-breadcrumb>
                <a-breadcrumb-item>
                  <a class="breadcrumb-a" @click="goToRootFolder">知识库文档</a>
                </a-breadcrumb-item>
                <a-breadcrumb-item v-if="groupId === -1">全部分组</a-breadcrumb-item>
                <a-breadcrumb-item v-else>{{ currentGroupItem.group_name }}</a-breadcrumb-item>
              </a-breadcrumb>
            </template>
            <template v-else>
              <span class="title-text">知识库文档</span>
            </template>
          </div>
        </div>
        <div class="list-tools">
          <div class="tools-items">
            <div class="tool-item">
              <div class="view-toggle">
                <div
                  class="toggle-seg"
                  :class="{ active: viewMode === 'folder' }"
                  @click="switchToFolderView"
                  aria-label="文件夹视图"
                >
                  <svg-icon
                    class="toggle-icon"
                    :name="viewMode === 'folder' ? 'folder-icon-active' : 'folder-icon'"
                  ></svg-icon>
                </div>
                <div
                  class="toggle-seg"
                  :class="{ active: viewMode === 'list' }"
                  @click="switchToListView"
                  aria-label="缩略图视图"
                >
                  <svg-icon
                    class="toggle-icon"
                    :name="viewMode === 'list' ? 'thumbnail-icon-active' : 'thumbnail-icon'"
                  ></svg-icon>
                </div>
              </div>
            </div>
            <div class="tool-item" v-if="viewMode === 'folder' && inFolderDetail">
              <div class="sort-btn" ref="sortBtnRef">
                <div class="sort-left">
                  <span class="sort-text">{{ folderSort.order_field === 'total_hits' ? '合计' : folderSort.order_field === 'today_hits' ? '今日' : folderSort.order_field === 'yesterday_hits' ? '昨日' : '默认' }}</span>
                  <svg-icon
                    v-if="!folderSort.order_field"
                    key="default"
                    class="sort-icon active"
                    name="sort-icon"
                  />
                  <svg-icon
                    v-else-if="folderSort.order_type === 'DESC'"
                    key="desc"
                    class="sort-icon active"
                    name="sort-down-icon"
                    @click.stop="onClickSortIcon"
                  />
                  <svg-icon
                    v-else
                    key="asc"
                    class="sort-icon active"
                    name="sort-top-icon"
                    @click.stop="onClickSortIcon"
                  />
                </div>
                <a-dropdown v-model:open="sortOpen" trigger="click">
                  <span class="sort-arrow" :class="{ active: sortOpen }" @click.prevent="toggleSortDropdown">
                    <DownOutlined />
                  </span>
                  <template #overlay>
                    <a-menu :selectedKeys="[(folderSort.order_field || 'default')]">
                      <a-menu-item :key="'default'" @click="setFolderSort('default')">默认</a-menu-item>
                      <a-menu-item :key="'total_hits'" @click="setFolderSort('total_hits')">合计</a-menu-item>
                      <a-menu-item :key="'today_hits'" @click="setFolderSort('today_hits')">今日</a-menu-item>
                      <a-menu-item :key="'yesterday_hits'" @click="setFolderSort('yesterday_hits')">昨日</a-menu-item>
                    </a-menu>
                  </template>
                </a-dropdown>
              </div>
            </div>
            <a-flex align="center" class="tool-item custom-select-box" v-show="false">
              <span>嵌入模型：</span>
              <model-select
                modelType="TEXT EMBEDDING"
                :isOffline="false"
                :modeName="modelForm.use_model"
                :modeId="modelForm.model_config_id"
                style="width: 240px"
                @change="onChangeModel"
                @loaded="onVectorTextModelLoaded"
              />
            </a-flex>
          </div>
          <div class="right-box">
            <div class="status-box" v-if="libraryInfo.type == 0">
              <div class="status-item">
                <div class="status-label">已学习：</div>
                <div class="status-content">{{ count_data.learned_count }}</div>
              </div>
              <div class="status-item" v-if="count_data.learned_wait_count > 0">
                <div class="status-label">，未学习：</div>
                <div class="status-content content-tip">{{ count_data.learned_wait_count }}</div>
              </div>
              <div class="status-item" v-if="count_data.learned_err_count > 0">
                <div class="status-label">，学习失败：</div>
                <div class="status-content content-tip">{{ count_data.learned_err_count }}</div>
              </div>
            </div>
            <div class="status-select-box" v-if="libraryInfo.type == 0">
              <a-select
                v-model:value="queryParams.status"
                style="width: 150px"
                placeholder="请选择"
                @change="handleChangeStatus"
              >
                <a-select-option v-for="item in statusList" :key="item.id" :value="item.id">{{
                  item.name
                }}</a-select-option>
              </a-select>
            </div>
            <div class="tool-item">
              <a-input
                style="width: 140px"
                v-model:value="queryParams.file_name"
                placeholder="文档名称搜索"
                @change="onSearch"
              >
                <template #suffix>
                  <SearchOutlined @click="onSearch" style="color: rgba(0, 0, 0, 0.25)" />
                </template>
              </a-input>
            </div>
            <a-flex align="center" v-if="neo4j_status" class="tool-item custom-select-box pd-5-8">
              <ModelSelect modelType="LLM" v-show="false" @loaded="onVectorModelLoaded" />
              <span>生成知识图谱：</span>
              <a-switch
                v-model:checked="createGraphSwitch"
                @change="createGraphSwitchChange"
                checked-children="开"
                un-checked-children="关"
              />
            </a-flex>
            <a-button v-if="props.type != 3" @click="showMetaModal(1)">元数据 <SettingOutlined/></a-button>
            <div class="tool-item">
              <a-dropdown :trigger="['hover']" overlayClassName="add-dropdown-btn">
                <template #overlay>
                  <a-menu @click="handleMenuClick">
                    <a-menu-item :key="1">
                      <div class="dropdown-btn-menu">
                        <a-flex class="title-block" :gap="4">
                          <svg-icon name="doc-icon"></svg-icon>
                          <div class="title">本地文档</div>
                        </a-flex>
                        <div class="desc" v-if="libraryInfo.type == 2">
                          上传本地 docx/csv/xlsx 等格式文件
                        </div>
                        <div class="desc" v-else>
                          上传本地 pdf/docx/ofd/txt/md/xlsx/csv/html 等格式文件
                        </div>
                      </div>
                    </a-menu-item>
                    <a-menu-item :key="2" v-if="libraryInfo.type != 2">
                      <div class="dropdown-btn-menu">
                        <a-flex class="title-block" :gap="4">
                          <svg-icon name="link-icon"></svg-icon>
                          <div class="title">在线数据</div>
                        </a-flex>
                        <div class="desc">获取在线网页内容</div>
                      </div>
                    </a-menu-item>
                    <a-menu-item :key="3">
                      <div class="dropdown-btn-menu">
                        <a-flex class="title-block" :gap="4">
                          <svg-icon name="cu-doc-icon"></svg-icon>
                          <div class="title">自定义文档</div>
                        </a-flex>
                        <div class="desc">自定义一个空文档，手动添加内容</div>
                      </div>
                    </a-menu-item>
                    <a-menu-item :key="4" v-if="libraryInfo.type == 0">
                      <div class="dropdown-btn-menu">
                        <a-flex class="title-block" :gap="4">
                          <svg-icon name="feishu-doc"></svg-icon>
                          <div class="title">飞书知识库</div>
                        </a-flex>
                        <div class="desc">在线获取飞书知识库下docx格式云文档内容</div>
                      </div>
                    </a-menu-item>
                    <!-- 添加分组 -->
                    <a-menu-item :key="5" v-if="viewMode === 'folder'">
                      <div class="dropdown-btn-menu">
                        <a-flex class="title-block" :gap="4">
                          <svg-icon name="group-default-icon"></svg-icon>
                          <div class="title">添加分组</div>
                        </a-flex>
                        <div class="desc">便于文档更好的归类</div>
                      </div>
                    </a-menu-item>
                  </a-menu>
                </template>
                <a-button v-if="type != 3" type="primary">
                  <template #icon>
                    <PlusOutlined />
                  </template>
                  <span>添加内容</span>
                </a-button>
              </a-dropdown>
            </div>
          </div>
        </div>
        <page-alert style="margin-bottom: 16px" title="使用说明" v-if="libraryInfo.type == 0 && viewMode==='list'">
          <div>
            <p>
              1、如果单次上传一个文档，上传成功后，系统会自动学习；如果单次上传多个文档，上传成功后，需要手动点击文档后面"学习"进行学习；如果解析失败，支持重新学习。
            </p>
            <p>2、未学习的文档数据不会被检索到。</p>
          </div>
        </page-alert>
        <!-- 文件夹视图 -->
        <div class="folder-view" v-if="viewMode==='folder'">
          <div class="empty-folder" v-if="inFolderDetail && folderItems.length === 0">
            <div class="empty-inner">
              <img :src="emptyDocumentIcon" alt="" class="empty-img" />
              <div class="empty-text">暂无文档</div>
            </div>
          </div>
          <div class="folder-grid" v-else>
            <div
              v-for="item in folderItems"
              :key="`${item.__type}-${item.id}`"
              class="folder-card"
              :class="[`card-${item.__type}`]"
              @click="onCardClick(item)"
            >
              <div class="folder-thumb" :class="['status-' + (item.cover_status || 'success')]">
                <template v-if="item.__type === 'group'">
                  <svg-icon name="group-icon" class="folder-type-icon"></svg-icon>
                </template>
                <template v-else>
                  <img v-if="item.cover_src" :src="item.cover_src" alt="" class="thumb-img" />
                  <img class="thumb-img" v-else src="@/assets/img/default-bg.png">
                  <div class="status-badge" v-if="item.status_label && item.cover_status !== 'success'">
                    <a-tooltip v-if="item.cover_status === 'error' && item.origin.errmsg && item.origin.errmsg != 'success'">
                      <template #title>
                        <span>{{ item.origin.errmsg }}</span>
                      </template>
                      <svg-icon style="color: rgba(0, 0, 0, 0.4); font-size: 16px; vertical-align: sub;" name="tip-icon-two"></svg-icon>
                      {{ item.status_label }}
                    </a-tooltip>
                    <svg-icon v-if="item.cover_status === 'error' && !item.origin.errmsg" style="color: rgba(0, 0, 0, 0.4); font-size: 16px; vertical-align: sub;" name="tip-icon-two"></svg-icon>
                    <svg-icon v-if="item.cover_status === 'wait'" style="color: rgba(0, 0, 0, 0.4); font-size: 16px; vertical-align: sub;" name="time-default-icon"></svg-icon>
                    <span v-if="item.origin.errmsg == 'success' && item.cover_status === 'loading' && getLearnPercent(item.origin) !== null" class="learn-percent">{{ getLearnPercent(item.origin) }}%</span>
                    <span v-if="!item.origin.errmsg || item.origin.errmsg == 'success'" style="margin-left: 2px;">{{ item.status_label }}</span>
                  </div>
                </template>
              </div>
              <div class="folder-meta">
                <template v-if="item.__type === 'group'">
                  <div class="name zm-line1">{{ item.group_name }}</div>
                </template>
                <template v-else>
                  <a-popover :title="null" v-if="item.origin.doc_type == 2">
                    <template #content>
                      原链接：<a :href="item.origin.doc_url" target="_blank">{{ item.origin.doc_url }} </a>
                      <CopyOutlined v-copy="`${item.origin.doc_url}`" style="margin-left: 4px; cursor: pointer" />
                    </template>
                    <div class="name name-2line">
                      {{ item.title }}
                    </div>
                  </a-popover>
                  <div class="name name-2line" v-else>{{ item.title }}</div>
                </template>
                <div class="desc">
                  <template v-if="item.__type === 'group'">
                    <span>{{ item.total || '0' }}个知识库</span>
                    <span>{{ formatCardTime(item.last_update_ts) }}</span>
                  </template>
                  <template v-else>
                    <span class="file-ext">
                      <svg-icon
                        :name="['docx','excel','html','pdf','txt','xlsx','ofd','csv','md'].some(t => String(item.origin.file_ext || '').includes(t)) ? item.origin.file_ext : 'file_ext'"
                        class="file-icon"
                      />
                      {{ item.origin.file_ext || '--' }}
                    </span>
                    <span>{{ formatCardTime(item.update_ts) }}</span>
                  </template>
                </div>
              </div>
              <a-dropdown v-if="item.__type === 'group' && item.id > 0">
                <div class="card-more" @click.stop>
                  <EllipsisOutlined />
                </div>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="openGroupModal({ id: item.id, group_name: item.group_name })">
                      <a-flex :gap="8" align="center"><svg-icon name="edit"></svg-icon><div>重命名</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="handleDelGroup({ id: item.id, group_name: item.group_name })">
                      <a-flex :gap="8" align="center"><svg-icon name="delete"></svg-icon><div>删除</div></a-flex>
                    </a-menu-item>
                  </a-menu>
                </template>
              </a-dropdown>
              <a-dropdown v-if="item.__type === 'file'">
                <div class="card-more" @click.stop>
                  <EllipsisOutlined />
                </div>
                <template #overlay>
                  <a-menu>
                    <a-menu-item @click="showBaseData(item.origin)">
                      <a-flex :gap="8" align="center"><svg-icon name="base-data-icon"></svg-icon><div>基础数据</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="handlePreview(item.origin)">
                      <a-flex :gap="8" align="center"><svg-icon name="preview-icon"></svg-icon><div>预览</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="toReSegmentationPage(item.origin)">
                      <a-flex :gap="8" align="center"><SyncOutlined class="action-icon" /><div>重新分段</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="handleOpenRenameModal(item.origin)">
                      <a-flex :gap="8" align="center"><svg-icon name="edit"></svg-icon><div>重命名</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="openEditGroupModal(item.origin)">
                      <a-flex :gap="8" align="center"><svg-icon name="group-default-icon"></svg-icon><div>修改分组</div></a-flex>
                    </a-menu-item>
                    <a-menu-item @click="handleDownload(item.origin)">
                      <a-flex :gap="8" align="center"><svg-icon name="down-file"></svg-icon><div>下载文档</div></a-flex>
                    </a-menu-item>
                    <a-popconfirm title="确定要删除吗?" placement="topRight" trigger="click" @confirm="onDelete(item.origin)">
                      <a-menu-item>
                        <a-flex :gap="8" align="center"><svg-icon name="delete"></svg-icon><div>删除</div></a-flex>
                      </a-menu-item>
                    </a-popconfirm>
                  </a-menu>
                </template>
              </a-dropdown>
            </div>
            <div ref="folderScrollSentinel" class="scroll-sentinel"></div>
          </div>
          <BaseDataDrawer
            v-model:open="baseDataOpen"
            :record="baseDataRecord"
            @close="baseDataOpen=false"
            @editOnline="handleEditOnlineDoc"
          />
        </div>
        <!-- 列表视图 -->
        <div class="list-content" v-else>
          <a-table
            @resizeColumn="handleResizeColumn"
            :columns="columns"
            :data-source="fileList"
            :scroll="{ x: 1000 }"
            row-key="id"
            :pagination="{
              current: queryParams.page,
              total: queryParams.total,
              pageSize: queryParams.size,
              showQuickJumper: true,
              showSizeChanger: true,
              pageSizeOptions: ['10', '20', '50', '100']
            }"
            @change="onTableChange"
          >
            <template #headerCell="{column}">
              <template v-if="column.key.indexOf('meta_') > -1">
                <a-tooltip :title="column.title?.length > 8 ? column.title : ''">
                  <div class="zm-line1">{{column.title}}</div>
                </a-tooltip>
              </template>
              <template v-else-if="column.key === 'file_name'">
                文档名称 (共{{ queryParams.total || 0 }}个)
              </template>
              <template v-else-if="column.key === 'graph_entity_count'">
                <div class="title-box">
                  <div class="zm-line1">实体数</div>
                  <a-tooltip title="知识图谱实体数">
                    <QuestionCircleOutlined />
                  </a-tooltip>
                </div>
              </template>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'file_name'">
                <div class="doc-name-td">
                  <!-- 根据file_ext显示图标 -->
                  <a-tooltip :title="`.${record.file_ext}`">
                    <svg-icon
                      :name="['docx','excel','html','pdf','txt','xlsx','ofd','csv','md'].some(t => String(record.file_ext || '').includes(t)) ? record.file_ext : 'file_ext'"
                      class="file-icon"
                    />
                  </a-tooltip>
                  <a-popover :title="null" v-if="record.doc_type == 2">
                    <template #content>
                      原链接：<a :href="record.doc_url" target="_blank">{{ record.doc_url }} </a>
                      <CopyOutlined
                        v-copy="`${record.doc_url}`"
                        style="margin-left: 4px; cursor: pointer"
                      />
                    </template>
                    <a @click="handlePreview(record, { open_new: true })">
                      <a-tooltip :title="getTooltipTitle((['5','6','7'].includes(record.status) ? record.doc_url : record.file_name), record, 14, 2, 12)" placement="top">
                        <span class="doc-name-text" :ref="el => setDescRef(el, record)">
                          <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                          <span v-else>{{ record.file_name }}</span>
                        </span>
                      </a-tooltip>
                    </a>
                  </a-popover>
                  <a @click="handlePreview(record, { open_new: true })" v-else>
                    <a-tooltip :title="getTooltipTitle((['5','6','7'].includes(record.status) ? record.doc_url : record.file_name), record, 14, 2, 12)" placement="top">
                      <span class="doc-name-text" :ref="el => setDescRef(el, record)">
                        <span v-if="['5', '6', '7'].includes(record.status)">{{ record.doc_url }}</span>
                        <span v-else>{{ record.file_name }}</span>
                      </span>
                    </a-tooltip>
                  </a>
                  <div v-if="record.doc_type == 2 && record.remark" class="url-remark">
                    备注：{{ record.remark }}
                  </div>
                </div>
              </template>
              <template v-if="column.key === 'graph_entity_count'">
                <a @click="toGraph(record)">{{ record.graph_entity_count }}</a>
              </template>
              <!-- 状态 -->
              <!-- <template v-if="column.key === 'status'">
                <div class="status-icon-row">
                  {{ record.status }}
                  <a-tooltip title="学习完成" v-if="record.status == 2">
                    <svg-icon name="accomplish-icon" class="status-icon" />
                  </a-tooltip>
                  <a-tooltip title="学习失败" v-if="[3, 7, 8].includes(Number(record.status))">
                    <svg-icon
                      name="fail-icon"
                      class="status-icon"
                    />
                  </a-tooltip>
                  <a-tooltip title="转换中" v-if="record.status == 0">
                    <svg-icon name="time-icon" class="status-icon" />
                  </a-tooltip>
                </div>
              </template> -->
              <template v-if="column.key === 'status'">
                <template v-if="record.file_ext == 'pdf' && record.pdf_parse_type >= 2">
                  <div class="pdf-progress-box" v-if="record.status == 0">
                    <div class="progress-title">
                      <span class="status-box"><LoadingOutlined />文档解析中</span>
                      <a @click="handleCancelOcrPdf(record)">取消</a>
                    </div>
                    <div class="progress-bar">
                      <a-progress
                        size="small"
                        class="progress-bar-box"
                        :percent="parseInt((record.ocr_pdf_index / record.ocr_pdf_total) * 100)"
                        :show-info="false"
                      />
                      <div class="num-box">
                        {{ record.ocr_pdf_index }} / {{ record.ocr_pdf_total }}
                      </div>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <span class="status-tag running" v-if="record.status == 0"
                    ><a-spin size="small" /> 转换中</span
                  >
                </template>

                <span class="status-tag running" v-if="record.status == 1"
                  ><a-spin size="small" /> 学习中</span
                >

                <span class="status-tag complete" v-if="record.status == 2"
                  ><CheckCircleFilled /> 学习完成</span
                >

                <a-tooltip placement="top" v-if="record.status == 3">
                  <template #title>
                    <span>{{ record.errmsg }}</span>
                  </template>
                  <span>
                    <span class="status-tag status-error"><CloseCircleFilled /> 转换失败</span>
                    <a class="ml8" v-if="libraryInfo.type == 2" @click="handlePreview(record)"
                      >学习</a
                    >
                  </span>
                </a-tooltip>
                <a-tooltip placement="top" v-if="record.status == 8">
                  <template #title>
                    <span>{{ record.errmsg }}</span>
                  </template>
                  <span>
                    <span class="status-tag status-error"><CloseCircleFilled /> 转化异常</span>
                  </span>
                </a-tooltip>
                <template v-if="record.status == 4">
                  <span class="status-tag"><ClockCircleFilled /> 待学习</span>
                  <!-- 产品说去掉这个学习，因为和重新分段走的一样的逻辑 -->
                  <!-- <a class="ml8" @click="handlePreview(record)">学习</a> -->
                </template>
                <template v-if="record.status == 5">
                  <span class="status-tag"><ClockCircleFilled /> 待获取</span>
                </template>
                <span class="status-tag running" v-if="record.status == 6"
                  ><a-spin size="small" /> 获取中</span
                >
                <a-tooltip placement="top" v-if="record.status == 7">
                  <template #title>
                    <span>{{ record.errmsg }}</span>
                  </template>
                  <span class="status-tag error"><CloseCircleFilled /> 获取失败</span>
                </a-tooltip>
                <template v-if="record.status == 9">
                  <span class="status-tag cancel"><ExclamationCircleOutlined /> 取消解析</span>
                </template>

                <span class="status-tag subning" v-if="record.status == 10">
                  <a-spin size="small" />
                  正在分段
                </span>
              </template>
              <template v-if="column.key === 'graph_status'">
                <!--0待生成 1排队中 2生成完成 3生成失败 4生成中 5部分成功-->
                <template v-if="record.graph_status == 0">
                  <span class="status-tag"><ClockCircleFilled /> 待生成</span>
                  <a class="ml8" @click="createGraphTask(record)">生成</a>
                </template>
                <span v-else-if="record.graph_status == 1" class="status-tag running"
                  ><HourglassFilled /> 排队中</span
                >
                <span v-else-if="record.graph_status == 2" class="status-tag complete"
                  ><CheckCircleFilled /> 生成完成</span
                >
                <template v-else-if="record.graph_status == 3">
                  <span class="status-tag error"><CloseCircleFilled /> 生成失败</span>
                  <a class="ml8" @click="createGraphTask(record)">生成</a>
                  <a-tooltip v-if="record.graph_err_msg" :title="record.graph_err_msg">
                    <div class="zm-line1 reason-text">原因：{{ record.graph_err_msg }}</div>
                  </a-tooltip>
                </template>
                <span v-else-if="record.graph_status == 4" class="status-tag running"
                  ><a-spin size="small" /> 生成中</span
                >
                <template v-else-if="record.graph_status == 5">
                  <span class="status-tag warning"><CheckCircleFilled /> 部分成功</span>
                  <div class="reason-text">
                    失败数：{{ record.graph_err_count || 0 }}
                    <a class="ml8" @click="handlePreview(record, { graph_status: 3 })">详情</a>
                  </div>
                </template>
              </template>
              <template v-if="column.key === 'file_size'">
                <span>
                  {{ record.doc_type == 3 ? '-' : record.file_size_str }}
                  /
                  {{ record.status == 0 || record.status == 1 ? '-' : record.paragraph_count }}
                </span>
              </template>
              <template v-if="column.key === 'update_info'">
                <div class="text-block">
                  <span class="time-content-box">
                    {{ formatUpdateShort(record.update_time) }}
                    <a-popover :title="null" placement="top" overlayClassName="update-popover">
                      <template #content>
                        <div class="update-tooltip">
                          <div class="update-tooltip-title">
                            更新频率：<span class="update-tooltip-value">{{ formatUpdateFrequency(record) }}</span>
                            <div class="update-icon-box">
                              <svg-icon
                                name="edit"
                                class="edit-icon"
                              @click.stop="handleEditOnlineDoc(record, true)"
                              ></svg-icon>
                            </div>
                          </div>
                          <div class="update-tooltip-time">
                            更新时间：<span class="update-tooltip-value">{{ record.update_time }}</span>
                            <a-popconfirm title="确认更新?" @confirm="handleUpdataDoc(record)">
                              <a class="ml4 btn-hover-block">更新</a>
                            </a-popconfirm>
                          </div>
                        </div>
                      </template>
                      <span class="update-trigger-text">
                        <template v-if="record.doc_type == 2 && record.doc_auto_renew_frequency">
                          <span class="ml4">/</span>
                          <svg-icon
                            class="update-icon"
                            :name="record.doc_auto_renew_frequency == 1 ? 'forbid-update-icon' : 'update-icon'"
                          ></svg-icon>
                        </template>
                      </span>
                    </a-popover>
                  </span>
                </div>
              </template>
              <template v-if="column.key === 'action'">
                <a-flex :gap="8">
                  <a-tooltip>
                    <template #title>重新分段</template>
                    <SyncOutlined class="btn-hover-block" @click="toReSegmentationPage(record)" />
                  </a-tooltip>
                  <a-dropdown>
                    <div class="table-btn" @click.prevent>
                      <MoreOutlined class="btn-hover-block" />
                    </div>
                    <template #overlay>
                      <a-menu>
                        <a-menu-item v-if="record.doc_type == 2">
                          <div @click="handleEditOnlineDoc(record)">编辑</div>
                        </a-menu-item>
                        <a-menu-item
                          :disabled="record.status == 6 || record.status == 7 || record.status == 0"
                        >
                          <div @click="handlePreview(record)">预览</div>
                        </a-menu-item>
                        <a-menu-item @click="handleOpenRenameModal(record)">
                          <div>重命名</div>
                        </a-menu-item>
                        <a-popconfirm title="确定要删除吗?" @confirm="onDelete(record)" placement="topRight">
                          <a-menu-item>
                            <span>删除</span>
                          </a-menu-item>
                        </a-popconfirm>
                        <a-menu-item @click="openEditGroupModal(record)">
                          <div>修改分组</div>
                        </a-menu-item>
                        <a-menu-item @click="handleDownload(record)">
                          <div>下载文档</div>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </a-flex>
              </template>
              <template v-if="column.key.indexOf('meta_') > -1 && record[column.key]">
                <a-tooltip :title="record[column.key].length > 8 ? record[column.key] : ''">
                  <div class="zm-line1">{{record[column.key]}}</div>
                </a-tooltip>
              </template>
            </template>
          </a-table>
        </div>
      </cu-scroll>
    </div>

    <a-modal
      v-model:open="addFileState.open"
      :confirm-loading="addFileState.confirmLoading"
      :maskClosable="false"
      title="上传文档"
      @ok="handleSaveFiles"
      @cancel="handleCloseFileUploadModal"
      width="740px"
    >
      <a-form class="mt24">
        <a-form-item required label="">
          <div class="upload-file-box">
            <UploadFilesInput
              :type="libraryInfo.type"
              :maxCount="20"
              v-model:value="addFileState.fileList"
              @change="onFilesChange"
            />
          </div>
        </a-form-item>
        <a-form-item v-if="existPdfFile" required label="PDF解析模式" class="select-card-main">
          <div class="select-card-box">
            <div
              v-for="item in PDF_PARSE_MODE"
              :key="item.key"
              :class="[
                'select-card-item',
                { active: addFileState.pdf_parse_type == item.key },
                { disabled: item.key == 4 && !ali_ocr_switch }
              ]"
              @click.stop="pdfParseTypeChange(item)"
            >
              <svg-icon class="check-arrow" name="check-arrow-filled"></svg-icon>
              <div class="card-title">
                <svg-icon :name="item.icon" style="font-size: 16px"></svg-icon>
                {{ item.title }}
                <div class="card-switch-box" v-if="item.key == 4 && !ali_ocr_switch">未开启</div>
              </div>
              <div class="card-desc">{{ item.desc }}</div>
              <div class="card-switch" v-if="item.key == 4 && !ali_ocr_switch">
                未开启阿里云OCR
                <div class="card-switch-btn" @click.stop="onGoSwitch">去开启</div>
              </div>
            </div>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
    <QaUploadModal @ok="getData" ref="qaUploadModalRef" />
    <a-modal
      v-model:open="addUrlState.open"
      :confirm-loading="addUrlState.confirmLoading"
      :maskClosable="false"
      title="添加在线数据"
      width="746px"
      @ok="handleSaveUrl"
      @cancel="handleCloseUrlModal"
    >
      <a-form
        class="url-add-form"
        layout="vertical"
        ref="urlFormRef"
        :model="addUrlState"
        :rules="addUrlState.rules"
      >
        <a-form-item name="urls" label="网页链接">
          <a-textarea
            style="height: 120px"
            v-model:value="addUrlState.urls"
            placeholder="请输入网页链接,形式：一行标题一行网页链接"
          />
        </a-form-item>
        <a-form-item name="doc_auto_renew_frequency" label="更新频率" required>
          <a-select v-model:value="addUrlState.doc_auto_renew_frequency" style="width: 100%">
            <a-select-option :value="1">不自动更新</a-select-option>
            <a-select-option :value="2">每天</a-select-option>
            <a-select-option :value="3">每3天</a-select-option>
            <a-select-option :value="4">每7天</a-select-option>
            <a-select-option :value="5">每30天</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item
          v-if="addUrlState.doc_auto_renew_frequency > 1"
          name="doc_auto_renew_minute"
          label="更新时间"
        >
          <a-time-picker
            valueFormat="HH:mm"
            v-model:value="addUrlState.doc_auto_renew_minute"
            format="HH:mm"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <AddCustomDocument
      :libraryInfo="libraryInfo"
      :group_id="groupId"
      @ok="onSearch"
      ref="addCustomDocumentRef"
    ></AddCustomDocument>
    <RenameModal @ok="onSearch" ref="renameModalRef" />
    <OpenGrapgModal @ok="handleOpenGrapgOk" ref="openGrapgModalRef" />
    <GuideLearningTips ref="guideLearningTipsRef" />
    <EditOnlineDoc @ok="getData" ref="editOnlineDocRef" />
    <AddGroup group_type="1" ref="addGroupRef" @ok="initData" />
    <EditGroup :libraryId="libraryId" :sense="1" ref="editGroupRef" @ok="initData" />
    <AddFeishuDocument ref="feishuRef" :libraryId="libraryId" @ok="initData"/>
    <MetadataManageModal ref="metaRef" :library-id="libraryId" @change="initData"/>
  </div>
</template>

<script setup>
import { useStorage } from '@/hooks/web/useStorage'
import { reactive, ref, toRaw, onUnmounted, onMounted, computed, createVNode, nextTick } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute, useRouter } from 'vue-router'
import {
  PlusOutlined,
  SearchOutlined,
  CheckCircleFilled,
  CloseCircleFilled,
  ClockCircleFilled,
  MoreOutlined,
  HourglassFilled,
  ExclamationCircleOutlined,
  SyncOutlined,
  CopyOutlined,
  EllipsisOutlined,
  LoadingOutlined,
  DownOutlined,
  SettingOutlined,
  QuestionCircleOutlined
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  getLibraryFileList,
  delLibraryFile,
  addLibraryFile,
  editLibrary,
  createGraph,
  manualCrawl,
  cancelOcrPdf,
  getLibraryGroup,
  deleteLibraryGroup,
  sortLibararyListGroup
} from '@/api/library'
import { formatFileSize, setDescRef, getTooltipTitle } from '@/utils/index'
import UploadFilesInput from '../add-library/components/upload-input.vue'
import { transformUrlData } from '@/utils/validate.js'
import AddCustomDocument from './components/add-custom-document.vue'
import ModelSelect from '@/components/model-select/model-select.vue'
import OpenGrapgModal from '@/views/library/library-details/components/open-grapg-modal.vue'
import QaUploadModal from './components/qa-upload-modal.vue'
import RenameModal from './components/rename-modal.vue'
import PageAlert from '@/components/page-alert/page-alert.vue'
import GuideLearningTips from './components/guide-learning-tips.vue'
import EditOnlineDoc from './components/edit-online-doc.vue'
import BaseDataDrawer from './components/base-data-drawer.vue'
import { useUserStore } from '@/stores/modules/user'
import defaultCover from '@/assets/img/default-bg.png'
import emptyDocumentIcon from '@/assets/svg/empty-document.svg'
import { useLibraryStore } from '@/stores/modules/library'
import { convertTime } from '@/utils/index'
import { useCompanyStore } from '@/stores/modules/company'
import AddGroup from './qa-knowledge-document/components/add-group.vue'
import EditGroup from './qa-knowledge-document/components/edit-group.vue'
import Draggable from 'vuedraggable'
import AddFeishuDocument from "./components/add-feishu-document.vue";
import MetadataManageModal from "@/views/library/library-details/components/metadata-manage-modal.vue";

const { setStorage } = useStorage('localStorage')

const libraryStore = useLibraryStore()
const { changeGraphSwitch } = libraryStore

const companyStore = useCompanyStore()
const neo4j_status = computed(() => {
  return companyStore.companyInfo?.neo4j_status == 'true'
})
const ali_ocr_switch = computed(() => {
  return companyStore.companyInfo?.ali_ocr_switch == '1'
})

const props = defineProps({
  type: {
    type: [Number, String],
    default: () => 2
  },
  library_id: {
    type: [Number, String],
    default: () => ''
  }
})

const libraryId = computed(() => {
  return props.library_id || query.id
})

const isHiddenGroup = ref(localStorage.getItem('document_group_hide_key') == 1)

const groupSearchKey = ref('')
const groupLists = ref([])
const groupId = ref(-1)

const currentGroupItem = computed(() => {
  return groupLists.value.find((item) => item.id == groupId.value) || {}
})

const filterGroupLists = computed(() => {
  return groupLists.value.filter((item) => item.group_name.includes(groupSearchKey.value))
})

const userStore = useUserStore()
const sourceTypeMap = {
  1: '本地文档',
  2: '在线文档',
  3: '自定义文档',
  4: '手工新增问答',
  5: '导入问答'
}
const PDF_PARSE_MODE = [
  {
    key: 2,
    title: '图文OCR解析',
    icon: 'pdf-ocr',
    desc: '通过OCR文字识别提取pdf文件内容，可以兼容扫描件，但是解析速度较慢。'
  },
  {
    key: 1,
    title: '纯文本解析',
    icon: 'pdf-text',
    desc: '只提取pdf中的文字内容，如果文档为扫描件可能提取不到内容。'
  },
  {
    key: 3,
    title: '图文解析',
    icon: 'pdf-img',
    desc: '提取PDF文档中的图片和文字'
  },
  {
    key: 4,
    title: '阿里云OCR解析',
    icon: 'ali-ocr',
    desc: '通过阿里云文档智能接口解析提取图片和文字'
  }
]
const rotue = useRoute()
const router = useRouter()
const query = rotue.query

const isLoading = ref(false)
const guideLearningTipsRef = ref(null)
const openGrapgModalRef = ref(null)
const feishuRef = ref(null)
const createGraphSwitch = ref(false)
const libraryInfo = ref({
  library_intro: '',
  library_name: '',
  use_model: '',
  is_offline: null,
  type: 0
})
const state = reactive({
  selectedRowKeys: []
})

const modelForm = reactive({
  use_model: '',
  model_config_id: ''
})

const onChangeModel = (val, option) => {
  let new_use_model = option.modelName
  let new_model_config_id = option.modelId

  if (fileList.value.length > 0) {
    Modal.confirm({
      title: `确定切换模型为${new_use_model}吗？`,
      content: '切换后，所有学习文档将自动重新学习。',
      onOk() {
        modelForm.use_model = new_use_model
        modelForm.model_config_id = new_model_config_id
        saveLibraryConfig()
      },
      onCancel() {}
    })
  } else {
    modelForm.use_model = new_use_model
    modelForm.model_config_id = new_model_config_id
    saveLibraryConfig()
  }
}

const saveLibraryConfig = (showSuccessTip = true, callback = null) => {
  editLibrary({
    ...toRaw(libraryInfo.value),
    use_model: modelForm.use_model,
    model_config_id: modelForm.model_config_id
  }).then(() => {
    libraryInfo.value.use_model = modelForm.use_model
    libraryInfo.value.model_config_id = modelForm.model_config_id

    typeof callback === 'function' && callback()
    if (showSuccessTip) {
      message.success('保存成功')
    }
  })
}

const vectorModelList = ref([])

const onVectorModelLoaded = (list) => {
  vectorModelList.value = list
}

const vectorModelTextList = ref([])

const onVectorTextModelLoaded = (list) => {
  vectorModelTextList.value = list

  setDefaultModel()
}

const setDefaultModel = () => {
  // 防止没有数据时，切换模型报错
  if (!libraryInfo.value.id) {
    setTimeout(() => {
      setDefaultModel()
    }, 100)

    return
  }

  if (vectorModelTextList.value.length > 0 && !libraryInfo.value.use_model) {
    // 遍历查找chatwiki模型
    let modelConfig = null
    let model = null

    for (let item of vectorModelTextList.value) {
      if (item.model_define === 'chatwiki') {
        modelConfig = item
        for (let child of modelConfig.children) {
          if (child.name === 'text-embedding-v2') {
            model = child
            break
          }
        }
        break
      }
    }

    if (!modelConfig) {
      modelConfig = vectorModelTextList.value[0]
      model = modelConfig.children[0]
    }

    if (modelConfig) {
      modelForm.use_model = model.name
      modelForm.model_config_id = model.model_config_id

      saveLibraryConfig(false)
    }
  }
}

const fileList = ref([])

const viewMode = ref('list')
const sortBtnRef = ref(null)

function switchToFolderView() {
  viewMode.value = 'folder'
  localStorage.setItem(getViewModeStorageKey(), 'folder')
  // 重置查询参数与状态
  groupId.value = -1
  inFolderDetail.value = false
  queryParams.page = 1
  queryParams.size = 50
  queryParams.file_name = undefined
  queryParams.status = ''
  queryParams.sort_field = ''
  queryParams.sort_type = ''
  folderSort.order_field = ''
  folderSort.order_type = 'DESC'
  state.selectedRowKeys = []
  sortOpen.value = false
  if (!isHiddenGroup.value) {
    handleChangeHideStatus()
  }
  onSearch()
}
function switchToListView() {
  viewMode.value = 'list'
  localStorage.setItem(getViewModeStorageKey(), 'list')
  queryParams.size = 10
  // 重置查询参数与状态
  groupId.value = -1
  queryParams.page = 1
  queryParams.file_name = undefined
  queryParams.status = ''
  queryParams.sort_field = ''
  queryParams.sort_type = ''
  folderSort.order_field = ''
  folderSort.order_type = 'DESC'
  state.selectedRowKeys = []
  sortOpen.value = false
  if (isHiddenGroup.value) {
    handleChangeHideStatus()
  }
  onSearch()
}
const sortOpen = ref(false)
const folderSort = reactive({
  order_field: '',
  order_type: 'DESC'
})
function setFolderSort(field) {
  if (field === 'default') {
    folderSort.order_field = ''
    folderSort.order_type = 'DESC'
    queryParams.sort_field = ''
    queryParams.sort_type = ''
  } else {
    folderSort.order_field = field
    folderSort.order_type = 'DESC'
    queryParams.sort_field = field
    queryParams.sort_type = 'DESC'
  }
  queryParams.page = 1
  sortOpen.value = false
  getData()
  getGroupLists()
}
function onClickSortIcon() {
  if (!folderSort.order_field) return
  folderSort.order_type = folderSort.order_type === 'DESC' ? 'ASC' : 'DESC'
  queryParams.sort_type = folderSort.order_type
  queryParams.page = 1
  getData()
}
function toggleSortDropdown() {
  sortOpen.value = !sortOpen.value
}
function getStatusClass(item) {
  const s = Number(item?.status ?? -1)
  if ([0, 1, 5, 6, 10].includes(s)) return 'loading'
  if ([4].includes(s)) return 'wait'
  if ([3, 7, 8, 9].includes(s)) return 'error'
  return 'success'
}
const learnProgressTicker = ref(0)
let learnProgressTimer = null
const learnStartTS = {}
function getLearnPercent(origin) {
  try {
    const _ = learnProgressTicker.value
    const a = Number(origin?.ocr_pdf_index ?? 0)
    const b = Number(origin?.ocr_pdf_total ?? 0)
    if (b > 0 && a >= 0) {
      return Math.max(0, Math.min(100, parseInt((a / b) * 100)))
    }
    const id = String(origin?.id || '0')
    if (!learnStartTS[id]) learnStartTS[id] = Date.now()
    const elapsedSec = (Date.now() - learnStartTS[id]) / 1000
    const stepPerSec = 20 + (Number(id) % 6) * 5
    const percent = Math.floor(10 + elapsedSec * stepPerSec)
    return Math.max(10, Math.min(95, percent))
  } catch {}
  return null
}

const groupCards = computed(() => {
  const lists = groupLists.value || []
  const data = fileList.value || []
  const timeField = props.type == 3 ? 'official_article_update_time' : 'update_time'
  const pickLatest = (items) => {
    if (!items.length) return null
    const ordered = items
      .filter(i => !!i[timeField])
      .sort((a, b) => (String(a[timeField]).localeCompare(String(b[timeField]))))
    return ordered.length ? ordered[ordered.length - 1] : items[items.length - 1]
  }
  return lists.map(g => {
    const items = g.id === -1 ? data : data.filter(d => Number(d.group_id || 0) === Number(g.id))
    const latest = pickLatest(items)
    return {
      id: g.id,
      group_name: g.group_name,
      total: g.total,
      last_update: latest ? latest[timeField] || '' : '',
      last_update_ts: latest && latest[timeField] ? dayjs(latest[timeField]).valueOf() : 0,
      cover_status: latest ? getStatusClass(latest) : 'success',
      cover_src: latest ? getFileCover(latest) : ''
    }
  })
})
const visibleGroupCards = computed(() => {
  const arr = [...groupCards.value]
  const { order_field, order_type } = folderSort
  if (order_field === 'update_time' || order_field === 'library_num') {
    arr.sort((a, b) => {
      let av = 0
      let bv = 0
      if (order_field === 'update_time') {
        av = a.last_update_ts || 0
        bv = b.last_update_ts || 0
      } else if (order_field === 'library_num') {
        av = Number(a.total || 0)
        bv = Number(b.total || 0)
      }
      return order_type === 'ASC' ? av - bv : bv - av
    })
  }
  return arr
})

const folderFiles = computed(() => {
  const data = fileList.value || []
  const timeField = props.type == 3 ? 'official_article_update_time' : 'update_time'
  return data.map(rec => {
    const ts = rec[timeField] ? dayjs(rec[timeField]).valueOf() : 0
    return {
      __type: 'file',
      id: rec.id,
      title: rec.file_name || rec.doc_url || '',
      update_ts: ts,
      cover_src: getFileCover(rec) || '',
      cover_status: getStatusClass(rec),
      status_label: getStatusLabel(rec),
      origin: rec
    }
  })
})

const rootUnGroupFiles = computed(() => {
  const data = fileList.value || []
  const timeField = props.type == 3 ? 'official_article_update_time' : 'update_time'
  return data
    .filter(rec => Number(rec.group_id || 0) <= 0)
    .map(rec => {
      const ts = rec[timeField] ? dayjs(rec[timeField]).valueOf() : 0
      return {
        __type: 'file',
        id: rec.id,
        title: rec.file_name || rec.doc_url || '',
        update_ts: ts,
        cover_src: getFileCover(rec) || '',
        cover_status: getStatusClass(rec),
        status_label: getStatusLabel(rec),
        origin: rec
      }
    })
})

const folderItems = computed(() => {
  if (!inFolderDetail.value) {
    return [
      ...visibleGroupCards.value.map(g => ({ __type: 'group', ...g })),
      ...rootUnGroupFiles.value
    ]
  }
  return folderFiles.value
})

function onCardClick(item) {
  if (item.__type === 'group') {
    inFolderDetail.value = true
    handleChangeGroup({ id: item.id, group_name: item.group_name, total: item.total })
  } else {
    handlePreview(item.origin)
  }
}

function formatCardTime(ts) {
  if (!ts) return '--'
  const d = dayjs(ts)
  const now = dayjs()
  if (d.isSame(now, 'day')) return d.format('HH:mm:ss')
  if (now.diff(d, 'year') >= 1) return d.format('YYYY/MM/DD')
  return d.format('MM/DD')
}

function getFileCover(rec) {
  const status = getStatusClass(rec)
  const src =
    rec.thumb_path ||
    rec.cover_url ||
    rec.cover ||
    rec.article_cover_url ||
    ''
  if (src) return src
  if (status === 'wait' || status === 'loading' || status === 'error') return defaultCover
  return ''
}

function getStatusLabel(record) {
  const s = Number(record?.status ?? -1)
  if (s === 0) return '转换中'
  if (s === 1) return '学习中'
  if (s === 2) return '学习完成'
  if (s === 3) return '转换失败'
  if (s === 4) return '待学习'
  if (s === 5) return '待获取'
  if (s === 6) return '获取中'
  if (s === 7) return '获取失败'
  if (s === 8) return '转化异常'
  if (s === 9) return '取消解析'
  if (s === 10) return '正在分段'
  return ''
}

const baseDataOpen = ref(false)
const baseDataRecord = ref(null)
function showBaseData(rec) {
  baseDataRecord.value = rec
  baseDataOpen.value = true
}

function formatUpdateShort(timeStr) {
  if (!timeStr) return '--'
  const d = dayjs(timeStr)
  return d.isValid() ? d.format('YY-MM-DD') : timeStr
}

function formatUpdateFrequency(record) {
  const freq = Number(record?.doc_auto_renew_frequency || 0)
  if (!freq) return '--'
  let first = ''
  if (freq === 1) { return '不自动更新' }
  if (freq === 2) { first = '每天' }
  if (freq === 3) { first = '每3天' }
  if (freq === 4) { first = '每7天' }
  if (freq === 5) { first = '每30天' }
  const minute = Number(record?.doc_auto_renew_minute || 0)
  if (minute > 0) {
    let time = convertTime(minute)
    if (typeof time === 'string' && /^0\d:\d{2}$/.test(time)) {
      time = time.slice(1)
    }
    return `${first}${time}更新`
  }
  return '更新'
}

const count_data = reactive({
  learned_err_count: 0, // 失败数
  learned_count: 0, // 成功数
  learned_wait_count: 0 // 待学习
})

const queryParams = reactive({
  library_id: libraryId.value,
  file_name: undefined,
  status: '', // 2:成功 3:全部失败 8:部分失败 4:待学习
  page: 1,
  size: 10,
  total: 0,
  sort_field: '',
  sort_type: '',
})

const statusList = [
  {
    id: '',
    name: '全部'
  },
  {
    id: '4',
    name: '待学习'
  },
  {
    id: '8',
    name: '部分失败'
  },
  {
    id: '3',
    name: '全部失败'
  },
  {
    id: '2',
    name: '成功'
  }
]

const columns = ref([])
const columnsDefault = [
  {
    title: '文档名称',
    dataIndex: 'file_name',
    key: 'file_name',
    resizable: true,
    minWidth: 180,
    maxWidth: 600,
    width: 300
  },
  {
    title: '大小/分段',
    dataIndex: 'file_size_str',
    key: 'file_size',
    resizable: true,
    minWidth: 120,
    maxWidth: 360,
    width: 140
  },
  {
    title: '合计',
    dataIndex: 'total_hits',
    key: 'total_hits',
    resizable: true,
    minWidth: 80,
    maxWidth: 200,
    width: 80,
    sorter: true
  },
  {
    title: '昨日',
    dataIndex: 'yesterday_hits',
    key: 'yesterday_hits',
    resizable: true,
    minWidth: 80,
    maxWidth: 200,
    width: 80,
    sorter: true
  },
  {
    title: '今日',
    dataIndex: 'today_hits',
    key: 'today_hits',
    resizable: true,
    minWidth: 80,
    maxWidth: 200,
    width: 80,
    sorter: true
  },
  {
    title: '知识图谱',
    dataIndex: 'graph_status',
    key: 'graph_status',
    resizable: true,
    minWidth: 120,
    maxWidth: 400,
    width: 200
  },
  {
    title: '知识图谱实体数',
    dataIndex: 'graph_entity_count',
    key: 'graph_entity_count',
    resizable: true,
    minWidth: 120,
    maxWidth: 320,
    width: 160
  },
  {
    title: '状态',
    dataIndex: 'status',
    key: 'status',
    resizable: true,
    minWidth: 130,
    maxWidth: 200,
    width: 130,
  },
  {
    title: '更新时间/频率',
    dataIndex: 'update_info',
    key: 'update_info',
    resizable: true,
    minWidth: 130,
    maxWidth: 600,
    width: 130
  },
  {
    title: '操作',
    dataIndex: 'action',
    key: 'action',
    fixed: 'right',
    width: 100
  }
]

const getColumnsWidthStorageKey = () => `qa_document_columns_widths_${libraryId.value || 0}`
const getViewModeStorageKey = () => `qa_document_view_mode_${libraryId.value || 0}`
const getSavedColumnsWidths = () => {
  try {
    const s = localStorage.getItem(getColumnsWidthStorageKey()) || '{}'
    const obj = JSON.parse(s)
    return typeof obj === 'object' && obj ? obj : {}
  } catch {
    return {}
  }
}
const applySavedWidths = (cols) => {
  const saved = getSavedColumnsWidths()
  return cols.map(c => {
    const w = saved[c.key]
    if (typeof w === 'number') {
      let width = w
      if (typeof c.minWidth === 'number') {
        width = Math.max(c.minWidth, width)
      }
      if (typeof c.maxWidth === 'number') {
        width = Math.min(c.maxWidth, width)
      }
      return { ...c, width }
    }
    return c
  })
}
const saveColumnWidth = (key, width) => {
  const saved = getSavedColumnsWidths()
  saved[key] = width
  localStorage.setItem(getColumnsWidthStorageKey(), JSON.stringify(saved))
}

const handleChangeStatus = (item) => {
  onSearch()
}

const onTableChange = (pagination, _ , sorter) => {
  console.log(sorter,'==')
  queryParams.page = pagination.current
  queryParams.size = pagination.pageSize
  queryParams.sort_field = ''
  queryParams.sort_type = ''
  if(sorter.order && sorter.field){
    queryParams.sort_field = sorter.field
    queryParams.sort_type = sorter.order == 'ascend' ? 'ASC' : 'DESC'
  }
  getData()
}

const onSearch = () => {
  queryParams.page = 1
  getData(false)
  getGroupLists()
}

let isLast = false

const onDelete = ({ id }) => {
  if (fileList.value.length == 1) {
    isLast = true
  }

  delLibraryFile({ id }).then(() => {
    if (isLast && queryParams.page > 1) {
      queryParams.page--
    }

    getData()
    getGroupLists()
    message.success('删除成功')
  })
}

const handleBatchDelete = () => {
  if (state.selectedRowKeys.length == 0) {
    return message.error('请选择要删除的文档')
  }
  Modal.confirm({
    title: '删除确认',
    icon: createVNode(ExclamationCircleOutlined),
    content: createVNode(
      'div',
      {
        style: 'color: #8c8c8c;'
      },
      '确定要删除选择的文档吗?'
    ),
    onOk() {
      delLibraryFile({ id: state.selectedRowKeys.join(',') }).then(() => {
        getData()
        message.success('批量删除成功')
      })
    },
    onCancel() {}
  })
}

const handlePreview = (record, params = {}) => {
  console.log(rotue)
  let routeName = rotue.name
  if (record.status == '4') {
    if (routeName == 'libraryConfig') {
      window.open(
        '/#/library/document-segmentation?document_id=' + record.id + '&page=' + queryParams.page
      )
      return
    }
    if (params.open_new === true) {
      let queryUrl = router.resolve({
        path: '/library/document-segmentation',
        query: { document_id: record.id, page: queryParams.page }
      })
      window.open(queryUrl.href, '_blank')
      return
    }
    return router.push({
      path: '/library/document-segmentation',
      query: { document_id: record.id, page: queryParams.page }
    })
  }
  if (record.status == '3' && libraryInfo.value.type != 2) {
    return message.error('学习失败,不可预览')
  }
  if (record.status == '0') {
    return message.error('转换中,稍候可预览')
  }
  if (record.status == '1') {
    return message.error('学习中,不可预览')
  }
  if (record.status == '6') {
    return message.error('获取中,不可预览')
  }
  if (record.status == '7') {
    return message.error('获取失败,不可预览')
  }
  if (record.status == '10') {
    return message.error('正在分段,不可预览')
  }

  if (params.open_new === true || routeName == 'libraryConfig') {
    const url = router.resolve({ name: 'libraryPreview', query: { id: record.id, ...params } })
    window.open(url.href, '_blank')
    return
  }
  router.push({ name: 'libraryPreview', query: { id: record.id, ...params } })
}

const getData = (append = false) => {
  let params = toRaw(queryParams)
  if (params.status == 0) {
    params.status = ''
  }
  params.group_id = groupId.value
  isLoading.value = true
  getLibraryFileList(params)
    .then((res) => {
      let info = res.data.info

      if (!modelForm.use_model && modelForm.use_model != info.use_model) {
        modelForm.use_model = info.use_model
        modelForm.model_config_id = info.model_config_id
      }

      libraryInfo.value = { ...info }
      createGraphSwitch.value = info.graph_switch == 1
      if (!append) {
        if (info.graph_switch == '0' || !neo4j_status.value) {
          columns.value = columnsDefault.filter(
            (item) => !['graph_status', 'graph_entity_count'].includes(item.key)
          )
        } else {
          columns.value = columnsDefault
        }
        if (props.type != 3) {
          columns.value = columns.value.filter(i => i.key != 'official_article_update_time')
        }
        columns.value = applySavedWidths(columns.value)
      }

      let list = res.data.list || []
      let countData = res.data.count_data || {}
      if (!append) {
        let metaList = list?.[0]?.meta_list || []
        let metaCols = []
        metaList.forEach(item => {
          if (columns.value.findIndex(i => i.dataIndex == item.key) > -1) return
          metaCols.push({
            title: item.name,
            key: `meta_${item.key}`,
            dataIndex: `meta_${item.key}`,
            width: 160,
          })
        })
        columns.value.splice(columns.value.length - 1, 0, ...metaCols)
      }

      queryParams.total = res.data.total

      count_data.learned_count = countData.learned_count
      count_data.learned_err_count = countData.learned_err_count
      count_data.learned_wait_count = countData.learned_wait_count

      let needRefresh = false
      const mapped = list.map((item) => {
        // , '4' 是待学习，如果加进去会一直刷新状态不会改变
        if (['1', '6', '0', '5'].includes(item.status)) {
          needRefresh = true
        }
        item.file_size_str = formatFileSize(item.file_size)
        item.update_time = dayjs(item.update_time * 1000).format('YYYY-MM-DD HH:mm')
        item.official_article_update_time = dayjs(item.official_article_update_time * 1000).format('YYYY-MM-DD HH:mm')
        item.doc_last_renew_time_desc =
          item.doc_last_renew_time > 0
            ? dayjs(item.doc_last_renew_time * 1000).format('YYYY-MM-DD HH:mm')
            : '--'
        if (Array.isArray(item.meta_list)) {
          item.meta_list.forEach(i => {
            if (i.type == 1 && i.value > 0) {
              i.value = dayjs(i.value * 1000).format('YYYY-MM-DD HH:mm')
            }
            if (i.key == 'source') {
              i.value = sourceTypeMap[i.value]
            }
            item[`meta_${i.key}`] = i.value
          })
        }
        return item
      })
      fileList.value = append ? [...fileList.value, ...mapped] : mapped
      needRefresh && timingRefreshStatus()
      !needRefresh && clearInterval(timingRefreshStatusTimer.value)
    })
    .finally(() => {
      isLoading.value = false
    })
}

const timingRefreshStatusTimer = ref(null)
const timingRefreshStatus = () => {
  clearInterval(timingRefreshStatusTimer.value)
  timingRefreshStatusTimer.value = setInterval(() => {
    getData()
  }, 1000 * 5)
}

const addCustomDocumentRef = ref(null)

const handleMenuClick = (e) => {
  if (vectorModelTextList.value.length == 0) {
    Modal.confirm({
      title: `请先到模型管理中添加嵌入模型？`,
      content:
        '知识库学习需要使用到嵌入模型，请在系统管理-模型管理中添加。推荐使用通义千问、openai或者火山引擎的嵌入模型。',
      okText: '去添加',
      onOk() {
        router.push({ path: '/user/model' })
      }
    })
    return
  }

  let { key } = e

  if (key == 1) {
    handleOpenFileUploadModal()
  }
  if (key == 2) {
    handleOpenUrlModal()
  }
  if (key == 3) {
    addCustomDocumentRef.value.add()
  }
  if (key == 4) {
    feishuRef.value.show()
  }
}
const addUrlState = reactive({
  open: false,
  urls: '',
  library_id: libraryId.value,
  doc_auto_renew_frequency: 1,
  confirmLoading: false,
  doc_auto_renew_minute: '',
  rules: {
    urls: [
      {
        message: '请输入网页地址',
        required: true
      },
      {
        validator: (_rule, value) => {
          if (transformUrlData(value) === false) {
            return Promise.reject(new Error('网页地址不合法'))
          }
          return Promise.resolve()
        }
      }
    ]
  }
})

const handleOpenUrlModal = () => {
  addUrlState.open = true
  addUrlState.confirmLoading = false
  addUrlState.urls = ''
  addUrlState.doc_auto_renew_minute = '02:00'
  addUrlState.doc_auto_renew_frequency = 1
}
const urlFormRef = ref(null)
const handleSaveUrl = () => {
  // 保存本地内容
  urlFormRef.value
    .validate()
    .then(() => {
      addUrlState.confirmLoading = true
      addLibraryFile({
        library_id: addUrlState.library_id,
        urls: JSON.stringify(transformUrlData(addUrlState.urls)),
        doc_auto_renew_frequency: addUrlState.doc_auto_renew_frequency,
        doc_auto_renew_minute: convertTime(addUrlState.doc_auto_renew_minute),
        doc_type: 2,
        group_id: groupId.value,
      }).then(() => {
        addUrlState.open = false
        addUrlState.confirmLoading = false
        onSearch()
      })
    })
    .catch(() => {
      addUrlState.confirmLoading = false
    })
}

const handleCloseUrlModal = () => {
  addUrlState.open = false
}

const addFileState = reactive({
  open: false,
  fileList: [],
  confirmLoading: false,
  pdf_parse_type: 1 //1纯文本解析，2ocr解析
})

const existPdfFile = computed(
  () => addFileState.fileList.filter((i) => i.type === 'application/pdf').length > 0
)

const qaUploadModalRef = ref(null)
const handleOpenFileUploadModal = () => {
  if (libraryInfo.value.type == 2) {
    qaUploadModalRef.value.show()
    return
  }
  addFileState.fileList = []
  addFileState.open = true
}

const handleCloseFileUploadModal = () => {
  addFileState.fileList = []
}

const onFilesChange = (files) => {
  addFileState.fileList = files
}

const handleSaveFiles = () => {
  // 提交后，需要显示引导学习的弹窗提示。如果勾选了“不再显示”，以后批量上传后不再显示弹窗提示。
  if (
    guideLearningTipsRef.value &&
    !userStore.getGuideLearningTips &&
    addFileState.fileList.length > 1
  ) {
    guideLearningTipsRef.value.show()
  }

  if (addFileState.fileList.length == 0) {
    message.error('请选择文件')
    return
  }

  addFileState.confirmLoading = true

  let formData = new FormData()

  formData.append('library_id', queryParams.library_id)
  formData.append('group_id', groupId.value || 0)
  let isTableType = false
  addFileState.fileList.forEach((file) => {
    if (file.name.includes('.xlsx') || file.name.includes('.csv')) {
      isTableType = true
    }
    formData.append('library_files', file)
  })
  if (existPdfFile.value) {
    formData.append('pdf_parse_type', addFileState.pdf_parse_type)
  }
  addLibraryFile(formData)
    .then((res) => {
      getData()
      addFileState.open = false
      addFileState.fileList = []
      addFileState.confirmLoading = false
      if (isTableType && res.data.file_ids.length == 1) {
        router.push('/library/document-segmentation?document_id=' + res.data.file_ids[0])
      }
    })
    .catch(() => {
      addFileState.confirmLoading = false
    })
}

const renameModalRef = ref(null)
const handleOpenRenameModal = (record) => {
  renameModalRef.value.show(record)
}

const toReSegmentationPage = (record) => {
  router.push({
    path: '/library/document-segmentation',
    query: {
      document_id: record.id
    }
  })
}

const createGraphTask = (record) => {
  createGraph({ id: record.id }).then(() => {
    getData()
  })
}

const createGraphSwitchChange = () => {
  if (createGraphSwitch.value) {
    createGraphSwitch.value = false
    let data = {
      graph_model_config_id: libraryInfo.value.graph_model_config_id,
      graph_use_model: libraryInfo.value.graph_use_model
    }
    if (
      (!data.graph_model_config_id || !data.graph_use_model) &&
      vectorModelList.value.length > 0
    ) {
      let modelConfig = vectorModelList.value[0]
      let model = modelConfig.children[0]
      data.graph_use_model = model.name
      data.graph_model_config_id = model.model_config_id
    }
    openGrapgModalRef.value.show(data)
  } else {
    libraryInfo.value.graph_switch = 0
    saveLibraryConfig(false, () => {
      getData()
      changeGraphSwitch(0)
    })
  }
}

const handleOpenGrapgOk = (data) => {
  createGraphSwitch.value = true
  libraryInfo.value.graph_switch = 1
  libraryInfo.value.graph_model_config_id = data.graph_model_config_id
  libraryInfo.value.graph_use_model = data.graph_use_model
  saveLibraryConfig(false, () => {
    getData()
    changeGraphSwitch(1)
  })
}

const pdfParseTypeChange = (item) => {
  if (item.key == 4 && !ali_ocr_switch.value) {
    return false
  }
  addFileState.pdf_parse_type = item.key
}

const onGoSwitch = () => {
  window.open('#/user/aliocr')
}

const editOnlineDocRef = ref(null)
const handleEditOnlineDoc = (record, fromUpdateInfo = false) => {
  editOnlineDocRef.value.show({
    ...record,
    from_update_info: fromUpdateInfo
  })
}
const handleUpdataDoc = (record) => {
  manualCrawl({
    id: record.id
  }).then((res) => {
    message.success('更新成功')
    getData()
  })
}

const handleCancelOcrPdf = (record) => {
  Modal.confirm({
    title: `取消确认？`,
    content: '确认取消该文档解析',
    okText: '确定',
    onOk() {
      cancelOcrPdf({
        id: record.id
      }).then((res) => {
        message.success('取消成功')
        getData()
      })
    }
  })
}

const toGraph = (record) => {
  setStorage('graph:autoOpenFileId', record.id)

  router.push({
    path: '/library/details/knowledge-graph',
    query: {
      id: libraryId.value
    }
  })
}

const handleDownload = (record) => {
  let targetUrl = `/manage/downloadLibraryFile?id=${record.id}&token=${userStore.getToken}`
  var aTag = document.createElement('a')
  aTag.href = targetUrl
  aTag.style.display = 'none'
  aTag.click()
}
const handleChangeHideStatus = () => {
  isHiddenGroup.value = !isHiddenGroup.value
  localStorage.setItem('document_group_hide_key', isHiddenGroup.value ? 1 : 2)
}

const getGroupLists = () => {
  getLibraryGroup({
    library_id: libraryId.value,
    group_type: 1,
    order_field: folderSort.order_field,
    order_type: folderSort.order_type
  }).then((res) => {
    let lists = res.data || []
    let allTotal = lists.reduce((total, item) => {
      return total + +item.total
    }, 0)
    groupLists.value = [
      {
        group_name: '全部分组',
        id: -1,
        total: allTotal
      },
      ...res.data
    ]
  })
}

getGroupLists()

const addGroupRef = ref(null)
const openGroupModal = (data) => {
  addGroupRef.value.show({
    ...data,
    library_id: libraryId.value
  })
}

const editGroupRef = ref(null)
const openEditGroupModal = (data) => {
  editGroupRef.value.show({
    ...data,
    sense: 1
  })
}

const initData = () => {
  getGroupLists()
  getData()
}

const inFolderDetail = ref(false)
const handleChangeGroup = (item) => {
  groupId.value = item.id
  onSearch()
}
const goToRootFolder = () => {
  inFolderDetail.value = false
  groupId.value = -1
  onSearch()
}
const folderScrollSentinel = ref(null)
const folderScrollObserver = ref(null)
function initFolderInfiniteScroll() {
  try {
    folderScrollObserver.value && folderScrollObserver.value.disconnect()
    folderScrollObserver.value = new IntersectionObserver((entries) => {
      const ent = entries[0]
      if (ent && ent.isIntersecting) {
        tryLoadMore()
      }
    }, { root: null, threshold: 0.1 })
    nextTick(() => {
      if (folderScrollSentinel.value) {
        folderScrollObserver.value.observe(folderScrollSentinel.value)
      }
    })
  } catch {}
}
function tryLoadMore() {
  if (viewMode.value !== 'folder') return
  const hasMore = queryParams.page * queryParams.size < (queryParams.total || 0)
  if (!hasMore || isLoading.value) return
  queryParams.page++
  getData(true)
}

const handleDelGroup = (item) => {
  Modal.confirm({
    title: `确认删除分组${item.group_name}`,
    icon: createVNode(ExclamationCircleOutlined),
    content: '',
    okText: '确认',
    okType: 'danger',
    cancelText: '取消',
    onOk() {
      deleteLibraryGroup({
        id: item.id,
        sense: 1
      }).then(() => {
        message.success('删除成功')
        getGroupLists()
        if (groupId.value == item.id) {
          groupId.value = -1
          onSearch()
        }
      })
    }
  })
}

const GROUP_BOX_WIDTH_KEY = 'qa_document_group_box_width'
const minGroupBoxWidth = 180
const maxGroupBoxWidth = 256

const groupBoxWidth = ref(minGroupBoxWidth)
const isDragging = ref(false) // 新增

onMounted(() => {
  const width = parseInt(localStorage.getItem(GROUP_BOX_WIDTH_KEY))
  if (width && width >= minGroupBoxWidth && width <= maxGroupBoxWidth) {
    groupBoxWidth.value = width
  }
  const savedViewMode = localStorage.getItem(getViewModeStorageKey())
  if (savedViewMode === 'folder' || savedViewMode === 'list') {
    viewMode.value = savedViewMode
    queryParams.size = savedViewMode === 'folder' ? 50 : 10
  } else {
    localStorage.setItem(getViewModeStorageKey(), 'list')
  }
  checkFeishuCallback()
  initFolderInfiniteScroll()
  learnProgressTimer = setInterval(() => {
    learnProgressTicker.value++
  }, 1500)
})

const checkFeishuCallback = () => {
  const {
    user_access_token,
    feishu_app_id,
    feishu_app_secret,
    ...other
  } = query
  if (user_access_token && feishu_app_id && feishu_app_secret) {
    router.replace({query: other})
    feishuRef.value.show({
      user_access_token,
      feishu_app_id,
      feishu_app_secret,
    })
  }
}

let isResizing = false
let startX = 0
let startWidth = 0

const handleResizeMouseDown = (e) => {
  if (isHiddenGroup.value) return
  isResizing = true
  isDragging.value = true // 新增
  startX = e.clientX
  startWidth = groupBoxWidth.value
  document.body.style.cursor = 'col-resize'
  document.addEventListener('mousemove', handleResizing)
  document.addEventListener('mouseup', handleResizeMouseUp)
}

const handleResizing = (e) => {
  if (!isResizing) return
  let newWidth = startWidth + (e.clientX - startX)
  newWidth = Math.max(minGroupBoxWidth, Math.min(maxGroupBoxWidth, newWidth))
  groupBoxWidth.value = newWidth
}

const handleResizeMouseUp = () => {
  if (!isResizing) return
  isResizing = false
  isDragging.value = false // 新增
  localStorage.setItem(GROUP_BOX_WIDTH_KEY, groupBoxWidth.value)
  document.body.style.cursor = ''
  document.removeEventListener('mousemove', handleResizing)
  document.removeEventListener('mouseup', handleResizeMouseUp)
}

const checkMove = (e) => {
  // 只允许id>0的项目拖拽
  return e.draggedContext.element.id > 0 && e.relatedContext.element.id > 0
}

const handleDragEnd = async () => {
  try {
    // 过滤掉"全部分组"(-1)
    const sortList = groupLists.value
      .filter((item) => item.id > 0)
      .map((item, index) => ({
        id: item.id,
        sort: groupLists.value.length - index
      }))

    // 调用API保存排序
    await sortLibararyListGroup({
      library_id: libraryId.value,
      sense: 1,
      sort_group: JSON.stringify(sortList)
    })
    message.success('排序已保存')
  } catch (error) {
    console.error('排序保存失败:', error)
    // 恢复原顺序
    getGroupLists()
  }
}

const metaRef = ref(null)

const showMetaModal = (type) => {
  metaRef.value.show(type == 2)
}


onMounted(() => {
  if (query.page) {
    queryParams.page = +query.page
  }
  getData()
})

const handleResizeColumn = (w, col) => {
  let width = w
  if (typeof col.minWidth === 'number') {
    width = Math.max(col.minWidth, width)
  }
  if (typeof col.maxWidth === 'number') {
    width = Math.min(col.maxWidth, width)
  }
  col.width = width
  columns.value = columns.value.map(c => (c.key === col.key ? { ...c, width } : c))
  saveColumnWidth(col.key, width)
}

onUnmounted(() => {
  timingRefreshStatusTimer.value && clearInterval(timingRefreshStatusTimer.value)
  learnProgressTimer && clearInterval(learnProgressTimer)
  folderScrollObserver.value && folderScrollObserver.value.disconnect()
  folderScrollObserver.value && folderScrollObserver.value.disconnect()
})
</script>

<style lang="less" scoped>
.details-library-page {
  height: 100%;
  display: flex;
  overflow: hidden;
}
.group-header-title {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  .title {
    color: #00000073;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
    .title-text {
      color: #262626;
      font-size: 16px;
      font-style: normal;
      font-weight: 600;
      line-height: 24px;
    }
  }
}

.group-content-box {
  position: relative;
  width: 180px;
  border-right: 1px solid #d9d9d9;
  padding-top: 24px;
  padding-left: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition:
    width 0.3s cubic-bezier(0.65, 0, 0.35, 1),
    transform 0.3s cubic-bezier(0.65, 0, 0.35, 1),
    padding 0.3s cubic-bezier(0.65, 0, 0.35, 1),
    border-right-width 0.3s cubic-bezier(0.65, 0, 0.35, 1);
  &.no-transition {
    transition: none !important;
  }
  &.collapsed {
    width: 0 !important;
    min-width: 0 !important;
    max-width: 0 !important;
    padding-left: 0;
    padding-right: 0;
    border-right-width: 0;
    transform: translateX(-100%);
  }

  .group-head-box {
    padding-right: 24px;
  }
  .head-title {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .head-btn-box {
    display: flex;
    gap: 16px;
    margin-top: 16px;
    div {
      flex: 1;
    }
  }
  .search-box {
    margin-top: 16px;
  }
}

.flex-between-box {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.classify-box {
  flex: 1;
  overflow: hidden;
  font-size: 14px;
  .classify-item.sortable-chosen {
    background: #f0f7ff;
  }
  .classify-item.sortable-ghost {
    opacity: 0.5;
    background: #e6f7ff;
  }
  .drag-btn {
    margin-right: 8px;
    cursor: grab;
    opacity: 0;
    transition: opacity 0.2s;
  }
  .classify-item {
    height: 32px;
    padding: 0 8px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 4px;
    cursor: pointer;
    border-radius: 6px;
    color: #595959;
    transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
    user-select: none;

    .classify-title {
      flex: 1;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }
    .num {
      display: block;
    }

    .btn-box {
      display: none;
    }
    &:hover {
      background: #f2f4f7;
      .num {
        display: none;
      }
      .num.num-block {
        display: block;
      }
      .btn-box {
        display: block;
      }
      .drag-btn {
        opacity: 1;
      }
    }
    &.active {
      color: #2475fc;
      background: #e6efff;
    }
  }
}
.empty-folder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 260px;
  border-radius: 8px;
  .empty-inner {
    text-align: center;
    .empty-img {
      width: 88px;
      height: 88px;
      opacity: 0.8;
    }
    .empty-text {
      margin-top: 12px;
      color: #8c8c8c;
    }
  }
}

.scroll-sentinel {
  height: 1px;
}

.hover-btn-wrap {
  width: fit-content;
  height: 24px;
  border-radius: 6px;
  padding: 0 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.645, 0.045, 0.355, 1);
  &:hover {
    background: #e4e6eb;
  }
}

.pover-group-content-box {
  width: 232px;
  padding: 4px;
  overflow: hidden;
  .head-title {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 16px;
    font-weight: 600;
    color: #262626;
  }
  .head-btn-box {
    display: flex;
    gap: 16px;
    margin-top: 16px;
    div {
      flex: 1;
    }
  }
  .search-box {
    margin-top: 16px;
  }
}

.classify-scroll-box {
  max-height: 400px;
  min-height: 180px;
  margin-top: 4px;
  overflow: hidden;
  overflow-y: auto;
  /* 整个页面的滚动条 */
  &::-webkit-scrollbar {
    width: 6px; /* 垂直滚动条宽度 */
    height: 6px; /* 水平滚动条高度 */
  }

  /* 滚动条轨道 */
  &::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 10px;
  }

  /* 滚动条滑块 */
  &::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 10px;
    transition: background 0.3s ease;
  }

  /* 滚动条滑块悬停状态 */
  &::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  /* 滚动条角落 */
  &::-webkit-scrollbar-corner {
    background: #f1f1f1;
  }
}

.resize-bar {
  position: absolute;
  top: 0;
  right: 0;
  width: 4px;
  height: 100%;
  cursor: col-resize;
  z-index: 10;
  background: transparent;
  transition: background 0.2s;
}
.resize-bar:hover {
  background: #2475fc;
}
.group-content-box.collapsed .resize-bar {
  display: none;
}

.main-content-box {
  flex: 1;
  height: 100%;
  overflow: hidden;
  padding: 24px 0 0 24px;
  font-size: 14px;
  line-height: 22px;
  transition: all 0.3s ease;
}

.doc-name-td {
  word-break: break-all;
  display: flex;
  align-items: center;

  .file-icon {
    color: white;
    font-size: 24px;
    margin-right: 12px;
  }
}
.doc-name-text {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
.url-remark {
  color: #8c8c8c;
  margin-top: 2px;
  flex-basis: 100%;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}
.padding-0 {
  padding: 0;
}
.test-menu-icon {
  color: #fff;
}
.library-name {
  height: 38px;
  line-height: 38px;
  padding: 0 16px;
  margin-bottom: 16px;
  font-size: 14px;
  font-weight: 600;
  color: #262626;
  border-radius: 2px;
  background-color: #f2f4f7;
  display: flex;
  align-items: center;
  .anticon-edit {
    margin-left: 8px;
    color: #8c8c8c;
    cursor: pointer;
  }
}
.between-content-box {
  display: flex;
  flex: 1;
  overflow: hidden;
  .left-menu-box {
    width: 232px;
    margin-right: 24px;
  }
  .right-content-box {
    flex: 1;
    overflow: hidden;
  }
}

.menu-item {
  width: 232px;
  height: 56px;
  display: flex;
  align-items: center;
  padding: 0 24px;
  background: #f5f5f5;
  border-radius: 2px;
  margin-bottom: 16px;
  cursor: pointer;
  &.active {
    background: #e6efff;
    border: 1px solid #2475fc;
  }
  .title {
    color: #242933;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
    margin-left: 8px;
  }
}

.list-tools {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-tools {
  margin-bottom: 16px;
  .tools-items {
    display: flex;
    align-items: center;
    .tool-item {
      margin-right: 8px;

      .toggle-btn {
        display: flex;
        padding: 12px;
        justify-content: center;
        align-items: center;
        gap: 8px;
        border-radius: 12px;
        background: #FFF;

        .toggle-icon {
          font-size: 32px;
        }
      }
    }
  }
}

.list-content {
  .text-block {
    color: #595959;
  }
  .c-gray {
    color: #8c8c8c;
  }
  .time-content-box {
    display: flex;
    color: #595959;
    align-items: center;
  }
  .btn-hover-block {
    height: 24px;
    display: flex;
    align-items: center;
    padding: 0 8px;
    cursor: pointer;
    width: fit-content;
    border-radius: 6px;
    transition: all 0.2s cubic-bezier(0.645, 0.045, 0.355, 1);
    &:hover {
      background: #e4e6eb;
    }
  }
  .status-tag {
    display: inline-block;
    height: 24px;
    line-height: 24px;
    padding: 0 6px;
    border-radius: 2px;
    font-size: 14px;
    font-weight: 500;
    text-align: center;
    color: #595959;
    background-color: #edeff2;

    &.running {
      color: #2475fc;
      background-color: #e8effc;
    }

    &.subning {
      color: #6524fc;
      background-color: #eae0ff;
    }

    &.complete {
      color: #21a665;
      background: #e8fcf3;
    }

    &.error {
      cursor: pointer;
      color: #fb363f;
      background-color: #f5c6c8;
    }

    &.warning {
      cursor: pointer;
      background: #faebe6;
      color: #ed744a;
    }

    &.status-learning {
      color: #2475fc;
      // background-color: #e8effc;
    }

    &.status-complete {
      color: #3a4559;
      background-color: #edeff2;
    }

    &.status-error {
      cursor: pointer;
      color: #fb363f;
      // background-color: #f5c6c8;
    }
    &.warning {
      cursor: pointer;
      // background: #faebe6;
      color: #ed744a;
    }
    &.cancel {
      background: #fae4dc;
      color: #ed744a;
    }
  }
  .update-icon {
    width: 16px;
    height: 16px;
    margin-left: 4px;
  }
  .update-trigger-text {
    cursor: pointer;
  }
  .status-icon-row {
    display: inline-flex;
    align-items: center;
    gap: 8px;
  }
  .status-icon {
    color: white;
    font-size: 16px;
  }
}
.view-toggle {
  display: inline-flex;
  align-items: center;
  border: 2px solid #EDEFF2;
  border-radius: 8px;
  overflow: hidden;
  box-sizing: border-box;
  height: 32px;
  background: #EDEFF2;
  .toggle-seg {
    width: 28px;
    height: 28px;
    padding: 2px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border-radius: 6px;
    &:not(:last-child) {
      border-right: 1px solid #EDEFF2;
    }
    .toggle-icon {
      font-size: 16px;
    }
    &.active {
      background: #fff;
      width: 28px;
      height: 28px;

      .toggle-icon {
        color: #2475fc;
      }
    }
  }
}

.sort-btn {
  width: 94px;
  height: 34px;
  padding: 5px 0px 5px 5px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border: 1px solid #EDEFF2;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  .sort-left {
    display: flex;
    padding: 1px 4px;
    align-items: center;
    gap: 8px;
    border-radius: 6px;
    background: #EDEFF2;
  }
  .sort-text {
    color: #595959;
    font-size: 14px;
    font-style: normal;
    font-weight: 400;
    line-height: 22px;
  }
  .sort-icon {
    width: 16px;
    height: 16px;
    transition: color 0.3s ease, transform 0.3s ease;
    &.active {
      color: #2475fc;
    }
  }
  .sort-arrow {
    color: #595959;
    font-size: 12px;
    transition: transform 0.3s ease;
  }
}

.sort-arrow.active {
  transform: rotate(180deg);
}

.sort-dropdown {
  width: 100%;
}
.sort-dropdown .ant-dropdown-menu {
  width: 100%;
  min-width: auto;
}

.pdf-progress-box {
  .progress-title {
    display: flex;
    align-items: center;
    gap: 8px;
    white-space: nowrap;
    .status-box {
      width: fit-content;
      padding: 0 6px;
      display: flex;
      align-items: center;
      gap: 2px;
      height: 22px;
      border-radius: 6px;
      background: #e8effc;
      color: #2475fc;
      font-weight: 500;
    }
  }
  .progress-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    line-height: 20px;
    .ant-progress-line {
      margin: 0;
    }
    .progress-bar-box {
      flex: 1;
    }
    .num-box {
      font-size: 12px;
      color: #8c8c8c;
    }
  }
}
//.upload-file-box {
//  padding: 30px 0;
//}
.ml8 {
  margin-left: 8px;
}
.ml4 {
  margin-left: 4px;
}
.url-add-form {
  margin-top: 24px;
}

.add-dropdown-btn.ant-dropdown {
  .ant-dropdown-menu {
    padding: 0;
    border-radius: 0;
    ::v-deep(.ant-dropdown-menu-item) {
      padding: 12px 16px;
    }
  }
}

.dropdown-btn-menu {
  .title-block {
    color: #262626;
    font-size: 14px;
    font-weight: 600;
    line-height: 22px;
  }
  .desc {
    color: #8c8c8c;
    font-size: 14px;
    line-height: 22px;
  }
}
.table-btn {
  cursor: pointer;
  &:hover {
    color: #2475fc;
  }
}
.custom-select-box {
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  > span {
    white-space: nowrap;
  }
  :deep(.ant-select-selector) {
    border: none !important;
    padding-left: 0 !important;
    height: unset !important;
  }
  padding-left: 8px;
}
.pd-5-8 {
  padding: 5px 8px;
}
.reason-text {
  color: #8c8c8c;
  font-size: 12px;
  line-height: 24px;
}

.title-box {
  display: flex;
  align-items: center;
  gap: 4px;
}

.select-card-box {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
  margin-top: 8px;
  .select-card-item {
    min-height: 105px;
    width: calc(50% - 8px);
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
      gap: 4px;
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
      line-height: 22px;
      font-size: 14px;
      color: #595959;
    }

    .card-switch {
      display: flex;
      gap: 4px;
      color: #8c8c8c;
      font-size: 14px;
      line-height: 22px;

      .card-switch-btn {
        cursor: pointer;
        color: #2475fc;

        &:hover {
          opacity: 0.8;
        }
      }
    }

    .card-switch-box {
      width: 52px;
      height: 22px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 10px;
      border-radius: 6px;
      border: 1px solid #edb8a6;
      background: #ffece5;
      color: #ed744a;
      font-size: 12px;
      font-style: normal;
      font-weight: 400;
      line-height: 20px;
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
.mt24 {
  margin-top: 24px;
}

.right-box {
  display: flex;
  gap: 16px;
}

.status-box {
  display: flex;
  align-items: center;

  .status-item {
    display: flex;
    align-items: center;

    .status-label {
      color: #595959;
    }

    .content-tip {
      color: red;
    }
  }
}

.select-card-main {
  :deep(.ant-row) {
    display: block;
  }
}

.subning {
  :deep(.ant-spin-dot-item) {
    background-color: #6524fc !important;
  }
}
.folder-view {
  padding-bottom: 24px;
  .folder-grid {
    display: grid;
    grid-template-columns: repeat(9, 1fr);
    gap: 24px;
  }
  .folder-card {
    position: relative;
    border: 1px solid #f0f0f0;
    border-radius: 12px;
    overflow: hidden;
    background: #fff;
    cursor: pointer;
    transition: box-shadow 0.2s ease, transform 0.2s ease;
    width: 100%;
    height: 182px;
    &:hover {
      box-shadow: 0 6px 16px rgba(0,0,0,0.08);
      transform: translateY(-2px);
    }
    .card-more {
      position: absolute;
      right: 7px;
      top: 8px;
      color: #fff;
      display: none;
      padding: 4px;
      justify-content: center;
      align-items: center;
      gap: 4px;
      border-radius: 6px;
      background: #00000040;
    }
    &:hover .card-more { display: inline-flex; }
  }
  .folder-thumb {
    position: relative;
    width: calc(100% - 16px);
    margin: 8px;
    height: 86px;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fff;

    .folder-type-icon {
      font-size: 68px;
    }
    .thumb-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    .type-tag {
      position: absolute;
      left: 8px;
      top: 8px;
      padding: 2px 6px;
      border-radius: 4px;
      font-size: 12px;
      line-height: 20px;
      color: #fff;
      background: rgba(36, 117, 252, 0.9);
    }
    .status-badge {
      position: absolute;
      left: 8px;
      top: 8px;
      padding: 2px 6px;
      border-radius: 6px;
      font-size: 12px;
      line-height: 20px;
      color: #fff;
      background: rgba(0,0,0,0.4);
    }
    .learn-percent {
      margin-right: 4px;
    }
    &.status-loading {
      background: linear-gradient(180deg, #f0f7ff 0%, #ffffff 100%);
    }
  }
  .card-file {
    .folder-thumb {
      border: 1px solid #F0F0F0;
    }
  }
  .folder-meta {
    padding: 0 8px 8px;
    .name {
      width: 100%;
      height: 44px;
      color: #242933;
      text-align: center;
      font-size: 14px;
      font-style: normal;
      font-weight: 400;
      line-height: 22px;
      margin-bottom: 8px;
    }
    .name-2line {
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .desc {
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 12px;
      color: #8c8c8c;
      line-height: 20px;
      .file-icon {
        font-size: 14px;
        vertical-align: -2px;
        margin-right: 4px;
      }
    }
  }
}
@media screen and (max-width: 1900px) {
  .folder-view {
    .folder-grid {
      grid-template-columns: repeat(7, 1fr);
    }
  }
}
</style>

<style lang="less">
.ant-popover-inner-content {
  .update-tooltip {
    .update-tooltip-title {
      display  : flex;
      align-items: center;
      margin-bottom: 4px;
    }
    .update-icon-box {
      display: flex;
      padding: 4px;
      justify-content: center;
      align-items: center;
      gap: 4px;
      border-radius: 6px;
      background: #E4E6EB;
      margin-left: 8px;
    }
    .edit-icon {
      width: 14px;
      height: 14px;
      cursor: pointer;
    }
  }
  .update-tooltip-value {
    color: #595959;
  }
  .update-trigger-text {
    cursor: pointer;
  }
}
.main-content-box .group-header-title .breadcrumb-a:hover {
  background-color: transparent !important;
}
</style>
