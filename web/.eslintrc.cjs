module.exports = {
    "env": {
        "browser": true,
        "node": true,
        "es2021": true
    },
    "parser": "vue-eslint-parser",
    "parserOptions": {
        "ecmaVersion": "latest",
        "parser": "@typescript-eslint/parser",
        "sourceType": "module",
        "ecmaFeatures": {
            "jsx": true
        }
    },
    "extends": [
        "eslint:recommended",
        "plugin:vue/vue3-recommended",
        "plugin:@typescript-eslint/recommended",
        "eslint-config-prettier",
        "./.eslintrc-auto-import.json",
        "plugin:prettier/recommended",
    ],
    "plugins": ["vue", "@typescript-eslint", "prettier"],
    "rules": {
        // eslint (http://eslint.cn/docs/rules)
        "no-var": "error", // 要求使用 let 或 const 而不是 var
        "no-undef": "off",
        "no-multiple-empty-lines": ["error", { "max": 1 }], // 不允许多个空行
        "prefer-const": "off", // 使用 let 关键字声明但在初始分配后从未重新分配的变量，要求使用 const
        "no-use-before-define": "off", // 禁止在 函数/类/变量 定义之前使用它们

        // typeScript (https://typescript-eslint.io/rules)
        "@typescript-eslint/no-unused-vars": [
            "off",
            {
                "argsIgnorePattern": "^_",
                "varsIgnorePattern": "^_"
            }
        ], // 禁止定义未使用的变量
        "@typescript-eslint/prefer-ts-expect-error": "error", // 禁止使用 @ts-ignore
        "@typescript-eslint/ban-ts-comment": "error", // 禁止 @ts-<directive> 使用注释或要求在指令后进行描述
        "@typescript-eslint/no-inferrable-types": "off", // 可以轻松推断的显式类型可能会增加不必要的冗长
        "@typescript-eslint/no-namespace": "off", // 禁止使用自定义 TypeScript 模块和命名空间
        "@typescript-eslint/no-explicit-any": "off", // 禁止使用 any 类型
        "@typescript-eslint/ban-types": "off", // 禁止使用特定类型
        "@typescript-eslint/no-var-requires": "off", // 允许使用 require() 函数导入模块
        "@typescript-eslint/no-empty-function": "off", // 禁止空函数
        "@typescript-eslint/no-non-null-assertion": "off", // 不允许使用后缀运算符的非空断言(!)

        // vue (https://eslint.vuejs.org/rules)
        "vue/script-setup-uses-vars": "error", // 防止<script setup>使用的变量<template>被标记为未使用，此规则仅在启用该no-unused-vars规则时有效
        "vue/v-slot-style": "error", // 强制执行 v-slot 指令样式
        "vue/no-mutating-props": "error", // 不允许改变组件 prop
        "vue/custom-event-name-casing": "error", // 为自定义事件名称强制使用特定大小写
        "vue/html-closing-bracket-newline": "error", // 在标签的右括号之前要求或禁止换行
        "vue/attribute-hyphenation": "error", // 对模板中的自定义组件强制执行属性命名样式：my-prop="prop"
        "vue/v-on-event-hyphenation": "off", // 关闭 eslint 中对 vue.js 时间名只用连接符的特殊要求
        "vue/attributes-order": "off", // vue api使用顺序，强制执行属性顺序
        "vue/no-v-html": "off", // 禁止使用 v-html
        "vue/require-default-prop": "off", // 此规则要求为每个 prop 为必填时，必须提供默认值
        "vue/multi-word-component-names": "off", // 要求组件名称始终为 “-” 链接的单词
        "prettier/prettier": 0, // 开启 prettier 格式化规则校验提示
        "vue/comment-directive": "off"
    }
}
