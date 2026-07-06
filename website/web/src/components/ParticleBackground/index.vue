<template>
  <div class="particle-background">
    <canvas ref="canvasRef" class="particle-canvas"></canvas>
  </div>
</template>

<script setup lang="ts">
interface Particle {
  x: number
  y: number
  vx: number
  vy: number
  size: number
  opacity: number
  color: string
}

const canvasRef = ref<HTMLCanvasElement>()
const particles = ref<Particle[]>([])
const animationId = ref<number>()

const colors = ['#0066ff', '#00d4ff', '#7c3aed']

const createParticle = (canvas: HTMLCanvasElement): Particle => {
  return {
    x: Math.random() * canvas.width,
    y: Math.random() * canvas.height,
    vx: (Math.random() - 0.5) * 0.5,
    vy: (Math.random() - 0.5) * 0.5,
    size: Math.random() * 2 + 1,
    opacity: Math.random() * 0.5 + 0.1,
    color: colors[Math.floor(Math.random() * colors.length)]
  }
}

const updateParticle = (particle: Particle, canvas: HTMLCanvasElement) => {
  particle.x += particle.vx
  particle.y += particle.vy

  if (particle.x < 0 || particle.x > canvas.width) particle.vx *= -1
  if (particle.y < 0 || particle.y > canvas.height) particle.vy *= -1

  particle.opacity += (Math.random() - 0.5) * 0.01
  particle.opacity = Math.max(0.1, Math.min(0.6, particle.opacity))
}

const drawParticle = (ctx: CanvasRenderingContext2D, particle: Particle) => {
  ctx.save()
  ctx.globalAlpha = particle.opacity
  ctx.fillStyle = particle.color
  ctx.beginPath()
  ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2)
  ctx.fill()
  ctx.restore()
}

const drawConnections = (ctx: CanvasRenderingContext2D, particles: Particle[]) => {
  for (let i = 0; i < particles.length; i++) {
    for (let j = i + 1; j < particles.length; j++) {
      const dx = particles[i].x - particles[j].x
      const dy = particles[i].y - particles[j].y
      const distance = Math.sqrt(dx * dx + dy * dy)

      if (distance < 100) {
        ctx.save()
        ctx.globalAlpha = (100 - distance) / 100 * 0.2
        ctx.strokeStyle = '#0066ff'
        ctx.lineWidth = 0.5
        ctx.beginPath()
        ctx.moveTo(particles[i].x, particles[i].y)
        ctx.lineTo(particles[j].x, particles[j].y)
        ctx.stroke()
        ctx.restore()
      }
    }
  }
}

const animate = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  const ctx = canvas.getContext('2d')
  if (!ctx) return

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  particles.value.forEach(particle => {
    updateParticle(particle, canvas)
    drawParticle(ctx, particle)
  })

  drawConnections(ctx, particles.value)

  animationId.value = requestAnimationFrame(animate)
}

const resizeCanvas = () => {
  const canvas = canvasRef.value
  if (!canvas) return

  canvas.width = window.innerWidth
  canvas.height = window.innerHeight

  // 重新创建粒子
  particles.value = []
  const particleCount = Math.floor((canvas.width * canvas.height) / 15000)
  for (let i = 0; i < particleCount; i++) {
    particles.value.push(createParticle(canvas))
  }
}

onMounted(() => {
  resizeCanvas()
  animate()

  window.addEventListener('resize', resizeCanvas)
})

onUnmounted(() => {
  if (animationId.value) {
    cancelAnimationFrame(animationId.value)
  }
  window.removeEventListener('resize', resizeCanvas)
})
</script>

<style scoped lang="scss">
.particle-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.particle-canvas {
  width: 100%;
  height: 100%;
}
</style>
