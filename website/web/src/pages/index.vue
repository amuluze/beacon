<script setup lang="ts">
import type { StatisticsUpdateParams } from '~/interface/statistics'

import useStore from '@/store'
import { useApi } from '~/composables/useAPI'
import { carouselImages } from '~/config/carousel'

const loading = ref(false)
const store = useStore()
const downloadCount = computed(() => store.statistics.getVisitCount)

// 添加动画状态
const isVisible = ref(false)
const buttonsLoaded = ref(false)

// 添加安装按钮加载状态
const copyingInstall = ref(false)
const downloadingInstall = ref(false)
const officialBaseURL = 'https://official.beacon.amuluze.com'
const installCommand = `curl -fsSL ${officialBaseURL}/download/install.sh -o install.sh && sh install.sh`

onMounted(async () => {
    if (store.statistics.getVisitCount === 0) {
        const { statistics } = useApi()
        const { data } = await statistics.statisticQuery()
        console.log('data', data)
        store.statistics.setStatistic(data.times)
        store.statistics.setID(data.id)
    }
    // 统计数据通过 computed 自动响应，无需手动赋值
    loading.value = true

    // 添加进入动画
    setTimeout(() => {
        isVisible.value = true
    }, 100)

    setTimeout(() => {
        buttonsLoaded.value = true
    }, 500)
})

async function refreshInstallStatistic() {
    const { statistics } = useApi()
    const params: StatisticsUpdateParams = {
        id: store.statistics.id,
    }

    try {
        await statistics.statisticUpdate(params)
        const { data } = await statistics.statisticQuery()
        store.statistics.setStatistic(data.times)
        store.statistics.setID(data.id)
    }
    catch (error) {
        console.error('统计更新失败:', error)
    }
}

async function copyInstallCommand() {
    if (copyingInstall.value)
        return // 防止重复点击

    copyingInstall.value = true
    try {
        await navigator.clipboard.writeText(installCommand)
        ElMessage.success('安装命令已复制')
    }
    catch (error) {
        console.error('复制安装命令失败:', error)
        ElMessage.error('复制失败，请手动复制文档中的安装命令')
    }
    finally {
        setTimeout(() => {
            copyingInstall.value = false
        }, 500)
    }
}

async function downloadInstallScript() {
    if (downloadingInstall.value)
        return // 防止重复点击

    downloadingInstall.value = true
    try {
        window.open('/download/install.sh')
        setTimeout(async () => {
            await refreshInstallStatistic()
        }, 100)
    }
    finally {
        setTimeout(() => {
            downloadingInstall.value = false
        }, 1000)
    }
}

const toGithub = () => window.open('https://github.com/amuluze/amprobe', '_blank')

const description = ref('轻量级主机及容器监控管理工具')

const colShow = ref(true)
function updateColShow() {
    console.log('window', window)
    if (typeof window !== 'undefined') {
        colShow.value = window.innerWidth > 992
    }
}
onMounted(() => {
    updateColShow()
    window.addEventListener('resize', updateColShow)
})
onUnmounted(() => {
    if (typeof window !== 'undefined') {
        window.removeEventListener('resize', updateColShow)
    }
})
// SEO 优化
useHead({
    title: 'Beacon - 开源轻量级主机监控工具',
    meta: [
        {
            name: 'description',
            content: '轻量级主机及容器监控管理工具，实时监控服务器资源使用情况，提供 Docker 容器管理功能',
        },
        {
            name: 'keywords',
            content: '主机监控,Docker管理,容器管理,服务器监控,开源工具',
        },
    ],
})
</script>

