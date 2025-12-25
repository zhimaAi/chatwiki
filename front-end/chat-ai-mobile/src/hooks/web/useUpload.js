import { ref, onUnmounted } from 'vue'
import { showToast } from 'vant';
import { uploadFile } from '@/api/app'
import { getUuid } from '@/utils/index'

// 真实上传任务
const uploadFileFetch = ({file, fileList, category, extraData}) => {
  return new Promise((resolve, reject) => {
    uploadFile({
      file: file.originFile,
      category: category,
      extraData: extraData
    }, (progressEvent) => {
      const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
      
      // 更新文件列表中的进度
      const currentFile = fileList.value.find(f => f.uid === file.uid);
      if (currentFile) {
        const updatedFile = { ...currentFile, percent: percentCompleted };
        const index = fileList.value.findIndex(f => f.uid === file.uid);
        if (index !== -1) {
          fileList.value.splice(index, 1, updatedFile);
        }
      }
    }).then((res) => {
      // 上传成功
      const currentFile = fileList.value.find(f => f.uid === file.uid);
      if (currentFile) {
        const updatedFile = { 
          ...currentFile, 
          percent: 100, 
          status: 'done', 
          url: res.data.link || '' 
        };
        const index = fileList.value.findIndex(f => f.uid === file.uid);
        if (index !== -1) {
          fileList.value.splice(index, 1, updatedFile);
        }
        resolve(updatedFile);
      }
    }).catch((error) => {
      // 上传失败
      const currentFile = fileList.value.find(f => f.uid === file.uid);
      if (currentFile) {
        const updatedFile = { ...currentFile, status: 'error' };
        const index = fileList.value.findIndex(f => f.uid === file.uid);
        if (index !== -1) {
          fileList.value.splice(index, 1, updatedFile);
        }
      }
      reject(error);
    });
  });
}

export function useUpload(options) {
  const {
    limit = 10,
    maxSize = 10, // MB
    accept = 'image/*',
    multiple = true,
    category,
    fileList = ref([]),
    extraData = {}
  } = options || {}

  if(!category){
    throw new Error('category is required')
  }

  let inputEl = null;

  const handleFileChange = (event) => {
    let files = Array.from(event.target.files)
    const remainingSlots = limit - fileList.value.length

    if (remainingSlots <= 0) {
      showToast(`最多只能上传 ${limit} 张图片`)
      return
    }

    if (files.length > remainingSlots) {
      files = files.slice(0, remainingSlots)
      showToast(`最多只能上传 ${limit} 张图片，已为您选择前 ${remainingSlots} 张`)
    }

    files.forEach((file) => {
      const isLtMaxSize = file.size / 1024 / 1024 < maxSize;
      if (!isLtMaxSize) {
        showToast(`文件大小不能超过 ${maxSize}MB!`);
        return;
      }
      // 检查文件类型
      const fileType = file.type;
      const acceptedTypes = accept.split(',').map(t => t.trim());
      const isTypeValid = acceptedTypes.some(acceptedType => {
        if (acceptedType.endsWith('/*')) {
          return fileType.startsWith(acceptedType.slice(0, -1));
        }
        return fileType === acceptedType;
      });

      if (!isTypeValid) {
        showToast(`不支持的文件类型: ${fileType}`);
        return;
      }

      const fileItem = {
        uid: file.uid || getUuid(16),
        name: file.name,
        status: 'uploading',
        percent: 0,
        originFile: file,
        url: URL.createObjectURL(file) // 生成临时的 URL
      }

      fileList.value.push(fileItem)

      uploadFileFetch({file: fileItem, fileList, category, extraData: extraData})
        .then((res) => {
          const index = fileList.value.findIndex((item) => item.uid === res.uid)
          if (index !== -1) {
            fileList.value.splice(index, 1, res)
          }
        })
        .catch(() => {
          const index = fileList.value.findIndex((item) => item.uid === fileItem.uid)
          if (index !== -1) {
            fileList.value[index].status = 'error'
          }
        })
    })

    // 清空 input 的 value，以便可以再次选择相同的文件
    event.target.value = ''
  }

  const createInputElement = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = accept;
    input.multiple = multiple;
    input.style.display = 'none';
    input.onchange = handleFileChange;
    document.body.appendChild(input);
    return input;
  }

  const openFileDialog = () => {
    if (fileList.value.length >= limit) {
      showToast(`最多只能上传 ${limit} 张图片`);
      return;
    }
    if (!inputEl) {
      inputEl = createInputElement();
    }
    inputEl.click();
  }

  onUnmounted(() => {
    if (inputEl && document.body.contains(inputEl)) {
      document.body.removeChild(inputEl);
      inputEl = null;
    }
  });

  return {
    openFileDialog,
  }
}