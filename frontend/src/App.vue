<template>
  <el-config-provider :locale="elLocale" :message="{ max: 3 }">
    <router-view />
  </el-config-provider>
  <DeveloperTerminal
    v-model:visible="terminalVisible"
    :history="terminalHistory"
    :loading="terminalLoading"
    @submit="processCommand"
  />
</template>

<script setup lang="ts">
import { ref, provide, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElConfigProvider } from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import DeveloperTerminal from '@/components/DeveloperTerminal.vue'
import { useUser } from '@/composables/useUser'

const { locale } = useI18n()

const elLocale = computed(() => (locale.value === 'zh' ? zhCn : en))

const user = useUser()

// Provide refreshUser for permission refresh after role changes
provide('refreshUser', user.refreshUser)

const applyTheme = () => {
  // 仅通过 el-theme-dark class 控制深色模式（与 Element Plus 联动）
  if (localStorage.getItem('isDark') === 'true') {
    document.documentElement.classList.add('el-theme-dark')
  } else {
    document.documentElement.classList.remove('el-theme-dark')
  }
}

// ========== Easter Eggs ==========
const terminalVisible = ref(false)
const terminalHistory = ref<{ type: 'input' | 'output' | 'error' | 'success' | 'art'; content: string }[]>([])
const terminalLoading = ref(false)

const welcomeLines = [
  { type: 'art' as const, content: '\n  ╔══════════════════════════════════════╗' },
  { type: 'art' as const, content: '  ║     DataRegistry Developer Terminal   ║' },
  { type: 'art' as const, content: '  ║     Type "help" for commands          ║' },
  { type: 'art' as const, content: '  ╚══════════════════════════════════════╝\n' },
]

const asciiArts = {
  bigCat: `
      /\\_/\\
     ( o.o )
      > ^ <
     /|   |\\
    (_|   |_)  ~meow!~
   /\\_/\\  /\\_/\\
  ( @.@ ) ( ~.~ )
   > ^ <    > * <
  /|   |\\/|   |\\
 (_|   |_(_|   |_)
    MEOW!
`,
  cat: `
  /\\_/\\
 ( o.o )
  > ^ <
 /|   |\\
(_|   |_)
`,
  coffee: `
       ( (
        ) )
      ........
      |      |]
      \\      /
       \`----'
`,
  love: `
   ******   ******
  ******** ********
 *******************
  *****************
   ***************
    *************
     ***********
      *********
       *******
        *****
         ***
          *
`,
  dino: `
         __
        / _)
   _.--/  /
  /  _ '--'|
  \`\_/ )   |
   / /_|   |
  /  \\__/  |
  \\        /
   \\_____/
`,
}

const matrixChars = '日ﾊﾐﾋｰｳｼﾅﾓﾆｻﾜﾂｵﾘｱﾎﾃﾏｹﾒｴｶｷﾑﾕﾗｾﾈｽﾀﾇﾍ01234567890'.split('')

