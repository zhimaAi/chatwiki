import dayjs from 'dayjs';
 
const timeFormat = (el, bindings) => {
  // 1.获取时间，并且转化成毫秒
  let timestamp = el.textContent
  if (timestamp.length && timestamp.indexOf(':') > -1) {
    return
  }
  if(timestamp.length === 10){
    timestamp = timestamp * 1000
  }

  timestamp = Number(timestamp)

  // 2.获取传入的参数
  let value = bindings.value;
  if(!value){
    value = "YYYY-MM-DD HH:mm:ss"
  }
 
  const date = dayjs(timestamp);
  const today = dayjs().startOf('day');
  const isToday = date.isSame(today, 'day');
  
  const format = isToday ? 'HH:mm' : value;
  el.textContent = date.format(format);
};
 
export const timeFormatter = {
  beforeMount(el, binding) {
    timeFormat(el, binding);
  },
  updated(el, binding) {
    timeFormat(el, binding);
  }
};
