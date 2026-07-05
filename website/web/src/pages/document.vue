<script setup lang="ts">
// 设置页面标题
useHead({
    title: 'Beacon 使用手册 - 轻量级主机监控及 Docker 容器管理工具',
})
</script>

<template>
    <div class="am-document">
        <div class="am-document__header">
            <h1 class="am-document__title">
                Beacon 2.0 使用手册
            </h1>
            <p class="am-document__subtitle">
                轻量级主机监控及 Docker 容器管理工具
            </p>
        </div>

        <div class="am-document__content">
            <!-- 第一部分：安装使用 -->
            <section class="am-section">
                <div class="am-section__header">
                    <h2 class="am-section__title">
                        <span class="am-section__number">01</span>
                        安装使用
                    </h2>
                </div>

                <div class="am-section__content">
                    <div class="am-install">
                        <h3>一键安装</h3>
                        <div class="am-code-block">
                            <pre><code>
# 下载并执行安装脚本
curl -fsSL https://official.beacon.amuluze.com/download/install.sh -o install.sh
sh install.sh

# 非交互安装示例
INSTALL_DIR=/data/beacon \
BEACON_HTTP_PORT=1443 \
BEACON_CONTROL_PORT=17000 \
BEACON_PUBLIC_BASE_URL=http://服务器IP:1443 \
sh install.sh
                            </code></pre>
                        </div>

                        <h3>安装脚本会做什么</h3>
                        <ul>
                            <li>提示输入安装目录，默认目录为 <code>/data/beacon</code>；</li>
                            <li>从官网接口下载 <code>compose.yaml</code> 到安装目录；</li>
                            <li>在同一目录生成并可编辑 <code>.env</code>，用于配置镜像、端口、公开访问地址、Agent 安装 Token 等环境变量；</li>
                            <li>执行 <code>docker compose pull</code> 和 <code>docker compose up -d</code> 启动 Beacon 容器。</li>
                        </ul>

                        <h3>手动安装</h3>
                        <div class="am-code-block">
                            <pre><code>
mkdir -p /data/beacon
cd /data/beacon

curl -fsSL https://official.beacon.amuluze.com/download/compose.yaml -o compose.yaml

cat > .env <<'EOF'
BEACON_IMAGE=registry.cn-hangzhou.aliyuncs.com/amuluze/beacon:latest
BEACON_CONTAINER_NAME=beacon
BEACON_HTTP_PORT=1443
BEACON_CONTROL_PORT=17000
BEACON_DATA_DIR=./data
BEACON_LOG_DIR=./logs
BEACON_DB_NAME=/app/data/beacon
BEACON_PUBLIC_BASE_URL=http://服务器IP:1443
BEACON_AGENT_INSTALL_TOKEN=请替换为随机长字符串
BEACON_AUTH_SIGNING_KEY=请替换为随机长字符串
EOF

docker compose up -d
                            </code></pre>
                        </div>

                        <p>
                            安装目录中的
                            <code>.env</code>
                            是容器部署的主要配置入口，修改后需要在同一目录执行
                            <code>docker compose up -d</code>
                            重新应用。
                        </p>
                        <p>
                            启动 Beacon 前，请确认：
                        </p>
                        <ul>
                            <li>
                                已安装 Docker，并支持 <code>docker compose</code> 命令；
                            </li>
                            <li>
                                <code>BEACON_HTTP_PORT</code> 和 <code>BEACON_CONTROL_PORT</code> 未被占用；
                            </li>
                            <li>
                                <code>BEACON_PUBLIC_BASE_URL</code> 能被浏览器和 Agent 节点访问；
                            </li>
                            <li>
                                <code>BEACON_AGENT_INSTALL_TOKEN</code> 和 <code>BEACON_AUTH_SIGNING_KEY</code> 已替换为随机长字符串。
                            </li>
                        </ul>
                        <h3>运维命令</h3>
                        <div class="am-code-block">
                            <pre><code>
# 进入安装目录
cd /data/beacon

# 启动或更新服务
docker compose up -d

# 查看服务状态和日志
docker compose ps
docker compose logs -f

# 停止服务
docker compose down
                            </code></pre>
                        </div>

                        <h3>访问服务</h3>
                        <p>启动 <code>beacon</code> 容器后，可通过浏览器访问 <code>BEACON_PUBLIC_BASE_URL</code> 或 <code>http://服务器IP:1443</code> 进行管理。</p>
                    </div>
                </div>
            </section>

            <!-- 第二部分：常见问题 -->
            <section class="am-section">
                <div class="am-section__header">
                    <h2 class="am-section__title">
                        <span class="am-section__number">02</span>
                        常见问题
                    </h2>
                </div>

                <div class="am-section__content">
                    <div class="am-faq">
                        <div class="am-faq__item">
                            <h3 class="am-faq__question">
                                Q: 如何备份和恢复配置？
                            </h3>
                            <div class="am-faq__answer">
                                <p><strong>A:</strong> 配置备份方法：</p>
                                <ul>
                                    <li>部署配置位置：<code>/data/beacon/.env</code> 和 <code>/data/beacon/compose.yaml</code></li>
                                    <li>数据目录：<code>/data/beacon/data</code></li>
                                    <li>
                                        备份命令：<code>tar -czf beacon-backup.tar.gz /data/beacon</code>
                                    </li>
                                    <li>恢复时解压到对应目录并执行 <code>docker compose up -d</code></li>
                                </ul>
                            </div>
                        </div>

                        <div class="am-faq__item">
                            <h3 class="am-faq__question">
                                Q: 如何更新到最新版本？
                            </h3>
                            <div class="am-faq__answer">
                                <p><strong>A:</strong> 更新步骤：</p>
                                <ul>
                                    <li>进入安装目录：<code>cd /data/beacon</code></li>
                                    <li>更新前建议备份 <code>.env</code>、<code>compose.yaml</code> 和数据目录</li>
                                    <li>执行 <code>docker compose pull && docker compose up -d</code></li>
                                </ul>
                            </div>
                        </div>

                        <div class="am-faq__item">
                            <h3 class="am-faq__question">
                                Q: 初始化用户有哪些?密码是什么?
                            </h3>
                            <div class="am-faq__answer">
                                <p><strong>A:</strong> 初始化用户及密码:</p>
                                <ul>
                                    <li> 管理用户 admin 密码 admin123 </li>
                                    <li> 普通用户 beacon 密码 123456 </li>
                                    <li>管理员用户能够对普通用户进行管理</li>
                                </ul>
                            </div>
                        </div>

                        <div class="am-faq__item">
                            <h3 class="am-faq__question">
                                Q: 如何获取技术支持？
                            </h3>
                            <div class="am-faq__answer">
                                <p><strong>A:</strong> 获取帮助的途径：</p>
                                <ul>
                                    <li>
                                        查看 <a href="https://github.com/amuluze/amprobe" target="_blank">GitHub 项目</a> 的
                                        Issues
                                    </li>
                                    <li>发送邮件至：314901758@qq.com</li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<style scoped lang="scss">
