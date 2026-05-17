<script setup>
import { Video, Tv, CheckCircle2, Download } from 'lucide-vue-next'

defineProps({
  activeJob: {
    type: Object,
    default: null
  },
  previewVideoUrl: {
    type: String,
    default: null
  }
})
</script>

<template>
  <div class="flex flex-col gap-6 items-center w-full">
    
    <!-- Minimalist Device Frame -->
    <div class="relative w-[280px] h-[500px] rounded-[32px] border-[6px] border-zinc-900 bg-[#050505] shadow-2xl overflow-hidden flex-shrink-0">
      
      <!-- Top Speaker / Notch Placeholder -->
      <div class="absolute top-0 inset-x-0 flex justify-center z-30">
        <div class="w-32 h-5 bg-zinc-900 rounded-b-xl flex items-end justify-center pb-1.5">
          <div class="w-12 h-1 bg-zinc-800 rounded-full"></div>
        </div>
      </div>

      <!-- Simulator Contents -->
      <div class="absolute inset-0 flex flex-col relative w-full h-full bg-[#050505]">
        
        <!-- Processing State -->
        <div v-if="activeJob && activeJob.status === 'processing'" class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center z-20 bg-[#050505]">
          <div class="relative mb-6">
            <div class="h-12 w-12 rounded-full border-2 border-zinc-800 border-t-zinc-300 animate-spin"></div>
            <Video class="h-5 w-5 text-zinc-400 absolute top-3.5 left-3.5" />
          </div>
          <h4 class="font-medium text-xs text-zinc-300 mb-1">Rendering Video</h4>
          <p class="text-[10px] text-zinc-500 mb-4 font-mono truncate w-full">{{ activeJob.currentStep }}</p>
          
          <!-- Progress Bar -->
          <div class="w-full bg-zinc-900 h-1 rounded-full overflow-hidden">
            <div class="bg-zinc-300 h-full transition-all duration-500" :style="{ width: activeJob.progress + '%' }"></div>
          </div>
          <span class="text-[9px] text-zinc-600 font-mono mt-2">{{ activeJob.progress }}%</span>
        </div>

        <!-- Ready State / Generic Video Preview -->
        <div v-else-if="previewVideoUrl" class="absolute inset-0 z-10 flex flex-col bg-black">
          <video :src="previewVideoUrl" autoplay loop muted class="h-full w-full object-cover"></video>
          
          <!-- Sleek Generic Overlay -->
          <div class="absolute bottom-0 inset-x-0 p-4 bg-gradient-to-t from-black/80 to-transparent flex flex-col gap-2 font-sans pointer-events-none">
             <div class="w-full text-center px-3 py-1.5 bg-black/50 rounded-lg backdrop-blur-md border border-white/10 mb-2">
                <p class="text-xs font-semibold text-white tracking-wide">
                  "Click the button to deploy."
                </p>
             </div>
             <div>
                <h5 class="text-[11px] font-semibold text-white">Generated Video Preview</h5>
                <p class="text-[9px] text-zinc-400 mt-0.5 line-clamp-2">
                  {{ activeJob?.data?.caption || 'Seamless auto-generated walkthrough.' }}
                </p>
             </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="absolute inset-0 flex flex-col items-center justify-center p-6 text-center text-zinc-500">
          <Tv class="h-8 w-8 text-zinc-800 mb-3" />
          <h4 class="font-medium text-xs text-zinc-400">Preview Panel</h4>
          <p class="text-[10px] leading-relaxed mt-1.5 text-zinc-600">
            Awaiting input...
          </p>
        </div>
        
      </div>
    </div>

    <!-- Download Assets Panel -->
    <div v-if="activeJob && activeJob.status === 'completed'" class="w-full max-w-[280px] bg-[#0D0D0D] border border-zinc-800 rounded-xl p-4 flex flex-col gap-3">
      <h4 class="text-[11px] font-semibold text-zinc-300 flex items-center gap-1.5 border-b border-zinc-800/50 pb-2">
        <CheckCircle2 class="h-3.5 w-3.5 text-emerald-500" />
        Assets Ready
      </h4>
      <div class="flex flex-col gap-1.5">
        <div class="flex items-center justify-between p-2 rounded bg-zinc-900 border border-zinc-800/50 hover:bg-zinc-800 transition group cursor-pointer">
          <span class="text-[10px] text-zinc-400 font-mono group-hover:text-zinc-300">video.mp4</span>
          <a :href="activeJob.data.video_url" download class="text-zinc-500 hover:text-white" target="_blank">
            <Download class="h-3.5 w-3.5" />
          </a>
        </div>
        <div class="flex items-center justify-between p-2 rounded bg-zinc-900 border border-zinc-800/50 hover:bg-zinc-800 transition group cursor-pointer">
          <span class="text-[10px] text-zinc-400 font-mono group-hover:text-zinc-300">subtitles.srt</span>
          <a :href="activeJob.data.subtitle_url" download class="text-zinc-500 hover:text-white" target="_blank">
            <Download class="h-3.5 w-3.5" />
          </a>
        </div>
      </div>
    </div>

  </div>
</template>
