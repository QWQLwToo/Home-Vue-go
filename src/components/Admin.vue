<template>
  <div class="admin-wrapper">
    <!-- 侧边栏 -->
    <aside class="admin-sidebar">
      <div class="sidebar-header">
        <h1 class="sidebar-logo">
          <i class="fas fa-cog"></i>
          <span>管理面板</span>
        </h1>
      </div>
      
      <nav class="sidebar-nav">
        <button
          v-for="tab in tabs"
          :key="tab"
          @click="handleTabChange(tab)"
          :class="['nav-item', { active: activeTab === tab }]"
        >
          <i :class="getTabIcon(tab)"></i>
          <span>{{ tab }}</span>
        </button>
      </nav>

      <div class="sidebar-footer">
        <div
          class="user-footer-wrapper"
          @mouseenter="userMenuOpen = true"
          @mouseleave="userMenuOpen = false"
        >
          <button class="user-avatar-btn" type="button" title="用户菜单">
            <i class="fas fa-user-circle"></i>
          </button>
          <div class="user-menu" :class="{ open: userMenuOpen }">
            <div class="user-menu-header">
              <i class="fas fa-user"></i>
              <span>用户菜单</span>
            </div>
            <button
              type="button"
              class="user-menu-item"
              @click="showChangePasswordModal = true; userMenuOpen = false"
            >
              <i class="fas fa-key"></i>
              <span>修改密码</span>
            </button>
          </div>
        </div>
        <button @click="logout" class="logout-icon-btn" type="button" title="退出登录">
          <i class="fas fa-sign-out-alt"></i>
        </button>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="admin-main">
      <!-- 顶部栏 -->
      <header class="admin-topbar">
        <h2 class="page-title">
          <i :class="getTabIcon(activeTab)"></i>
          <span>{{ activeTab }}</span>
        </h2>
        <div class="topbar-actions">
          <button v-if="activeTab === '站点配置'" @click="saveSiteConfig" class="btn-save">
            <i class="fas fa-save"></i>
            <span>保存配置</span>
          </button>
          <button v-if="activeTab === '站点管理'" @click="showSiteForm = true" class="btn-primary">
            <i class="fas fa-plus"></i>
            <span>添加站点</span>
          </button>
          <button v-if="activeTab === '联系方式管理'" @click="showContactForm = true" class="btn-primary">
            <i class="fas fa-plus"></i>
            <span>添加联系方式</span>
          </button>
        </div>
      </header>

      <!-- 内容区域 -->
      <div class="admin-content">
        <!-- 加载状态 -->
        <div v-if="loading" class="loading-overlay">
          <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
            <p>加载中...</p>
          </div>
        </div>

        <!-- 仪表盘 -->
        <div v-if="activeTab === '仪表盘' && !loading" class="content-section">
          <Dashboard />
        </div>

        <!-- 站点配置 -->
        <div v-if="activeTab === '站点配置' && !loading" class="content-section">
          <!-- 配置子标签 -->
          <div class="config-subtabs">
            <button
              v-for="subTab in configSubTabs"
              :key="subTab"
              @click="activeConfigSubTab = subTab"
              :class="['subtab-btn', { active: activeConfigSubTab === subTab }]"
            >
              <i :class="getConfigSubTabIcon(subTab)"></i>
              <span>{{ subTab }}</span>
            </button>
          </div>

          <!-- 基础信息 -->
          <div v-if="activeConfigSubTab === '基础信息'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-info-circle"></i> 基础信息</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>站点名称</label>
                    <input v-model="siteConfig.siteName" type="text" placeholder="输入站点名称" />
                    <small class="form-hint">此名称将显示在网站各个位置</small>
                  </div>
                  <div class="form-group">
                    <label>站点URL</label>
                    <input v-model="siteConfig.siteURL" type="url" placeholder="https://example.com" />
                    <small class="form-hint">网站的完整URL地址</small>
                  </div>
                  <div class="form-group">
                    <label>站点描述</label>
                    <textarea v-model="siteConfig.siteDescription" rows="3" placeholder="输入站点描述"></textarea>
                    <small class="form-hint">网站的描述信息，用于SEO和社交媒体分享</small>
                  </div>
                  <div class="form-group">
                    <label>站点关键词</label>
                    <input v-model="siteConfig.siteKeywords" type="text" placeholder="关键词1, 关键词2, 关键词3" />
                    <small class="form-hint">
                      <i class="fas fa-info-circle"></i>
                      多个关键词请使用<strong>英文逗号</strong>（,）分隔，例如：个人主页,Vue3,前端开发
                    </small>
                  </div>
                  <div class="form-group">
                    <label>底部版权年份</label>
                    <div class="color-input-group">
                      <input
                        v-model="siteConfig.footerYearStart"
                        type="text"
                        inputmode="numeric"
                        pattern="[0-9]*"
                        placeholder="起始年份，例如：2020"
                      />
                      <span>-</span>
                      <input
                        v-model="siteConfig.footerYearEnd"
                        type="text"
                        inputmode="numeric"
                        pattern="[0-9]*"
                        placeholder="结束年份（可选）"
                      />
                    </div>
                    <small class="form-hint">
                      <i class="fas fa-info-circle"></i>
                      只填写起始年份时，底部显示为「© 起始年份」；同时填写起止年份且不同，则显示为「© 起始年份-结束年份」。
                    </small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 用户信息 -->
          <div v-if="activeConfigSubTab === '用户信息'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-user"></i> 用户信息</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>用户名</label>
                    <input v-model="siteConfig.userName" type="text" placeholder="输入用户名" />
                    <small class="form-hint">显示在主页上的用户名</small>
                  </div>
                  <div class="form-group">
                    <label>头像</label>
                    <IconSelector
                      v-model="siteConfig.profileImageURL"
                      default-icon-path=""
                    />
                    <small class="form-hint">支持上传、URL或使用默认图标</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 图标配置 -->
          <div v-if="activeConfigSubTab === '图标配置'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-image"></i> 图标配置</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>站点图标</label>
                    <IconSelector v-model="siteConfig.siteIcon" default-icon-path="/favicon.ico" />
                    <small class="form-hint">网站的主要图标</small>
                  </div>
                  <div class="form-group">
                    <label>网页图标（Favicon）</label>
                    <IconSelector v-model="siteConfig.favicon" default-icon-path="/favicon.ico" />
                    <small class="form-hint">浏览器标签页显示的图标</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 备案信息 -->
          <div v-if="activeConfigSubTab === '备案信息'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-certificate"></i> 备案信息</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>ICP备案号</label>
                    <input v-model="siteConfig.icpNumber" type="text" placeholder="暂未填写" />
                    <small class="form-hint">ICP备案号，将显示在网站底部。留空或填写"暂未填写"则不显示</small>
                  </div>
                  <div class="form-group">
                    <label>公安备案号</label>
                    <input v-model="siteConfig.policeNumber" type="text" placeholder="暂未填写" />
                    <small class="form-hint">公安备案号，将显示在网站底部。留空或填写"暂未填写"则不显示</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 前端配置 -->
          <div v-if="activeConfigSubTab === '前端配置'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-code"></i> 前端配置</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>网页标题</label>
                    <input v-model="siteConfig.pageTitle" type="text" placeholder="例如：梦爱吃鱼 - 个人主页" />
                    <small class="form-hint">浏览器标签页显示的标题</small>
                  </div>
                  <div class="form-group">
                    <label>图标库CDN地址</label>
                    <input v-model="siteConfig.iconLibrary" type="text" placeholder="//lib.baomitu.com/font-awesome/6.5.0/css/all.min.css" />
                    <small class="form-hint">Font Awesome图标库地址，用于显示图标</small>
                  </div>
                  <div class="form-group">
                    <label>字体库CDN地址</label>
                    <input v-model="siteConfig.fontLibrary" type="text" placeholder="https://s1.hdslb.com/bfs/static/jinkela/long/font/regular.css" />
                    <small class="form-hint">自定义字体库地址，留空则使用默认字体</small>
                  </div>
                  <div class="form-group">
                    <label class="switch-label">
                      <span>显示停留时间</span>
                      <label class="switch">
                        <input type="checkbox" v-model="siteConfig.showVisitTimer" />
                        <span class="slider"></span>
                      </label>
                    </label>
                    <small class="form-hint">控制主页是否显示停留时间和日期信息</small>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 轮换文本配置 -->
          <div v-if="activeConfigSubTab === '轮换文本'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card rotating-text-card">
                <div class="card-header">
                  <h3><i class="fas fa-text-height"></i> 轮换文本配置</h3>
                  <button @click="saveRotatingTexts" class="btn-save">
                    <i class="fas fa-save"></i>
                    <span>保存文本配置</span>
                  </button>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>轮换文本列表</label>
                    <small class="form-hint">最多可配置8条文本，这些文本将在主页上轮换显示</small>
                    <div class="rotating-texts-list">
                      <div
                        v-for="(text, index) in rotatingTexts"
                        :key="index"
                        class="rotating-text-item"
                      >
                        <div class="text-item-header">
                          <span class="text-index">文本 {{ index + 1 }}</span>
                          <button
                            v-if="rotatingTexts.length > 1"
                            @click="removeRotatingText(index)"
                            class="btn-remove-text"
                            type="button"
                            title="删除"
                          >
                            <i class="fas fa-times"></i>
                          </button>
                        </div>
                        <textarea
                          v-model="rotatingTexts[index]"
                          rows="2"
                          :placeholder="`输入第 ${index + 1} 条轮换文本...`"
                          class="text-input"
                        ></textarea>
                      </div>
                      <button
                        v-if="rotatingTexts.length < 8"
                        @click="addRotatingText"
                        class="btn-add-text"
                        type="button"
                      >
                        <i class="fas fa-plus"></i>
                        添加文本（{{ rotatingTexts.length }}/8）
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 统计配置 -->
          <div v-if="activeConfigSubTab === '统计配置'" class="config-tab-content">
            <div class="config-grid">
              <div class="config-card">
                <div class="card-header">
                  <h3><i class="fas fa-chart-line"></i> 统计配置</h3>
                </div>
                <div class="card-body">
                  <div class="form-group">
                    <label>Umami统计脚本地址</label>
                    <input v-model="siteConfig.umamiScript" type="text" placeholder="umami.is/tongji.js" />
                    <small class="form-hint">留空则不启用统计</small>
                  </div>
                  <div class="form-group">
                    <label>Umami统计网站ID</label>
                    <input v-model="siteConfig.umamiWebsiteId" type="text" placeholder="请填写正确的umami统计id" />
                    <small class="form-hint">需要先配置Umami统计脚本地址</small>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 站点管理 -->
        <div v-if="activeTab === '站点管理' && !loading" class="content-section">
          <div class="data-table-card">
            <div class="table-header">
              <div class="table-info">
                <span>共 {{ sites.length }} 个站点</span>
              </div>
            </div>
            <div class="table-wrapper">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>名称</th>
                    <th>URL</th>
                    <th>图标</th>
                    <th>排序</th>
                    <th class="actions-col">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="site in sites" :key="site.id">
                    <td>{{ site.id }}</td>
                    <td>
                      <div class="cell-content">
                        <i :class="site.icon" v-if="site.icon"></i>
                        <span>{{ site.name }}</span>
                      </div>
                    </td>
                    <td>
                      <a :href="site.url" target="_blank" class="link-cell">
                        {{ site.url }}
                        <i class="fas fa-external-link-alt"></i>
                      </a>
                    </td>
                    <td>
                      <code class="icon-code">{{ site.icon }}</code>
                    </td>
                    <td>
                      <span class="badge">{{ site.sortOrder }}</span>
                    </td>
                    <td class="actions-cell">
                      <button @click="editSite(site)" class="btn-icon btn-edit" title="编辑">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button @click="deleteSite(site.id)" class="btn-icon btn-delete" title="删除">
                        <i class="fas fa-trash"></i>
                      </button>
                    </td>
                  </tr>
                  <tr v-if="sites.length === 0">
                    <td colspan="6" class="empty-state">
                      <i class="fas fa-inbox"></i>
                      <p>暂无站点，点击上方"添加站点"按钮创建</p>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>

        <!-- 联系方式管理 -->
        <div v-if="activeTab === '联系方式管理' && !loading" class="content-section">
          <div class="data-table-card">
            <div class="table-header">
              <div class="table-info">
                <span>共 {{ contacts.length }} 个联系方式</span>
              </div>
            </div>
            <div class="table-wrapper">
              <table>
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>类型</th>
                    <th>图标</th>
                    <th>URL/二维码</th>
                    <th>悬停颜色</th>
                    <th>排序</th>
                    <th class="actions-col">操作</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="contact in contacts" :key="contact.id">
                    <td>{{ contact.id }}</td>
                    <td>
                      <span class="type-badge" :style="{ backgroundColor: contact.hoverColor + '20', color: contact.hoverColor }">
                        {{ contact.type }}
                      </span>
                    </td>
                    <td>
                      <div class="icon-preview">
                        <i :class="contact.icon" :style="{ color: contact.hoverColor }"></i>
                      </div>
                    </td>
                    <td>
                      <div class="url-cell">
                        <a v-if="contact.url" :href="contact.url" target="_blank" class="link-cell">
                          {{ contact.url }}
                          <i class="fas fa-external-link-alt"></i>
                        </a>
                        <span v-else-if="contact.qrCode" class="qr-indicator">
                          <i class="fas fa-qrcode"></i>
                          二维码
                        </span>
                        <span v-else class="text-muted">-</span>
                      </div>
                    </td>
                    <td>
                      <div class="color-preview">
                        <span class="color-dot" :style="{ backgroundColor: contact.hoverColor }"></span>
                        <code>{{ contact.hoverColor }}</code>
                      </div>
                    </td>
                    <td>
                      <span class="badge">{{ contact.sortOrder }}</span>
                    </td>
                    <td class="actions-cell">
                      <button @click="editContact(contact)" class="btn-icon btn-edit" title="编辑">
                        <i class="fas fa-edit"></i>
                      </button>
                      <button @click="deleteContact(contact.id)" class="btn-icon btn-delete" title="删除">
                        <i class="fas fa-trash"></i>
                      </button>
                    </td>
                  </tr>
                  <tr v-if="contacts.length === 0">
                    <td colspan="7" class="empty-state">
                      <i class="fas fa-inbox"></i>
                      <p>暂无联系方式，点击上方"添加联系方式"按钮创建</p>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- 站点表单模态框 -->
    <Transition name="modal">
      <div v-if="showSiteForm" class="modal-overlay" @click="showSiteForm = false">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3>
              <i class="fas" :class="editingSite ? 'fa-edit' : 'fa-plus'"></i>
              {{ editingSite ? '编辑站点' : '添加站点' }}
            </h3>
            <button @click="cancelSiteForm" class="modal-close">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>名称 <span class="required">*</span></label>
              <input v-model="siteForm.name" type="text" placeholder="输入站点名称" required />
            </div>
            <div class="form-group">
              <label>URL <span class="required">*</span></label>
              <input v-model="siteForm.url" type="url" placeholder="https://example.com" required />
            </div>
            <div class="form-group">
              <label>图标类名 <span class="required">*</span></label>
              <div class="icon-input-wrapper">
                <input 
                  v-model="siteForm.icon" 
                  type="text" 
                  placeholder="点击右侧按钮选择图标" 
                  required 
                  readonly
                  class="icon-input"
                />
                <button @click="openIconPicker('site')" type="button" class="icon-picker-btn">
                  <i class="fas fa-icons"></i>
                  选择图标
                </button>
              </div>
              <small class="form-hint">点击"选择图标"按钮从图标库中选择，或手动输入Font Awesome图标类名</small>
              <div v-if="siteForm.icon" class="icon-preview-inline">
                <i :class="siteForm.icon"></i>
                <code>{{ siteForm.icon }}</code>
              </div>
            </div>
            <div class="form-group">
              <label>排序</label>
              <input v-model.number="siteForm.sortOrder" type="number" min="0" />
              <small class="form-hint">数字越小越靠前</small>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="cancelSiteForm" class="btn-secondary">
              <i class="fas fa-times"></i>
              <span>取消</span>
            </button>
            <button @click="saveSite" class="btn-primary">
              <i class="fas fa-save"></i>
              <span>保存</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 联系方式表单模态框 -->
    <Transition name="modal">
      <div v-if="showContactForm" class="modal-overlay" @click="showContactForm = false">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h3>
              <i class="fas" :class="editingContact ? 'fa-edit' : 'fa-plus'"></i>
              {{ editingContact ? '编辑联系方式' : '添加联系方式' }}
            </h3>
            <button @click="cancelContactForm" class="modal-close">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label>类型 <span class="required">*</span></label>
              <select v-model="contactForm.type" required>
                <option value="Email">Email</option>
                <option value="Github">Github</option>
                <option value="支付宝">支付宝</option>
                <option value="微信">微信</option>
                <option value="其他">其他</option>
              </select>
            </div>
            <div class="form-group">
              <label>图标类名 <span class="required">*</span></label>
              <div class="icon-input-wrapper">
                <input 
                  v-model="contactForm.icon" 
                  type="text" 
                  placeholder="点击右侧按钮选择图标" 
                  required 
                  readonly
                  class="icon-input"
                />
                <button @click="openIconPicker('contact')" type="button" class="icon-picker-btn">
                  <i class="fas fa-icons"></i>
                  选择图标
                </button>
              </div>
              <small class="form-hint">点击"选择图标"按钮从图标库中选择，或手动输入Font Awesome图标类名</small>
              <div v-if="contactForm.icon" class="icon-preview-inline">
                <i :class="contactForm.icon"></i>
                <code>{{ contactForm.icon }}</code>
              </div>
            </div>
            <div class="form-group" v-if="contactForm.type === 'Email' || contactForm.type === 'Github'">
              <label>URL <span class="required">*</span></label>
              <input v-model="contactForm.url" type="text" :placeholder="contactForm.type === 'Email' ? 'mailto:example@domain.com' : 'https://...'" required />
              <small class="form-hint" v-if="contactForm.type === 'Email'">Email必须使用mailto:格式</small>
            </div>
            <div class="form-group" v-if="contactForm.type === '支付宝' || contactForm.type === '微信'">
              <label>二维码图片 <span class="required">*</span></label>
              <IconSelector
                v-model="contactForm.qrCode"
                default-icon-path=""
              />
            </div>
            <div class="form-group">
              <label>悬停颜色</label>
              <div class="color-input-group">
                <input v-model="contactForm.hoverColor" type="color" />
                <input v-model="contactForm.hoverColor" type="text" placeholder="#000000" />
              </div>
            </div>
            <div class="form-group">
              <label>排序</label>
              <input v-model.number="contactForm.sortOrder" type="number" min="0" />
              <small class="form-hint">数字越小越靠前</small>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="cancelContactForm" class="btn-secondary">
              <i class="fas fa-times"></i>
              <span>取消</span>
            </button>
            <button @click="saveContact" class="btn-primary">
              <i class="fas fa-save"></i>
              <span>保存</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 图标选择器模态框 -->
    <Transition name="modal">
      <div v-if="showIconPicker" class="modal-overlay" @click.self="showIconPicker = false">
        <div class="modal-content icon-picker-modal" @click.stop>
          <div class="modal-header">
            <h3>
              <i class="fas fa-icons"></i>
              选择图标
            </h3>
            <button @click="showIconPicker = false" class="modal-close">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="modal-body icon-picker-modal-body">
            <IconPicker 
              :modelValue="iconPickerTarget === 'site' ? siteForm.icon : contactForm.icon"
              @close="showIconPicker = false"
              @update:modelValue="handleIconSelect"
            />
          </div>
        </div>
      </div>
    </Transition>

    <!-- 修改密码模态框 -->
    <Transition name="modal">
      <div v-if="showChangePasswordModal" class="modal-overlay" @click.self="closePasswordModal">
        <div class="modal-content password-modal" @click.stop>
          <div class="modal-header">
            <div class="modal-header-content">
              <div class="modal-icon-wrapper">
                <i class="fas fa-key"></i>
              </div>
              <div>
                <h3>修改密码</h3>
                <p class="modal-subtitle">为了您的账户安全，请定期更新密码</p>
              </div>
            </div>
            <button @click="closePasswordModal" class="modal-close">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="modal-body">
            <div class="password-form">
              <div class="form-group">
                <label>
                  <i class="fas fa-lock"></i>
                  当前密码 <span class="required">*</span>
                </label>
                <div class="password-input-wrapper">
                  <input
                    v-model="passwordForm.oldPassword"
                    :type="showOldPassword ? 'text' : 'password'"
                    placeholder="请输入当前密码"
                    required
                    class="password-input"
                    @focus="passwordFormErrors.oldPassword = ''"
                  />
                  <button
                    type="button"
                    @click="showOldPassword = !showOldPassword"
                    class="password-toggle"
                    :title="showOldPassword ? '隐藏密码' : '显示密码'"
                  >
                    <i :class="showOldPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
                  </button>
                </div>
                <small v-if="passwordFormErrors.oldPassword" class="form-error">
                  <i class="fas fa-exclamation-circle"></i>
                  {{ passwordFormErrors.oldPassword }}
                </small>
              </div>

              <div class="form-group">
                <label>
                  <i class="fas fa-key"></i>
                  新密码 <span class="required">*</span>
                </label>
                <div class="password-input-wrapper">
                  <input
                    v-model="passwordForm.newPassword"
                    :type="showNewPassword ? 'text' : 'password'"
                    placeholder="至少8位字符，建议包含字母、数字和特殊字符"
                    required
                    class="password-input"
                    @input="validatePassword"
                    @focus="passwordFormErrors.newPassword = ''"
                  />
                  <button
                    type="button"
                    @click="showNewPassword = !showNewPassword"
                    class="password-toggle"
                    :title="showNewPassword ? '隐藏密码' : '显示密码'"
                  >
                    <i :class="showNewPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
                  </button>
                </div>
                <small class="form-hint">密码至少8位字符，建议包含大小写字母、数字和特殊字符</small>
                <div v-if="passwordStrength" class="password-strength">
                  <div class="strength-bar" :class="passwordStrength.level">
                    <div class="strength-fill" :style="{ width: passwordStrength.percent + '%' }"></div>
                  </div>
                  <div class="strength-info">
                    <span class="strength-text" :class="passwordStrength.level">
                      <i :class="getStrengthIcon(passwordStrength.level)"></i>
                      {{ passwordStrength.text }}
                    </span>
                    <span class="strength-tips" v-if="passwordStrength.tips">
                      {{ passwordStrength.tips }}
                    </span>
                  </div>
                </div>
                <small v-if="passwordFormErrors.newPassword" class="form-error">
                  <i class="fas fa-exclamation-circle"></i>
                  {{ passwordFormErrors.newPassword }}
                </small>
              </div>

              <div class="form-group">
                <label>
                  <i class="fas fa-check-double"></i>
                  确认新密码 <span class="required">*</span>
                </label>
                <div class="password-input-wrapper">
                  <input
                    v-model="passwordForm.confirmPassword"
                    :type="showConfirmPassword ? 'text' : 'password'"
                    placeholder="请再次输入新密码"
                    required
                    class="password-input"
                    :class="{ 'input-error': passwordMismatch && passwordForm.confirmPassword }"
                    @focus="passwordFormErrors.confirmPassword = ''"
                  />
                  <button
                    type="button"
                    @click="showConfirmPassword = !showConfirmPassword"
                    class="password-toggle"
                    :title="showConfirmPassword ? '隐藏密码' : '显示密码'"
                  >
                    <i :class="showConfirmPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
                  </button>
                </div>
                <small v-if="passwordMismatch && passwordForm.confirmPassword" class="form-error">
                  <i class="fas fa-exclamation-circle"></i>
                  两次输入的密码不一致
                </small>
                <small v-else-if="passwordForm.confirmPassword && !passwordMismatch" class="form-success">
                  <i class="fas fa-check-circle"></i>
                  密码匹配
                </small>
                <small v-if="passwordFormErrors.confirmPassword" class="form-error">
                  <i class="fas fa-exclamation-circle"></i>
                  {{ passwordFormErrors.confirmPassword }}
                </small>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button @click="closePasswordModal" class="btn-secondary">
              <i class="fas fa-times"></i>
              <span>取消</span>
            </button>
            <button
              @click="changePassword"
              class="btn-primary"
              :disabled="!canSubmitPassword || changingPassword"
            >
              <i v-if="changingPassword" class="fas fa-spinner fa-spin"></i>
              <i v-else class="fas fa-save"></i>
              <span>{{ changingPassword ? '修改中...' : '确认修改' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminAPI } from '../api'
import IconSelector from './IconSelector.vue'
import IconPicker from './IconPicker.vue'
import Dashboard from './Dashboard.vue'
import { loadAndApplyFrontendConfig } from '../utils/frontendConfig'

const route = useRoute()
const router = useRouter()

const tabs = ['仪表盘', '站点配置', '站点管理', '联系方式管理']
const activeTab = ref('仪表盘')
const configSubTabs = ['基础信息', '用户信息', '图标配置', '备案信息', '前端配置', '轮换文本', '统计配置']
const activeConfigSubTab = ref('基础信息')
const loading = ref(false)
const sites = ref([])
const contacts = ref([])
const siteConfig = ref({
  siteName: '',
  siteURL: '',
  siteIcon: '',
  siteDescription: '',
  siteKeywords: '',
  userName: '',
  profileImageURL: '',
  icpNumber: '',
  policeNumber: '',
  pageTitle: '',
  favicon: '',
  umamiScript: '',
  umamiWebsiteId: '',
  iconLibrary: '',
  fontLibrary: '',
  // 底部版权年份
  footerYearStart: '',
  footerYearEnd: '',
  // 显示停留时间
  showVisitTimer: true,
})

const rotatingTexts = ref([
  "你好鸭，欢迎来到我的主页！！",
  "随时可以联系我，期待与你交流。",
  "愿你历尽千帆，归来仍是少年。",
  "梦想还是要有的，万一实现了呢？",
  "I hope you have a happy day every day."
])

const showSiteForm = ref(false)
const showContactForm = ref(false)
const showChangePasswordModal = ref(false)
const showIconPicker = ref(false)
const iconPickerTarget = ref(null) // 'site' 或 'contact'
const editingSite = ref(null)
const editingContact = ref(null)

// 修改密码相关
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const showOldPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)
const passwordStrength = ref(null)
const changingPassword = ref(false)
const passwordFormErrors = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const siteForm = ref({
  name: '',
  url: '',
  icon: '',
  sortOrder: 0,
})

