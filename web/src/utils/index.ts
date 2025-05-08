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

/**
 * @description 生成唯一 uuid
 */
export function generateUUID() {
    let uuid = ''
    for (let i = 0; i < 32; i++) {
        const random = (Math.random() * 16) | 0
        if (i === 8 || i === 12 || i === 16 || i === 20)
            uuid += '-'
        uuid += (i === 12 ? 4 : i === 16 ? (random & 3) | 8 : random).toString(16)
    }
    return uuid
}
