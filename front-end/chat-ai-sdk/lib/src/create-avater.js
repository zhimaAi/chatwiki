function loadImage(src) {
  return new Promise((resolve, reject) => {
    let img = new Image();
    img.src = src;
    img.onload = () => resolve(img);
    img.onerror = reject;
  });
}

// type = 1 （简洁）大小固定为 50*50
async function createAvatarByType1(config) {
  let img = document.createElement("img");

  img.src = config.buttonIcon;

  img.style.width = 50 + "px";
  img.style.height = 50 + "px";

  return img;
}
// type = 2 （加文案）
async function createAvatarByType2(config) {
  let el = document.createElement("div");
  el.className = "chat-wiki-avatar_type2";

  let img = document.createElement("img");

  img.src = config.buttonIcon;
  img.className = "chat-wiki-avatar_type2_icon";
  el.appendChild(img);

  let text = document.createElement("div");
  text.className = "chat-wiki-avatar_type2_text";
  text.innerText = config.buttonText;

  el.appendChild(text);

  return el;
}
// type = 3 （自定义）大小等于图片的大小
async function createAvatarByType3(config) {
  try {
    // 等待图片加载完成
    const img = await loadImage(config.buttonIcon);
    img.style.width = img.width + "px";
    img.style.height = img.height + "px";
    

    return img;
  } catch (e) {
    console.log(e);
  }
}

async function createAvatar(config) {
  if (config.displayType === 3) {
    return createAvatarByType3(config);
  } else if (config.displayType === 2) {
    return createAvatarByType2(config);
  } else {
    return createAvatarByType1(config);
  }
}

export default createAvatar;
