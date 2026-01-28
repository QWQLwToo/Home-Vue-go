<template>
  <div class="content">
    <div class="user-profile-container">
      <div class="user-profile-image" v-motion-pop>
        <img :src="profileImage" alt="头像" @click.stop="toggleInfo">
        <span class="status-ball"></span>
      </div>
      <div class="user-name" v-motion-slide-left>
        <h1>Hi,</h1>
        <h1>I'm <span class="name-style">{{ userName }}</span></h1>
      </div>
    </div>
    <div class="description">
      <p ref="descriptionElement"></p>
    </div>
    <div class="contact-section" v-motion-pop>
      <template v-for="contact in contacts" :key="contact.type">
        <a v-if="contact.url" :href="contact.url" target="_blank" class="contact-item" :style="{ '--hover-color': contact.hoverColor }">
          <i :class="contact.icon"></i>
          <span class="tooltip">{{ contact.type }}</span>
        </a>
        <span v-else @click="toggleQRCode(contact.qrCode)" class="contact-item" :style="{ '--hover-color': contact.hoverColor }">
          <i :class="contact.icon"></i>
          <span class="tooltip">{{ contact.type }}</span>
        </span>
      </template>
      <span class="contact-item" @click="toggleDarkMode" :style="{ '--hover-color': isDarkMode ? '#ffcc00' : '#666' }">
        <i :class="darkModeIconClass"></i>
        <span class="tooltip">{{ isDarkMode ? '浅色' : '深色' }}</span>
      </span>
    </div>
    <Website /> 
    <!-- 使用v-if确保组件完全从DOM中移除，包括所有class -->
    <VisitTimer v-if="showVisitTimer" :key="showVisitTimer ? 'visit-timer-show' : 'visit-timer-hide'" />

    <Transition name="fade">
      <div v-if="showAbout" class="overlay" @click="showAbout = false">
        <div class="modal-content">
          <AboutPage @close="showAbout = false" />
        </div>
      </div>
    </Transition>

    <Transition name="fade">
      <div v-if="showQR" class="overlay" @click="hideQRCode">
        <div class="modal-content">
          <img :src="qrCodeSrc" alt="QR Code" class="qr-image" @click.stop>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { getContacts, getSiteConfig } from '../api';
import api from '../api';
import Website from './Website.vue';
import AboutPage from './AboutPage.vue';
import VisitTimer from './VisitTimer.vue';
import Typed from 'typed.js';

const contacts = ref([]);
const showQR = ref(false);
const showAbout = ref(false);
const qrCodeSrc = ref('');
const profileImage = ref('');
const userName = ref('');
const siteConfig = ref({});
const descriptionElement = ref(null);
const showVisitTimer = ref(true);

const loadData = async () => {
  try {
    const [contactsRes, configRes] = await Promise.all([
      getContacts(),
      getSiteConfig(),
    ]);
    contacts.value = contactsRes.data;
    siteConfig.value = configRes.data;
    userName.value = configRes.data.userName || import.meta.env.VITE_APP_USER_NAME || '用户';
    profileImage.value = configRes.data.profileImageURL || import.meta.env.VITE_APP_PROFILE_IMAGE_URL || '';
    // 确保正确读取 showVisitTimer，明确处理 false 值
    const timerValue = configRes.data.showVisitTimer;
    const newValue = timerValue !== undefined && timerValue !== null 
      ? Boolean(timerValue) 
      : true;
    // 使用赋值触发响应式更新
    const oldValue = showVisitTimer.value;
    showVisitTimer.value = newValue;
    console.log('加载配置 - showVisitTimer:', { 
      oldValue, 
      newValue, 
      rawValue: timerValue, 
      type: typeof timerValue,
      isFalse: timerValue === false 
    });
    
    // 如果从true变为false，确保DOM被移除
    if (oldValue && !newValue) {
      await nextTick();
      const timerElements = document.querySelectorAll('.visit-timer-container, .visit-timer');
      if (timerElements.length > 0) {
        console.warn('检测到需要移除的visit-timer元素，数量:', timerElements.length);
        timerElements.forEach(el => {
          console.log('移除元素:', el);
          el.remove();
        });
      }
    }
  } catch (error) {
    console.error('加载数据失败:', error);
    // 如果API失败，使用环境变量或默认值
    userName.value = import.meta.env.VITE_APP_USER_NAME || '用户';
    profileImage.value = import.meta.env.VITE_APP_PROFILE_IMAGE_URL || '';
    showVisitTimer.value = true; // 默认显示
  }
};

