import { fetchEventSource, EventStreamContentType } from './fetch';
const isBrowser = typeof window !== 'undefined' && window.document;
if (isBrowser) {
    window.FetchEventSource = {
        fetchEventSource,
        EventStreamContentType,
    };
}
export { fetchEventSource, EventStreamContentType };
//# sourceMappingURL=browser.js.map