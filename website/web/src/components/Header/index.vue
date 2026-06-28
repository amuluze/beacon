<script setup lang="ts">
function toDevelopment() {
    window.open('https://amprobe.amuluze.com/', '_blank')
}
</script>

<template>
    <div class="am-container">
        <div class="am-logo">
            <RouterLink to="/" class="am-logo__link">
                <span>Amprobe</span>
            </RouterLink>
        </div>

        <nav class="am-menu">
            <RouterLink to="/document" class="am-menu__link">
                <span>使用手册</span>
            </RouterLink>
            <div class="am-menu__button">
                <el-button type="primary" size="default" @click="toDevelopment">
                    立即体验
                </el-button>
            </div>
        </nav>
    </div>
</template>

<style scoped lang="scss">
@include b(container) {
  position: fixed;
  z-index: 999;
  width: 100%;
  /* 修改为 100% 填满窗口 */
  left: 0;
  right: 0;
  height: 70px;

  display: flex;
  align-items: center;
  justify-content: space-between;

  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  background: linear-gradient(to right, #ffffff, #f8f9fa);
  transition: all 0.3s ease;

  @media (max-width: 768px) {
    padding: 0 20px;
    height: 60px;
  }
}

@include b(logo) {
  font-size: 26px;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  font-weight: 700;
  color: #0f2b5e;
  cursor: pointer;
  transition: transform 0.3s ease;
  padding-left: 84px;

  &:hover {
    transform: scale(1.05);
  }

  @include e(link) {
    text-decoration: none;
    color: inherit;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: -4px;
      left: 0;
      width: 0;
      height: 2px;
      background-color: #1a56db;
      transition: width 0.3s ease;
    }

    &:hover::after {
      width: 100%;
    }
  }

  @media (max-width: 768px) {
    font-size: 22px;
  }
}

@include b(menu) {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-end;
  gap: 28px;
  min-width: 200px;
  padding-right: 84px;

  font-size: 16px;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  font-weight: 600;
  color: #1e293b;

  @include e(link) {
    text-decoration: none;
    color: #4b5563;
    padding: 8px 0;
    position: relative;
    transition: color 0.3s ease;
    display: flex;
    align-items: center;
    margin-right: auto;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      width: 0;
      height: 2px;
      background-color: #1a56db;
      transition: width 0.3s ease;
    }

    &:hover {
      color: #1a56db;
    }

    &:hover::after {
      width: 100%;
    }
  }

  @include e(button) {
    display: flex;
    align-items: center;
    margin-left: auto;

    :deep(.el-button) {
      background-color: #0f2b5e;
      border-color: #0f2b5e;
      transition: all 0.3s ease;
      height: 40px;
      padding: 0 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 50px;
      font-weight: 600;

      &:hover {
        background-color: #1e40af;
        border-color: #1e40af;
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(26, 86, 219, 0.3);
      }
    }
  }

  @media (max-width: 768px) {
    gap: 16px;
    font-size: 14px;
    min-width: 150px;

    @include e(button) {
      :deep(.el-button) {
        height: 36px;
        padding: 0 16px;
        font-size: 14px;
      }
    }
  }
}
</style>