const contactForm = ref({
  type: 'Email',
  icon: '',
  url: '',
  qrCode: '',
  hoverColor: '#000000',
  sortOrder: 0,
})

// 左下角用户菜单状态
const userMenuOpen = ref(false)

// 打开图标选择器
const openIconPicker = (target) => {
  // target: 'site' 或 'contact'
  iconPickerTarget.value = target
  showIconPicker.value = true
}

// 处理从图标选择器返回的图标
const handleIconSelect = (icon) => {
  if (iconPickerTarget.value === 'site') {
    siteForm.value.icon = icon
  } else if (iconPickerTarget.value === 'contact') {
    contactForm.value.icon = icon
  }
}

const getTabIcon = (tab) => {
  const icons = {
    '仪表盘': 'fas fa-chart-bar',
    '站点配置': 'fas fa-cog',
    '站点管理': 'fas fa-link',
    '联系方式管理': 'fas fa-address-book'
  }
  return icons[tab] || 'fas fa-circle'
}

const getConfigSubTabIcon = (subTab) => {
  const icons = {
    '基础信息': 'fas fa-info-circle',
    '用户信息': 'fas fa-user',
    '图标配置': 'fas fa-image',
    '备案信息': 'fas fa-certificate',
    '前端配置': 'fas fa-code',
    '轮换文本': 'fas fa-text-height',
    '统计配置': 'fas fa-chart-line'
  }
  return icons[subTab] || 'fas fa-circle'
}

