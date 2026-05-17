import { ref, onMounted, onUnmounted } from 'vue'

// ========== Easter Egg State ==========
const terminalVisible = ref(false)
const konamiTriggered = ref(false)
const logoClickCount = ref(0)
const logoClickTimer = ref<ReturnType<typeof setTimeout> | null>(null)

// ========== Terminal History ==========
export interface TerminalLine {
  type: 'input' | 'output' | 'error' | 'success' | 'art'
  content: string
}

const terminalHistory = ref<TerminalLine[]>([])
const terminalInput = ref('')
const terminalLoading = ref(false)

// ========== ASCII Art ==========
const catArt = `
  /\\_/\\
 ( o.o )
  > ^ <
 /|   |\\
(_|   |_)
`

const bigCatArt = `
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
`

const nyanArt = [
`
  ████
 █▀█▀█░█
 █ █ █ █
 █▄█▄█░█
  ████
`,
`
   ████
  █▀█▀█░█
  █ █ █ █
  █▄█▄█░█
   ████
   ███
   █
   ████
   ██████
  █▀█▀█░█
  █ █ █ █
  █▄█▄█░█
   ████
`,
`
    ████
   █▀█▀█░█
   █ █ █ █
   █▄█▄█░█
    ████
    ███
     █
    ████
   ██████
  █▀█▀█░█
  █ █ █ █
  █▄█▄█░█
   ████
`]

const coffeeArt = `
       ( (
        ) )
      ........
      |      |]
      \\      /
       \`----'
`

const loveArt = `
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
`

const dinoArt = `
         __
        / _)
   _.--/  /
  /  _ '--'|
  \`_/ )   |
   / /_|   |
  /  \__/  |
  \        /
   \_____/
`

const hackArt = [
'Initializing hack.exe...',
'Bypassing firewall... ██████████ 100%',
'Accessing mainframe... ██████████ 100%',
'Decrypting passwords... ██████████ 100%',
'Embedding backdoor... ██████████ 100%',
'Covering tracks... ██████████ 100%',
'> ACCESS GRANTED <',
''
]

const matrixChars = '日ﾊﾐﾋｰｳｼﾅﾓﾆｻﾜﾂｵﾘｱﾎﾃﾏｹﾒｴｶｷﾑﾕﾗｾﾈｽﾀﾇﾍ01234567890'.split('')

const welcomeLines: TerminalLine[] = [
  { type: 'art', content: '\n  ╔══════════════════════════════════════╗' },
  { type: 'art', content: '  ║     DataRegistry Developer Terminal   ║' },
  { type: 'art', content: '  ║     Type "help" for commands          ║' },
  { type: 'art', content: '  ╚══════════════════════════════════════╝\n' },
]

