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