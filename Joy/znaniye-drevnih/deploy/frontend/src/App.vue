<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const ws = ref<WebSocket | null>(null)
const wsStatus = ref<'connecting' | 'open' | 'closed'>('connecting')

const awaitingStart = ref(true)
const mode = ref<'idle' | 'dialog' | 'task'>('idle')
const currentDialog = ref<DialogMessage | null>(null)
const currentPhrase = ref<PhraseMessage | null>(null)
const currentTask = ref<TaskMessage | null>(null)
const lastResult = ref<ResultMessage | null>(null)
const gameStatus = ref<GameStatus | null>(null)
const flagValue = ref<string | null>(null)
const hasSeenZuck = ref(false)

const code = ref('')

const bgmRef = ref<HTMLAudioElement | null>(null)
const audioEnabled = ref(false)
const voiceAudio = new Audio()
let phraseToken = 0
let voiceToken = 0
const BGM_BASE_VOLUME = 0.6
const BGM_DUCK_VOLUME = 0.33

const wsUrl = computed(() => {
  const env = import.meta.env.VITE_WS_URL
  if (env) return env
  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  return `${protocol}://${window.location.host}/ws`
})

const activeRole = computed<Role | null>(() => {
  if (currentDialog.value) return currentDialog.value.role
  if (currentPhrase.value) return currentPhrase.value.role
  return null
})

const hideSlavDuringPhrase = computed(
  () => currentPhrase.value && activeRole.value === 'zuck'
)
const hideZuckDuringPhrase = computed(
  () => currentPhrase.value && activeRole.value === 'slav'
)

const showStartOverlay = computed(() => awaitingStart.value && !gameStatus.value)
const showOutcomeOverlay = computed(() => gameStatus.value && gameStatus.value !== 'win')

const outcomeTitle = computed(() => {
  switch (gameStatus.value) {
    case 'timeout':
      return 'Время ожидания ответа вышло'
    case 'wrong':
      return 'Ответ неверен'
    case 'error':
      return 'Произошла ошибка'
    case 'win':
      return 'Победа!'
    default:
      return ''
  }
})

const outcomeHint = computed(() => {
  switch (gameStatus.value) {
    case 'timeout':
      return 'Соберись и начни заново.'
    case 'wrong':
      return 'Попытайся снова, славянин.'
    case 'error':
      return 'Связь с Перуном дала сбой.'
    case 'win':
      return 'Свершилось. Сохрани знамение.'
    default:
      return ''
  }
})

function connectWs() {
  if (ws.value) {
    ws.value.close()
  }
  wsStatus.value = 'connecting'
  const socket = new WebSocket(wsUrl.value)
  ws.value = socket

  socket.onopen = () => {
    wsStatus.value = 'open'
  }
  socket.onclose = () => {
    wsStatus.value = 'closed'
  }
  socket.onerror = () => {
    wsStatus.value = 'closed'
  }
  socket.onmessage = (event) => {
    handleServerMessage(event.data)
  }
}

function handleServerMessage(raw: string) {
  let msg: ServerMessage
  try {
    msg = JSON.parse(raw)
  } catch {
    return
  }

  switch (msg.type) {
    case 'dialog':
      currentDialog.value = msg
      currentPhrase.value = null
      currentTask.value = null
      lastResult.value = null
      mode.value = 'dialog'
      if (msg.role === 'zuck') hasSeenZuck.value = true
      playVoice(msg.audio)
      return
    case 'phrase':
      currentPhrase.value = msg
      flagValue.value = msg.flag ?? flagValue.value
      if (msg.role === 'zuck') hasSeenZuck.value = true
      playPhrase(msg)
      return
    case 'task':
      currentTask.value = msg
      currentDialog.value = null
      lastResult.value = null
      mode.value = 'task'
      return
    case 'result':
      lastResult.value = msg
      if (msg.success) {
        code.value = ''
      }
      return
    case 'game':
      gameStatus.value = msg.status
      awaitingStart.value = true
      if (msg.status !== 'win') {
        currentDialog.value = null
        currentPhrase.value = null
      }
      mode.value = 'idle'
      return
  }
}

