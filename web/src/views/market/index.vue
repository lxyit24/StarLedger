<template>
  <div class="market-page">
    <div class="page-header"><h2>模块市场</h2><p class="subtitle">按需启用功能模块，定制您的专属管理系统</p></div>
    <el-row :gutter="20">
      <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="mod in modules" :key="mod.name">
        <el-card class="module-card" :class="{ 'is-disabled': !mod.enabled }">
          <div class="module-header">
            <el-icon :size="36" class="module-icon"><component :is="mod.icon" /></el-icon>
            <div class="module-title-area">
              <h3>{{ mod.display_name }}</h3>
              <el-tag v-if="mod.is_core" type="danger" size="small">必选</el-tag>
              <el-tag v-else-if="mod.enabled" type="success" size="small">已启用</el-tag>
              <el-tag v-else type="info" size="small">未启用</el-tag>
            </div>
          </div>
          <p class="module-desc">{{ mod.description }}</p>
          <div class="module-actions">
            <el-switch v-if="!mod.is_core" :model-value="mod.enabled" @change="(val: boolean) => handleToggle(mod, val)" :loading="mod._loading" active-text="启用" inactive-text="停用" />
            <span v-else class="core-label">核心模块，始终启用</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-empty v-if="modules.length === 0 && !loading" description="暂无可用模块" />
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { listModules, enableModule, disableModule } from '../../api/market'
import type { ModuleInfo } from '../../api/market'
import { useAppStore } from '../../stores/app'
const modules = ref<(ModuleInfo & { _loading?: boolean })[]>([])
const loading = ref(false)
const appStore = useAppStore()
async function fetchModules() { loading.value = true; try { const res = await listModules(); modules.value = res.data || [] } catch { ElMessage.error('获取模块列表失败') } finally { loading.value = false } }
async function handleToggle(mod: ModuleInfo & { _loading?: boolean }, enabled: boolean) { mod._loading = true; try { if (enabled) { await enableModule(mod.name); mod.enabled = true; ElMessage.success(`${mod.display_name} 已启用`) } else { await disableModule(mod.name); mod.enabled = false; ElMessage.success(`${mod.display_name} 已停用`) } await appStore.fetchModules() } catch { ElMessage.error('操作失败') } finally { mod._loading = false } }
onMounted(fetchModules)
</script>
<style scoped>
.market-page { padding: 20px; }
.page-header { margin-bottom: 24px; }
.page-header h2 { margin: 0 0 8px; font-size: 22px; }
.subtitle { color: #909399; margin: 0; }
.module-card { margin-bottom: 20px; transition: all 0.3s; }
.module-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1); }
.module-card.is-disabled { opacity: 0.7; }
.module-header { display: flex; align-items: flex-start; gap: 12px; margin-bottom: 12px; }
.module-icon { color: #409eff; flex-shrink: 0; }
.module-title-area h3 { margin: 0 0 6px; font-size: 16px; }
.module-desc { color: #606266; font-size: 13px; line-height: 1.5; margin: 0 0 16px; }
.module-actions { border-top: 1px solid #f0f0f0; padding-top: 12px; }
.core-label { color: #909399; font-size: 13px; }
</style>
