import { createPinia } from 'pinia'
import { setActivePinia } from 'pinia'

// Create a fresh Pinia instance for each test
beforeEach(() => {
  setActivePinia(createPinia())
})
