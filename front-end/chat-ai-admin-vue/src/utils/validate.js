export function validatePassword(password) {
  if (password.length < 6 || password.length > 32) {
    return false; // 密码长度不在6-32位之间  
  }

  // 定义用于检查不同类型字符的正则表达式  
  const hasLetter = /[a-zA-Z]/.test(password);
  const hasNumber = /[0-9]/.test(password);
  const hasSpecialChar = /[!@#$%^&*()\-_=+{};:,<.>]/.test(password); // 这里只列出了一些常见的特殊字符，你可以根据需要添加更多  

  // 确保密码包含至少两种不同类型的字符  
  const hasTwoTypes = (
    (hasLetter && hasNumber) ||
    (hasLetter && hasSpecialChar) ||
    (hasNumber && hasSpecialChar)
  );

  return hasTwoTypes;
}


export function isValidURL(url) {
  try {
    new URL(url);
    return true;
  } catch (err) {
    return false;
  }
}

export function transformUrlData(str) {
  const isValidUrl = (string) => {
    try {
      new URL(string);
      return true;
    } catch (_) {
      return false;
    }
  }
  str = str.trim()
  var lines = str.split(/\r?\n|\r/).map(line => line.trim()).filter(line => line.length > 0)
  // console.log(lines)
  let result = []
  let item = {}
  for (let i = 0; i < lines.length; i++) {
    let line = lines[i]
    if (isValidUrl(line)) {
      item.url = line
      result.push(item)
      item = {}
    } else {
      if (!item.remark) {
        item.remark = line
      } else {
        if (!item.url) {
          return false
        }
      }
    }
  }
  return result
}


export function isNumberOrNumberString(value) {
  // 检查是否为number类型  
  if (typeof value === 'number') {
    // 如果不是NaN或Infinity（可选，根据你的需求）  
    return Number.isFinite(value);
  }
  // 检查是否为可以转换为有限数字的字符串  
  if (typeof value === 'string') {
    // 尝试将字符串转换为数字  
    const num = Number(value);
    // 检查转换后的值是否为有限数字且原始字符串不是空字符串（可选）  
    return !Number.isNaN(num) && Number.isFinite(num) && value !== '';
  }
  // 如果不是number类型也不是字符串类型，则不是数字或数字字符串  
  return false;
}  