const handleTabChange = (tab) => {
  activeTab.value = tab
  // 切换到站点配置时，重置子标签为第一个
  if (tab === '站点配置') {
    activeConfigSubTab.value = '基础信息'
  }
  // 更新路由参数
  updateRoute()
}

const updateRoute = () => {
  const params = {}
  if (activeTab.value !== '仪表盘') {
    params.tab = activeTab.value
  }
  if (activeTab.value === '站点配置' && activeConfigSubTab.value !== '基础信息') {
    params.subTab = activeConfigSubTab.value
  }
  router.replace({ 
    path: '/admin',
    query: params 
  })
}

// 监听子标签变化
watch(activeConfigSubTab, () => {
  if (activeTab.value === '站点配置') {
    updateRoute()
  }
})

const loadData = async () => {
  loading.value = true
  try {
    const [sitesRes, contactsRes, configRes, textsRes] = await Promise.all([
      adminAPI.getSites(),
      adminAPI.getContacts(),
      adminAPI.getSiteConfig(),
      adminAPI.getRotatingTexts().catch(() => ({ data: { texts: [] } })),
    ])
    sites.value = sitesRes.data
    contacts.value = contactsRes.data
    Object.assign(siteConfig.value, configRes.data)
    // 确保 showVisitTimer 正确赋值，处理 false 值
    if (configRes.data.showVisitTimer !== undefined) {
      siteConfig.value.showVisitTimer = Boolean(configRes.data.showVisitTimer)
    }
    console.log('Admin 加载配置 - showVisitTimer:', siteConfig.value.showVisitTimer, '原始值:', configRes.data.showVisitTimer)
    if (textsRes.data?.texts && textsRes.data.texts.length > 0) {
      rotatingTexts.value = textsRes.data.texts
    }
  } catch (error) {
    console.error('加载数据失败:', error)
    showNotification('加载数据失败: ' + (error.response?.data?.error || error.message), 'error')
  } finally {
    loading.value = false
  }
}

