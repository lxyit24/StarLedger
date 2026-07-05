<template>
  <div class="landing-page">
    <!-- Header -->
    <header class="header">
      <div class="container header-inner">
        <div class="logo">
          <span class="logo-icon">⭐</span>
          <span class="logo-text">星账 StarLedger</span>
        </div>
        <nav class="nav">
          <a href="#features">功能特性</a>
          <a href="#pricing">价格方案</a>
          <a href="#about">关于我们</a>
        </nav>
        <div class="header-actions">
          <el-button text @click="$router.push('/login')">登录</el-button>
          <el-button type="primary" @click="$router.push('/login')">免费试用</el-button>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <section class="hero">
      <div class="container">
        <h1 class="hero-title">一站式业财管理<br/>让经营更简单</h1>
        <p class="hero-subtitle">
          专为 IDC/服务器租赁/域名交易/云服务等行业打造的 SaaS 财税业务管理平台<br/>
          从账单到报表，星账帮您管好每一笔
        </p>
        <div class="hero-actions">
          <el-button type="primary" size="large" @click="$router.push('/login')">免费开始使用</el-button>
          <el-button size="large" @click="scrollTo('features')">了解更多</el-button>
        </div>
        <div class="hero-stats">
          <div class="stat-item">
            <div class="stat-value">500+</div>
            <div class="stat-label">企业用户</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">99.9%</div>
            <div class="stat-label">系统可用性</div>
          </div>
          <div class="stat-item">
            <div class="stat-value">50万+</div>
            <div class="stat-label">月处理账单</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="features">
      <div class="container">
        <h2 class="section-title">强大的功能模块</h2>
        <p class="section-desc">覆盖业务全流程，按需启用，灵活扩展</p>
        <div class="feature-grid">
          <div v-for="feature in features" :key="feature.title" class="feature-card">
            <div class="feature-icon">
              <el-icon :size="32"><component :is="feature.icon" /></el-icon>
            </div>
            <h3>{{ feature.title }}</h3>
            <p>{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Pricing Section -->
    <section id="pricing" class="pricing">
      <div class="container">
        <h2 class="section-title">灵活的价格方案</h2>
        <p class="section-desc">从个人到企业，总有一款适合您</p>
        <div class="pricing-grid">
          <div v-for="plan in plans" :key="plan.name" class="pricing-card" :class="{ featured: plan.featured }">
            <div v-if="plan.featured" class="pricing-badge">推荐</div>
            <h3 class="plan-name">{{ plan.name }}</h3>
            <div class="plan-price">
              <span v-if="plan.price === '免费'" class="price-free">免费</span>
              <template v-else>
                <span class="price-amount">{{ plan.price }}</span>
                <span class="price-unit">/月</span>
              </template>
            </div>
            <p class="plan-desc">{{ plan.desc }}</p>
            <ul class="plan-features">
              <li v-for="f in plan.features" :key="f">
                <el-icon color="#67c23a"><Check /></el-icon>
                {{ f }}
              </li>
            </ul>
            <el-button :type="plan.featured ? 'primary' : 'default'" size="large" class="plan-btn">
              {{ plan.price === '免费' ? '免费开始' : '立即订阅' }}
            </el-button>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="cta">
      <div class="container">
        <h2>准备好提升您的业财管理效率了吗？</h2>
        <p>14 天免费试用，无需绑定支付方式</p>
        <el-button type="primary" size="large" @click="$router.push('/login')">立即免费试用</el-button>
      </div>
    </section>

    <!-- Footer -->
    <footer class="footer">
      <div class="container footer-inner">
        <div class="footer-brand">
          <div class="logo">
            <span class="logo-icon">⭐</span>
            <span class="logo-text">星账 StarLedger</span>
          </div>
          <p>一站式业财管理平台</p>
        </div>
        <div class="footer-links">
          <div class="link-group">
            <h4>产品</h4>
            <a href="#features">功能介绍</a>
            <a href="#pricing">价格方案</a>
            <a href="#">更新日志</a>
          </div>
          <div class="link-group">
            <h4>支持</h4>
            <a href="#">帮助中心</a>
            <a href="#">API 文档</a>
            <a href="#">联系我们</a>
          </div>
          <div class="link-group">
            <h4>法律</h4>
            <router-link to="/terms">服务协议</router-link>
            <router-link to="/privacy">隐私政策</router-link>
            <a href="#">SLA 协议</a>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <div class="container">
          <p>© {{ new Date().getFullYear() }} 星账 StarLedger. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Check, Document, Monitor, Tickets, List, DataAnalysis, ShoppingCart, Connection } from '@element-plus/icons-vue'

const features = [
  { title: '账单管理', icon: Document, desc: '账单创建、支付跟踪、发票管理、逾期提醒，支持批量操作和自动对账' },
  { title: '服务器租赁', icon: Monitor, desc: '服务器资源登记、租赁合同关联、到期提醒、资源利用率统计' },
  { title: '合同管理', icon: Tickets, desc: '合同创建/审批/归档、状态管理、文件上传、到期续约提醒' },
  { title: '任务协作', icon: List, desc: '任务创建/分配/跟踪、看板视图、成员指派、工时统计' },
  { title: '数据报表', icon: DataAnalysis, desc: '收支趋势、账单分布、服务器成本分析、任务统计，多维度数据洞察' },
  { title: '模块市场', icon: ShoppingCart, desc: '按需启用/停用模块，灵活扩展，满足不同规模企业需求' },
  { title: '多租户隔离', icon: Connection, desc: '严格的数据隔离机制，支持个人/团队/企业多种租户类型' },
]