<template>
    <div class="am-container">
        <div class="am-introduction" :class="{ 'is-visible': isVisible }">
            <!-- 添加装饰性背景元素 -->
            <div class="am-introduction__bg-decoration">
                <div class="floating-circle circle-1" />
                <div class="floating-circle circle-2" />
                <div class="floating-circle circle-3" />
            </div>

            <div class="am-introduction__content">
                <div class="am-introduction__title">
                    <span class="title-highlight">Beacon</span>
                    <div class="title-subtitle">
                        开源，轻量，现代化
                    </div>
                </div>

                <div class="am-introduction__description">
                    <Marquee :text="description" :time="240" />
                </div>

                <div class="am-introduction__features">
                    <div class="feature-tag">
                        <Icon name="mdi:docker" />
                        <span>Docker 管理</span>
                    </div>
                    <div class="feature-tag">
                        <Icon name="mdi:monitor-dashboard" />
                        <span>实时监控</span>
                    </div>
                    <div class="feature-tag">
                        <Icon name="mdi:cloud-outline" />
                        <span>轻量部署</span>
                    </div>
                </div>

                <div class="am-introduction__download" :class="{ 'buttons-loaded': buttonsLoaded }">
                    <el-row :gutter="20" justify="center" align="middle">
                        <el-col class="am-introduction__download--middle" :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
                            <el-button
                                class="download-btn download-btn--primary"
                                size="large"
                                :loading="copyingInstall"
                                :disabled="copyingInstall"
                                @click="copyInstallCommand"
                            >
                                <Icon v-if="!copyingInstall" name="mdi:content-copy" class="btn-icon" />
                                <span>{{ copyingInstall ? '复制中...' : '复制安装命令' }}</span>
                                <div class="btn-shine" />
                            </el-button>
                        </el-col>
                        <el-col class="am-introduction__download--middle" :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
                            <el-button
                                class="download-btn download-btn--secondary"
                                size="large"
                                :loading="downloadingInstall"
                                :disabled="downloadingInstall"
                                @click="downloadInstallScript"
                            >
                                <Icon v-if="!downloadingInstall" name="mdi:download" class="btn-icon" />
                                <span>{{ downloadingInstall ? '下载中...' : '下载安装脚本' }}</span>
                                <div class="btn-shine" />
                            </el-button>
                        </el-col>
                        <el-col class="am-introduction__download--middle" :xl="8" :lg="8" :md="8" :sm="24" :xs="24">
                            <el-button class="download-btn download-btn--github" size="large" @click="toGithub">
                                <Icon name="mdi:github" class="btn-icon" />
                                <span>GitHub</span>
                                <div class="star-count">
                                    ★ 400+
                                </div>
                                <div class="btn-shine" />
                            </el-button>
                        </el-col>
                    </el-row>
                </div>

                <div v-show="loading" class="am-introduction__statistics">
                    <div class="statistics-content">
                        <Icon name="mdi:download-circle" class="statistics-icon" />
                        <span class="statistics-text">累计获取</span>
                        <span class="statistics-number">{{ downloadCount?.toLocaleString() }}</span>
                        <span class="statistics-text">次</span>
                    </div>
                </div>
            </div>
        </div>

        <!-- 轮播图部分 - 优化尺寸 -->
        <div class="am-carousel">
            <el-carousel :motion-blur="true" trigger="click" :autoplay="false" height="600px">
                <el-carousel-item v-for="(item, index) in carouselImages" :key="index">
                    <div class="carousel-image-container">
                        <el-image fit="contain" :src="item.url" alt="" />
                    </div>
                </el-carousel-item>
            </el-carousel>
        </div>

        <!-- 功能特性部分保持原有结构，样式会在CSS中优化 -->
        <el-row :gutter="20" justify="center" align="middle" class="am-feature am-feature--enhanced bg-gray-1">
            <el-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__image-col">
                <div class="am-feature__image-wrapper">
                    <img src="/images/docker.png" alt="Docker管理" class="am-feature__image">
                    <div class="am-feature__image-overlay">
                        <Icon name="mdi:docker" class="overlay-icon" />
                    </div>
                </div>
            </el-col>
            <el-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__content-col">
                <div class="am-feature__content">
                    <div class="am-feature__header">
                        <div class="am-feature__title">
                            容器管理
                        </div>
                        <div class="am-feature__description">
                            全面的Docker容器管理解决方案
                        </div>
                    </div>

                    <div class="am-feature__list">
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>查看 Docker 版本信息</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>Docker 镜像源设置</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>容器的创建、启动、停止、重启、删除，查看容器日志</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>查看容器运行状态，包括各容器的 CPU、内存使用情况</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>镜像的导入、导出、删除，虚悬镜像清理</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>网络的创建、删除、查看网络状态</span>
                        </div>
                    </div>
                </div>
            </el-col>
        </el-row>

        <!-- 主机管理部分 -->
        <el-row
            :gutter="20" justify="center" align="middle"
            class="am-feature am-feature--enhanced am-feature--reverse"
        >
            <el-col v-show="!colShow" :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__image-col">
                <div class="am-feature__image-wrapper">
                    <img src="/images/audit.png" alt="主机管理" class="am-feature__image">
                    <div class="am-feature__image-overlay">
                        <Icon name="mdi:monitor-dashboard" class="overlay-icon" />
                    </div>
                </div>
            </el-col>
            <el-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__content-col">
                <div class="am-feature__content">
                    <div class="am-feature__header">
                        <div class="am-feature__title">
                            主机管理
                        </div>
                        <div class="am-feature__description">
                            实时监控和管理您的主机系统
                        </div>
                    </div>

                    <div class="am-feature__list">
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>查看主机名称、启动时间、发行版本、内核版本、系统类型</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>查看主机 CPU 使用率、内存使用率、磁盘使用率、网络流量</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>系统时间、系统时区设置、重启、关机</span>
                        </div>
                    </div>
                </div>
            </el-col>
            <el-col v-show="colShow" :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__image-col">
                <div class="am-feature__image-wrapper">
                    <img src="/images/audit.png" alt="主机管理" class="am-feature__image">
                    <div class="am-feature__image-overlay">
                        <Icon name="mdi:monitor-dashboard" class="overlay-icon" />
                    </div>
                </div>
            </el-col>
        </el-row>

        <!-- 用户管理部分 -->
        <el-row :gutter="20" justify="center" align="middle" class="am-feature am-feature--enhanced bg-gray-1">
            <el-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__image-col">
                <div class="am-feature__image-wrapper">
                    <img src="/images/host.png" alt="用户管理" class="am-feature__image">
                    <div class="am-feature__image-overlay">
                        <Icon name="mdi:account-group" class="overlay-icon" />
                    </div>
                </div>
            </el-col>
            <el-col :xl="8" :lg="12" :md="12" :sm="24" :xs="24" class="am-feature__content-col hidden-sm-and-down">
                <div class="am-feature__content">
                    <div class="am-feature__header">
                        <div class="am-feature__title">
                            用户管理
                        </div>
                        <div class="am-feature__description">
                            完善的用户权限和角色管理系统
                        </div>
                    </div>

                    <div class="am-feature__list">
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>用户查看、新建、编辑、删除（admin 管理用户禁止编辑）</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>用户角色列表、角色接口权限查看</span>
                        </div>
                        <div class="am-feature__subtitle am-feature__subtitle--enhanced">
                            <div class="feature-check">
                                <Icon name="mdi:check-circle" />
                            </div>
                            <span>API 接口列表</span>
                        </div>
                    </div>
                </div>
            </el-col>
        </el-row>
    </div>