const addRotatingText = () => {
  if (rotatingTexts.value.length < 8) {
    rotatingTexts.value.push('')
  }
}

const removeRotatingText = (index) => {
  if (rotatingTexts.value.length > 1) {
    rotatingTexts.value.splice(index, 1)
  }
}

const saveRotatingTexts = async () => {
  try {
    // 过滤空文本
    const filteredTexts = rotatingTexts.value.filter(text => text.trim() !== '')
    if (filteredTexts.length === 0) {
      showNotification('至少需要保留一条文本', 'error')
      return
    }
    
    await adminAPI.updateRotatingTexts(filteredTexts)
    rotatingTexts.value = filteredTexts
    showNotification('轮换文本配置保存成功', 'success')
  } catch (error) {
    showNotification('保存失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

const showNotification = (message, type = 'success') => {
  // 移除已存在的通知
  const existingNotifications = document.querySelectorAll('.top-notification')
  existingNotifications.forEach(n => n.remove())
  
  // 顶部呼吸弹窗通知
  const notification = document.createElement('div')
  notification.className = 'top-notification'
  
  // 根据类型设置样式和图标
  const typeConfig = {
    success: {
      bgColor: 'rgba(76, 175, 80, 0.95)',
      icon: 'fa-check-circle',
      borderColor: '#4caf50'
    },
    error: {
      bgColor: 'rgba(244, 67, 54, 0.95)',
      icon: 'fa-exclamation-circle',
      borderColor: '#f44336'
    },
    warning: {
      bgColor: 'rgba(255, 152, 0, 0.95)',
      icon: 'fa-exclamation-triangle',
      borderColor: '#ff9800'
    },
    info: {
      bgColor: 'rgba(33, 150, 243, 0.95)',
      icon: 'fa-info-circle',
      borderColor: '#2196f3'
    }
  }
  
  const config = typeConfig[type] || typeConfig.success
  
  // 先设置初始样式（完全隐藏）
  notification.style.cssText = `
    position: fixed;
    top: 0;
    left: 50%;
    transform: translate(-50%, -100%);
    background: ${config.bgColor};
    color: white;
    padding: 18px 40px;
    border-radius: 0 0 16px 16px;
    box-shadow: 0 8px 32px rgba(0,0,0,0.3);
    z-index: 10000;
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 15px;
    font-weight: 500;
    border-bottom: 4px solid ${config.borderColor};
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    min-width: 320px;
    max-width: 90vw;
    text-align: center;
    justify-content: center;
    opacity: 0;
    pointer-events: none;
  `
  
  notification.innerHTML = `
    <i class="fas ${config.icon}" style="font-size: 20px;"></i>
    <span style="white-space: nowrap;">${message}</span>
  `
  
  document.body.appendChild(notification)
  
  // 强制重排，确保初始状态生效
  void notification.offsetHeight
  
  // 应用动画 - 使用双重requestAnimationFrame确保动画正确触发
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      notification.style.animation = 'slideDownAndBreathe 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards'
      notification.style.opacity = '1'
      notification.style.pointerEvents = 'auto'
    })
  })
  
  // 3秒后淡出并移除
  setTimeout(() => {
    notification.style.animation = 'slideUpAndFade 0.4s cubic-bezier(0.4, 0, 0.2, 1) forwards'
    setTimeout(() => {
      if (notification.parentNode) {
        notification.remove()
      }
    }, 400)
  }, 3000)
}