// ========== Command Processor ==========
const processCommand = async (cmd: string) => {
  const c = cmd.trim().toLowerCase()
  terminalHistory.value.push({ type: 'input', content: `$ ${cmd}` })

  switch (c) {
    case '':
      break
    case 'help':
      terminalHistory.value.push({ type: 'output', content: `
Available commands:
  help       - Show this help message
  clear      - Clear the terminal
  meow       - A cute cat appears!
  nyan       - Nyan cat animation
  coffee     - Take a coffee break
  matrix     - Enter the Matrix
  hack       - Ironic hacking sequence
  love       - Because you matter
  42         - The answer to everything
  whoami     - Who are you?
  stats      - System statistics
  dinosaur   - Rawr!
  date       - Current date and time
  echo       - Echo your message
  ping       - Ping the server
  sudo       - Superuser do
  cat        - More cats
  喵          - Secret Chinese cat
  喵呜        - Cat says meow in Chinese
`})
      break
    case 'clear':
      terminalHistory.value = []
      break
    case 'meow':
      terminalHistory.value.push({ type: 'art', content: bigCatArt })
      break
    case 'cat':
    case 'pets':
    case 'kitten':
      terminalHistory.value.push({ type: 'art', content: catArt })
      break
    case '喵':
    case '喵呜':
    case 'pet cat':
      terminalHistory.value.push({ type: 'art', content: catArt })
      terminalHistory.value.push({ type: 'success', content: '🐱 喵~ (The cat purrs happily!)' })
      break
    case 'nyan':
      for (let i = 0; i < 5; i++) {
        terminalHistory.value.push({ type: 'art', content: nyanArt[i % nyanArt.length] })
        await new Promise(r => setTimeout(r, 400))
      }
      terminalHistory.value.push({ type: 'success', content: '★ ~Nyan Cat mode activated~ ★' })
      break
    case 'coffee':
      terminalHistory.value.push({ type: 'art', content: coffeeArt })
      terminalHistory.value.push({ type: 'success', content: 'Coffee is brewing... ☕' })
      break
    case 'love':
      terminalHistory.value.push({ type: 'art', content: loveArt })
      terminalHistory.value.push({ type: 'success', content: 'Love is in the air! <3' })
      break
    case '42':
      terminalHistory.value.push({ type: 'output', content: 'The answer to the ultimate question of life, the universe, and everything.' })
      terminalHistory.value.push({ type: 'art', content: '          *  *  *' })
      break
    case 'dinosaur':
      terminalHistory.value.push({ type: 'art', content: dinoArt })
      terminalHistory.value.push({ type: 'success', content: 'RAWR! 🦖 (T-Rex mode)' })
      break
    case 'matrix':
      terminalHistory.value.push({ type: 'output', content: 'Entering the Matrix...' })
      for (let i = 0; i < 8; i++) {
        const line = Array.from({ length: 40 }, () =>
          matrixChars[Math.floor(Math.random() * matrixChars.length)]
        ).join('')
        terminalHistory.value.push({ type: 'success', content: line })
        await new Promise(r => setTimeout(r, 150))
      }
      terminalHistory.value.push({ type: 'success', content: 'Wake up, Neo...' })
      break
    case 'hack':
      terminalLoading.value = true
      for (const line of hackArt) {
        terminalHistory.value.push({ type: 'output', content: line })
        await new Promise(r => setTimeout(r, 300))
      }
      terminalLoading.value = false
      break
    case 'whoami':
      terminalHistory.value.push({ type: 'output', content: `User: ${localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).nickname || 'admin' : 'admin'}` })
      terminalHistory.value.push({ type: 'output', content: `Role: ${localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!).roleName || 'Administrator' : 'Administrator'}` })
      terminalHistory.value.push({ type: 'output', content: `Theme: ${localStorage.getItem('theme') || 'blue'}` })
      break
    case 'stats':
      terminalHistory.value.push({ type: 'output', content: `
╔═══════════════════════════════╗
║     SYSTEM STATISTICS          ║
╠═══════════════════════════════╣
║  CPU:    ████████░░  78%       ║
║  Memory: ██████████  96%       ║
║  Disk:   ██████░░░░  62%       ║
║  Uptime: ${new Date().toLocaleTimeString().padEnd(17)} ║
║  Status: ALL SYSTEMS GO        ║
╚═══════════════════════════════╝` })
      break
    case 'date':
      terminalHistory.value.push({ type: 'output', content: new Date().toString() })
      break
    case 'ping':
      terminalHistory.value.push({ type: 'output', content: `PING localhost: 56 bytes of data.` })
      await new Promise(r => setTimeout(r, 500))
      terminalHistory.value.push({ type: 'success', content: `64 bytes from localhost: icmp_seq=1 ttl=64 time=0.42ms` })
      await new Promise(r => setTimeout(r, 300))
      terminalHistory.value.push({ type: 'success', content: `64 bytes from localhost: icmp_seq=2 ttl=64 time=0.38ms` })
      await new Promise(r => setTimeout(r, 300))
      terminalHistory.value.push({ type: 'output', content: `--- localhost ping statistics ---` })
      terminalHistory.value.push({ type: 'output', content: `2 packets transmitted, 2 received, 0% packet loss` })
      break
    case 'sudo':
      terminalHistory.value.push({ type: 'error', content: 'Nice try. ☕' })
      break
    case 'echo':
      terminalHistory.value.push({ type: 'output', content: '' })
      break
    default:
      if (c.startsWith('echo ')) {
        terminalHistory.value.push({ type: 'output', content: cmd.slice(5) })
      } else if (c === 'cd' || c === 'ls' || c === 'pwd') {
        terminalHistory.value.push({ type: 'output', content: '/mnt/d/data-registry' })
      } else if (c === 'uname' || c === 'uname -a') {
        terminalHistory.value.push({ type: 'output', content: 'DataRegistry-OS 3.0.0-SAFE x86_64 GNU/Linux' })
      } else if (c.startsWith('cd ')) {
        terminalHistory.value.push({ type: 'output', content: `cd: ${cmd.slice(3)}: No such directory` })
      } else {
        terminalHistory.value.push({ type: 'error', content: `Command not found: ${cmd}. Type "help" for available commands.` })
      }
  }
}

