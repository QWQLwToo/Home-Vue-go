<template>
  <div class="container">
    <div class="swiper-container">
      <div class="swiper-wrapper">
        <div v-for="(siteChunk, index) in chunkedSites" :key="index" class="swiper-slide">
          <div class="site-grid">
            <div v-for="(site, i) in siteChunk" :key="i" class="site-box" @click="openLink(site.url)">
              <div class="site-content">
                <i :class="site.icon" aria-hidden="true"></i>
                <span class="site-name">{{ site.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="swiper-pagination"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import Swiper from 'swiper/bundle';
import 'swiper/swiper-bundle.css';
import { getSites } from '../api';

const sites = ref([]);
const chunkedSites = ref([]);

const loadSites = async () => {
  try {
    const res = await getSites();
    sites.value = res.data;
    // 将站点数据分块（每页6个）
    chunkedSites.value = sites.value.reduce((acc, site, index) => {
      const chunkIndex = Math.floor(index / 6);
      if (!acc[chunkIndex]) acc[chunkIndex] = [];
      acc[chunkIndex].push(site);
      return acc;
    }, []);
    
    // 如果数据加载后需要重新初始化Swiper
    if (chunkedSites.value.length > 0) {
      setTimeout(() => {
        initSwiper();
      }, 100);
    }
  } catch (error) {
    console.error('加载站点数据失败:', error);
    // 如果API失败，使用空数组
    chunkedSites.value = [];
  }
};

let swiperInstance = null;

const initSwiper = () => {
  if (swiperInstance) {
    swiperInstance.destroy();
  }
  swiperInstance = new Swiper('.swiper-container', {
    slidesPerView: 1,
    spaceBetween: 20,
    pagination: { el: '.swiper-pagination', clickable: true },
    mousewheel: true,
  });
};

const openLink = (url) => {
  if (url) window.open(url, '_blank');
};

onMounted(() => {
  loadSites();
});
</script>


<style scoped>
.container {
  max-width: 700px;
  width: 100%;
  margin: 30px 0 20px;
}

.swiper-container {
  overflow: hidden;
  padding: 10px;
}

.swiper-pagination {
  bottom: inherit;
}

.site-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 15px;
}

.site-box {
  padding: 30px;
  backdrop-filter: blur(10px);
  border-radius: var(--border-radius);
  border: 1px solid var(--border-color);
  background-color: rgba(var(--background-color-rgb), 0.2);
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 1px 8px var(--shadow-color);
  }
}

.site-content {
  display: flex;
  gap: 10px;
  justify-content: center;
  align-items: center;

  i {
    font-size: var(--icon-size);
  }
}

.site-name {
  margin: 0;
  font-size: 1.17em;
  font-weight: bold;
}

@media screen and (max-width: 768px) {
  .site-content {
    gap: 5px;
    flex-direction: column;
  }

  .site-box {
    padding: 15px;
    border-radius: 8px;
  }

  .site-name {
    font-size: 16px;
  }

  .site-content i {
    font-size: 18px;
  }
}

:deep(.swiper-pagination-bullet-active) {
  background: #8c8c8c94;
  width: 20px;
  border-radius: 5px;
}
</style>