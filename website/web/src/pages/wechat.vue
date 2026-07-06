<script setup lang="ts"></script>

<template>
    <div class="wechat-container">
        <div class="qrcode-card">
            <img src="/images/wechat.jpg" alt="微信公众号二维码" class="qrcode-image">
            <p class="qrcode-text">
                扫码关注获取最新动态
            </p>
        </div>
    </div>
</template>

<style scoped lang="scss">
.wechat-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  // 科技感背景
  background:
    radial-gradient(circle at center, rgba(0, 102, 255, 0.1) 0%, transparent 70%),
    linear-gradient(135deg, #0a0e27 0%, #1a1f3a 100%);
  padding: 2rem;

  .qrcode-card {
    @include glass-effect(0.05);
    padding: 2.5rem;
    border-radius: 16px;
    text-align: center;
    transition: all 0.3s ease;
    @include glow-effect(#0066ff, 0.1);

    &:hover {
      transform: translateY(-5px);
      @include glow-effect(#0066ff, 0.2);
    }

    .qrcode-image {
      width: 280px;
      height: 280px;
      border-radius: 12px;
      margin-bottom: 1.5rem;
      box-shadow: 0 8px 24px rgba(0, 102, 255, 0.2);
      position: relative;

      // 添加扫描线动画效果
      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 2px;
        background: linear-gradient(90deg, transparent, #0066ff, transparent);
        animation: scan 2s linear infinite;
      }
    }

    .qrcode-text {
      font-size: 1.4rem;
      color: #e2e8f0;
      font-weight: 600;
      letter-spacing: 1px;
      margin: 0;
      background: linear-gradient(135deg, #0066ff, #00d4ff);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }
}

@media (max-width: 768px) {
  .wechat-container {
    padding: 1rem;

    .qrcode-card {
      padding: 1.5rem;

      .qrcode-image {
        width: 220px;
        height: 220px;
      }

      .qrcode-text {
        font-size: 1.2rem;
      }
    }
  }
}
</style>