const saveSiteConfig = async () => {
  try {
    // 确保年份字段以字符串形式发送，避免后端解析错误
    // 确保 showVisitTimer 明确包含在 payload 中，即使是 false 值
    const payload = {
      ...siteConfig.value,
      footerYearStart:
        siteConfig.value.footerYearStart !== undefined && siteConfig.value.footerYearStart !== null
          ? String(siteConfig.value.footerYearStart).trim()
          : '',
      footerYearEnd:
        siteConfig.value.footerYearEnd !== undefined && siteConfig.value.footerYearEnd !== null
          ? String(siteConfig.value.footerYearEnd).trim()
          : '',
      // 明确包含 showVisitTimer，确保 false 值也能正确发送
      showVisitTimer: siteConfig.value.showVisitTimer !== undefined 
        ? Boolean(siteConfig.value.showVisitTimer) 
        : true,
    }
    
    console.log('保存配置 - showVisitTimer 值:', payload.showVisitTimer, '类型:', typeof payload.showVisitTimer)

    await adminAPI.updateSiteConfig(payload)
    showNotification('配置保存成功', 'success')
    
    // 重新加载配置以获取最新值
    await loadData()
    
    // 触发热重载
    await triggerHotReload()
  } catch (error) {
    showNotification('保存失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

// 热重载功能
const triggerHotReload = async () => {
  try {
    // 通知后端配置已更新
    await adminAPI.notifyConfigUpdate()
    
    // 重新加载前端配置
    await loadAndApplyFrontendConfig()
    
    // 发送消息给其他标签页和主页
    if (window.BroadcastChannel) {
      const channel = new BroadcastChannel('config-update')
      channel.postMessage({ type: 'config-updated' })
      channel.close()
    }
    
    // 如果当前不在管理页面，刷新页面
    if (window.location.pathname !== '/admin') {
      // 延迟刷新，让用户看到通知
      setTimeout(() => {
        window.location.reload()
      }, 1000)
    }
  } catch (error) {
    console.error('热重载失败:', error)
  }
}

// 监听配置更新消息
let broadcastChannel = null
if (window.BroadcastChannel) {
  broadcastChannel = new BroadcastChannel('config-update')
  broadcastChannel.onmessage = (event) => {
    if (event.data.type === 'config-updated') {
      loadAndApplyFrontendConfig()
    }
  }
}

const saveSite = async () => {
  try {
    if (editingSite.value) {
      await adminAPI.updateSite(editingSite.value.id, siteForm.value)
      showNotification('站点更新成功', 'success')
    } else {
      await adminAPI.createSite(siteForm.value)
      showNotification('站点创建成功', 'success')
    }
    await loadData()
    cancelSiteForm()
  } catch (error) {
    showNotification('保存失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

const editSite = (site) => {
  editingSite.value = site
  siteForm.value = {
    name: site.name,
    url: site.url,
    icon: site.icon,
    sortOrder: site.sortOrder,
  }
  showSiteForm.value = true
}

const deleteSite = async (id) => {
  if (!confirm('确定要删除这个站点吗？此操作不可恢复。')) return
  try {
    await adminAPI.deleteSite(id)
    showNotification('站点删除成功', 'success')
    await loadData()
  } catch (error) {
    showNotification('删除失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

const cancelSiteForm = () => {
  showSiteForm.value = false
  editingSite.value = null
  siteForm.value = { name: '', url: '', icon: '', sortOrder: 0 }
}

const saveContact = async () => {
  try {
    if (editingContact.value) {
      await adminAPI.updateContact(editingContact.value.id, contactForm.value)
      showNotification('联系方式更新成功', 'success')
    } else {
      await adminAPI.createContact(contactForm.value)
      showNotification('联系方式创建成功', 'success')
    }
    await loadData()
    cancelContactForm()
  } catch (error) {
    showNotification('保存失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

const editContact = (contact) => {
  editingContact.value = contact
  contactForm.value = { ...contact }
  showContactForm.value = true
}

const deleteContact = async (id) => {
  if (!confirm('确定要删除这个联系方式吗？此操作不可恢复。')) return
  try {
    await adminAPI.deleteContact(id)
    showNotification('联系方式删除成功', 'success')
    await loadData()
  } catch (error) {
    showNotification('删除失败: ' + (error.response?.data?.error || error.message), 'error')
  }
}

const cancelContactForm = () => {
  showContactForm.value = false
  editingContact.value = null
  contactForm.value = {
    type: 'Email',
    icon: '',
    url: '',
    qrCode: '',
    hoverColor: '#000000',
    sortOrder: 0,
  }
}

// 密码强度验证
const validatePassword = () => {
  const password = passwordForm.value.newPassword
  if (!password) {
    passwordStrength.value = null
    passwordFormErrors.value.newPassword = ''
    return
  }

  let strength = 0
  let text = '弱'
  let level = 'weak'
  let tips = ''

  // 长度检查
  if (password.length >= 8) strength++
  if (password.length >= 12) strength++
  
  // 字符类型检查
  const hasLower = /[a-z]/.test(password)
  const hasUpper = /[A-Z]/.test(password)
  const hasNumber = /\d/.test(password)
  const hasSpecial = /[^a-zA-Z0-9]/.test(password)
  
  if (hasLower && hasUpper) strength++
  if (hasNumber) strength++
  if (hasSpecial) strength++

  // 确定强度等级
  if (strength <= 2) {
    text = '弱'
    level = 'weak'
    tips = '建议添加数字、大小写字母或特殊字符'
  } else if (strength <= 4) {
    text = '中'
    level = 'medium'
    tips = '可以添加特殊字符或增加长度以提升安全性'
  } else {
    text = '强'
    level = 'strong'
    tips = '密码强度良好'
  }

  passwordStrength.value = {
    level,
    text,
    percent: Math.min((strength / 6) * 100, 100),
    tips
  }

  // 验证错误信息
  if (password.length < 8) {
    passwordFormErrors.value.newPassword = '密码长度至少8位'
  } else {
    passwordFormErrors.value.newPassword = ''
  }
}

// 获取强度图标
const getStrengthIcon = (level) => {
  const icons = {
    weak: 'fas fa-exclamation-circle',
    medium: 'fas fa-info-circle',
    strong: 'fas fa-check-circle'
  }
  return icons[level] || 'fas fa-circle'
}

// 关闭密码模态框
const closePasswordModal = () => {
  showChangePasswordModal.value = false
  // 重置表单
  passwordForm.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  passwordStrength.value = null
  passwordFormErrors.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
  showOldPassword.value = false
  showNewPassword.value = false
  showConfirmPassword.value = false
}

// 计算密码是否匹配
const passwordMismatch = computed(() => {
  return passwordForm.value.newPassword && 
         passwordForm.value.confirmPassword && 
         passwordForm.value.newPassword !== passwordForm.value.confirmPassword
})

// 是否可以提交密码修改
const canSubmitPassword = computed(() => {
  return passwordForm.value.oldPassword &&
         passwordForm.value.newPassword &&
         passwordForm.value.confirmPassword &&
         passwordForm.value.newPassword.length >= 8 &&
         !passwordMismatch.value
})

// 修改密码
const changePassword = async () => {
  // 重置错误信息
  passwordFormErrors.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }

  // 验证表单
  if (!passwordForm.value.oldPassword) {
    passwordFormErrors.value.oldPassword = '请输入当前密码'
    return
  }

  if (!passwordForm.value.newPassword) {
    passwordFormErrors.value.newPassword = '请输入新密码'
    return
  }

  if (passwordForm.value.newPassword.length < 8) {
    passwordFormErrors.value.newPassword = '密码长度至少8位'
    return
  }

  if (!passwordForm.value.confirmPassword) {
    passwordFormErrors.value.confirmPassword = '请确认新密码'
    return
  }

  if (passwordMismatch.value) {
    passwordFormErrors.value.confirmPassword = '两次输入的密码不一致'
    return
  }

  if (!canSubmitPassword.value) {
    showNotification('请填写完整信息并确保密码符合要求', 'error')
    return
  }

  changingPassword.value = true

  try {
    await adminAPI.changePassword(
      passwordForm.value.oldPassword,
      passwordForm.value.newPassword
    )
    showNotification('密码修改成功', 'success')
    closePasswordModal()
  } catch (error) {
    const errorMessage = error.response?.data?.error || error.message || '密码修改失败'
    showNotification('密码修改失败: ' + errorMessage, 'error')
    
    // 根据错误信息设置相应的错误提示
    if (errorMessage.includes('旧密码') || errorMessage.includes('密码错误')) {
      passwordFormErrors.value.oldPassword = '当前密码错误'
    }
  } finally {
    changingPassword.value = false
  }
}

const logout = () => {
  if (confirm('确定要退出登录吗？')) {
    localStorage.removeItem('token')
    window.location.href = '/login'
  }
}

onMounted(() => {
  // 从路由参数恢复页面状态
  const tabParam = route.query.tab
  const subTabParam = route.query.subTab
  
  if (tabParam && tabs.includes(tabParam)) {
    activeTab.value = tabParam
    if (tabParam === '站点配置' && subTabParam && configSubTabs.includes(subTabParam)) {
      activeConfigSubTab.value = subTabParam
    }
  }
  
  loadData()
})

onUnmounted(() => {
  if (broadcastChannel) {
    broadcastChannel.close()
  }
})
</script>

<style scoped>
/* 主容器 */
.admin-wrapper {
  display: flex;
  min-height: 100vh;
  background: var(--background-color);
  color: var(--text-color);
  font-family: 'HarmonyOS_Regular', 'MiSans', 'HarmonyOS Sans SC', sans-serif;
}

/* 侧边栏 */
.admin-sidebar {
  width: 260px;
  background: rgba(var(--background-color-rgb), 0.98);
  backdrop-filter: blur(10px);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  position: fixed;
  height: 100vh;
  left: 0;
  top: 0;
  z-index: 100;
  box-shadow: 2px 0 12px rgba(0, 0, 0, 0.08);
}

.sidebar-header {
  padding: 24px 20px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-color);
}

.sidebar-logo i {
  font-size: 1.5rem;
  color: var(--hover-link-color);
}

.sidebar-nav {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
}

.nav-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background: transparent;
  border: none;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
  font-size: 0.95rem;
  position: relative;
}

.nav-item i {
  width: 20px;
  text-align: center;
  font-size: 1.1rem;
}

.nav-item:hover {
  background: rgba(var(--background-color-rgb), 0.5);
  color: var(--hover-link-color);
}

.nav-item.active {
  background: linear-gradient(90deg, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.15) 0%, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.08) 100%);
  color: var(--hover-link-color);
  font-weight: 600;
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: var(--hover-link-color);
}

.sidebar-footer {
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
}

.user-footer-wrapper {
  position: relative;
  flex: 1;
}

.user-avatar-btn {
  width: 44px;
  height: 44px;
  border: 2px solid var(--border-color);
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.15) 0%, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1) 100%);
  color: var(--hover-link-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.user-avatar-btn:hover {
  background: linear-gradient(135deg, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.25) 0%, rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.15) 100%);
  border-color: var(--hover-link-color);
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.2);
}

.user-menu {
  position: absolute;
  bottom: 60px;
  left: 0;
  min-width: 200px;
  background: linear-gradient(135deg, rgba(var(--background-color-rgb), 0.98) 0%, rgba(var(--background-color-rgb), 0.95) 100%);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15), 0 4px 8px rgba(0, 0, 0, 0.1);
  opacity: 0;
  visibility: hidden;
  transform: translateX(-10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1000;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  overflow: hidden;
}

.user-menu.open {
  opacity: 1;
  visibility: visible;
  transform: translateX(0);
}

.user-menu-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--text-color);
  background: rgba(var(--background-color-rgb), 0.5);
}

.user-menu-header i {
  color: var(--hover-link-color);
  font-size: 1rem;
}

.user-menu-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  color: var(--text-color);
  cursor: pointer;
  transition: all 0.3s ease;
  text-align: left;
  font-size: 0.95rem;
  border-bottom: 1px solid transparent;
}

