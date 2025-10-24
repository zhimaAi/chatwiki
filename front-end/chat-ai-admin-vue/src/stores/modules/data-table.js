import { defineStore } from 'pinia'
import { getFormList, getFormFieldList } from '@/api/database/index'

export const useDataTableStore = defineStore('dataTable', {
  state: () => {
    return {
      // 表单列表缓存
      formList: [],
      // 表单字段列表缓存映射，key为form_id
      formFieldListMap: {},
      // 表单列表加载状态
      formListLoading: false,
      // 表单字段列表加载状态映射，key为form_id
      formFieldListLoadingMap: {},
      // 表单列表请求队列，存储等待结果的Promise resolve函数
      formListQueue: [],
      // 表单字段列表请求队列映射，key为form_id，value为Promise resolve函数数组
      formFieldListQueueMap: {},
    }
  },
  getters: {},
  actions: {
    /**
     * 获取表单列表（带队列机制）
     * @returns {Promise<Array>} 表单列表
     */
    async getFormList(){
      // 如果已有缓存数据，直接返回
      if(this.formList.length > 0){
        return [...this.formList]
      }
      
      // 如果正在加载中，将当前请求加入队列等待
      if(this.formListLoading){
        return new Promise((resolve) => {
          this.formListQueue.push(resolve)
        })
      }
      
      // 设置加载状态为true，防止重复请求
      this.formListLoading = true
      
      try{
        const res = await getFormList()
        this.formList = res.data || []
        
        // 请求成功后，处理队列中等待的所有请求
        const result = [...this.formList]
        this.formListQueue.forEach(resolve => resolve(result))
        this.formListQueue = []
        
        return result
      }catch(error){
        console.log(error)
        this.formList = []
        
        // 请求失败时，也要处理队列中等待的请求
        const result = [...this.formList]
        this.formListQueue.forEach(resolve => resolve(result))
        this.formListQueue = []
        
        return result
      }finally{
        // 无论成功失败，都要重置加载状态
        this.formListLoading = false
      }
    },
    /**
     * 获取表单字段列表（带队列机制）
     * @param {Object} params - 参数对象
     * @param {string} params.form_id - 表单ID
     * @returns {Promise<Array>} 表单字段列表
     */
    async getFormFieldList({form_id}){
      // 如果已有缓存数据，直接返回
      if(this.formFieldListMap[form_id] && this.formFieldListMap[form_id].length > 0){
        return [...this.formFieldListMap[form_id]]
      }
      
      // 如果该form_id正在加载中，将当前请求加入对应的队列等待
      if(this.formFieldListLoadingMap[form_id]){
        return new Promise((resolve) => {
          // 初始化该form_id的队列（如果不存在）
          if(!this.formFieldListQueueMap[form_id]){
            this.formFieldListQueueMap[form_id] = []
          }
          this.formFieldListQueueMap[form_id].push(resolve)
        })
      }
      
      // 设置该form_id的加载状态为true，防止重复请求
      this.formFieldListLoadingMap[form_id] = true
      
      try{
        const res = await getFormFieldList({form_id})
        this.formFieldListMap[form_id] = res.data || []
        
        // 请求成功后，处理该form_id队列中等待的所有请求
        const result = [...this.formFieldListMap[form_id]]
        if(this.formFieldListQueueMap[form_id]){
          this.formFieldListQueueMap[form_id].forEach(resolve => resolve(result))
          this.formFieldListQueueMap[form_id] = []
        }
        
        return result
      }catch(error){
        console.log(error)
        this.formFieldListMap[form_id] = []
        
        // 请求失败时，也要处理队列中等待的请求
        const result = [...this.formFieldListMap[form_id]]
        if(this.formFieldListQueueMap[form_id]){
          this.formFieldListQueueMap[form_id].forEach(resolve => resolve(result))
          this.formFieldListQueueMap[form_id] = []
        }
        
        return result
      }finally{
        // 无论成功失败，都要重置该form_id的加载状态
        this.formFieldListLoadingMap[form_id] = false
      }
    },
    /**
     * 刷新表单列表
     * 清空缓存并重新获取数据
     * @returns {Promise<Array>} 最新的表单列表
     */
    async refreshFormList(){
      // 清空缓存数据
      this.formList = []
      // 重置加载状态和队列
      this.formListLoading = false
      this.formListQueue = []

      let list = await this.getFormList()

      return list
    },
    
    /**
     * 刷新指定表单的字段列表
     * 清空缓存并重新获取数据
     * @param {Object} params - 参数对象
     * @param {string} params.form_id - 表单ID
     * @returns {Promise<Array>} 最新的表单字段列表
     */
    async refreshFormFieldList({form_id}){
      // 清空该form_id的缓存数据
      this.formFieldListMap[form_id] = []
      // 重置该form_id的加载状态和队列
      this.formFieldListLoadingMap[form_id] = false
      this.formFieldListQueueMap[form_id] = []
      
      let list = await this.getFormFieldList({form_id})

      return list
    }
  },
})