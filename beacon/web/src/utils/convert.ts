/**
 * @Author     : Amu
 * @Date       : 2024/3/12 13:17
 * @Description:
 */

export function convertBytesToReadable(bytes: number): string {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let unitIndex = 0
    while (bytes >= 1024 && unitIndex < units.length - 1) {
        bytes /= 1024
        unitIndex++
    }
    return `${bytes.toFixed(2)} ${units[unitIndex]}`
}