const predefinedDescriptions = ref([
  "你好鸭，欢迎来到我的主页！！",
  "随时可以联系我，期待与你交流。",
  "愿你历尽千帆，归来仍是少年。",
  "梦想还是要有的，万一实现了呢？",
  "I hope you have a happy day every day."
]);

let typedInstance = null;

const loadRotatingTexts = async () => {
  try {
    const res = await api.get('/rotating-texts');
    if (res.data?.texts && res.data.texts.length > 0) {
      predefinedDescriptions.value = res.data.texts;
    }
  } catch (error) {
    console.debug('加载轮换文本失败，使用默认文本:', error);
  }
};

const initializeTyped = () => {
  if (typedInstance) {
    typedInstance.destroy();
  }
  typedInstance = new Typed(descriptionElement.value, {
    strings: predefinedDescriptions.value,
    typeSpeed: 120,
    backSpeed: 80,
    showCursor: true,
    cursorChar: '|',
    loop: true,
  });
};

// 记录访问统计
const trackVisit = async () => {
  try {
    const path = window.location.pathname;
    const referer = document.referrer || '';
    await api.post('/track-visit', {
      path: path,
      referer: referer,
    });
  } catch (error) {
    // 静默失败，不影响用户体验
    console.debug('访问统计记录失败:', error);
  }
};

// 监听配置更新消息
let broadcastChannel = null
if (window.BroadcastChannel) {
  broadcastChannel = new BroadcastChannel('config-update')
  broadcastChannel.onmessage = async (event) => {
    if (event.data.type === 'config-updated') {
      console.log('收到配置更新消息，重新加载数据...')
      const oldTimerValue = showVisitTimer.value
      await loadData()
      await loadRotatingTexts()
      if (typedInstance) {
        typedInstance.destroy()
      }
      initializeTyped()
      // 停留时间显示状态会在loadData中更新
      // 使用 nextTick 确保 DOM 更新
      await nextTick()
      console.log('配置更新后 - showVisitTimer:', showVisitTimer.value, '之前的值:', oldTimerValue, 'DOM已更新')
      
      // 如果从true变为false，强制检查并移除残留的DOM元素
      if (oldTimerValue && !showVisitTimer.value) {
        console.log('showVisitTimer从true变为false，强制清理残留元素')
        const timerElements = document.querySelectorAll('.visit-timer-container, .visit-timer, [class*="visit-timer"]')
        if (timerElements.length > 0) {
          console.warn('发现残留的visit-timer相关元素，强制移除:', timerElements.length)
          timerElements.forEach(el => {
            console.log('移除残留元素:', el.className, el)
            el.remove()
          })
        }
        // 再次等待DOM更新
        await nextTick()
      }
    }
  }
}

// 监听showVisitTimer变化，确保DOM正确更新
watch(showVisitTimer, async (newValue, oldValue) => {
  console.log('showVisitTimer变化:', { oldValue, newValue })
  // 等待DOM更新完成
  await nextTick()
  // 如果设置为false，强制检查并移除任何残留的visit-timer元素
  if (!newValue) {
    const timerElements = document.querySelectorAll('.visit-timer-container, .visit-timer')
    if (timerElements.length > 0) {
      console.warn('发现残留的visit-timer元素，强制移除:', timerElements.length)
      timerElements.forEach(el => el.remove())
    }
  }
}, { immediate: false })

onMounted(async () => {
  await loadData();
  await loadRotatingTexts();
  initializeTyped();
  trackVisit(); // 记录访问
});

// 清理
onUnmounted(() => {
  if (broadcastChannel) {
    broadcastChannel.close()
  }
  if (typedInstance) {
    typedInstance.destroy()
  }
})

const toggleQRCode = (qrCode) => {
  qrCodeSrc.value = qrCode || '';
  showQR.value = !showQR.value;
};

const hideQRCode = () => {
  showQR.value = false;
};

const toggleInfo = () => {
  showAbout.value = !showAbout.value;
};

const isDarkMode = ref(false);
const darkModeIconClass = ref('fas fa-moon');

const toggleDarkMode = () => {
  isDarkMode.value = !isDarkMode.value;
  document.body.classList.toggle('dark-mode', isDarkMode.value);

  localStorage.setItem('darkMode', isDarkMode.value);

  darkModeIconClass.value = isDarkMode.value ? 'fas fa-sun' : 'fas fa-moon';
};

onMounted(() => {
  const savedDarkMode = localStorage.getItem('darkMode');
  if (savedDarkMode !== null) {
    isDarkMode.value = savedDarkMode === 'true';
    document.body.classList.toggle('dark-mode', isDarkMode.value);
  }

  darkModeIconClass.value = isDarkMode.value ? 'fas fa-sun' : 'fas fa-moon';
});
</script>

