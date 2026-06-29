import { beforeEach } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

// Create a fresh Pinia instance for each test
beforeEach(() => {
    setActivePinia(createPinia())
})
