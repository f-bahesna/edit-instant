<script setup>
import { Link as LinkIcon, Play, RefreshCw } from 'lucide-vue-next'

const props = defineProps({
  formData: {
    type: Object,
    required: true
  },
  activeJob: {
    type: Object,
    default: null
  },
  stylePresets: {
    type: Array,
    required: true
  },
  selectedVideoStyle: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:selectedVideoStyle', 'submit'])

const isProcessing = () => props.activeJob && props.activeJob.status === 'processing'
</script>

<template>
  <div class="bg-[#0D0D0D] border border-zinc-800 rounded-xl p-6">
    <h2 class="text-sm font-medium text-zinc-100 mb-5">Job Configuration</h2>

    <div class="flex flex-col gap-5">
      <!-- Title -->
      <div>
        <label class="block text-sm font-medium text-zinc-400 mb-1.5">Video Title</label>
        <input 
          v-model="formData.title" 
          type="text" 
          placeholder="e.g. How to deploy Next.js in 30 seconds" 
          class="w-full bg-[#111] border border-zinc-800 rounded-lg px-3 py-2 text-sm text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-zinc-600 focus:ring-1 focus:ring-zinc-600 transition"
        />
      </div>

      <!-- Product URL -->
      <div>
        <label class="block text-sm font-medium text-zinc-400 mb-1.5">Target URL</label>
        <div class="relative">
          <LinkIcon class="absolute left-3 top-2.5 h-4 w-4 text-zinc-500" />
          <input 
            v-model="formData.product_url" 
            type="url" 
            placeholder="https://vercel.com" 
            class="w-full bg-[#111] border border-zinc-800 rounded-lg pl-9 pr-3 py-2 text-sm text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-zinc-600 focus:ring-1 focus:ring-zinc-600 transition"
          />
        </div>
      </div>

      <!-- Script Idea -->
      <div>
        <label class="block text-sm font-medium text-zinc-400 mb-1.5">Script Idea (Optional)</label>
        <textarea 
          v-model="formData.script_idea" 
          rows="3" 
          placeholder="Focus on the dashboard UI and the seamless deployment flow..." 
          class="w-full bg-[#111] border border-zinc-800 rounded-lg px-3 py-2 text-sm text-zinc-200 placeholder-zinc-600 focus:outline-none focus:border-zinc-600 focus:ring-1 focus:ring-zinc-600 transition resize-none"
        ></textarea>
      </div>

      <!-- Presets Selection -->
      <div>
        <label class="block text-sm font-medium text-zinc-400 mb-2">Aesthetic Preset</label>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <button 
            v-for="preset in stylePresets" 
            :key="preset.id"
            @click="emit('update:selectedVideoStyle', preset.id)"
            class="p-3 rounded-lg border text-left transition flex flex-col gap-1"
            :class="selectedVideoStyle === preset.id ? 'border-zinc-400 bg-zinc-800/30 text-zinc-100' : 'border-zinc-800 bg-[#111] text-zinc-500 hover:border-zinc-700'"
          >
            <span class="text-xs font-medium flex items-center gap-1.5">
              <div class="h-1.5 w-1.5 rounded-full" :class="selectedVideoStyle === preset.id ? 'bg-zinc-300' : 'bg-zinc-700'"></div>
              {{ preset.name }}
            </span>
            <span class="text-[10px] leading-relaxed line-clamp-1">{{ preset.desc }}</span>
          </button>
        </div>
      </div>

      <!-- Trigger Button -->
      <button 
        @click="emit('submit')"
        :disabled="isProcessing()"
        class="w-full mt-2 bg-zinc-100 hover:bg-white text-zinc-900 disabled:opacity-50 disabled:cursor-not-allowed font-medium py-2 px-4 rounded-lg flex items-center justify-center gap-2 transition"
      >
        <Play v-if="!isProcessing()" class="h-4 w-4" />
        <RefreshCw v-else class="h-4 w-4 animate-spin" />
        {{ isProcessing() ? 'Executing Pipeline...' : 'Run Generation' }}
      </button>
    </div>
  </div>
</template>