.user-menu-item:last-child {
  border-bottom: none;
}

.user-menu-item:hover {
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1);
  color: var(--hover-link-color);
  padding-left: 20px;
}

.user-menu-item i {
  width: 20px;
  text-align: center;
  font-size: 1rem;
  color: var(--hover-link-color);
}

.logout-icon-btn {
  width: 44px;
  height: 44px;
  border: 2px solid rgba(244, 67, 54, 0.3);
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.15) 0%, rgba(244, 67, 54, 0.1) 100%);
  color: #f44336;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 1.3rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
}

.logout-icon-btn:hover {
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.25) 0%, rgba(244, 67, 54, 0.15) 100%);
  border-color: #f44336;
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(244, 67, 54, 0.3);
}

.change-password-btn {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(33, 150, 243, 0.1);
  color: #2196f3;
  border: 1px solid rgba(33, 150, 243, 0.3);
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.95rem;
}

.change-password-btn:hover {
  background: rgba(33, 150, 243, 0.2);
  border-color: rgba(33, 150, 243, 0.5);
}

.logout-btn {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
  border: 1px solid rgba(244, 67, 54, 0.3);
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 0.95rem;
}

.logout-btn:hover {
  background: rgba(244, 67, 54, 0.2);
  border-color: rgba(244, 67, 54, 0.5);
}

