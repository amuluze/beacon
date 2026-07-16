interface PageSeoOptions {
    title: string
    description: string
    path: string
}

const siteURL = 'https://help.beacon.amuluze.com'
const socialImage = `${siteURL}/images/beacon.png`

export function usePageSeo(options: PageSeoOptions) {
    const canonical = new URL(options.path, siteURL).toString()

    useSeoMeta({
        title: options.title,
        description: options.description,
        ogTitle: options.title,
        ogDescription: options.description,
        ogType: 'website',
        ogUrl: canonical,
        ogImage: socialImage,
        twitterCard: 'summary_large_image',
        twitterTitle: options.title,
        twitterDescription: options.description,
        twitterImage: socialImage,
    })
    useHead({
        link: [{ rel: 'canonical', href: canonical }],
    })
}