const plans = [
  {
    name: '免费版',
    price: '免费',
    desc: '适合个人用户和初创者',
    featured: false,
    features: ['账单管理（50笔/月）', '1 个用户', '基础数据报表', '社区支持'],
  },
  {
    name: '专业版',
    price: '¥99',
    desc: '适合中小团队',
    featured: true,
    features: ['账单管理（不限）', '服务器租赁 + 任务协作', '最多 10 个用户', '高级报表 + 数据导出', '邮件/站内信通知', '工单支持（24h响应）'],
  },
  {
    name: '企业版',
    price: '¥499',
    desc: '适合中大型企业',
    featured: false,
    features: ['全部模块', '不限用户数', '审批工作流', '多公司/多主体管理', '开放 API + Webhook', '专属客户经理', 'SLA 99.9% 保障'],
  },
]

function scrollTo(id: string) {
  document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
}
</script>

<style scoped>
.landing-page {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: #333;
}
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 24px;
}

/* Header */
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #eee;
  z-index: 100;
}
.header-inner {
  display: flex;
  align-items: center;
  height: 64px;
}
.logo {
  display: flex;
  align-items: center;
  gap: 8px;
}
.logo-icon { font-size: 24px; }
.logo-text { font-size: 18px; font-weight: 700; color: #303133; }
.nav {
  display: flex;
  gap: 32px;
  margin-left: 48px;
}
.nav a {
  color: #606266;
  text-decoration: none;
  font-size: 14px;
}
.nav a:hover { color: #409eff; }
.header-actions {
  margin-left: auto;
  display: flex;
  gap: 8px;
}

/* Hero */
.hero {
  padding: 140px 0 80px;
  text-align: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}
.hero-title {
  font-size: 48px;
  font-weight: 800;
  line-height: 1.3;
  margin-bottom: 24px;
}
.hero-subtitle {
  font-size: 18px;
  opacity: 0.9;
  line-height: 1.8;
  margin-bottom: 40px;
}
.hero-actions {
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 60px;
}
.hero-stats {
  display: flex;
  justify-content: center;
  gap: 80px;
}
.stat-item { text-align: center; }
.stat-value { font-size: 36px; font-weight: 700; }
.stat-label { font-size: 14px; opacity: 0.8; margin-top: 4px; }

/* Features */
.features {
  padding: 100px 0;
  background: #f9fafb;
}
.section-title {
  text-align: center;
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 12px;
}
.section-desc {
  text-align: center;
  color: #909399;
  font-size: 16px;
  margin-bottom: 60px;
}
.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
}
.feature-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
}
.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}
.feature-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}
.feature-card h3 {
  font-size: 18px;
  margin-bottom: 8px;
}
.feature-card p {
  color: #909399;
  font-size: 14px;
  line-height: 1.6;
}

/* Pricing */
.pricing {
  padding: 100px 0;
}
.pricing-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  max-width: 1000px;
  margin: 0 auto;
}
.pricing-card {
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 16px;
  padding: 40px 32px;
  text-align: center;
  position: relative;
  transition: transform 0.3s;
}
.pricing-card:hover { transform: translateY(-4px); }
.pricing-card.featured {
  border-color: #409eff;
  box-shadow: 0 8px 32px rgba(64, 158, 255, 0.15);
}
.pricing-badge {
  position: absolute;
  top: -12px;
  left: 50%;
  transform: translateX(-50%);
  background: #409eff;
  color: #fff;
  padding: 4px 16px;
  border-radius: 12px;
  font-size: 12px;
}
.plan-name { font-size: 20px; font-weight: 600; margin-bottom: 16px; }
.plan-price { margin-bottom: 12px; }
.price-free { font-size: 32px; font-weight: 700; color: #67c23a; }
.price-amount { font-size: 40px; font-weight: 700; color: #303133; }
.price-unit { font-size: 14px; color: #909399; }
.plan-desc { color: #909399; font-size: 14px; margin-bottom: 24px; }
.plan-features {
  list-style: none;
  padding: 0;
  margin: 0 0 32px;
  text-align: left;
}
.plan-features li {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  font-size: 14px;
  color: #606266;
}
.plan-btn { width: 100%; }

/* CTA */
.cta {
  padding: 80px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  text-align: center;
}
.cta h2 { font-size: 32px; margin-bottom: 12px; }
.cta p { font-size: 16px; opacity: 0.9; margin-bottom: 32px; }

/* Footer */
.footer {
  background: #303133;
  color: #c0c4cc;
  padding: 60px 0 0;
}
.footer-inner {
  display: flex;
  gap: 80px;
}
.footer-brand { flex: 1; }
.footer-brand .logo { margin-bottom: 12px; }
.footer-brand .logo-text { color: #fff; }
.footer-brand p { font-size: 14px; }
.footer-links {
  display: flex;
  gap: 60px;
}
.link-group h4 {
  color: #fff;
  font-size: 14px;
  margin-bottom: 16px;
}
.link-group a {
  display: block;
  color: #c0c4cc;
  text-decoration: none;
  font-size: 13px;
  padding: 4px 0;
}
.link-group a:hover { color: #409eff; }
.footer-bottom {
  margin-top: 40px;
  padding: 20px 0;
  border-top: 1px solid #4a4a4a;
  text-align: center;
  font-size: 13px;
}
</style>
