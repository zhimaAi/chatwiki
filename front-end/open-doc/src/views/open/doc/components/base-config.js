const basicConfig = {
  externals: {
    echarts: false,
    katex: false,
    MathJax: true,
  },
  isPreviewOnly: true,
  engine: {
    global: {
      // urlProcessor(url, srcType) {
      //   console.log(`url-processor`, url, srcType)
      //   return url
      // }
    },
    syntax: {
      image: {
        videoWrapper: (link, type, defaultWrapper) => {
          return defaultWrapper
        },
      },
      link: {
        /** 生成的<a>标签追加target属性的默认值 空：在<a>标签里不会追加target属性， _blank：在<a>标签里追加target="_blank"属性 */
        target: '_blank',
        /** 生成的<a>标签追加rel属性的默认值 空：在<a>标签里不会追加rel属性， nofollow：在<a>标签里追加rel="nofollow：在"属性*/
        rel: '',
      },
      autoLink: {
        /** 生成的<a>标签追加target属性的默认值 空：在<a>标签里不会追加target属性， _blank：在<a>标签里追加target="_blank"属性 */
        target: '_blank',
        /** 生成的<a>标签追加rel属性的默认值 空：在<a>标签里不会追加rel属性， nofollow：在<a>标签里追加rel="nofollow：在"属性*/
        rel: '',
        /** 是否开启短链接 */
        enableShortLink: false,
        /** 短链接长度 */
        shortLinkLength: 20,
      },
      header: {
        anchorStyle: 'none',
      },
      codeBlock: {
        theme: 'twilight',
        lineNumber: true, // 默认显示行号
        expandCode: true,
        copyCode: true,
        editCode: true,
        changeLang: true,
        customBtns: [],
      },
      table: {
        enableChart: false,
      },
      fontEmphasis: {
        allowWhitespace: false, // 是否允许首尾空格
      },
      strikethrough: {
        needWhitespace: false, // 是否必须有前后空格
      },
    },
  },
  multipleFileSelection: {
    video: false,
    audio: false,
    image: false,
    word: false,
    pdf: false,
    file: false,
  },
  toolbars: {
    toolbar: [
      'bold',
      'italic',
      {
        strikethrough: ['strikethrough', 'underline', 'sub', 'sup', 'ruby', 'customMenuAName'],
      },
      'size',
      '|',
      'color',
      'header',
      '|',
      'ol',
      'ul',
      'checklist',
      'panel',
      'justify',
      'detail',
      '|',
      'formula',
      {
        insert: [
          'image',
          'audio',
          'video',
          'link',
          'hr',
          'br',
          'code',
          'inlineCode',
          'formula',
          'toc',
          'table',
          'pdf',
          'word',
          'file',
        ],
      },
      'graph',
      'togglePreview',
      'shortcutKey',
    ],
    toolbarRight: ['fullScreen', '|', 'export', 'changeLocale', 'wordCount'],
    bubble: [
      'bold',
      'italic',
      'underline',
      'strikethrough',
      'sub',
      'sup',
      'quote',
      'ruby',
      '|',
      'size',
      'color',
      'theme',
    ], // array or false
    sidebar: ['mobilePreview', 'copy'], // 'theme'
    toc: {
      updateLocationHash: false, // 要不要更新URL的hash
      defaultModel: 'full', // pure: 精简模式/缩略模式，只有一排小点； full: 完整模式，会展示所有标题
    },
    // config: {
    //   publish: [
    //     {
    //       name: '微信公众号',
    //       key: 'wechat',
    //       icon: `data:image/svg+xml;charset=utf8,%3Csvg xmlns='http://www.w3.org/2000/svg' width='80' height='80' viewBox='0 0 80 80'%3E  %3Cg fill='none'%3E    %3Cpath fill='%23FFF' d='M0 0h80v80H0z' opacity='0'/%3E    %3Cpath fill='%2307C160' d='M60.962 22.753c-7.601-2.567-18.054-2.99-27.845 4.49-5.423 4.539-9.56 10.715-10.675 18.567-2.958-3.098-5.025-7.995-5.58-11.706-.806-5.403.483-10.82 4.311-15.45C26.906 11.724 34.577 10 39.6 10c9.57.001 18.022 5.882 21.363 12.753zm7.64 11.78c7.516 9.754 5.441 24.73-5.1 32.852-2.618 2.018-5.67 3.198-8.651 4.024a26.067 26.067 0 0 0 5.668-9.54c4.613-13.806-2.868-28.821-16.708-33.536-.3-.102-.601-.191-.903-.282 9.348-3.467 19.704-1.292 25.694 6.482zM39.572 59.37c6.403 0 11.474-1.49 16.264-5.013-.124 1.993-.723 4.392-1.271 5.805-4.509 11.633-17.56 16.676-31.238 12.183C11.433 68.438 4.145 54.492 7.475 42.851c.893-3.12 1.805-5.26 3.518-7.953 1.028 7.504 5.7 14.803 12.511 19.448.518.35.872.932.901 1.605a2.4 2.4 0 0 1-.08.653l-1.143 5.19c-.052.243-.142.499-.13.752.023.56.495.997 1.053.973.22-.01.395-.1.576-.215l6.463-4.143c.486-.312 1.007-.513 1.587-.538a3.03 3.03 0 0 1 .742.067c1.96.438 3.996.68 6.1.68z'/%3E  %3C/g%3E%3C/svg%3E`,
    //       serviceUrl: 'http://localhost:3001',
    //       injectPayload: {
    //         thumb_media_id: 'ft7IwCi1eukC6lRHzmkYuzeMmVXWbU3JoipysW2EZamblyucA67wdgbYTix4X377',
    //         author: 'Cherry Markdown',
    //       },
    //     }
    //   ],
    // },
  },

  editor: {
    autoSave2Textarea: true,
    defaultModel: 'edit&preview',
    showFullWidthMark: true, // 是否高亮全角符号 ·|￥|、|：|“|”|【|】|（|）|《|》
    showSuggestList: true, // 是否显示联想框
  },
  themeSettings: {
    mainTheme: 'ant-design',
    themeList: [{ className: 'ant-design', label: 'ant-design' }],
    codeBlockTheme: 'one-light',
  },
}

export default basicConfig
