import { fetchEventSource, FetchEventSourceInit, EventStreamContentType } from './fetch';
import { EventSourceMessage } from './parse';

// 导出到全局对象
declare global {
    interface Window {
        FetchEventSource: {
            fetchEventSource: typeof fetchEventSource;
            // 移除 EventSourceMessage，因为它只是类型定义而非值
            EventStreamContentType: typeof EventStreamContentType;
            // 移除 FetchEventSourceInit，因为它只是类型定义而非值
        }
    }
}

// 确保在浏览器环境中运行
const isBrowser = typeof window !== 'undefined' && window.document;

if (isBrowser) {
    window.FetchEventSource = {
        fetchEventSource,
        EventStreamContentType,
    } as const;  // 使用 as const 确保类型不可变
}

export {
    fetchEventSource,
    FetchEventSourceInit,
    EventStreamContentType,
    EventSourceMessage
};