async function playPhrase(msg: PhraseMessage) {
  const token = ++phraseToken
  await playVoice(msg.audio)
  if (phraseToken !== token) return
  if (msg.flag) return
  setTimeout(() => {
    if (phraseToken === token) {
      currentPhrase.value = null
    }
  }, 450)
}

function playVoice(base64: string) {
  if (!base64) return Promise.resolve()
  return new Promise<void>((resolve) => {
    const token = ++voiceToken
    if (bgmRef.value) {
      bgmRef.value.volume = BGM_DUCK_VOLUME
    }
    voiceAudio.pause()
    voiceAudio.currentTime = 0
    voiceAudio.volume = 1
    voiceAudio.onended = () => {
      if (token === voiceToken && bgmRef.value) {
        bgmRef.value.volume = BGM_BASE_VOLUME
      }
      resolve()
    }
    voiceAudio.onerror = () => {
      if (token === voiceToken && bgmRef.value) {
        bgmRef.value.volume = BGM_BASE_VOLUME
      }
      resolve()
    }
    voiceAudio.src = `data:audio/mpeg;base64,${base64}`
    voiceAudio.play().catch(() => {
      if (token === voiceToken && bgmRef.value) {
        bgmRef.value.volume = BGM_BASE_VOLUME
      }
      resolve()
    })
  })
}

function tryEnableAudio() {
  if (!bgmRef.value) return
  bgmRef.value.volume = BGM_BASE_VOLUME
  bgmRef.value
    .play()
    .then(() => {
      audioEnabled.value = true
    })
    .catch(() => {
      audioEnabled.value = false
    })
}

function resetUiForStart() {
  awaitingStart.value = false
  gameStatus.value = null
  currentDialog.value = null
  currentPhrase.value = null
  currentTask.value = null
  lastResult.value = null
  flagValue.value = null
  hasSeenZuck.value = false
  mode.value = 'idle'
  code.value = ''
}

function sendStart() {
  resetUiForStart()
  tryEnableAudio()
  sendMessage({ type: 'start' })
}

function sendContinue() {
  if (currentDialog.value) {
    voiceAudio.pause()
    voiceAudio.currentTime = 0
    voiceToken++
    if (bgmRef.value) {
      bgmRef.value.volume = BGM_BASE_VOLUME
    }
  }
  sendMessage({ type: 'continue' })
}

function sendSubmit() {
  if (!code.value.trim()) return
  const normalized = normalizeUtf8(code.value)
  sendMessage({ type: 'submit', code: normalized })
}

function sendMessage(payload: Record<string, unknown>) {
  const json = JSON.stringify(payload)
  const encoded = new TextEncoder().encode(json)
  ws.value?.send(encoded)
}

function normalizeUtf8(input: string) {
  const hasCyrillic = /[\\u0400-\\u04FF]/.test(input)
  const looksMojibake = /[ÐÑ]/.test(input)
  if (hasCyrillic || !looksMojibake) return input
  try {
    const bytes = Uint8Array.from(input, (c) => c.charCodeAt(0))
    return new TextDecoder('utf-8').decode(bytes)
  } catch {
    return input
  }
}

onMounted(() => {
  connectWs()
  tryEnableAudio()
})

onBeforeUnmount(() => {
  ws.value?.close()
})

type Role = 'slav' | 'zuck'

type DialogMessage = {
  type: 'dialog'
  role: Role
  audio: string
  text: string
}

type PhraseMessage = {
  type: 'phrase'
  role: Role
  audio: string
  text: string
  flag?: string
}

type TaskMessage = {
  type: 'task'
  text: string
}