// ========== Konami Code ==========
const konamiSequence = [
  'ArrowUp', 'ArrowUp',
  'ArrowDown', 'ArrowDown',
  'ArrowLeft', 'ArrowRight',
  'ArrowLeft', 'ArrowRight',
  'KeyB', 'KeyA'
]
let konamiIndex = 0

const handleKonamiKey = (e: KeyboardEvent) => {
  if (e.code === konamiSequence[konamiIndex]) {
    konamiIndex++
    if (konamiIndex === konamiSequence.length) {
      konamiIndex = 0
      triggerKonami()
    }
  } else {
    konamiIndex = 0
  }
}

// ========== Konami Effect ==========
const konamiCanvas = ref<HTMLCanvasElement | null>(null)
const konamiActive = ref(false)
let konamiAnimationId: number | null = null

const triggerKonami = () => {
  if (konamiActive.value) return
  konamiActive.value = true

  // Fullscreen canvas overlay
  const canvas = document.createElement('canvas')
  canvas.id = 'konami-canvas'
  canvas.style.cssText = 'position:fixed;top:0;left:0;width:100vw;height:100vh;z-index:99999;background:rgba(0,0,0,0.95)'
  document.body.appendChild(canvas)
  konamiCanvas.value = canvas

  const ctx = canvas.getContext('2d')!
  canvas.width = window.innerWidth
  canvas.height = window.innerHeight

  const particles: { x: number; y: number; vx: number; vy: number; color: string; size: number; life: number }[] = []
  const emojis = ['⭐', '💫', '🌟', '✨', '🎉', '🎊', '🔥', '💥', '🚀', '🌈', '🦄', '🐱', '❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '🤍']

  const spawn = () => {
    for (let i = 0; i < 8; i++) {
      particles.push({
        x: Math.random() * canvas.width,
        y: -20,
        vx: (Math.random() - 0.5) * 3,
        vy: Math.random() * 4 + 2,
        color: emojis[Math.floor(Math.random() * emojis.length)],
        size: Math.random() * 24 + 12,
        life: 1
      })
    }
  }

  let lastTime = 0
  const animate = (time: number) => {
    const delta = time - lastTime
    lastTime = time

    ctx.fillStyle = 'rgba(0,0,0,0.08)'
    ctx.fillRect(0, 0, canvas.width, canvas.height)

    if (delta < 100) spawn()

    for (let i = particles.length - 1; i >= 0; i--) {
      const p = particles[i]
      p.x += p.vx
      p.y += p.vy
      p.vy += 0.05
      p.life -= 0.003

      if (p.life <= 0 || p.y > canvas.height) {
        particles.splice(i, 1)
        continue
      }

      ctx.globalAlpha = p.life
      ctx.font = `${p.size}px serif`
      ctx.fillText(p.color, p.x, p.y)
    }
    ctx.globalAlpha = 1

    // Show message
    ctx.fillStyle = '#fff'
    ctx.font = 'bold 32px "DM Sans", sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText('🎉 KONAMI CODE ACTIVATED! 🎉', canvas.width / 2, 60)
    ctx.font = '16px "DM Sans", sans-serif'
    ctx.fillStyle = '#aaa'
    ctx.fillText('You found the secret!', canvas.width / 2, 90)

    if (particles.length > 0 || konamiActive.value) {
      konamiAnimationId = requestAnimationFrame(animate)
    }
  }

  konamiAnimationId = requestAnimationFrame(animate)

  setTimeout(() => {
    cleanupKonami()
  }, 8000)

  canvas.addEventListener('click', cleanupKonami, { once: true })
}

const cleanupKonami = () => {
  if (konamiAnimationId) {
    cancelAnimationFrame(konamiAnimationId)
    konamiAnimationId = null
  }
  konamiActive.value = false
  const canvas = document.getElementById('konami-canvas')
  if (canvas) canvas.remove()
  konamiCanvas.value = null
}

// ========== Logo Click Easter Egg ==========
const triggerLogoClickEasterEgg = () => {
  logoClickCount.value++
  if (logoClickTimer.value) clearTimeout(logoClickTimer.value)
  logoClickTimer.value = setTimeout(() => {
    logoClickCount.value = 0
  }, 1500) // 1.5s window for 7 clicks

  if (logoClickCount.value >= 7) {
    logoClickCount.value = 0
    if (logoTimer.value) clearTimeout(logoTimer.value)
    triggerLogoAnimation()
  }
}

// ========== Logo Animation ==========
const logoAnimating = ref(false)
const logoTimer = ref<ReturnType<typeof setTimeout> | null>(null)

const triggerLogoAnimation = () => {
  if (logoAnimating.value) return
  logoAnimating.value = true

  // Create floating quote bubble
  const bubble = document.createElement('div')
  bubble.id = 'logo-easter-egg'
  bubble.style.cssText = `
    position: fixed; top: 70px; left: 50%; transform: translateX(-50%);
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white; padding: 16px 28px; border-radius: 16px;
    font-family: 'DM Sans', sans-serif; font-size: 15px; font-weight: 600;
    box-shadow: 0 8px 32px rgba(102, 126, 234, 0.4);
    z-index: 99999; white-space: nowrap;
    animation: logoEasterPop 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
  `
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
  bubble.textContent = quotes[Math.floor(Math.random() * quotes.length)]
  document.body.appendChild(bubble)

  // Add animation keyframes
  if (!document.getElementById('logo-easter-style')) {
    const style = document.createElement('style')
    style.id = 'logo-easter-style'
    style.textContent = `
      @keyframes logoEasterPop {
        0% { opacity: 0; transform: translateX(-50%) scale(0.5) rotate(-5deg); }
        100% { opacity: 1; transform: translateX(-50%) scale(1) rotate(0deg); }
      }
      @keyframes logoEasterFade {
        0% { opacity: 1; transform: translateX(-50%) scale(1); }
        100% { opacity: 0; transform: translateX(-50%) translateY(-20px) scale(0.8); }
      }
    `
    document.head.appendChild(style)
  }

  logoTimer.value = setTimeout(() => {
    bubble.style.animation = 'logoEasterFade 0.4s ease forwards'
    setTimeout(() => bubble.remove(), 400)
    logoAnimating.value = false
  }, 3500)
}

// ========== Console ASCII Art ==========
const printConsoleArt = () => {
  const art = `
%c
    ╔═══════════════════════════════════════════╗
    ║   ██╗     ███████╗ ██████╗ ██████╗██████╗ ║
    ║   ██║     ██╔════╝██╔════╝██╔════╝██╔══██╗║
    ║   ██║     █████╗  ██║     ██║     ██████╔╝║
    ║   ██║     ██╔══╝  ██║     ██║     ██╔══██╗║
    ║   ███████╗███████╗╚██████╗╚██████╗██║  ██║║
    ║   ╚══════╝╚══════╝ ╚═════╝ ╚═════╝╚═╝  ╚═╝║
    ║          %c🐱 DataRegistry Console 🐱%c            ║
    ╚═══════════════════════════════════════════╝

    %c👋 Welcome, developer! You found the secret console!
    %c💡 Try pressing %c\\\`%c (backtick) for a surprise...
    %c🎮 Or try the Konami Code: ↑↑↓↓←→←→BA
  `

  console.log(art,
    'color: #667eea', 'color: #764ba2', 'color: #667eea',
    'color: #a0aec0', 'color: #a0aec0', 'color: #f0abfc', 'color: #a0aec0',
    'color: #a0aec0'
  )

  console.log('%c🐱 meow!', 'color: #f0abfc; font-size: 16px;')
}

// ========== Global Keyboard Handler ==========
const handleGlobalKey = (e: KeyboardEvent) => {
  // Don't trigger in input/textarea/select
  const tag = (e.target as HTMLElement).tagName
  if (['INPUT', 'TEXTAREA', 'SELECT'].includes(tag)) return

  // Backtick to toggle terminal
  if (e.code === 'Backquote' || e.key === '`') {
    e.preventDefault()
    toggleTerminal()
    return
  }

  // Konami code
  handleKonamiKey(e)
}

// ========== Terminal Functions ==========
const toggleTerminal = () => {
  terminalVisible.value = !terminalVisible.value
  if (terminalVisible.value) {
    // Init welcome
    if (terminalHistory.value.length === 0) {
      terminalHistory.value.push(...welcomeLines)
    }
  }
}

const submitCommand = (cmd: string) => {
  if (!cmd.trim()) return
  processCommand(cmd)
  terminalInput.value = ''
}

// ========== Export ==========
export function useEasterEggs() {
  onMounted(() => {
    document.addEventListener('keydown', handleGlobalKey)
    printConsoleArt()
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleGlobalKey)
    if (logoTimer.value) clearTimeout(logoTimer.value)
  })

  return {
    // Terminal
    terminalVisible,
    terminalHistory,
    terminalInput,
    terminalLoading,
    toggleTerminal,
    submitCommand,
    // Konami
    konamiActive,
    triggerKonami,
    cleanupKonami,
    // Logo click
    triggerLogoClickEasterEgg,
    logoAnimating,
  }
}
