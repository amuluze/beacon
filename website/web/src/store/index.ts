/**
 * @Author     : Amu
 * @Date       : 2024/05/21 23:26:15
 * @Description:
 */

import { useStatisticsStore } from './modules/statistics'

function useStore() {
    return {
        statistics: useStatisticsStore(),
    }
}

export default useStore
