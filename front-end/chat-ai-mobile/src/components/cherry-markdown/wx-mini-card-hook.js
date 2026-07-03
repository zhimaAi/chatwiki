import CherryEngine from 'cherry-markdown/dist/cherry-markdown.engine.core'

const wxMiniCardHook = CherryEngine.createSyntaxHook(
  'wxMiniCard',
  CherryEngine.constants.HOOKS_TYPE_LIST.PAR,
  {
    needCache: true,
    beforeMakeHtml(str) {
      return str.replace(this.RULE.reg, (whole, json) => {
        try {
          const card = JSON.parse(json)
          const sign = this.$engine.md5(whole)
          const lines = this.getLineCount(whole)
          const enc = (v) => encodeURIComponent(v || '')
          const esc = (v) =>
            String(v || '')
              .replace(/&/g, '&amp;')
              .replace(/</g, '&lt;')
              .replace(/>/g, '&gt;')
              .replace(/"/g, '&quot;')
          const iconSvg =
            '<svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg"><path d="M10.5004 8C10.9948 8 11.4782 7.8534 11.8893 7.5787C12.3004 7.30395 12.6209 6.9135 12.8101 6.4567C12.9993 5.9999 13.0488 5.49725 12.9523 5.0123C12.8559 4.5273 12.6178 4.08185 12.2681 3.73225C11.9185 3.3826 11.4731 3.1445 10.9881 3.04805C10.5031 2.95155 10.0005 3.0011 9.54366 3.1903C9.08686 3.3795 8.69642 3.69995 8.42172 4.1111C8.14697 4.5222 8.00037 5.00555 8.00037 5.5V10.5C8.00037 10.9944 7.85377 11.4778 7.57907 11.8889C7.30432 12.3 6.91392 12.6205 6.45707 12.8097C6.00027 12.9989 5.49762 13.0484 5.01267 12.9519C4.52767 12.8555 4.08222 12.6174 3.73262 12.2677C3.38297 11.9181 3.14487 11.4727 3.04842 10.9877C2.95192 10.5027 3.00147 10.0001 3.19067 9.5433C3.37987 9.0865 3.70032 8.69605 4.11147 8.42135C4.52257 8.1466 5.00592 8 5.50037 8" stroke="#BFBFBF" stroke-linecap="round" stroke-linejoin="round"/></svg>'
          const coverHtml = card.thumb_url
            ? `<div class="wx-mini-card-cover"><img src="${esc(card.thumb_url)}" alt="" /></div>`
            : ''
          const result = `<div class="wx-mini-card" data-title="${enc(card.title)}" data-appid="${enc(card.appid)}" data-page-path="${enc(card.page_path)}" data-thumb-url="${enc(card.thumb_url)}" data-sign="${sign}" data-lines="${lines}"><div class="wx-mini-card-title">${esc(card.title)}</div>${coverHtml}<div class="wx-mini-card-footer">${iconSvg}<span>小程序</span></div></div>`
          return this.pushCache(result, sign, lines)
        } catch (e) {
          return whole
        }
      })
    },
    makeHtml(str) {
      return str
    },
    rule() {
      return { reg: /\[wx_mini_card\](\{[\s\S]*?\})\[\/wx_mini_card\]/g }
    }
  }
)

export default wxMiniCardHook
