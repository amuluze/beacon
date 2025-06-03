import antfu from '@antfu/eslint-config'

// https://ithelp.ithome.com.tw/articles/10349428
export default antfu(
    {
        formatters: {
            /**
             * Format CSS, LESS, SCSS files, also the `<style>` blocks in Vue
             * By default uses Prettier
             */
            css: true,
            /**
             * Format HTML files
             * By default uses Prettier
             */
            html: true,
            /**
             * Format Markdown files
             * Supports Prettier and dprint
             * By default uses Prettier
             */
            markdown: 'prettier',
        },
        typescript: {
            tsconfigPath: 'tsconfig.json',
        },
        vue: true,
        unocss: true,
        ignores: ['**/fixtures'],
    },
    {
        rules: {
            'no-console': 'off',
            'style/no-tabs': 'off',
            'perfectionist/sort-imports': 'off',
            'ts/restrict-template-expressions': 'off',
            'ts/no-unsafe-function-type': 'off',
            'ts/no-unsafe-call': 'off',
            'ts/no-unsafe-argument': 'off',
            'ts/no-unsafe-member-access': 'off',
            'ts/no-unsafe-assignment': 'off',
            'ts/no-unsafe-return': 'off',
            'ts/strict-boolean-expressions': 'off',
            'ts/explicit-function-return-type': 'off',
            'ts/ban-ts-comment': 'off',
            'style/indent': ['error', 4],
            'jsonc/indent': ['error', 4],
            'vue/html-indent': ['error', 4],
            'vue/script-indent': ['error', 4],
            'vue/style-indent': ['error', 4],
        },
    },
    {
        files: ['**/*.vue'],
        rules: {
            // https://eslint.vuejs.org/rules/script-indent
            'vue/style-indent': ['error', 4], // 确保<style>标签内容缩进为4个空格
            'vue/html-indent': ['error', 4], // 确保HTML标签的缩进为4个空格
            'vue/script-indent': ['error', 4], // 确保<script>标签的缩进为4个空格
            'style/indent': 'off',
            'vue/operator-linebreak': ['error', 'before'],
            // 修改这个规则，允许多行属性的标签在同一行关闭
            'vue/html-closing-bracket-newline': ['error', {
                singleline: 'never',
                multiline: 'never',  // 改为 'never'，允许多行属性的标签在同一行关闭
                selfClosingTag: {
                    singleline: 'never',
                    multiline: 'never',  // 改为 'never'
                },
            }],
            'vue/html-self-closing': ['error', {
                html: {
                    void: 'always',
                    normal: 'always',
                    component: 'always',
                },
                svg: 'always',
                math: 'always',
            }],
        },
    },
)