/* 主内容区 */
.admin-main {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  overflow: hidden;
}

.admin-topbar {
  background: rgba(var(--background-color-rgb), 0.98);
  backdrop-filter: blur(10px);
  padding: 20px 32px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 10;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.page-title {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-color);
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title i {
  color: var(--hover-link-color);
  font-size: 1.3rem;
}

.topbar-actions {
  display: flex;
  gap: 12px;
}

.admin-content {
  flex: 1;
  padding: 32px;
  overflow-y: auto;
  overflow-x: hidden;
  position: relative;
  height: calc(100vh - 80px);
  max-height: calc(100vh - 80px);
  -webkit-overflow-scrolling: touch;
  overscroll-behavior: contain;
}

.content-section {
  max-width: 1600px;
  margin: 0 auto;
  width: 100%;
}

/* 配置子标签 */
.config-subtabs {
  display: flex;
  gap: 0;
  margin-bottom: 32px;
  border-bottom: 2px solid var(--border-color);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  position: relative;
}

.subtab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 28px;
  background: transparent;
  border: none;
  border-bottom: 3px solid transparent;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.6);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  font-size: 15px;
  font-weight: 400;
  white-space: nowrap;
  font-family: inherit;
}

.subtab-btn::before {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, var(--hover-link-color), #ffd700);
  transform: scaleX(0);
  transform-origin: center;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.subtab-btn:hover {
  color: var(--text-color);
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.05);
}

.subtab-btn:hover::before {
  transform: scaleX(0.5);
}

.subtab-btn.active {
  color: var(--hover-link-color);
  font-weight: 500;
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.08);
}

.subtab-btn.active::before {
  transform: scaleX(1);
}

.subtab-btn i {
  font-size: 0.95rem;
  transition: transform 0.3s ease;
}

.subtab-btn.active i {
  transform: scale(1.1);
  color: var(--hover-link-color);
}

.config-tab-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 配置卡片网格 */
.config-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
  margin-bottom: 32px;
}

.config-card {
  background: rgba(var(--background-color-rgb), 0.85);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
}

.rotating-text-card .card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.config-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12), 0 4px 8px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
  border-color: var(--hover-link-color);
}

.card-header {
  padding: 20px 24px;
  background: rgba(var(--background-color-rgb), 0.8);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-color);
}

.card-header i {
  color: var(--hover-link-color);
}

.card-body {
  padding: 24px;
}

/* 表单组 */
.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  font-size: 0.9rem;
  color: var(--text-color);
}

.required {
  color: #f44336;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  background: rgba(var(--background-color-rgb), 0.8);
  color: var(--text-color);
  font-size: 0.95rem;
  transition: all 0.3s ease;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--hover-link-color);
  box-shadow: 0 0 0 3px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.form-hint {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-top: 6px;
  font-size: 0.85rem;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.6);
  font-style: italic;
  line-height: 1.4;
}

.form-hint i {
  color: var(--hover-link-color);
  margin-top: 2px;
  flex-shrink: 0;
}

.form-hint strong {
  color: var(--hover-link-color);
  font-weight: 600;
}

.color-input-group {
  display: flex;
  gap: 12px;
  align-items: center;
}

.color-input-group input[type="color"] {
  width: 60px;
  height: 40px;
  padding: 2px;
  border-radius: var(--border-radius);
  cursor: pointer;
}

.color-input-group input[type="text"] {
  flex: 1;
}

/* 开关样式 */
.switch-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  font-weight: 500;
  font-size: 0.9rem;
  color: var(--text-color);
}

.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 26px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(var(--text-color-rgb, 51, 51, 51), 0.3);
  transition: 0.3s;
  border-radius: 26px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

.switch input:checked + .slider {
  background-color: var(--hover-link-color);
}

.switch input:checked + .slider:before {
  transform: translateX(24px);
}

.switch input:focus + .slider {
  box-shadow: 0 0 0 3px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.2);
}

/* 按钮样式 */
.btn-primary,
.btn-secondary,
.btn-save {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: var(--border-radius);
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  transition: all 0.3s ease;
  font-family: inherit;
}

.btn-primary,
.btn-save {
  background: var(--hover-link-color);
  color: #333;
  box-shadow: 0 2px 8px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.3);
}

.btn-primary:hover,
.btn-save:hover {
  background: #ffd700;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.4);
}

.btn-secondary {
  background: rgba(var(--text-color-rgb, 51, 51, 51), 0.1);
  color: var(--text-color);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: rgba(var(--text-color-rgb, 51, 51, 51), 0.2);
}

/* 数据表格卡片 */
.data-table-card {
  background: rgba(var(--background-color-rgb), 0.85);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08), 0 2px 4px rgba(0, 0, 0, 0.04);
}

.table-header {
  padding: 20px 24px;
  background: rgba(var(--background-color-rgb), 0.8);
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-info {
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.7);
  font-size: 0.9rem;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: rgba(var(--background-color-rgb), 0.4);
}

th {
  padding: 16px 24px;
  text-align: left;
  font-weight: 600;
  font-size: 0.9rem;
  color: var(--text-color);
  border-bottom: 2px solid var(--border-color);
  white-space: nowrap;
}

.actions-col {
  text-align: center;
  width: 120px;
}

td {
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-color);
}

tbody tr {
  transition: background 0.2s ease;
}

tbody tr:hover {
  background: rgba(var(--background-color-rgb), 0.4);
}

.cell-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cell-content i {
  font-size: 1.2rem;
  color: var(--hover-link-color);
}

.link-cell {
  color: var(--hover-link-color);
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: color 0.3s ease;
  max-width: 300px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.link-cell:hover {
  color: #ffd700;
  text-decoration: underline;
}

.link-cell i {
  font-size: 0.85rem;
  opacity: 0.7;
}

.icon-code {
  background: rgba(var(--background-color-rgb), 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 0.85rem;
  font-family: 'Courier New', monospace;
  color: var(--hover-link-color);
}

.type-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.icon-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  font-size: 1.2rem;
}

.color-preview {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-dot {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid var(--border-color);
  flex-shrink: 0;
}

.color-preview code {
  background: rgba(var(--background-color-rgb), 0.8);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 0.85rem;
  font-family: 'Courier New', monospace;
}

.badge {
  display: inline-block;
  padding: 4px 10px;
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.2);
  color: var(--hover-link-color);
  border-radius: 12px;
  font-size: 0.85rem;
  font-weight: 500;
}

.qr-indicator {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: var(--hover-link-color);
  font-size: 0.9rem;
}

.text-muted {
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
}

.actions-cell {
  text-align: center;
}

.btn-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border: none;
  border-radius: var(--border-radius);
  cursor: pointer;
  transition: all 0.3s ease;
  margin: 0 4px;
  font-size: 0.9rem;
}

.btn-edit {
  background: rgba(33, 150, 243, 0.1);
  color: #2196f3;
}

.btn-edit:hover {
  background: rgba(33, 150, 243, 0.2);
  transform: scale(1.1);
}

.btn-delete {
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
}

.btn-delete:hover {
  background: rgba(244, 67, 54, 0.2);
  transform: scale(1.1);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 16px;
  opacity: 0.3;
}

.empty-state p {
  margin: 0;
  font-size: 0.95rem;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: rgba(var(--background-color-rgb), 0.98);
  border-radius: var(--border-radius);
  border: 1px solid var(--border-color);
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  overflow: hidden;
}

.modal-header {
  padding: 24px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(var(--background-color-rgb), 0.8);
}