</template>

<style scoped lang="scss">
@include b(container) {
  width: 100%;
  height: 100%;
  margin-top: 72px;
}

@include b(introduction) {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  // 优化背景渐变，与Header/Footer配色保持一致
  background: linear-gradient(135deg, #0f2b5e 0%, #0f256c 50%, #0f2b5e 100%);
  overflow: hidden;

  // 进入动画
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1);

  &.is-visible {
    opacity: 1;
    transform: translateY(0);
  }

  // 背景装饰
  &__bg-decoration {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    pointer-events: none;
    z-index: 1;
  }

  &__content {
    position: relative;
    z-index: 2;
    text-align: center;
    max-width: 800px;
    padding: 0 20px;
  }

  &__badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    // 使用与Header一致的白色背景和深蓝色边框
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(10px);
    border: 2px solid #1a56db;
    border-radius: 50px;
    padding: 8px 20px;
    color: #0f2b5e;
    font-size: 14px;
    font-weight: 600;
    margin-bottom: 32px;
    animation: fadeInUp 0.8s ease-out 0.2s both;
    box-shadow: 0 4px 12px rgba(15, 43, 94, 0.2);

    .badge-icon {
      color: #ffd700;
      animation: sparkle 2s ease-in-out infinite;
    }
  }

  &__title {
    margin-bottom: 32px;
    animation: fadeInUp 0.8s ease-out 0.4s both;

    .title-highlight {
      display: block;
      color: white;
      font-weight: 800;
      font-size: clamp(48px, 8vw, 80px);
      line-height: 1.1;
      text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
      margin-bottom: 16px;
      // 优化渐变色彩，使用品牌色调
      background: linear-gradient(45deg, #ffffff, #f8f9fa, #e0e7ff);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }

    .title-subtitle {
      color: rgba(255, 255, 255, 0.95);
      font-weight: 400;
      font-size: clamp(18px, 3vw, 24px);
      letter-spacing: 2px;
    }
  }

  &__description {
    margin-bottom: 40px;
    animation: fadeInUp 0.8s ease-out 0.6s both;

    :deep(.marquee-text) {
      color: rgba(255, 255, 255, 0.95);
      font-size: 20px;
      font-weight: 300;
    }
  }

  &__features {
    display: flex;
    justify-content: center;
    gap: 16px;
    margin-bottom: 48px;
    animation: fadeInUp 0.8s ease-out 0.8s both;

    @media (max-width: 768px) {
      flex-direction: column;
      align-items: center;
    }

    .feature-tag {
      display: flex;
      align-items: center;
      gap: 8px;
      // 使用与Header一致的白色背景
      background: rgba(255, 255, 255, 0.9);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(26, 86, 219, 0.3);
      border-radius: 25px;
      padding: 10px 16px;
      color: #0f2b5e;
      font-size: 14px;
      font-weight: 600;
      transition: all 0.3s ease;
      box-shadow: 0 2px 8px rgba(15, 43, 94, 0.1);

      &:hover {
        background: rgba(255, 255, 255, 1);
        transform: translateY(-2px);
        border-color: #1a56db;
        box-shadow: 0 4px 12px rgba(15, 43, 94, 0.2);
      }

      .nuxt-icon {
        font-size: 18px;
        color: #1a56db;
      }
    }
  }

  &__download {
    margin-bottom: 48px;

    &.buttons-loaded {
      .download-btn {
        animation: slideInUp 0.6s ease-out both;

        &:nth-child(1) {
          animation-delay: 0.1s;
        }

        &:nth-child(2) {
          animation-delay: 0.2s;
        }

        &:nth-child(3) {
          animation-delay: 0.3s;
        }
      }
    }

    @include m(middle) {
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 16px;

      @media (max-width: 768px) {
        margin-bottom: 12px;
      }
    }
  }

  &__statistics {
    animation: fadeInUp 0.8s ease-out 1.2s both;

    .statistics-content {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      // 使用与Header一致的白色背景
      background: rgba(255, 255, 255, 0.9);
      backdrop-filter: blur(10px);
      border: 1px solid rgba(26, 86, 219, 0.3);
      border-radius: 50px;
      padding: 16px 24px;
      color: #0f2b5e;
      box-shadow: 0 4px 12px rgba(15, 43, 94, 0.15);

      .statistics-icon {
        font-size: 24px;
        color: #1a56db;
      }

      .statistics-text {
        font-size: 16px;
        font-weight: 500;
        color: #4b5563;
      }

      .statistics-number {
        font-size: 24px;
        font-weight: 700;
        color: #0f2b5e;
        text-shadow: none;
      }
    }
  }
}

// 浮动装饰圆圈 - 科技感重构
.floating-circle {
  position: absolute;
  border-radius: 50%;
  // 使用品牌色调的半透明效果
  background: rgba(26, 86, 219, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);

  &.circle-1 {
    width: 200px;
    height: 200px;
    top: 10%;
    left: 10%;
    animation: float 6s ease-in-out infinite;
  }

  &.circle-2 {
    width: 150px;
    height: 150px;
    top: 60%;
    right: 15%;
    animation: float 8s ease-in-out infinite reverse;
  }

  &.circle-3 {
    width: 100px;
    height: 100px;
    bottom: 20%;
    left: 20%;
    animation: float 10s ease-in-out infinite;
  }
}

// 下载按钮样式 - 科技感重构
.download-btn {
  position: relative;
  width: 200px !important;
  height: 50px;
  border: none !important;
  border-radius: 25px !important;
  font-size: 16px !important;
  font-weight: 600 !important;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  cursor: pointer;

  .btn-icon {
    margin-right: 8px;
    font-size: 18px;
    transition: transform 0.3s ease;
  }

  .btn-shine {
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
    transition: left 0.6s ease;
  }

  &:hover {
    transform: translateY(-3px) !important;
    box-shadow: 0 10px 25px rgba(15, 43, 94, 0.3) !important;

    .btn-icon {
      transform: scale(1.1);
    }

    .btn-shine {
      left: 100%;
    }
  }

  &--primary {
    // 使用品牌主色调
    background: linear-gradient(135deg, #0f2b5e, #1a56db) !important;
    color: white !important;

    &:hover {
      background: linear-gradient(135deg, #1e40af, #2563eb) !important;
    }
  }

  &--secondary {
    // 使用品牌辅助色调
    background: linear-gradient(135deg, #1a56db, #3b82f6) !important;
    color: white !important;

    &:hover {
      background: linear-gradient(135deg, #2563eb, #60a5fa) !important;
    }
  }

  &--github {
    // 保持GitHub风格但调整阴影
    background: linear-gradient(135deg, #24292e, #1a1e22) !important;
    color: white !important;
    display: flex !important;
    align-items: center !important;
    justify-content: center !important;
    gap: 8px !important;

    .star-count {
      background: rgba(255, 255, 255, 0.2);
      border-radius: 12px;
      padding: 2px 8px;
      font-size: 12px;
      font-weight: 700;
    }

    &:hover {
      background: linear-gradient(135deg, #1a1e22, #0d1117) !important;
      box-shadow: 0 10px 25px rgba(36, 41, 46, 0.4) !important;
    }
  }
}

// 功能特性部分样式优化
@include b(feature) {
  padding: 100px 64px;
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(15, 43, 94, 0.2), transparent);
  }

  &.bg-gray-1 {
    background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  }

  &--enhanced {
    &:hover {
      .am-feature__image {
        transform: translateY(-10px) scale(1.05);
      }
    }
  }

  &--reverse {
    background: linear-gradient(135deg, #ffffff 0%, #f8fafc 100%);

    @media (min-width: 1200px) {
      .am-feature__content-col {
        order: 1;
      }

      .am-feature__image-col {
        order: 2;
      }
    }
  }

  @include e(image-col) {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 20px;
  }

  @include e(content-col) {
    display: flex;
    align-items: center;
    padding: 20px;
  }

  @include e(image-wrapper) {
    position: relative;
    display: inline-block;
    border-radius: 20px;
    overflow: hidden;
    box-shadow: 0 20px 40px rgba(15, 43, 94, 0.1);
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      box-shadow: 0 30px 60px rgba(15, 43, 94, 0.15);

      .am-feature__image-overlay {
        opacity: 1;
      }
    }
  }

  @include e(image) {
    height: 400px;
    width: 100%;
    object-fit: cover;
    border-radius: 20px;
    transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);

    @media (max-width: 768px) {
      height: 300px;
    }
  }

  @include e(image-overlay) {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(135deg, rgba(15, 43, 94, 0.8), rgba(26, 86, 219, 0.6));
    display: flex;
    align-items: center;
    justify-content: center;
    opacity: 0;
    transition: all 0.3s ease;
    border-radius: 20px;

    .overlay-icon {
      font-size: 60px;
      color: white;
      animation: pulse 2s ease-in-out infinite;
    }
  }

  @include e(content) {
    max-width: 500px;
    width: 100%;
  }

  @include e(header) {
    margin-bottom: 32px;
    text-align: left;

    @media (max-width: 768px) {
      text-align: center;
      margin-bottom: 24px;
    }
  }

  @include e(icon-badge) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, #0f2b5e, #1a56db);
    border-radius: 16px;
    margin-bottom: 16px;
    box-shadow: 0 8px 24px rgba(15, 43, 94, 0.3);

    .nuxt-icon {
      font-size: 28px;
      color: white;
    }

    &--secondary {
      background: linear-gradient(135deg, #1a56db, #3b82f6);
    }

    &--tertiary {
      background: linear-gradient(135deg, #3b82f6, #60a5fa);
    }
  }

  @include e(title) {
    color: #0f2b5e;
    font-size: 42px;
    font-weight: 800;
    line-height: 1.2;
    margin-bottom: 12px;
    position: relative;
    font-family:
      system-ui,
      -apple-system,
      sans-serif;

    @media (max-width: 768px) {
      font-size: 32px;
    }

    @media (max-width: 480px) {
      font-size: 28px;
    }
  }

  @include e(description) {
    color: #64748b;
    font-size: 18px;
    font-weight: 400;
    line-height: 1.6;
    margin-bottom: 24px;

    @media (max-width: 768px) {
      font-size: 16px;
    }
  }

  @include e(list) {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  @include e(subtitle) {
    color: #475569;
    font-size: 16px;
    line-height: 1.6;
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 16px 20px;
    border-radius: 12px;
    transition: all 0.3s ease;
    border-left: 3px solid transparent;
    background: rgba(255, 255, 255, 0.7);
    backdrop-filter: blur(10px);
    border: 1px solid rgba(15, 43, 94, 0.1);

    &--enhanced {
      &:hover {
        background: rgba(255, 255, 255, 0.9);
        border-left-color: #1a56db;
        transform: translateX(8px);
        box-shadow: 0 4px 12px rgba(15, 43, 94, 0.1);
      }
    }

    .feature-check {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      background: linear-gradient(135deg, #10b981, #059669);
      border-radius: 50%;
      flex-shrink: 0;
      margin-top: 2px;

      .nuxt-icon {
        color: white;
        font-size: 14px;
      }
    }

    span {
      flex: 1;
      font-weight: 500;
    }
  }

  // 响应式设计
  @media (max-width: 768px) {
    padding: 60px 32px;

    @include e(image-col) {
      margin-bottom: 20px;
    }

    @include e(content-col) {
      text-align: center;
    }
  }

  @media (max-width: 480px) {
    padding: 40px 16px;
  }
}

// 动画定义
@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 1;
  }

  50% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}

// 轮播图样式优化
@include b(carousel) {
  background: linear-gradient(135deg, #f8fafc 0%, #f1f5f9 100%);
  padding: 40px 0;

  .el-carousel {
    border-radius: 16px;
    overflow: hidden;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    max-width: 1000px;
    margin: 0 auto;

    @media (max-width: 768px) {
      max-width: 90%;
      border-radius: 12px;
    }
  }

  .carousel-image-container {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    padding: 20px;
    background: white;

    @media (max-width: 768px) {
      padding: 15px;
    }
  }

  .el-image {
    max-height: 500px;
    max-width: 90%;
    border-radius: 12px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;

    &:hover {
      transform: scale(1.02);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
    }

    @media (max-width: 768px) {
      max-height: 350px;
      max-width: 95%;
      border-radius: 8px;
    }

    @media (max-width: 480px) {
      max-height: 280px;
    }
  }

  // 轮播图指示器样式优化
  :deep(.el-carousel__indicators) {
    bottom: 15px;

    .el-carousel__indicator {
      .el-carousel__button {
        background-color: rgba(15, 43, 94, 0.3);
        border-radius: 6px;
        width: 30px;
        height: 6px;
        transition: all 0.3s ease;
      }

      &.is-active .el-carousel__button {
        background-color: #0f2b5e;
        width: 40px;
      }
    }
  }

  // 轮播图箭头样式优化
  :deep(.el-carousel__arrow) {
    background-color: rgba(255, 255, 255, 0.9);
    color: #0f2b5e;
    border: 1px solid rgba(15, 43, 94, 0.1);
    width: 40px;
    height: 40px;
    border-radius: 50%;
    transition: all 0.3s ease;

    &:hover {
      background-color: #0f2b5e;
      color: white;
      transform: scale(1.1);
    }

    @media (max-width: 768px) {
      width: 35px;
      height: 35px;
    }
  }
}

@media (max-width: 768px) {
  @include b(introduction) {
    min-height: 90vh;
    padding: 40px 0;

    &__content {
      padding: 0 16px;
    }

    &__badge {
      font-size: 12px;
      padding: 6px 16px;
    }

    &__features {
      gap: 12px;
      margin-bottom: 32px;

      .feature-tag {
        font-size: 12px;
        padding: 8px 12px;
      }
    }

    &__download {
      margin-bottom: 32px;
    }
  }

  .download-btn {
    width: 180px !important;
    height: 45px;
    font-size: 14px !important;

    .btn-icon {
      font-size: 16px;
    }
  }

  @include b(feature) {
    padding: 40px 20px;

    img {
      height: 250px;
      margin-bottom: 24px;
    }

    @include e(title) {
      font-size: 28px;
      margin-bottom: 24px;
    }

    @include e(subtitle) {
      font-size: 14px;
      padding: 10px 12px;
      margin-top: 12px;
    }
  }
}

@media (max-width: 480px) {
  .floating-circle {
    display: none;
  }

  @include b(introduction) {
    &__features {
      .feature-tag {
        width: 100%;
        justify-content: center;
      }
    }
  }

  .download-btn {
    width: 160px !important;
    height: 42px;
    font-size: 13px !important;
  }
}
</style>