.am-document {
  max-width: 1200px;
  margin: 0 auto;
  padding: 70px 20px;
  font-family:
    system-ui,
    -apple-system,
    sans-serif;
  line-height: 1.6;
  color: #333;

  @media (max-width: 768px) {
    padding: 20px 16px;
  }

  &__header {
    text-align: center;
    margin-bottom: 60px;
    padding-bottom: 30px;
    border-bottom: 2px solid #f0f0f0;
  }

  &__title {
    font-size: 48px;
    font-weight: 700;
    color: #0f2b5e;
    margin-bottom: 16px;

    @media (max-width: 768px) {
      font-size: 36px;
    }
  }

  &__subtitle {
    font-size: 20px;
    color: #666;
    font-weight: 400;

    @media (max-width: 768px) {
      font-size: 18px;
    }
  }

  &__content {
    display: flex;
    flex-direction: column;
    gap: 80px;

    @media (max-width: 768px) {
      gap: 60px;
    }
  }
}

.am-section {
  &__header {
    margin-bottom: 40px;
  }

  &__title {
    display: flex;
    align-items: center;
    gap: 16px;
    font-size: 36px;
    font-weight: 700;
    color: #0f2b5e;
    margin-bottom: 8px;

    @media (max-width: 768px) {
      font-size: 28px;
    }
  }

  &__number {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 60px;
    height: 60px;
    background: linear-gradient(135deg, #0f2b5e, #1e4a8c);
    color: white;
    border-radius: 50%;
    font-size: 24px;
    font-weight: 700;

    @media (max-width: 768px) {
      width: 50px;
      height: 50px;
      font-size: 20px;
    }
  }

  &__content {
    background: #fff;
    border-radius: 12px;
    padding: 40px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
    border: 1px solid #f0f0f0;

    @media (max-width: 768px) {
      padding: 24px;
    }
  }
}

.am-install {
  h3 {
    font-size: 24px;
    color: #0f2b5e;
    margin: 32px 0 16px 0;
    font-weight: 600;

    &:first-child {
      margin-top: 0;
    }
  }

  h4 {
    font-size: 18px;
    color: #333;
    margin: 20px 0 12px 0;
    font-weight: 600;
  }

  ul,
  ol {
    margin: 16px 0;
    padding-left: 24px;

    li {
      margin: 8px 0;
      color: #555;
    }
  }

  code {
    background: #f8f9fa;
    padding: 2px 6px;
    border-radius: 4px;
    font-family: 'Monaco', 'Consolas', monospace;
    font-size: 14px;
    color: #e83e8c;
  }
}

.am-code-block {
  margin: 24px 0;

  pre {
    background: #2d3748;
    color: #e2e8f0;
    padding: 20px;
    border-radius: 8px;
    overflow-x: auto;
    font-family: 'Monaco', 'Consolas', monospace;
    font-size: 14px;
    line-height: 1.5;
    margin: 12px 0;

    code {
      background: none;
      padding: 0;
      color: inherit;
    }
  }
}

.am-table {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  th,
  td {
    padding: 12px 16px;
    text-align: left;
    border-bottom: 1px solid #f0f0f0;
  }

  th {
    background: #f8f9fa;
    font-weight: 600;
    color: #333;
  }

  td {
    color: #555;

    code {
      background: #f8f9fa;
      padding: 2px 6px;
      border-radius: 4px;
      font-family: 'Monaco', 'Consolas', monospace;
      font-size: 13px;
    }
  }

  tr:last-child {
    td {
      border-bottom: none;
    }
  }
}

.am-faq {
  &__item {
    margin-bottom: 32px;
    padding: 24px;
    background: #f8f9fa;
    border-radius: 8px;
    border-left: 4px solid #0f2b5e;

    &:last-child {
      margin-bottom: 0;
    }
  }

  &__question {
    font-size: 20px;
    color: #0f2b5e;
    margin-bottom: 16px;
    font-weight: 600;

    @media (max-width: 768px) {
      font-size: 18px;
    }
  }

  &__answer {
    color: #555;

    p {
      margin: 12px 0;

      strong {
        color: #333;
      }
    }

    ul {
      margin: 12px 0;
      padding-left: 24px;

      li {
        margin: 8px 0;

        code {
          background: #e9ecef;
          padding: 2px 6px;
          border-radius: 4px;
          font-family: 'Monaco', 'Consolas', monospace;
          font-size: 13px;
          color: #495057;
        }
      }
    }

    a {
      color: #0f2b5e;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }
}
</style>
