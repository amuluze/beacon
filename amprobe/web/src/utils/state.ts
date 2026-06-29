/**
 * @Author     : Amu
 * @Date       : 2024/03/28 00:24:39
 * @Description:
 */
export function getTimeState() {
    const timeNow = new Date()
    const hours = timeNow.getHours()

    if (hours >= 6 && hours <= 10)
        return `早上好 ⛅`
    if (hours >= 10 && hours <= 14)
        return `中午好 🌞`
    if (hours >= 14 && hours <= 18)
        return `下午好 🌞`
    if (hours >= 18 && hours <= 24)
        return `晚上好 🌛`
    if (hours >= 0 && hours <= 6)
        return `凌晨好 🌛`
}
