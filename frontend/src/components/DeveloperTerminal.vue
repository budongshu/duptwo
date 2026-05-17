<template>
  <Teleport to="body">
    <transition name="terminal-fade">
      <div v-if="visible" class="terminal-overlay" @click.self="close">
        <div class="terminal-window">
          <!-- Title Bar -->
          <div class="terminal-titlebar">
            <div class="terminal-dots">
              <span class="dot dot--red" @click="close"></span>
              <span class="dot dot--yellow"></span>
              <span class="dot dot--green"></span>
            </div>
            <span class="terminal-title">DataRegistry Terminal</span>
            <button class="terminal-close" @click="close">×</button>
          </div>

          <!-- Output Area -->
          <div ref="outputEl" class="terminal-output">
            <div
              v-for="(line, idx) in history"
              :key="idx"
              class="terminal-line"
              :class="`terminal-line--${line.type}`"
            >
              <pre>{{ line.content }}</pre>
            </div>
            <!-- Loading indicator -->
            <div v-if="loading" class="terminal-loading">
              <span class="loading-dot"></span>
              <span class="loading-dot"></span>
              <span class="loading-dot"></span>
            </div>
          </div>

          <!-- Input Area -->
          <div class="terminal-input-row">
            <span class="terminal-prompt">$</span>
            <input
              ref="inputEl"
              v-model="inputVal"
              class="terminal-input"
              placeholder="Type a command..."
              autocomplete="off"
              spellcheck="false"
              @keydown.enter="handleSubmit"
              @keydown.up="scrollHistoryUp"
              @keydown.down="scrollHistoryDown"
              @keydown.tab.prevent="autoComplete"
            />
          </div>

          <!-- Hint -->
          <div class="terminal-hint">
            Press <kbd>`</kbd> or <kbd>Esc</kbd> to close · Type <code>help</code> for commands
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onMounted } from 'vue'
import type { TerminalLine } from '@/composables/useEasterEggs'

const props = defineProps<{
  visible: boolean
  history: TerminalLine[]
  loading: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void
  (e: 'submit', cmd: string): void
}>()

const outputEl = ref<HTMLElement | null>(null)
const inputEl = ref<HTMLInputElement | null>(null)
const inputVal = ref('')
const historyIndex = ref(-1)

const close = () => {
  emit('update:visible', false)
  inputVal.value = ''
  historyIndex.value = -1
}

// Auto-focus input when terminal opens
watch(() => props.visible, (v) => {
  if (v) {
    nextTick(() => inputEl.value?.focus())
  }
})

// Auto-scroll to bottom
watch(() => props.history.length, () => {
  nextTick(() => {
    if (outputEl.value) {
      outputEl.value.scrollTop = outputEl.value.scrollHeight
    }
  })
})

const handleSubmit = () => {
  if (!inputVal.value.trim()) return
  emit('submit', inputVal.value)
  historyIndex.value = -1
  inputVal.value = ''
}

const scrollHistoryUp = () => {
  // Simple command history navigation could be added here
  historyIndex.value++
}

const scrollHistoryDown = () => {
  historyIndex.value--
  if (historyIndex.value < 0) historyIndex.value = -1
}

const autoComplete = () => {
  const commands = ['help', 'clear', 'meow', 'nyan', 'coffee', 'matrix', 'hack', 'love', '42', 'whoami', 'stats', 'dinosaur', 'date', 'ping', 'cat', '喵', '喵呜']
  const match = commands.find(c => c.startsWith(inputVal.value.toLowerCase()))
  if (match) inputVal.value = match
}

// Keyboard shortcut to close
onMounted(() => {
  const handleKey = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && props.visible) {
      close()
    }
  }
  document.addEventListener('keydown', handleKey)
})
</script>

<style scoped>
.terminal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  z-index: 99998;
  display: flex;
  align-items: flex-end;
  justify-content: center;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.terminal-window {
  width: 100%;
  max-width: 860px;
  height: 520px;
  background: #0d1117;
  border-radius: 12px 12px 0 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-bottom: none;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 -8px 40px rgba(0, 0, 0, 0.6);
}

.terminal-titlebar {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: #161b22;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  gap: 12px;
}

.terminal-dots {
  display: flex;
  gap: 6px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  cursor: pointer;
  transition: opacity 0.2s;
}
.dot:hover { opacity: 0.7; }
.dot--red { background: #ff5f57; }
.dot--yellow { background: #febc2e; }
.dot--green { background: #28c840; }

.terminal-title {
  flex: 1;
  text-align: center;
  font-family: 'DM Sans', sans-serif;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: 500;
}

.terminal-close {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.3);
  font-size: 20px;
  cursor: pointer;
  padding: 0 4px;
  line-height: 1;
  transition: color 0.2s;
}
.terminal-close:hover { color: rgba(255, 255, 255, 0.7); }

.terminal-output {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
  font-size: 13px;
  line-height: 1.6;
  scrollbar-width: thin;
  scrollbar-color: rgba(255,255,255,0.1) transparent;
}
.terminal-output::-webkit-scrollbar { width: 4px; }
.terminal-output::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 2px; }

.terminal-line {
  margin-bottom: 2px;
}

.terminal-line pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: inherit;
}

.terminal-line--input pre {
  color: #e6edf3;
}
.terminal-line--input pre::before {
  content: '$ ';
  color: #3fb950;
}

.terminal-line--output pre { color: #8b949e; }
.terminal-line--error pre { color: #f85149; }
.terminal-line--success pre { color: #3fb950; }
.terminal-line--art pre {
  color: #58a6ff;
  text-shadow: 0 0 8px rgba(88, 166, 255, 0.3);
}

.terminal-loading {
  display: flex;
  gap: 4px;
  padding: 4px 0;
}
.loading-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #58a6ff;
  animation: terminalPulse 1.2s ease-in-out infinite;
}
.loading-dot:nth-child(2) { animation-delay: 0.2s; }
.loading-dot:nth-child(3) { animation-delay: 0.4s; }
@keyframes terminalPulse {
  0%, 100% { opacity: 0.3; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

.terminal-input-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  background: #161b22;
  gap: 8px;
}

.terminal-prompt {
  color: #3fb950;
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  font-weight: bold;
}

.terminal-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: #e6edf3;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  caret-color: #58a6ff;
}
.terminal-input::placeholder { color: rgba(255,255,255,0.2); }

.terminal-hint {
  padding: 6px 16px 10px;
  background: #161b22;
  font-family: 'DM Sans', sans-serif;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.2);
  text-align: center;
}
.terminal-hint kbd {
  background: rgba(255,255,255,0.08);
  padding: 1px 6px;
  border-radius: 4px;
  font-family: inherit;
  border: 1px solid rgba(255,255,255,0.1);
}
.terminal-hint code {
  color: rgba(255,255,255,0.35);
}

/* Transition */
.terminal-fade-enter-active { transition: opacity 0.2s ease; }
.terminal-fade-leave-active { transition: opacity 0.15s ease; }
.terminal-fade-enter-from,
.terminal-fade-leave-to { opacity: 0; }
.terminal-fade-enter-active .terminal-window {
  animation: terminalSlideUp 0.25s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}
@keyframes terminalSlideUp {
  from { transform: translateY(40px) scale(0.95); }
  to { transform: translateY(0) scale(1); }
}
</style>