<style scoped>
.content {
  flex: 1;
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
  gap: 30px;
  margin-top: 20px;

  .user-profile-container {
    display: flex;
    align-items: center;
    gap: 30px;
  }

  .user-profile-image {
    display: flex;
    border-radius: 50%;
    box-shadow: 0 2px 8px var(--shadow-color);
    padding: 5px;
    border: 3px solid var(--border-color);
    position: relative;

    img {
      width: 150px;
      height: 150px;
      border-radius: 50%;
      background-size: cover;
      background-position: center;
    }

    .status-ball {
      position: absolute;
      background: #00c800;
      width: 2em;
      height: 2em;
      border-radius: 20px;
      border: 3px solid #eee;
      bottom: 5px;
      right: 15px;
      display: flex;
      justify-content: center;
      align-items: center;
      transition: all 0.3s ease;
      z-index: 1;
      cursor: pointer;
      overflow: hidden;

      &::before {
        content: "在线中";
        color: #00c800;
        opacity: 0;
        transition: opacity 0.3s ease-in-out, color 0.1s ease-in-out;
      }

      &:hover {
        width: 4.5em;
        height: 2em;
      }

      &:hover::before {
        opacity: 1;
        color: #eee;
      }
    }
  }

  .user-name {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    font-size: 1.3em;

    h1 {
      margin: 0;
    }
  }

  .name-style {
    position: relative;

    &:before {
      position: absolute;
      border-radius: 5px;
      bottom: 0;
      left: 50%;
      transform: translate(-50%);
      z-index: -1;
      content: "";
      background: #ffcc00ad;
      height: 30%;
      width: 110%;
      transition: height 0.3s ease-in-out;
    }
    &:hover::before {
      height: 60%;
    }
  }

  .description {
    display: flex;
    min-height: 32px;
    width: 100%;
    max-width: 500px;
    font-family: 'Georgia', serif;
    font-size: 1.2rem;
    white-space: nowrap;
    text-overflow: ellipsis;
    align-items: center;
    justify-content: center;
    transition: all 0.3s ease-in-out;

    &::before,
    &::after {
      content: '"';
      font-size: 1.5em;
      color: #999;
      margin: 0 10px;
    }

    p {
      margin: 0;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  .contact-section {
    display: flex;
    justify-content: center;
    gap: 20px;
    padding: 5px 10px;
    border: 1px solid transparent;
    border-radius: var(--border-radius);
    transition: all 0.3s ease-in-out;

    .contact-item {
      color: var(--text-color);
      font-size: var(--icon-size);
      cursor: pointer;
      transition: transform 0.3s ease-in-out, color 0.3s ease-in-out;
      position: relative;

      .fas.fa-moon {
        width: 20px;
        height: 20px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
      }

      &:hover {
        transform: translateY(-5px) rotate(10deg);
        color: var(--hover-color);

        .tooltip {
          opacity: 1;
          transform: translate(-50%, 0);
        }
      }

      .tooltip {
        position: absolute;
        bottom: 100%;
        left: 50%;
        transform: translate(-50%, 10px);
        opacity: 0;
        transition: opacity 0.3s ease, transform 0.3s ease;
        white-space: nowrap;
        pointer-events: none;
      }
    }

    &:hover {
      backdrop-filter: blur(10px);
      border: 1px solid var(--border-color);
      box-shadow: 0 2px 8px var(--shadow-color);
      background-color: rgba(var(--background-color-rgb), 0.2);
    }
  }

  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .fade-enter-active,
  .fade-leave-active {
    transition: all 0.3s ease-out;
    
    .modal-content {
      transition: all 0.3s ease-out;
    }
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
    
    .modal-content {
      transform: translateY(30px) scale(0.8);
      opacity: 0;
    }
  }

  .fade-enter-to,
  .fade-leave-from {
    opacity: 1;
    
    .modal-content {
      transform: translateY(0) scale(1);
      opacity: 1;
    }
  }

  .qr-image {
    width: 300px;
    height: 300px;
    background: white;
    padding: 20px;
    border-radius: var(--border-radius);
    box-shadow: 0 4px 8px var(--shadow-color);
    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
    &:hover {
      transform: scale(1.03) translateY(-5px);
      box-shadow: 0 15px 30px -10px rgba(0, 0, 0, 0.2);
    }
  }
}

@media screen and (max-width: 768px) {
  .content {
    gap: 15px;
  }
  .content .user-profile-container {
    flex-direction: column;
    gap: 0;
  }

  h1 {
    font-size: 1.5em;
  }
}
</style>