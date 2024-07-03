export function getSdkCode({ pc_domain, robot_key }) {
  var text = `<script>
  (function (z, h, i, m, a, j, s) {
    z[m] = z[m] || function () {
      (z[m].a = z[m].a || []).push(arguments);
    };
    (j = h.createElement(i)), (s = h.getElementsByTagName(i)[0]);
    j.async = true;
    j.charset = "UTF-8";
    j.setAttribute(
      "data-json",
      JSON.stringify({
        robot_key: "${robot_key}",
        language: "zh-CN",
      })
    );
    j.id = "ai_chat_js";
    j.src = '${pc_domain}/sdk/ai-chat-sdk.umd.cjs';
    s.parentNode.insertBefore(j, s);
  })(window, document, "script", "_ai_chat");
</script>`

  return text
}