.modal-header h3 {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 10px;
  color: var(--text-color);
}

.modal-close {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-color);
  cursor: pointer;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  font-size: 1.1rem;
}

.modal-close:hover {
  background: rgba(var(--text-color-rgb, 51, 51, 51), 0.1);
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: 20px 24px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  background: rgba(var(--background-color-rgb), 0.8);
}

/* 动画 */
@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

@keyframes slideOut {
  from {
    transform: translateX(0);
    opacity: 1;
  }
  to {
    transform: translateX(100%);
    opacity: 0;
  }
}


.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-content,
.modal-leave-active .modal-content {
  transition: transform 0.3s ease, opacity 0.3s ease;
}

.modal-enter-from .modal-content,
.modal-leave-to .modal-content {
  transform: scale(0.9);
  opacity: 0;
}

/* 密码输入相关样式 */
.password-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input {
  width: 100%;
  padding: 12px 48px 12px 16px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.6);
  color: var(--text-color);
  font-size: 15px;
  transition: all 0.3s ease;
  outline: none;
}

.password-input:focus {
  border-color: #007aff;
  background: rgba(var(--background-color-rgb), 0.8);
  box-shadow: 0 0 0 4px rgba(0, 122, 255, 0.1);
}

.password-input.input-error {
  border-color: #f44336;
}

.password-toggle {
  position: absolute;
  right: 12px;
  background: transparent;
  border: none;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.3s ease;
}

.password-toggle:hover {
  color: var(--text-color);
}

.password-strength {
  margin-top: 8px;
}

.strength-bar {
  height: 4px;
  background: rgba(var(--text-color-rgb, 51, 51, 51), 0.1);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 4px;
}

.strength-fill {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
}

.strength-bar.weak .strength-fill {
  background: #f44336;
}

.strength-bar.medium .strength-fill {
  background: #ff9800;
}

.strength-bar.strong .strength-fill {
  background: #4caf50;
}

.strength-text {
  font-size: 12px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.6);
}

.form-error {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  color: #f44336;
  font-size: 12px;
}

.required {
  color: #f44336;
  margin-left: 4px;
}

/* 修改密码模态框特殊样式 */
.password-modal {
  max-width: 520px;
}

.modal-header-content {
  display: flex;
  align-items: center;
  gap: 16px;
}

.modal-icon-wrapper {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #2196f3, #1976d2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
}

.modal-subtitle {
  margin: 4px 0 0 0;
  font-size: 13px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.6);
  font-weight: 400;
}

.password-form {
  padding: 8px 0;
}

.password-form .form-group {
  margin-bottom: 24px;
}

.password-form .form-group label {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-weight: 500;
  color: var(--text-color);
}

.password-form .form-group label i {
  color: #2196f3;
  font-size: 14px;
}

.form-success {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 4px;
  color: #4caf50;
  font-size: 12px;
}

.strength-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 6px;
}

.strength-text {
  font-size: 12px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 4px;
}

.strength-text.weak {
  color: #f44336;
}

.strength-text.medium {
  color: #ff9800;
}

.strength-text.strong {
  color: #4caf50;
}

.strength-tips {
  font-size: 11px;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.5);
}

/* 图标选择器相关样式 */
.icon-input-wrapper {
  display: flex;
  gap: 8px;
  align-items: stretch;
}

.icon-input {
  flex: 1;
  padding: 12px 16px;
  border: 2px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.6);
  color: var(--text-color);
  font-size: 14px;
  transition: all 0.3s ease;
}

.icon-input:focus {
  outline: none;
  border-color: #007aff;
  background: rgba(var(--background-color-rgb), 0.8);
}

.icon-picker-btn {
  padding: 12px 20px;
  border: 2px solid var(--hover-link-color);
  border-radius: 8px;
  background: var(--hover-link-color);
  color: #333;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.icon-picker-btn:hover {
  background: #ffd700;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(255, 204, 0, 0.3);
}

.icon-preview-inline {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 8px;
  padding: 8px 12px;
  background: rgba(var(--background-color-rgb), 0.5);
  border-radius: 6px;
  border: 1px solid var(--border-color);
}

.icon-preview-inline i {
  font-size: 20px;
  color: var(--hover-link-color);
}

.icon-preview-inline code {
  background: rgba(var(--background-color-rgb), 0.8);
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-family: 'Courier New', monospace;
  color: var(--text-color);
}

.icon-picker-modal {
  max-width: 900px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.icon-picker-modal-body {
  padding: 0;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .config-grid {
    grid-template-columns: 1fr;
  }
  
  .admin-sidebar {
    width: 220px;
  }
  
  .admin-main {
    margin-left: 220px;
  }
}

@media (max-width: 768px) {
  .admin-sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }
  
  .admin-sidebar.open {
    transform: translateX(0);
  }
  
  .admin-main {
    margin-left: 0;
  }
  
  .admin-topbar {
    padding: 16px 20px;
  }
  
  .admin-content {
    padding: 20px;
  }
  
  .config-grid {
    grid-template-columns: 1fr;
  }
  
  .table-wrapper {
    overflow-x: scroll;
  }
  
  table {
    min-width: 800px;
  }
}

/* 加载状态 */
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(var(--background-color-rgb), 0.8);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.loading-spinner {
  text-align: center;
  color: var(--text-color);
}

.loading-spinner i {
  font-size: 3rem;
  color: var(--hover-link-color);
  margin-bottom: 16px;
}

.loading-spinner p {
  margin: 0;
  font-size: 1rem;
  color: rgba(var(--text-color-rgb, 51, 51, 51), 0.7);
}

/* 深色模式适配 */
.dark-mode .admin-wrapper {
  filter: brightness(0.7);
}

.dark-mode .config-card,
.dark-mode .data-table-card,
.dark-mode .modal-content {
  background: rgba(var(--background-color-rgb), 0.8);
}

.dark-mode .nav-item.active {
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.15);
}

/* 轮换文本配置样式 */
.rotating-texts-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
}

.rotating-text-item {
  background: rgba(var(--background-color-rgb), 0.5);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  padding: 16px;
  transition: all 0.3s ease;
}

.rotating-text-item:hover {
  border-color: var(--hover-link-color);
  background: rgba(var(--background-color-rgb), 0.7);
}

.text-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.text-index {
  font-weight: 600;
  color: var(--text-color);
  font-size: 0.9rem;
}

.btn-remove-text {
  width: 28px;
  height: 28px;
  border: 1px solid rgba(244, 67, 54, 0.3);
  border-radius: 6px;
  background: rgba(244, 67, 54, 0.1);
  color: #f44336;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  font-size: 0.85rem;
}

.btn-remove-text:hover {
  background: rgba(244, 67, 54, 0.2);
  border-color: #f44336;
  transform: scale(1.1);
}

.text-input {
  width: 100%;
  padding: 10px 14px;
  border: 1.5px solid var(--border-color);
  border-radius: 8px;
  background: rgba(var(--background-color-rgb), 0.8);
  color: var(--text-color);
  font-size: 0.95rem;
  font-family: inherit;
  resize: vertical;
  transition: all 0.3s ease;
}

.text-input:focus {
  outline: none;
  border-color: var(--hover-link-color);
  box-shadow: 0 0 0 3px rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1);
  background: rgba(var(--background-color-rgb), 1);
}

.btn-add-text {
  padding: 12px 20px;
  border: 2px dashed var(--border-color);
  border-radius: 10px;
  background: rgba(var(--background-color-rgb), 0.5);
  color: var(--text-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.3s ease;
  font-size: 0.95rem;
  font-weight: 500;
}

.btn-add-text:hover {
  border-color: var(--hover-link-color);
  background: rgba(var(--hover-link-color-rgb, 255, 204, 0), 0.1);
  color: var(--hover-link-color);
  transform: translateY(-2px);
}
</style>