const processCommand = async (cmd: string) => {
  const c = cmd.trim().toLowerCase()
  terminalHistory.value.push({ type: 'input', content: `$ ${cmd}` })

  switch (c) {
    case '': break
    case 'help':
      terminalHistory.value.push({ type: 'output', content: `
Available commands:
  help       Show this help message
  clear      Clear the terminal
  meow       A cute cat appears!
  nyan       Nyan cat animation 🐱
  coffee     Take a coffee break ☕
  matrix     Enter the Matrix 🌐
  hack       Ironic hacking sequence 💻
  love       Because you matter ❤️
  42         The answer to everything
  whoami     Who are you?
  stats      System statistics
  dinosaur   Rawr! 🦖
  date       Current date and time
  ping       Ping the server
  cat        More cats!
  喵          Secret Chinese cat
  喵呜        Cat says meow 🐱
` })
      break
    case 'clear': terminalHistory.value = []; break
    case 'meow': terminalHistory.value.push({ type: 'art', content: asciiArts.bigCat }); break
    case 'cat': case 'pets': case 'kitten': terminalHistory.value.push({ type: 'art', content: asciiArts.cat }); break
    case '喵': case '喵呜': case 'pet cat':
      terminalHistory.value.push({ type: 'art', content: asciiArts.cat })
      terminalHistory.value.push({ type: 'success', content: '🐱 喵~ (The cat purrs happily!)' })
      break
    case 'nyan':
      const frames = ['  ████\n █▀█▀█░█\n █ █ █ █\n █▄█▄█░█\n  ████', '   ████\n  █▀█▀█░█\n  █ █ █ █\n  █▄█▄█░█\n   ████\n   ███\n   █\n   ████\n   ██████', '    ████\n   █▀█▀█░█\n   █ █ █ █\n   █▄█▄█░█\n    ████\n    ███\n     █\n    ████']
      for (let i = 0; i < 5; i++) {
        terminalHistory.value.push({ type: 'art', content: frames[i % 3] })
        await new Promise(r => setTimeout(r, 400))
      }
      terminalHistory.value.push({ type: 'success', content: '★ ~Nyan Cat mode activated~ ★' })
      break
    case 'coffee': terminalHistory.value.push({ type: 'art', content: asciiArts.coffee }); terminalHistory.value.push({ type: 'success', content: 'Coffee is brewing... ☕' }); break
    case 'love': terminalHistory.value.push({ type: 'art', content: asciiArts.love }); terminalHistory.value.push({ type: 'success', content: 'Love is in the air! <3' }); break
    case '42':
      terminalHistory.value.push({ type: 'output', content: 'The answer to the ultimate question of life, the universe, and everything.' })
      terminalHistory.value.push({ type: 'art', content: '          *  *  *' })
      break
    case 'dinosaur': terminalHistory.value.push({ type: 'art', content: asciiArts.dino }); terminalHistory.value.push({ type: 'success', content: 'RAWR! 🦖 (T-Rex mode)' }); break
    case 'matrix':
      terminalHistory.value.push({ type: 'output', content: 'Entering the Matrix...' })
      for (let i = 0; i < 8; i++) {
        const line = Array.from({ length: 40 }, () => matrixChars[Math.floor(Math.random() * matrixChars.length)]).join('')
        terminalHistory.value.push({ type: 'success', content: line })
        await new Promise(r => setTimeout(r, 150))
      }
      terminalHistory.value.push({ type: 'success', content: 'Wake up, Neo...' })
      break
    case 'hack':
      terminalLoading.value = true
      for (const line of ['Initializing hack.exe...', 'Bypassing firewall... ██████████ 100%', 'Accessing mainframe... ██████████ 100%', 'Decrypting passwords... ██████████ 100%', 'Embedding backdoor... ██████████ 100%', 'Covering tracks... ██████████ 100%', '> ACCESS GRANTED <', '']) {
        terminalHistory.value.push({ type: 'output', content: line })
        await new Promise(r => setTimeout(r, 300))
      }
      terminalLoading.value = false
      break
    case 'whoami':
      try {
        const userData = JSON.parse(localStorage.getItem('user') || '{}')
        terminalHistory.value.push({ type: 'output', content: `User: ${userData.nickname || userData.username || 'admin'}` })
        terminalHistory.value.push({ type: 'output', content: `Role: ${userData.roleName || 'Administrator'}` })
        terminalHistory.value.push({ type: 'output', content: `Theme: ${localStorage.getItem('theme') || 'blue'}` })
      } catch {
        terminalHistory.value.push({ type: 'output', content: 'User: admin' })
        terminalHistory.value.push({ type: 'output', content: 'Role: Administrator' })
      }
      break
    case 'stats':
      terminalHistory.value.push({ type: 'output', content: `\n╔═══════════════════════════════╗\n║     SYSTEM STATISTICS          ║\n╠═══════════════════════════════╣\n║  CPU:    ████████░░  78%       ║\n║  Memory: ██████████  96%       ║\n║  Disk:   ██████░░░░  62%       ║\n║  Uptime: ${new Date().toLocaleTimeString().padEnd(17)} ║\n║  Status: ALL SYSTEMS GO        ║\n╚═══════════════════════════════╝` })
      break
    case 'date': terminalHistory.value.push({ type: 'output', content: new Date().toString() }); break
    case 'ping':
      terminalHistory.value.push({ type: 'output', content: 'PING localhost: 56 bytes of data.' })
      await new Promise(r => setTimeout(r, 500))
      terminalHistory.value.push({ type: 'success', content: '64 bytes from localhost: icmp_seq=1 ttl=64 time=0.42ms' })
      await new Promise(r => setTimeout(r, 300))
      terminalHistory.value.push({ type: 'success', content: '64 bytes from localhost: icmp_seq=2 ttl=64 time=0.38ms' })
      await new Promise(r => setTimeout(r, 300))
      terminalHistory.value.push({ type: 'output', content: '--- localhost ping statistics ---' })
      terminalHistory.value.push({ type: 'output', content: '2 packets transmitted, 2 received, 0% packet loss' })
      break
    case 'sudo': terminalHistory.value.push({ type: 'error', content: 'Nice try. ☕' }); break
    case 'echo': terminalHistory.value.push({ type: 'output', content: '' }); break
    default:
      if (c.startsWith('echo ')) {
        terminalHistory.value.push({ type: 'output', content: cmd.slice(5) })
      } else if (['cd', 'ls', 'pwd'].includes(c)) {
        terminalHistory.value.push({ type: 'output', content: '/mnt/d/data-registry' })
      } else {
        terminalHistory.value.push({ type: 'error', content: `Command not found: ${cmd}. Type "help" for commands.` })
      }
  }
}

