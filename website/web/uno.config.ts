import {
    defineConfig,
    presetAttributify,
    presetIcons,
    presetTypography,
    presetUno,
    presetWebFonts,
    transformerDirectives,
    transformerVariantGroup,
} from 'unocss'

// https://github.com/antfu/vitesse-nuxt/blob/main/uno.config.ts
export default defineConfig({
    shortcuts: [
        ['btn', 'px-4 py-1 rounded inline-block bg-[var(--primary)] text-white cursor-pointer hover:bg-[var(--color-brand-hover)] disabled:cursor-default disabled:opacity-50'],
        ['icon-btn', 'inline-block cursor-pointer select-none opacity-75 transition duration-200 ease-in-out hover:opacity-100 hover:text-[var(--primary)]'],
    ],
    presets: [
        presetUno(),
        presetAttributify(),
        presetIcons({
            scale: 1.2,
        }),
        presetTypography(),
        presetWebFonts({
            provider: 'none',
            fonts: {
                sans: 'Inter',
                mono: 'Geist Mono',
            },
        }),
    ],
    transformers: [
        transformerDirectives(),
        transformerVariantGroup(),
    ],
})
