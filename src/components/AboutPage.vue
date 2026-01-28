<template>
  <div class="about-page" @click.stop>
    <div class="about-modal">
      <div class="about-modal-content">
        <div class="tech-stack">
          <h3>使用的技术栈</h3>
          <ul class="tech-list">
            <li v-for="tech in techStack" :key="tech.name" :class="['tech-item', tech.name.toLowerCase()]">
              <i :class="tech.icon"></i>
              {{ tech.name }}
            </li>
          </ul>
        </div>
        <div class="github-info">
          <h3>开源地址</h3>
          <div class="github-links">
            <a href="https://github.com/JLinMr/Home-Vue" target="_blank" class="github-link">
              <i class="fab fa-github"></i>
              <div class="link-content">
                <span class="link-title">静态原项目</span>
                <span class="link-desc">Home-Vue</span>
              </div>
            </a>
            <a href="https://github.com/JLinMr/Home-Vue-go" target="_blank" class="github-link">
              <i class="fab fa-github"></i>
              <div class="link-content">
                <span class="link-title">动态现项目</span>
                <span class="link-desc">Home-Vue-go</span>
              </div>
            </a>
          </div>
        </div>
      </div>
      <button @click="closeModal" class="close-btn">
        <i class="fas fa-times"></i>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const emit = defineEmits(['close']);

const techStack = [
  // 前端技术栈
  { name: 'Vue3', icon: 'fab fa-vuejs' },
  { name: 'Vite', icon: 'fas fa-bolt' },
  { name: 'CSS3', icon: 'fab fa-css3-alt' },
  { name: 'HTML5', icon: 'fab fa-html5' },
  { name: 'JavaScript', icon: 'fab fa-js' },
  // 后端技术栈
  { name: 'Go', icon: 'fab fa-golang' },
  { name: 'Gin', icon: 'fas fa-server' },
  { name: 'SQLite', icon: 'fas fa-database' },
  { name: 'Ent', icon: 'fas fa-code-branch' },
  { name: 'JWT', icon: 'fas fa-key' }
];
const closeModal = () => emit('close');
</script>

<style scoped>
.about-page {
  width: 500px;
  backdrop-filter: blur(5px);
  background-color: rgba(var(--background-color-rgb), 0.9);
  padding: 40px;
  border-radius: var(--border-radius);
  position: relative;
  text-align: center;
  box-sizing: border-box;

  @media (max-width: 600px) {
    padding: 20px;
    width: 90%;
    margin: auto;
  }
}

h3 {
  border-bottom: 2px solid var(--border-color);
  padding-bottom: 10px;
}

.tech-stack ul,
.tech-list {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
  justify-content: center;
}

.tech-stack li,
.tech-item {
  padding: 10px 15px;
  border-radius: var(--border-radius);
  transition: all 0.3s ease;
  display: flex;
  gap: 5px;
  flex-direction: column;
}

.tech-stack li:hover,
.tech-item:hover {
  background-color: var(--hover-other-color);
  color: var(--hover-link-color);
}

.tech-stack li i,
.tech-item i {
  font-size: 1.5em;
}

@media (max-width: 600px) {
  .tech-stack li,
  .tech-item {
    font-size: 0.8em;
  }
}

.github-links {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
  width: 100%;
}

.github-link {
  display: flex;
  align-items: center;
  gap: 12px;
  background-color: #040404d0;
  color: white;
  padding: 15px 20px;
  border-radius: var(--border-radius);
  text-decoration: none;
  transition: all 0.3s ease;
  flex: 1;

  &:hover {
    background-color: #1a1a1a;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  i {
    font-size: 1.5em;
    flex-shrink: 0;
  }

  .link-content {
    display: flex;
    flex-direction: column;
    gap: 4px;
    text-align: left;
  }

  .link-title {
    font-size: 0.9em;
    opacity: 0.9;
  }

  .link-desc {
    font-size: 0.85em;
    opacity: 0.7;
    font-family: 'Courier New', monospace;
  }
}

@media (max-width: 600px) {
  .github-links {
    grid-template-columns: 1fr;
  }
}

.close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  background: none;
  border: none;
  font-size: 1.5em;
  cursor: pointer;
  color: #7f8c8d;
  transition: color 0.3s ease;

  &:hover {
    color: #c61b09;
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.about-modal-content > div {
  display: flex;
  flex-direction: column;
  animation: fadeIn 0.5s ease-out forwards;
  opacity: 0;

  &:nth-child(1) { animation-delay: 0.2s; }
  &:nth-child(2) { animation-delay: 0.3s; }
}
</style>