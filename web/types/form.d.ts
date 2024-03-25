declare namespace Form {
    type ItemType =
        | 'text'
        | 'password'
        | 'select'
        | 'datepicker'
        | 'datetimepicker'
        | 'radio'
        | 'checkbox'
        | 'textarea'
        | 'upload'

    interface FieldOptions {
        label?: string
        value?: string
        placeholder?: string
        format?: string
        valueFormat?: string
        data: Record<string, any>[]
    }

    interface FormItem {
        label?: string
        labelWidth?: string | number
        prop: string
        type?: ItemType
        value?: any
        rules?: Arrayable<RuleItem>
        placeholder?: string
        options?: FieldOptions
        readonly?: boolean // 是否只读
        showPassword?: boolean // 是否切换显示密码
        clearable?: boolean // 是否可清空
        disabled?: boolean // 是否禁用
        enterable?: boolean // 当为输入框时，是否启用回车触发提交功能
        multiple?: boolean
    }

    type RuleType =
        | 'string'
        | 'number'
        | 'boolean'
        | 'method'
        | 'regexp'
        | 'integer'
        | 'float'
        | 'array'
        | 'object'
        | 'enum'
        | 'date'
        | 'url'
        | 'hex'
        | 'email'
        | 'pattern'
        | 'any'

    interface RuleItem {
        type?: RuleType // default type is 'string'
        required?: boolean
        pattern?: RegExp | ((a?: string) => boolean)
        message?: string | ((a?: string) => string)
        min?: number // Range of type 'string' and 'array'
        max?: number // Range of type 'string' and 'array'
        len?: number // Length of type 'string' and 'array'
        enum?: Array<string | number | boolean | null | undefined> // possible values of type 'enum'
        trigger?: string | string[]
        fields?: Record<string, Rule> // ignore when without required
        defaultField?: Rule // 'object' or 'array' containing validation rules
    }

    type Rule = Partial<Record<string, Arrayable<RuleItem>>>
    type ButtonOptionsPosition = 'left' | 'right' | 'top'
    type ButtonOptionsSize = 'large' | 'small' | 'default'
    interface ButtonOptions {
        labelWidth?: string | number
        labelPosition?: ButtonOptionsPosition
        disabled?: boolean
        size?: ButtonOptionsSize
        showResetButton?: boolean // 是否展示重置按钮
        showCancelButton?: boolean // 是否展示取消按钮
        submitButtonText?: string
        resetButtonText?: string
        cancelButtonText?: string
        blockSubmitButton?: boolean // 提交按钮是否以块级按钮呈现
    }
}
