import type { Router } from 'vue-router';
import { NO_REDIRECT_WHITE_LIST } from '@/constants/index'
import { useChatStore } from '@/stores/modules/chat'
import type { Chat } from '@/stores/modules/chat'
import { useLocaleStoreWithOut } from '@/stores/modules/locale'
import { useLocale } from '@/hooks/web/useLocale'

// 辅助函数：从路由的查询参数中提取Chat数据  
function extractChatDataFromQuery(query: any): Partial<Chat> {  
    return {  
        openid: query.openid || '',  
        robot_key: query.robot_key || '',  
        avatar: query.avatar || '',  
        name: query.name || '',  
        nickname: query.nickname || '',  
        dialogue_id: query.dialogue_id || '',  
    };  
}  

// 辅助函数：检查是否需要重定向  
function shouldRedirect(to: any): boolean {  
    return NO_REDIRECT_WHITE_LIST.some((n) => n === to.name);  
}  

export function createRouterGuards(router: Router) {  
    const chatStore = useChatStore();
    const localeStore = useLocaleStoreWithOut();
    const { changeLocale } = useLocale();

    // 辅助函数：处理语言切换
    const handleLanguageSwitch = (robotInfo: any) => {
        let newLang = 'zh-CN'; // 默认语言

        // 从 external_config_h5 中获取语言配置
        if (robotInfo?.external_config_h5) {
            try {
                const externalConfig = JSON.parse(robotInfo.external_config_h5);
                newLang = externalConfig.lang || 'zh-CN';
            } catch (error) {
                console.error('Failed to parse external_config_h5:', error);
            }
        }

        // 检查是否需要切换语言
        if (
            newLang &&
            newLang !== localeStore.currentLocale.lang &&
            ['zh-CN', 'en-US'].includes(newLang)
        ) {
            changeLocale(newLang as any);
        }
    };

    // 导航守卫  
    router.beforeEach(async (to, from, next) => {  
        const data = extractChatDataFromQuery(to.query);  

        if (shouldRedirect(to)) {  
            next();  
        } else {  
            if (chatStore.robot.id) {  
                next();  
            } else {  
                try {  
                    const res = await chatStore.createChat(data as Chat);
                    
                    // 在 createChat 完成后处理语言切换
                    if (res?.data?.robot) {
                        handleLanguageSwitch(res.data.robot);
                    }
                    
                    next();  
                } catch (error) {  
                    // 处理错误，比如重定向到错误页面或显示错误消息  
                    console.error('Failed to get robot info:', error);  
                    // 可以选择性地调用 next(false) 来阻止导航，并显示错误页面  
                }  
            }  
        }  
    });  
}