declare namespace Form {
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
      | 'datetime'

    type Rule = Partial<Record<string, Arrayable<RuleItem>>>

    interface RuleItem {
        type?: RuleType // default type is 'string'
        required?: boolean
        message?: string | ((a?: string) => string)
        min?: number // Range of type 'string' and 'array'
        max?: number // Range of type 'string' and 'array'
        len?: number // Length of type 'string' and 'array'
        enum?: Array<string | number | boolean | null | undefined> // possible values of type 'enum'
        trigger?: string | string[]
        fields?: Record<string, Rule> // ignore when without required
        defaultField?: Rule // 'object' or 'array' containing validation rules
    }

    type ItemType = 'input'
      | 'password'
      | 'select'
      | 'datepicker'
      | 'timepicker'
      | 'datetimerange'
      | 'radio'
      | 'checkbox'
      | 'textarea'
      | 'upload'

    interface Options {
        label?: string
        value?: string
        placeholder?: string
        data: Record<string, any>[]
    }

    interface Item {
        label: string
        labelWidth?: number | string
        prop: string
        type?: ItemType
        value?: any
        rules?: RuleItem[]
        placeholder?: string
        options?: Options
        clearable?: boolean
        disabled?: boolean
        enterable?: boolean
        readonly?: boolean
        multiple?: boolean
    }
}
