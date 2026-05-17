<script setup>
import { Terminal } from 'lucide-vue-next'

defineProps({
  logs: {
    type: Array,
    required: true
  }
})
</script>

<template>
  <div class="bg-[#0D0D0D] border border-zinc-800 rounded-xl flex-1 flex flex-col min-h-[250px] overflow-hidden">
    <div class="px-4 py-2.5 border-b border-zinc-800 flex items-center justify-between bg-[#111]">
      <span class="flex items-center gap-2 text-xs font-mono text-zinc-400">
        <Terminal class="h-3.5 w-3.5 text-zinc-500" />
        Execution Logs
      </span>
      <span class="text-[10px] text-zinc-600 font-mono">stdout</span>
    </div>
    
    <div class="flex-1 p-4 font-mono text-[11px] overflow-y-auto max-h-[300px] flex flex-col gap-1.5 bg-[#0A0A0A]">
      <div v-if="logs.length === 0" class="text-zinc-600 flex items-center justify-center h-full text-xs">
        System idle. Ready for input.
      </div>
      <div 
        v-for="(log, idx) in logs" 
        :key="idx" 
        class="flex items-start gap-3 leading-relaxed"
        :class="{
          'text-red-400': log.type === 'error',
          'text-emerald-400': log.type === 'success',
          'text-zinc-300': log.type === 'info'
        }"
      >
        <span class="text-zinc-600 shrink-0 select-none">[{{ log.timestamp }}]</span>
        <div class="flex flex-col sm:flex-row sm:gap-2">
          <span class="font-semibold shrink-0" :class="{
            'text-blue-400': log.source === 'Playwright',
            'text-purple-400': log.source === 'ElevenLabs',
            'text-yellow-400': log.source === 'FFmpeg',
            'text-zinc-500': log.source === 'System' || log.source === 'Client' || log.source === 'API'
          }">[{{ log.source }}]</span>
          <span class="break-words">{{ log.message }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