// Global keyboard: backtick opens terminal, Konami code
const konamiSeq = ['ArrowUp','ArrowUp','ArrowDown','ArrowDown','ArrowLeft','ArrowRight','ArrowLeft','ArrowRight','KeyB','KeyA']
let konamiIdx = 0
let logoClicks = 0
let logoTimer: ReturnType<typeof setTimeout> | null = null
let logoAnimTimer: ReturnType<typeof setTimeout> | null = null

const printConsoleArt = () => {
  console.log(`%c
    ╔═══════════════════════════════════════════╗
    ║   ██╗     ███████╗ ██████╗ ██████╗██████╗ ║
    ║   ██║     ██╔════╝██╔════╝██╔════╝██╔══██╗║
    ║   ██║     █████╗  ██║     ██║     ██████╔╝║
    ║   ██║     ██╔══╝  ██║     ██║     ██╔══██╗║
    ║   ███████╗███████╗╚██████╗╚██████╗██║  ██║║
    ║   ╚══════╝╚══════╝ ╚═════╝ ╚═════╝╚═╝  ╚═╝║
    ║          🐱 DataRegistry Console 🐱            ║
    ╚═══════════════════════════════════════════╝
  `, 'color: #667eea')
  console.log('%c🐱 meow!', 'color: #f0abfc; font-size: 16px;')
  console.log('%c💡 Try pressing %c` (backtick)%c for a surprise...', 'color: #a0aec0', 'color: #f0abfc; font-weight: bold', 'color: #a0aec0')
  console.log('%c🎮 Konami Code: ↑↑↓↓←→←→BA', 'color: #a0aec0')
}

const handleGlobalKey = (e: KeyboardEvent) => {
  const tag = (e.target as HTMLElement).tagName
  if (['INPUT','TEXTAREA','SELECT'].includes(tag)) return

  // Backtick
  if (e.code === 'Backquote' || e.key === '`') {
    e.preventDefault()
    terminalVisible.value = !terminalVisible.value
    if (terminalVisible.value && terminalHistory.value.length === 0) {
      terminalHistory.value.push(...welcomeLines)
    }
    return
  }

  // Konami
  if (e.code === konamiSeq[konamiIdx]) {
    konamiIdx++
    if (konamiIdx === konamiSeq.length) {
      konamiIdx = 0
      triggerKonami()
    }
  } else {
    konamiIdx = 0
  }
}

