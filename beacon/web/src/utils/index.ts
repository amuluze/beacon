/**
 * @description 获取浏览器默认语言
 */
export function getBrowserLanguage(): string {
    const browserLanguage = navigator.language ? navigator.language : navigator.browserLanguage
    let defaultBrowserLanguage: string = ''

    if (['cn', 'zh', 'zh-cn'].includes(browserLanguage.toLowerCase())) {
        defaultBrowserLanguage = 'zh'
    }
    else {
        defaultBrowserLanguage = 'en'
    }
    return defaultBrowserLanguage
}
