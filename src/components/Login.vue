<template>
  <div class="login-container">
    <div class="login-background">
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
      </div>
    </div>
    <div class="login-box">
      <div class="login-header">
        <div class="login-icon">
          <i class="fas fa-lock"></i>
        </div>
        <h2>管理员登录</h2>
        <p class="login-subtitle">欢迎回来，请登录您的账户</p>
      </div>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <div class="input-wrapper">
            <i class="fas fa-user input-icon"></i>
            <input 
              v-model="username" 
              type="text" 
              placeholder="请输入用户名" 
              required 
              class="login-input"
            />
          </div>
        </div>
        <div class="form-group">
          <div class="input-wrapper">
            <i class="fas fa-lock input-icon"></i>
            <input 
              v-model="password" 
              type="password" 
              placeholder="请输入密码" 
              required 
              class="login-input"
            />
          </div>
        </div>
        <button type="submit" :disabled="loading" class="login-btn">
          <span v-if="!loading">
            <i class="fas fa-sign-in-alt"></i>
            登录
          </span>
          <span v-else>
            <i class="fas fa-spinner fa-spin"></i>
            登录中...
          </span>
        </button>
        <div v-if="error" class="error-message">
          <i class="fas fa-exclamation-circle"></i>
          {{ error }}
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { login } from '../api'

const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  try {
    const res = await login(username.value, password.value)
    localStorage.setItem('token', res.data.token)
    window.location.href = '/admin'
  } catch (err) {
    error.value = err.response?.data?.error || '登录失败'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: var(--background-color);
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, 
    rgba(0, 122, 255, 0.1) 0%, 
    rgba(255, 204, 0, 0.1) 50%, 
    rgba(0, 122, 255, 0.1) 100%);
  background-size: 200% 200%;
  animation: gradientShift 15s ease infinite;
  z-index: 0;
}

@keyframes gradientShift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  animation: float 20s infinite ease-in-out;
}

.shape-1 {
  width: 300px;
  height: 300px;
  background: var(--hover-link-color);
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 200px;
  height: 200px;
  background: #007aff;
  bottom: -50px;
  right: -50px;
  animation-delay: 5s;
}

.shape-3 {
  width: 150px;
  height: 150px;
  background: var(--hover-link-color);
  top: 50%;
  right: 10%;
  animation-delay: 10s;
}

@keyframes float {
  0%, 100% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(30px, -30px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
}

.login-box {
  position: relative;
  z-index: 1;
  background: rgba(var(--background-color-rgb), 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  padding: 48px 40px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
  width: 100%;
  max-width: 420px;
  animation: fadeInUp 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.login-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, #007aff, var(--hover-link-color));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(0, 122, 255, 0.3);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 4px 15px rgba(0, 122, 255, 0.3);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 6px 20px rgba(0, 122, 255, 0.4);
  }
}

.login-icon i {
  font-size: 28px;
  color: white;
}

.login-header h2 {
  margin: 0 0 8px 0;
  color: var(--text-color);
  font-size: 28px;
  font-weight: 600;
}

.login-subtitle {
  margin: 0;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.6);
  font-size: 14px;
}

.login-form {
  width: 100%;
}

.form-group {
  margin-bottom: 24px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
  font-size: 16px;
  z-index: 1;
  transition: color 0.3s ease;
}

.login-input {
  width: 100%;
  padding: 14px 16px 14px 48px;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  background: rgba(var(--background-color-rgb), 0.6);
  color: var(--text-color);
  font-size: 15px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  outline: none;
}

.login-input::placeholder {
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.4);
}

.login-input:focus {
  border-color: #007aff;
  background: rgba(var(--background-color-rgb), 0.8);
  box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.1);
  transform: translateY(-2px);
}

.login-input:focus + .input-icon,
.login-input:focus ~ .input-icon {
  color: #007aff;
}

.login-btn {
  width: 100%;
  padding: 16px;
  background: linear-gradient(135deg, #007aff, #0056b3);
  color: white;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  font-size: 16px;
  font-weight: 600;
  margin-top: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 15px rgba(0, 122, 255, 0.3);
  position: relative;
  overflow: hidden;
}

.login-btn::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
}

.login-btn:hover::before {
  width: 300px;
  height: 300px;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 122, 255, 0.4);
}

.login-btn:active {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.login-btn span {
  position: relative;
  z-index: 1;
}

.error-message {
  color: #f44336;
  margin-top: 16px;
  text-align: center;
  padding: 12px;
  background: rgba(244, 67, 54, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(244, 67, 54, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 14px;
  animation: shake 0.5s ease;
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-10px);
  }
  75% {
    transform: translateX(10px);
  }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-box {
    padding: 36px 24px;
    margin: 20px;
    max-width: calc(100% - 40px);
  }
  
  .login-header h2 {
    font-size: 24px;
  }
  
  .shape {
    display: none;
  }
}
</style>
