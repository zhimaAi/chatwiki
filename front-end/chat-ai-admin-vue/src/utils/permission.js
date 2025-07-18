import { usePermissionStore } from '@/stores/modules/permission'

/**
 * 字符权限校验
 * @param {Array} value 校验值
 * @returns {Boolean}
 */
export function checkPermi(value) {
  const permissionStore = usePermissionStore()
  const { permissionMap } = permissionStore

  if (value && value instanceof Array && value.length > 0) {
    const permissionDatas = value
    const all_permission = '*:*:*'
    const hasPermission = permissionDatas.some((permissionKey) => {
      return permissionKey == all_permission || permissionMap[permissionKey.replace(/\//g, '_')]
    })

    if (!hasPermission) {
      return false
    }

    return true
  } else {
    console.error(`need roles! Like checkPermi="['system:user:add','system:user:edit']"`)
    return false
  }
}

/**
 * 角色权限校验
 * @param {Array} value 校验值
 * @returns {Boolean}
 */
export function checkRole(value) {
    const permissionStore = usePermissionStore()
    const { role_permission, role_type } = permissionStore
    if(role_type == 1) {
      return true
    }
  if (value && value instanceof Array && value.length > 0) {
    const permissionRoles = value
    const super_admin = 'admin'
    const all_permission = '*:*:*'

    const hasRole = permissionRoles.some((role) => {
      return super_admin === role || role_permission.includes(role) || role == all_permission
    })

    if (!hasRole) {
      return false
    }
    return true
  } else {
    console.error(`need roles! Like checkRole="['admin','editor']"`)
    return false
  }
}

export function checkSystemPermisission(to) {
  const { role_permission } = usePermissionStore()
  if(to.name == 'userModel' && !role_permission.includes('ModelManage')){
    return 'AccountManage'
  }
}

export function checkRouterPermisission(value) {
  const { user_roles, role_permission } = usePermissionStore()
  if (user_roles == '1') {
    return true
  }
  if (value.includes('/robot/')) {
    // 机器人管理页面
    return role_permission.includes('RobotManage')
  }
  if (value.includes('/library/')) {
    // 知识库管理
    return role_permission.includes('LibraryManage')
  }
  if (value.includes('/user/clientDownload')) {
    // 客户端管理
    return role_permission.includes('ClientSideManage')
  }
  if (value.includes('/user/') && value != '/user/account') {
    // 知识库管理
    return role_permission.includes('SystemManage')
  }

  return true
}

// object_type 1:机器人 2:知识库 3:数据库

// operate_rights  4:可管理  2:可编辑 1:查看
export function getRobotPermission(id) {
  const {  permission_manage_data } = usePermissionStore()
  // 有权限的机器人
  let list = permission_manage_data.filter(item => item.object_type == 1)
  let allPermiston = list.filter(item => item.object_id == -1)
  if(allPermiston.length >  0){
    // 所有权限
    return 4
  }
  let currentItem = list.filter(item => item.object_id == id)
  if(currentItem.length == 0){
    // 没有任何权限
    return 0
  }
  return  currentItem[0].operate_rights
}


export function getLibraryPermission(id) {
  const {  permission_manage_data } = usePermissionStore()
  // 有权限的知识库
  let list = permission_manage_data.filter(item => item.object_type == 2)
  let allPermiston = list.filter(item => item.object_id == -1)
  if(allPermiston.length >  0){
    // 所有权限
    return 4
  }
  let currentItem = list.filter(item => item.object_id == id)
  if(currentItem.length == 0){
    // 没有任何权限
    return 0
  }
  return  currentItem[0].operate_rights
}

export function getDatabasePermission(id) {
  const {  permission_manage_data } = usePermissionStore()
  // 有权限的数据库
  let list = permission_manage_data.filter(item => item.object_type == 3)
  let allPermiston = list.filter(item => item.object_id == -1)
  if(allPermiston.length >  0){
    // 所有权限
    return 4
  }
  let currentItem = list.filter(item => item.object_id == id)
  if(currentItem.length == 0){
    // 没有任何权限
    return 0
  }
  return  currentItem[0].operate_rights
}