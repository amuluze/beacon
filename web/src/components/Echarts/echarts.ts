/**
 * @Author     : Amu
 * @Date       : 2024/3/10 16:37
 * @Description:
 */
// 引入 echarts 核心模块。
import * as echarts from 'echarts/core'

//引入柱状图和折线图组件。
import { BarChart, GaugeChart, LineChart } from 'echarts/charts'

// 引入标题、提示框、网格、数据集和数据转换器组件。
import {
    DatasetComponent,
    GridComponent,
    LegendComponent,
    TitleComponent,
    ToolboxComponent,
    TooltipComponent, // 数据集组件
    TransformComponent // 内置数据转换器组件 (filter, sort)
} from 'echarts/components'

//引入标签布局和通用过渡动画特性。
import { LabelLayout, UniversalTransition } from 'echarts/features'

// 引入 Canvas 渲染器。
import { CanvasRenderer } from 'echarts/renderers'

// 系列类型的定义后缀都为 SeriesOption
import type { BarSeriesOption, GaugeSeriesOption, LineSeriesOption } from 'echarts/charts'

// 组件类型的定义后缀都为 ComponentOption
import type { DatasetComponentOption, GridComponentOption, TitleComponentOption, TooltipComponentOption } from 'echarts/components'

import type { ComposeOption } from 'echarts/core'

// 通过 ComposeOption 来组合出一个只有必须组件和图表的 Option 类型
export type EChartsOption = ComposeOption<BarSeriesOption | LineSeriesOption | TitleComponentOption | TooltipComponentOption | GridComponentOption | DatasetComponentOption | GaugeSeriesOption>

/**
 注册必须的组件，包括标题、提示框、网格、数据集、数据转换器，
 以及柱状图、折线图、标签布局、通用过渡动画和 Canvas 渲染器。
 */
echarts.use([TitleComponent, TooltipComponent, LegendComponent, ToolboxComponent, GridComponent, DatasetComponent, TransformComponent, BarChart, LineChart, GaugeChart, LabelLayout, UniversalTransition, CanvasRenderer])
// 导出
export default echarts
