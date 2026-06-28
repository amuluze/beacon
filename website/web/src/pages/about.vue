<script setup lang="ts"></script>

<template>
    <div class="about-container">
        <div class="content">
            <h1>团队故事</h1>
            <el-divider />
            <p>亲爱的朋友，你好～我是阿慕。</p>
            <p>一个35岁的临“退”的程序员。一名见习独立开发者。</p>
            <p>
                35 岁的我，已经是第二份职业的第 7
                个年头了，程序员职业生涯的进度条已经快到头了，下一步的打算是成为一名自由职业者或者独立开发者，能够为自己打工，做自己喜欢的事，同时也能养家糊口。如果你跟我一样面临职业生涯窘境的，想探索和开启新的职业生涯，可以加我微信：Amuluze
                交流。
            </p>
            <p>
                Amprobe
                是我作为一名见习独立开发者的第一款开源产品，算是一个满足自己需求的练手项目，一方面是自己确实有这个需求，方便自己查看自己的云服务器的使用情况，另一方面也是想通过这个项目来锻炼自己独立开发的能力。
            </p>
            <p>
                Amprobe 是一款开源的云服务器监控及 Docker
                容器管理工具，一方面它能够帮助用户实时监控云服务器的使用情况，包括CPU、内存、磁盘等资源的占用情况。方便我们通过Amprobe来监控云服务器的使用情况，指导我们进行相应的优化和调整。
                另一方面，它也能够帮助我们监控和管理部署在云服务上的 docker
                容器，方便我们查看容器的状态，并进行相应的操作，比如启动、停止、重启等。
            </p>
            <p>
                Amprobe 目前还在不断的开发和完善中，如果你有什么需求或者建议，欢迎通过 GitHub Issues
                或者微信联系我（微信号： Amuluze），也可以通过关注公众号了解 Amprobe 的最新进展。
            </p>

            <p>阿慕</p>
            <p>2024年5月26日夜</p>
        </div>
    </div>
</template>

<style scoped lang="scss">
.about-container {
  overflow-x: hidden;
  min-height: 100vh;
  padding: 80px 10%;
  // 科技感背景
  background:
    radial-gradient(circle at center, rgba(0, 102, 255, 0.1) 0%, transparent 70%),
    linear-gradient(135deg, #0a0e27 0%, #1a1f3a 100%);
  position: relative;

  // 添加几何图案背景
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-image:
      linear-gradient(45deg, transparent 40%, rgba(0, 102, 255, 0.05) 50%, transparent 60%),
      linear-gradient(-45deg, transparent 40%, rgba(124, 58, 237, 0.05) 50%, transparent 60%);
    background-size: 60px 60px;
    z-index: 1;
  }

  .content {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px;
    position: relative;
    z-index: 2;
    @include glass-effect(0.05);
    border-radius: 16px;
    @include glow-effect(#0066ff, 0.1);

    // 添加渐变边框
    &::before {
      content: '';
      position: absolute;
      inset: 0;
      padding: 2px;
      background: linear-gradient(135deg, #0066ff, #00d4ff);
      border-radius: inherit;
      mask:
        linear-gradient(#fff 0 0) content-box,
        linear-gradient(#fff 0 0);
      mask-composite: exclude;
      -webkit-mask-composite: xor;
    }

    h1 {
      font-size: 2.5rem;
      color: #e2e8f0;
      margin-bottom: 1.5rem;
      position: relative;
      background: linear-gradient(135deg, #0066ff, #00d4ff);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;

      &::after {
        content: '';
        position: absolute;
        bottom: -0.5rem;
        left: 0;
        width: 60px;
        height: 3px;
        background: linear-gradient(135deg, #0066ff, #00d4ff);
        border-radius: 2px;
      }
    }

    p {
      font-size: 1.1rem;
      line-height: 1.8;
      color: #cbd5e1;
      margin-bottom: 1.2rem;
      text-align: justify;

      &:last-child {
        margin-bottom: 0;
        color: #94a3b8;
        font-style: italic;
      }
    }

    .el-divider {
      margin: 2rem 0;
      background: linear-gradient(90deg, transparent, rgba(0, 102, 255, 0.3), transparent);
      height: 1px;
      border: none;
    }
  }

  @media (max-width: 768px) {
    padding: 40px 5%;

    .content {
      padding: 30px;

      h1 {
        font-size: 2rem;
      }

      p {
        font-size: 1rem;
      }
    }
  }
}
</style>
