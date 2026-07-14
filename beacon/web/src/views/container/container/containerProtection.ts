const protectedContainerNames = new Set(['beacon', 'amprobe'])

export function isProtectedContainerName(containerName: string) {
    return protectedContainerNames.has(containerName)
}
