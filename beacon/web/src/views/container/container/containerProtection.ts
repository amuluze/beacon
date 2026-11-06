const protectedContainerNames = new Set(['beacon'])

// 与服务端 assertNotProtected 的判定对齐：先 trim 再转小写。
// 否则「Beacon」「 beacon 」这类名字会在界面放行，而服务端仍然拒绝并返回 403。
export function isProtectedContainerName(containerName: string) {
    return protectedContainerNames.has(containerName.trim().toLowerCase())
}
