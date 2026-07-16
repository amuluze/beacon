(function () {
    try {
        const stored = localStorage.getItem('beacon-theme')
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        if (stored === 'dark' || (stored === null && prefersDark))
            document.documentElement.classList.add('dark')
    }
    catch {}
})()
