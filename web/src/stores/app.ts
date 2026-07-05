import { defineStore } from 'pinia'
import { ref } from 'vue'
import { listModules } from '../api/market'
import type { ModuleInfo } from '../api/market'

export const useAppStore = defineStore('app', () => {
  const modules = ref<ModuleInfo[]>([])
  const loaded = ref(false)

  async function fetchModules() {
    try {
      const res = await listModules()
      modules.value = res.data || []
      loaded.value = true
    } catch {
      // Silently fail - sidebar will show all items
      loaded.value = false
    }
  }

  function isModuleEnabled(name: string): boolean {
    if (!loaded.value) return true // Show all if not loaded yet
    const mod = modules.value.find(m => m.name === name)
    if (!mod) return false
    return mod.enabled || mod.is_core
  }

  function getEnabledModules(): ModuleInfo[] {
    return modules.value.filter(m => m.enabled || m.is_core)
  }

  return {
    modules,
    loaded,
    fetchModules,
    isModuleEnabled,
    getEnabledModules,
  }
})
