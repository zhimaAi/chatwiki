"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EventStreamContentType = exports.fetchEventSource = void 0;
const fetch_1 = require("./fetch");
Object.defineProperty(exports, "fetchEventSource", { enumerable: true, get: function () { return fetch_1.fetchEventSource; } });
Object.defineProperty(exports, "EventStreamContentType", { enumerable: true, get: function () { return fetch_1.EventStreamContentType; } });
const isBrowser = typeof window !== 'undefined' && window.document;
if (isBrowser) {
    window.FetchEventSource = {
        fetchEventSource: fetch_1.fetchEventSource,
        EventStreamContentType: fetch_1.EventStreamContentType,
    };
}
//# sourceMappingURL=browser.js.map