module.exports = {
    // 继承推荐规范配置
    extends: [
        'stylelint-config-standard',
        'stylelint-config-recommended-scss',
        'stylelint-config-recommended-vue/scss',
        'stylelint-config-html/vue',
        'stylelint-config-recess-order'
    ],
    customSyntax: 'postcss-html',
    // 指定不同文件对应的解析器
    overrides: [
        {
            files: ['**/*.{vue,html}'],
            customSyntax: 'postcss-html'
        },
        {
            files: ['**/*.{css,scss}'],
            customSyntax: 'postcss-scss'
        }
    ],
    ignoreFiles: [
        '**/*.js',
        '**/*.jsx',
        '**/*.tsx',
        '**/*.ts',
        '**/*.json',
        '**/*.md',
        '**/*.yaml',
    ],
    // 自定义规则
    rules: {
        // 颜色值小写
        'color-hex-case': 'lower',
        // 缩进
        'indentation': 4,
        //
        'property-no-unknown': null,
        // 允许 global 、export 、v-deep等伪类
        'selector-pseudo-class-no-unknown': [
            true,
            {
                ignorePseudoClasses: ['global', 'export', 'v-deep', 'deep']
            }
        ]
    }
}
