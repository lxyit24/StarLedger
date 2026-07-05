import request from '../utils/request'

export interface ModuleInfo {
  name: string
  display_name: string
  description: string
  icon: string
  is_core: boolean
  enabled: boolean
}

export function listModules() {
  return request.get('/market/modules')
}

export function enableModule(name: string) {
  return request.put(`/market/modules/${name}/enable`)
}

export function disableModule(name: string) {
  return request.put(`/market/modules/${name}/disable`)
}
