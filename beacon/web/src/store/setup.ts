import { createPinia } from 'pinia'
import { createPersistedState } from 'pinia-plugin-persistedstate'
import type { App } from 'vue'

// Remove the legacy persisted user store once. Older releases stored access
// and refresh JWTs under this key; the Cookie-based session must not leave
// those credentials readable by JavaScript for the lifetime of the tab.
sessionStorage.removeItem('user')

const store = createPinia()
store.use(createPersistedState({
    storage: sessionStorage,
}))

export function setupStore(app: App<Element>) {
    app.use(store)
}

export { store }
