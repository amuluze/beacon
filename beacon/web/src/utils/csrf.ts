export const csrfCookieName = 'beacon_csrf_token'
export const csrfHeaderName = 'X-CSRF-Token'

export function getCSRFToken(): string {
    const prefix = `${csrfCookieName}=`
    const cookie = document.cookie
        .split(';')
        .map(part => part.trim())
        .find(part => part.startsWith(prefix))
    if (!cookie)
        return ''
    return decodeURIComponent(cookie.slice(prefix.length))
}

export function getCSRFHeaders(): Record<string, string> {
    const token = getCSRFToken()
    return token ? { [csrfHeaderName]: token } : {}
}