let konamiAnimId: number | null = null
const triggerKonami = () => {
  const canvas = document.createElement('canvas')
  canvas.style.cssText = 'position:fixed;top:0;left:0;width:100vw;height:100vh;z-index:99999;background:rgba(0,0,0,0.95)'
  document.body.appendChild(canvas)
  const ctx = canvas.getContext('2d')!
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight

  const particles: { x: number; y: number; vx: number; vy: number; emoji: string; size: number; life: number }[] = []
  const emojis = ['⭐','💫','🌟','✨','🎉','🎊','🔥','💥','🚀','🌈','🦄','🐱','❤️','🧡','💛','💚','💙','💜']

  ctx.textAlign = 'center'
  let lastT = 0
  const anim = (t: number) => {
    if (t - lastT < 16) { konamiAnimId = requestAnimationFrame(anim); return }
    lastT = t

    ctx.fillStyle = 'rgba(0,0,0,0.06)'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    if (t % 100 < 16) {
      for (let i = 0; i < 6; i++) particles.push({
        x: Math.random() * canvas.width,
        y: -20,
        vx: (Math.random() - 0.5) * 3,
        vy: Math.random() * 4 + 2,
        emoji: emojis[Math.floor(Math.random() * emojis.length)],
        size: Math.random() * 20 + 10,
        life: 1
      })
    }

    for (let i = particles.length - 1; i >= 0; i--) {
      const p = particles[i]
      p.x += p.vx; p.y += p.vy; p.vy += 0.04; p.life -= 0.004
      if (p.life <= 0) { particles.splice(i, 1); continue }
      ctx.globalAlpha = p.life
      ctx.font = `${p.size}px serif`
      ctx.fillText(p.emoji, p.x, p.y)
    }
    ctx.globalAlpha = 1

    ctx.fillStyle = '#fff'
    ctx.font = 'bold 28px "DM Sans", sans-serif'
    ctx.fillText('🎉 KONAMI CODE ACTIVATED! 🎉', canvas.width / 2, 60)
    ctx.font = '15px "DM Sans", sans-serif'
    ctx.fillStyle = '#aaa'
    ctx.fillText('You found the secret!', canvas.width / 2, 88)

    if (particles.length > 0) {
      konamiAnimId = requestAnimationFrame(anim)
    } else {
      canvas.remove()
    }
  }

  konamiAnimId = requestAnimationFrame(anim)
  canvas.addEventListener('click', () => {
    if (konamiAnimId) cancelAnimationFrame(konamiAnimId)
    canvas.remove()
  }, { once: true })
  setTimeout(() => { if (konamiAnimId) { cancelAnimationFrame(konamiAnimId); canvas.remove() } }, 8000)
}

// Expose for logo click from layout
const handleLogoClick = () => {
  logoClicks++
  if (logoTimer) clearTimeout(logoTimer)
  logoTimer = setTimeout(() => { logoClicks = 0 }, 1500)
  if (logoClicks >= 7) {
    logoClicks = 0
    if (logoAnimTimer) clearTimeout(logoAnimTimer)
    triggerLogoEaster()
  }
}

const triggerLogoEaster = () => {
  const quotes = [
    '"代码之道，在于折腾！"',
    '"Every bug is a feature in disguise."',
    '"Talk is cheap, show me the code!"',
    '"// TODO: infinity loop"',
    '"rm -rf /: Nice try 😏"',
    '"// This code works, trust me"',
    '"0x13AD: The secret to everything"',
    '"🐱 Cat detected: Meow!"',
  ]
  const bubble = document.createElement('div')
  bubble.style.cssText = `position:fixed;top:70px;left:50%;transform:translateX(-50%);background:linear-gradient(135deg,#667eea,#764ba2);color:white;padding:14px 26px;border-radius:14px;font-family:'DM Sans',sans-serif;font-size:14px;font-weight:600;box-shadow:0 8px 32px rgba(102,126,234,.4);z-index:99999;white-space:nowrap;animation:logoEasterPop .5s cubic-bezier(.34,1.56,.64,1) both;`
  bubble.textContent = quotes[Math.floor(Math.random() * quotes.length)]
  document.body.appendChild(bubble)
  if (!document.getElementById('logo-easter-style')) {
    const s = document.createElement('style')
    s.id = 'logo-easter-style'
    s.textContent = `@keyframes logoEasterPop{0%{opacity:0;transform:translateX(-50%) scale(.5) rotate(-5deg)}100%{opacity:1;transform:translateX(-50%) scale(1) rotate(0deg)}}@keyframes logoEasterFade{0%{opacity:1;transform:translateX(-50%) scale(1)}100%{opacity:0;transform:translateX(-50%) translateY(-20px) scale(.8)}}`
    document.head.appendChild(s)
  }
  logoAnimTimer = setTimeout(() => {
    bubble.style.animation = 'logoEasterFade .4s ease forwards'
    setTimeout(() => bubble.remove(), 400)
  }, 3500)
}

onMounted(() => {
  applyTheme()
  printConsoleArt()
  document.addEventListener('keydown', handleGlobalKey)
  user.initUser()
})

// Provide easter egg handlers for child components
provide('easterEggs', {
  handleLogoClick,
})

onUnmounted(() => {
  if (logoTimer) clearTimeout(logoTimer)
  if (logoAnimTimer) clearTimeout(logoAnimTimer)
})
</script>

<style>
/* 基础重置由 design-system.css 提供 */
html, body, #app {
  width: 100%;
  height: 100%;
}
</style>
