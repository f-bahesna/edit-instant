<script setup>
import { 
  Flame, 
  CheckCircle2, 
  XCircle, 
  AlertTriangle, 
  RefreshCw 
} from 'lucide-vue-next'

defineProps({
  apiStatus: {
    type: Object,
    required: true
  }
})

defineEmits(['refresh'])
</script>

<template>
  <header class="border-b border-zinc-800 bg-[#0A0A0A]/80 backdrop-blur-md sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-6 py-4 flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div class="h-8 w-8 rounded bg-zinc-100 flex items-center justify-center">
          <Flame class="h-5 w-5 text-zinc-900" />
        </div>
        <div>
          <h1 class="text-sm font-semibold tracking-tight text-zinc-100 flex items-center gap-2">
            tiktok-agent
            <span class="text-[10px] bg-zinc-800 text-zinc-400 px-1.5 py-0.5 rounded font-mono">v1.0</span>
          </h1>
        </div>
      </div>

      <!-- Connection Badges -->
      <div class="flex items-center gap-3 text-xs font-medium">
        <!-- Backend API Status -->
        <div class="flex items-center gap-2 px-2.5 py-1 rounded bg-zinc-900 border" :class="apiStatus.connected ? 'border-zinc-800 text-zinc-300' : 'border-red-900/50 text-red-400'">
          <div class="h-1.5 w-1.5 rounded-full" :class="apiStatus.connected ? 'bg-emerald-500' : 'bg-red-500'"></div>
          <span>API: {{ apiStatus.connected ? 'Online' : 'Offline' }}</span>
        </div>

        <!-- Postgres DB Status -->
        <div class="flex items-center gap-2 px-2.5 py-1 rounded bg-zinc-900 border" :class="apiStatus.dbStatus === 'Connected' ? 'border-zinc-800 text-zinc-300' : 'border-red-900/50 text-red-400'">
          <CheckCircle2 v-if="apiStatus.dbStatus === 'Connected'" class="h-3 w-3 text-emerald-500" />
          <XCircle v-else class="h-3 w-3 text-red-500" />
          <span>DB</span>
        </div>

        <!-- FFmpeg Status -->
        <div class="flex items-center gap-2 px-2.5 py-1 rounded bg-zinc-900 border" :class="apiStatus.ffmpegStatus === 'Available' ? 'border-zinc-800 text-zinc-300' : 'border-yellow-900/50 text-yellow-500'">
          <CheckCircle2 v-if="apiStatus.ffmpegStatus === 'Available'" class="h-3 w-3 text-emerald-500" />
          <AlertTriangle v-else class="h-3 w-3 text-yellow-500" />
          <span>FFmpeg</span>
        </div>

        <button @click="$emit('refresh')" class="p-1.5 rounded hover:bg-zinc-800 text-zinc-400 hover:text-zinc-100 transition-colors">
          <RefreshCw class="h-3.5 w-3.5" :class="{ 'animate-spin': apiStatus.loading }" />
        </button>
      </div>
    </div>
  </header>
</template>
