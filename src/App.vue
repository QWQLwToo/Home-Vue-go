<template>
  <router-view v-if="isAdminPage || isLoginPage" />
  <template v-else>
    <div class="background"></div>
    <Home />
    <footer>
      <span>© {{ footerYearText }} Made in <a href="/" target="_blank">{{ userName }}</a></span>
      <a v-if="icpNumber && icpNumber !== '暂未填写' && icpNumber.trim() !== ''" href="https://beian.miit.gov.cn/" target="_blank">{{ icpNumber }}</a>
      <a v-if="policenumber && policenumber !== '暂未填写' && policenumber.trim() !== ''" :href="`https://beian.mps.gov.cn/#/query/webSearch?police=${policenumber}`" target="_blank" class="police_link">
        <span class="police_img"></span> {{ policenumber }}
      </a>
    </footer>
  </template>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import Home from './components/Home.vue';
import { getSiteConfig } from './api';

const route = useRoute();
const userName = ref(import.meta.env.VITE_APP_USER_NAME || '用户');
const icpNumber = ref(import.meta.env.VITE_APP_ICP_NUMBER || '');
const policenumber = ref(import.meta.env.VITE_APP_POLICE_NUMBER || '');

// 底部年份（支持起止年份）
const footerYearStart = ref('');
const footerYearEnd = ref('');

const footerYearText = computed(() => {
  const currentYear = new Date().getFullYear().toString();
  const start = (footerYearStart.value || '').trim();
  const end = (footerYearEnd.value || '').trim();

  // 后台未配置时，默认显示当前年份
  if (!start && !end) {
    return currentYear;
  }

  // 只配置了起始年份
  if (start && !end) {
    return start;
  }

  // 起止年份都有且不同
  if (start && end && start !== end) {
    return `${start}~${end}`;
  }

  // 其他情况（例如起止相同），只显示一个年份
  return start || end || currentYear;
});

const isAdminPage = computed(() => route.path === '/admin');
const isLoginPage = computed(() => route.path === '/login');

const loadConfig = async () => {
  try {
    const res = await getSiteConfig();
    // 同步所有配置信息
    if (res.data.siteName) {
      // 站点名称可用于页面显示
    }
    if (res.data.siteURL) {
      // 站点URL可用于链接等
    }
    if (res.data.siteDescription) {
      // 站点描述已通过frontendConfig.js同步到meta标签
    }
    if (res.data.siteKeywords) {
      // 站点关键词已通过frontendConfig.js同步到meta标签
    }
    // 只有非空且不是"暂未填写"时才显示备案信息
    if (res.data.icpNumber && res.data.icpNumber !== '暂未填写' && res.data.icpNumber.trim() !== '') {
      icpNumber.value = res.data.icpNumber;
    } else {
      icpNumber.value = '';
    }
    if (res.data.policeNumber && res.data.policeNumber !== '暂未填写' && res.data.policeNumber.trim() !== '') {
      policenumber.value = res.data.policeNumber;
    } else {
      policenumber.value = '';
    }
    if (res.data.userName) userName.value = res.data.userName;

    // 底部年份配置
    if (typeof res.data.footerYearStart === 'string') {
      footerYearStart.value = res.data.footerYearStart;
    }
    if (typeof res.data.footerYearEnd === 'string') {
      footerYearEnd.value = res.data.footerYearEnd;
    }
  } catch (error) {
    console.error('加载配置失败:', error);
  }
};

onMounted(() => {
  if (!isAdminPage.value && !isLoginPage.value) {
    loadConfig();
  }
});
</script>