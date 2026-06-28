// eslint.config.js
import antfu from '@antfu/eslint-config'

export default antfu(
    {
        typescript: true,
        unocss: true,
        vue: true,
        ignores: ['**/fixtures', '.nuxt/'],
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
    },
    {
        rules: {
            'no-console': 'off',
            'node/prefer-global/process': 'off',
            'style/indent': ['error', 4],
            'jsonc/indent': ['error', 4],
            'vue/html-indent': ['error', 4],

        },
    },
)
