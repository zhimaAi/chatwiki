import { fetchEventSource, FetchEventSourceInit, EventStreamContentType } from './fetch';
import { EventSourceMessage } from './parse';
declare global {
    interface Window {
        FetchEventSource: {
            fetchEventSource: typeof fetchEventSource;
            EventStreamContentType: typeof EventStreamContentType;
        };
    }
}
export { fetchEventSource, FetchEventSourceInit, EventStreamContentType, EventSourceMessage };
