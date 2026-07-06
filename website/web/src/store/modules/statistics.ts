/**
 * @Author     : Amu
 * @Date       : 2024/05/21 23:26:15
 * @Description:
 */

interface StatisticsState {
    id: number
    visitCount: number // 更明确的字段名
}

export const useStatisticsStore = defineStore('statistics', {
    state: (): StatisticsState => ({
        id: 0,
        visitCount: 0,
    }),
    getters: {
        getVisitCount: state => state.visitCount,
    },
    actions: {
        setStatistic(statistic: number) {
            this.visitCount = statistic
        },
        setID(id: number) {
            this.id = id
        },
    },
    persist: {
        storage: piniaPluginPersistedstate.localStorage(),
    },
})
