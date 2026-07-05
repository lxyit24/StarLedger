<template>
  <div class="legal-page">
    <header class="legal-header">
      <div class="container header-inner">
        <router-link to="/" class="back-link">← 返回首页</router-link>
        <h1>服务协议</h1>
      </div>
    </header>
    <main class="container legal-content">
      <p class="update-time">最后更新时间：2024年1月1日</p>

      <section>
        <h2>一、总则</h2>
        <p>欢迎使用星账系统（以下简称"本平台"）。本服务协议（以下简称"本协议"）是您与星账系统运营方之间关于使用本平台服务所订立的协议。请您仔细阅读本协议的全部内容，如果您对本协议有任何疑问，您可以通过平台提供的联系方式与我们沟通。</p>
        <p>您注册并使用本平台服务，即表示您已阅读并同意受本协议的约束。如果您不同意本协议的任何内容，请勿使用本平台服务。</p>
      </section>

      <section>
        <h2>二、服务内容</h2>
        <p>本平台为您提供一站式业财管理 SaaS 服务，包括但不限于：</p>
        <ul>
          <li>账单管理：账单创建、支付跟踪、发票管理、逾期提醒等功能</li>
          <li>服务器租赁管理：服务器资源登记、租赁合同关联、到期提醒等功能</li>
          <li>合同管理：合同创建、审批、归档、状态管理等功能</li>
          <li>任务协作：任务创建、分配、跟踪、成员指派等功能</li>
          <li>数据报表：多维度数据统计分析、图表展示、数据导出等功能</li>
          <li>模块市场：按需启用/停用功能模块</li>
        </ul>
        <p>具体服务内容以您订购的套餐版本为准，不同版本的功能范围详见平台公示的价格方案。</p>
      </section>

      <section>
        <h2>三、账号注册与安全</h2>
        <p>1. 您需要提供真实、准确、完整的注册信息来创建账号，并对账号信息的真实性、完整性负责。</p>
        <p>2. 您应妥善保管账号和密码，对账号下发生的所有活动承担全部责任。</p>
        <p>3. 如发现账号被未经授权使用或存在安全漏洞，请立即通知我们。</p>
        <p>4. 我们有权对涉嫌违规的账号进行暂停或终止处理。</p>
      </section>

      <section>
        <h2>四、使用规范</h2>
        <p>您在使用本平台服务时，应遵守以下规范：</p>
        <ul>
          <li>遵守中华人民共和国法律法规</li>
          <li>不得利用本平台从事违法违规活动</li>
          <li>不得干扰或破坏本平台的正常运行</li>
          <li>不得侵犯其他用户的合法权益</li>
          <li>不得未经授权访问他人数据</li>
        </ul>
      </section>

      <section>
        <h2>五、数据安全与隐私</h2>
        <p>1. 我们重视您的数据安全，采用行业标准的加密技术保护您的数据。</p>
        <p>2. 多租户数据严格隔离，确保您的数据不会被其他用户访问。</p>
        <p>3. 我们会定期备份数据，但因不可抗力导致的数据丢失，我们将尽力恢复但不承担赔偿责任。</p>
        <p>4. 关于个人信息的收集、使用和保护，请参阅我们的《隐私政策》。</p>
      </section>

      <section>
        <h2>六、知识产权</h2>
        <p>1. 本平台的所有权及知识产权归星账系统运营方所有。</p>
        <p>2. 您在使用本平台过程中上传的数据，其知识产权归您所有。</p>
        <p>3. 未经书面许可，您不得复制、修改、分发本平台的任何部分。</p>
      </section>

      <section>
        <h2>七、服务变更与终止</h2>
        <p>1. 我们保留随时修改或中断服务的权利，但会提前通知您。</p>
        <p>2. 如您违反本协议，我们有权暂停或终止您的服务。</p>
        <p>3. 服务终止后，您的数据将保留 90 天，届时未续费的账号数据将被清除。</p>
      </section>

      <section>
        <h2>八、免责声明</h2>
        <p>1. 本平台按"现状"提供服务，不作任何明示或暗示的保证。</p>
        <p>2. 因不可抗力（包括但不限于自然灾害、政策变化、网络故障等）导致的服务中断，我们不承担责任。</p>
        <p>3. 我们的累计赔偿责任不超过您在过去 12 个月内实际支付的服务费用总额。</p>
      </section>

      <section>
        <h2>九、争议解决</h2>
        <p>本协议的解释和执行均适用中华人民共和国法律。因本协议产生的争议，双方应友好协商解决；协商不成的，任何一方可向本平台运营方所在地有管辖权的人民法院提起诉讼。</p>
      </section>

      <section>
        <h2>十、联系方式</h2>
        <p>如对本协议有任何疑问，请联系：</p>
        <p>邮箱：support@starledger.com</p>
        <p>地址：请联系我们获取最新地址信息</p>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.legal-page {
  min-height: 100vh;
  background: #f9fafb;
}
.legal-header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  padding: 24px 0;
}
.header-inner {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 24px;
}
.back-link {
  color: #409eff;
  text-decoration: none;
  font-size: 14px;
  display: inline-block;
  margin-bottom: 12px;
}
.legal-header h1 {
  font-size: 28px;
  margin: 0;
}
.legal-content {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 24px;
}
.update-time {
  color: #909399;
  font-size: 14px;
  margin-bottom: 32px;
}
section {
  margin-bottom: 32px;
}
section h2 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}
section p {
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
  margin-bottom: 8px;
}
section ul {
  padding-left: 24px;
  margin-bottom: 8px;
}
section li {
  font-size: 14px;
  line-height: 1.8;
  color: #606266;
}
</style>
