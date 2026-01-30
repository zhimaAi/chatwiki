import { ref, onMounted, onBeforeUnmount } from 'vue';
 
export function useWindowWidth(debounceDelay = 250) {
    const windowWidth = ref(window.innerWidth);
 
    const updateWidth = () => {
        windowWidth.value = window.innerWidth;
    };
 
    // 防抖函数
    const debounce = (func, delay) => {
        let timeoutId;
        return function(...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    };
 
    const debouncedUpdate = debounce(updateWidth, debounceDelay);
 
    onMounted(() => {
        window.addEventListener('resize', debouncedUpdate);
    });
 
    onBeforeUnmount(() => {
        window.removeEventListener('resize', debouncedUpdate);
    });
 
    return { windowWidth };
}