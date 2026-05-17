<script setup>
import { onMounted } from 'vue'
import { useVideoGeneration } from './composables/useVideoGeneration'

import AppHeader from './components/AppHeader.vue'
import VideoGeneratorForm from './components/VideoGeneratorForm.vue'
import AgentConsole from './components/AgentConsole.vue'
import DevicePreview from './components/DevicePreview.vue'

const {
  apiStatus,
  formData,
  generationLogs,
  activeJob,
  previewVideoUrl,
  selectedVideoStyle,
  stylePresets,
  checkSystemHealth,
  triggerVideoGeneration
} = useVideoGeneration()

onMounted(() => {
  checkSystemHealth()
})
</script>

<template>
  <div class="min-h-screen bg-[#0A0A0A] text-zinc-300 pb-16 font-sans">
    
    <AppHeader 
      :apiStatus="apiStatus" 
      @refresh="checkSystemHealth" 
    />

    <!-- Main Workspace -->
    <main class="max-w-[1100px] mx-auto px-6 mt-8 grid grid-cols-1 lg:grid-cols-12 gap-10">
      
      <!-- Left Column (7 cols) -->
      <section class="lg:col-span-7 flex flex-col gap-8">
        
        <VideoGeneratorForm 
          :formData="formData"
          :activeJob="activeJob"
          :stylePresets="stylePresets"
          v-model:selectedVideoStyle="selectedVideoStyle"
          @submit="triggerVideoGeneration"
        />

        <AgentConsole 
          :logs="generationLogs"
        />

      </section>

      <!-- Right Column (5 cols) -->
      <section class="lg:col-span-5 flex flex-col items-center pt-2">
        
        <DevicePreview 
          :activeJob="activeJob"
          :previewVideoUrl="previewVideoUrl"
        />

      </section>

    </main>
  </div>
</template>

