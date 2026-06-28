<script setup lang="ts">

</script>

<template>
    <footer class="am-container" role="contentinfo">
        <div class="am-container__top">
            <div class="am-description">
                <div class="am-description__title">
                    Amprobe
                </div>
                <div class="am-description__subtitle">
                    <span>轻量级主机监控及 Docker 容器管理工具</span>
                </div>
            </div>
            <div class="am-description">
                <div class="am-description__title">
                    关于我们
                </div>
                <div class="am-description__subtitle">
                    <NuxtLink to="/about" aria-label="了解团队故事">
                        团队故事
                    </NuxtLink>
                </div>
            </div>
            <div class="am-description">
                <div class="am-description__title">
                    关注我们
                </div>
                <div class="am-description__subtitle">
                    <NuxtLink to="/wechat" aria-label="查看微信二维码">
                        <img src="/images/wechat.png" alt="微信二维码" loading="lazy">
                    </NuxtLink>
                </div>
            </div>
        </div>
        <div class="am-link">
            <span class="am-link__label">友情链接：</span>
            <div class="am-link__items">
                <a
                    href="https://www.amuluze.com" target="_blank" rel="noopener noreferrer"
                    aria-label="访问amuluze网站"
                >amuluze</a>
                <span class="am-link__separator" aria-hidden="true">|</span>
                <a
                    href="https://www.coverpage.one" target="_blank" rel="noopener noreferrer"
                    aria-label="访问coverpage网站"
                >coverpage</a>
            </div>
        </div>
        <el-divider />
        <div class="am-icp">
            Copyright © 2024-present 阿慕 | 京ICP备2023024485号 | 京公网安备11011402053925号
        </div>
    </footer>
</template>

<style scoped lang="scss">
@include b(container) {
  min-height: 200px;
  background: #0f2b5e;
  color: rgba(255, 255, 255, 0.95);
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.12);

  @include e(top) {
    display: flex;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 40px;
    padding: 24px 84px;

    // 添加响应式布局
    @media (max-width: 768px) {
      flex-direction: column;
      gap: 24px;
      padding: 20px 16px;
      text-align: center;
    }

    @media (max-width: 480px) {
      padding: 16px 12px;
      gap: 20px;
    }
  }

  .el-divider {
    margin: 16px 0;
    background-color: rgba(255, 255, 255, 0.15);

    @media (max-width: 768px) {
      margin: 12px 0;
    }
  }
}

@include b(description) {
  @include e(title) {
    font-size: 24px;
    font-weight: 700;
    letter-spacing: 0.5px;
    font-family:
      system-ui,
      -apple-system,
      sans-serif;
    position: relative;
    display: inline-block;

    // 响应式字体大小
    @media (max-width: 768px) {
      font-size: 20px;
    }

    @media (max-width: 480px) {
      font-size: 18px;
    }

    &::after {
      content: '';
      position: absolute;
      bottom: -4px;
      left: 0;
      width: 0;
      height: 2px;
      background-color: rgba(255, 255, 255, 0.8);
      transition: width 0.3s ease;
    }

    &:hover::after {
      width: 100%;
    }
  }

  @include e(subtitle) {
    font-size: 16px;
    opacity: 0.85;
    transition: all 0.3s ease;
    font-family:
      system-ui,
      -apple-system,
      sans-serif;
    margin-top: 8px;

    @media (max-width: 768px) {
      font-size: 14px;
    }

    &:hover {
      opacity: 1;
    }

    a {
      color: inherit;
      text-decoration: none;
      position: relative;
      padding-bottom: 2px;

      &::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        width: 0;
        height: 1px;
        background-color: rgba(255, 255, 255, 0.8);
        transition: width 0.3s ease;
      }

      &:hover {
        text-decoration: none;
      }

      &:hover::after {
        width: 100%;
      }
    }

    img {
      width: 36px;
      height: 36px;
      transition: transform 0.3s ease;
      border-radius: 4px;

      @media (max-width: 480px) {
        width: 32px;
        height: 32px;
      }

      &:hover {
        transform: scale(1.05);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
      }
    }
  }
}

@include b(link) {
  font-size: 16px;
  text-align: left;
  display: flex;
  align-items: center;
  gap: 12px;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  padding: 0 84px;

  @media (max-width: 768px) {
    flex-direction: column;
    text-align: center;
    padding: 0 16px;
    gap: 8px;
  }

  @media (max-width: 480px) {
    padding: 0 12px;
    font-size: 14px;
  }

  @include e(label) {
    white-space: nowrap;
    font-weight: 600;
  }

  @include e(items) {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 16px;

    @media (max-width: 480px) {
      gap: 12px;
    }
  }

  @include e(separator) {
    opacity: 0.5;
    color: rgba(255, 255, 255, 0.6);
  }

  a {
    color: inherit;
    opacity: 0.85;
    transition: all 0.3s ease;
    text-decoration: none;
    position: relative;
    padding-bottom: 2px;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 0;
      height: 1px;
      background-color: rgba(255, 255, 255, 0.8);
      transition: width 0.3s ease;
    }

    &:hover {
      opacity: 1;
      transform: translateY(-1px);
    }

    &:hover::after {
      width: 100%;
    }
  }
}

@include b(icp) {
  text-align: center;
  font-size: 14px;
  opacity: 0.75;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  transition: opacity 0.3s ease;
  padding-bottom: 24px;

  @media (max-width: 768px) {
    padding: 0 16px 8px;
    font-size: 12px;
    line-height: 1.4;
  }

  @media (max-width: 480px) {
    padding: 0 12px 8px;
    font-size: 11px;
  }

  &:hover {
    opacity: 0.9;
  }
}
</style>