type ResultMessage = {
  type: 'result'
  couldRun: boolean
  containsNonRuChars: boolean
  success: boolean
}

type GameStatus = 'timeout' | 'wrong' | 'win' | 'error'

type GameMessage = {
  type: 'game'
  status: GameStatus
}

type ServerMessage = DialogMessage | PhraseMessage | TaskMessage | ResultMessage | GameMessage
</script>

<template>
  <div class="app">
    <audio ref="bgmRef" src="/audio/nueki_tolchonov_hit_fonk_extended_bgm.mp3" loop></audio>

    <header class="topbar">
      <div class="brand">
        <div class="brand__title">Знание древних</div>
      </div>
      <div class="status">
        <span class="status__dot" :class="`status__dot--${wsStatus}`"></span>
        <span class="status__text">{{ wsStatus === 'open' ? 'Связь с Перуном налажена' : 'Попытка достучаться до Перуна' }}</span>
      </div>
    </header>

    <main class="scene">
      <div class="stage">
        <img
          class="actor actor--left"
          :class="{
            active: activeRole === 'slav',
            idle: activeRole === 'zuck',
            hidden: !activeRole || hideSlavDuringPhrase,
          }"
          src="/images/slav.png"
          alt="Древний рус"
        />
        <img
          class="actor actor--right"
          :class="{
            active: activeRole === 'zuck',
            idle: activeRole === 'slav',
            hidden: !activeRole || !hasSeenZuck || hideZuckDuringPhrase,
          }"
          src="/images/zuck.png"
          alt="Змей"
        />

        <div v-if="currentPhrase" class="phrase">
          <div class="phrase__text">{{ currentPhrase.text }}</div>
          <div v-if="currentPhrase.flag" class="phrase__flag">{{ currentPhrase.flag }}</div>
        </div>

        <div v-if="mode === 'dialog'" class="dialog-overlay">
          <div class="dialog-overlay__text">
            {{ currentDialog?.text }}
          </div>
        </div>

        <button v-if="mode === 'dialog'" class="cta dialog-cta" @click="sendContinue">
          Продолжить
        </button>
      </div>

      <section class="task-panel" :class="{ disabled: mode !== 'task' }">
        <div class="task-panel__header">
          <div class="task-panel__title">Береста</div>
        </div>
        <div class="task-panel__text">
          {{ currentTask?.text || 'Ожидается задание...' }}
        </div>
        <div class="task-panel__editor">
          <textarea
            v-model="code"
            placeholder="Начертай код здесь..."
            :disabled="mode !== 'task'"
          ></textarea>
        </div>
        <div class="task-panel__actions">
          <button class="cta" :disabled="mode !== 'task'" @click="sendSubmit">
            Отправить
          </button>
          <div class="task-panel__note">
            {{ audioEnabled ? 'Звук славянский пробуждён' : 'Нажми «Начать», чтобы пробудить славянский звук' }}
          </div>
        </div>
      </section>
    </main>

    <transition name="veil">
      <div v-if="showStartOverlay" class="overlay">
        <div class="overlay__card">
          <div class="overlay__title">Начать испытание</div>
          <div class="overlay__text">
            Нажми, чтобы разжечь огонь и услышать голоса предков
          </div>
          <button class="cta cta--wide" @click="sendStart">Начать</button>
        </div>
      </div>
    </transition>

    <transition name="veil">
      <div v-if="showOutcomeOverlay" class="overlay">
        <div class="overlay__card">
          <div class="overlay__title">{{ outcomeTitle }}</div>
          <div class="overlay__text">{{ outcomeHint }}</div>
          <div v-if="gameStatus === 'win' && flagValue" class="overlay__flag">
            {{ flagValue }}
          </div>
          <button class="cta cta--wide" @click="sendStart">
            {{ gameStatus === 'win' ? 'Сыграть снова' : 'Попробовать снова' }}
          </button>
        </div>
      </div>
    </transition>
  </div>
</template>
