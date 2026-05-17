import { ref } from 'vue'

export function useVideoGeneration() {
  const apiStatus = ref({
    loading: true,
    connected: false,
    dbStatus: 'Disconnected',
    ffmpegStatus: 'Not Available',
    ffmpegInfo: 'Unknown',
  })

  const formData = ref({
    title: '',
    script_idea: '',
    product_url: '',
  })

  const generationLogs = ref([])
  const activeJob = ref(null)
  const previewVideoUrl = ref(null)
  const selectedVideoStyle = ref('modern')

  const stylePresets = [
    { id: 'modern', name: 'Modern SaaS', desc: 'Fast, high-contrast, professional zooms' },
    { id: 'energy', name: 'Energetic Creator', desc: 'CapCut style transitions, sound effects' },
    { id: 'minimalist', name: 'Clean Tutorial', desc: 'Beginner-friendly pacing, subtle background music' }
  ]

  function addLog(source, message, type = 'info') {
    const timestamp = new Date().toLocaleTimeString()
    generationLogs.value.unshift({ timestamp, source, message, type })
  }

  async function checkSystemHealth() {
    apiStatus.value.loading = true
    try {
      const res = await fetch('http://localhost:8080/api/health')
      if (res.ok) {
        const data = await res.json()
        apiStatus.value = {
          loading: false,
          connected: true,
          dbStatus: data.database,
          ffmpegStatus: data.ffmpeg,
          ffmpegInfo: data.ffmpeg_info,
        }
        addLog('System', 'Successfully connected to backend API.')
      } else {
        throw new Error('Non-200 response')
      }
    } catch (err) {
      apiStatus.value = {
        loading: false,
        connected: false,
        dbStatus: 'Disconnected',
        ffmpegStatus: 'Not Available',
        ffmpegInfo: err.message || 'Connection refused',
      }
      addLog('System', 'Error connecting to backend API.', 'error')
    }
  }

  async function triggerVideoGeneration() {
    if (!formData.value.title || !formData.value.product_url) {
      addLog('Client', 'Validation failed: Title and Product URL are required.', 'error')
      return
    }

    addLog('Client', `Initiating generation for: "${formData.value.title}"`)
    activeJob.value = {
      status: 'processing',
      progress: 10,
      currentStep: 'Analyzing website UI & layout...',
    }

    try {
      // Fake API call for the mockup
      addLog('API', `Server accepted job: job_test_123`)

      setTimeout(() => {
        activeJob.value.progress = 30
        activeJob.value.currentStep = 'Scraping target website with Playwright...'
        addLog('Playwright', 'Navigating to Product URL in headless mobile viewport...')
      }, 1500)

      setTimeout(() => {
        activeJob.value.progress = 55
        activeJob.value.currentStep = 'Recording user interaction flows...'
        addLog('Playwright', 'Simulating cursor path, scrolling dashboard elements...')
      }, 3500)

      setTimeout(() => {
        activeJob.value.progress = 75
        activeJob.value.currentStep = 'Generating synthetic AI Voice and Whisper subtitles...'
        addLog('ElevenLabs', 'Generated high-retention audio narration. Syncing SRT files...')
      }, 5500)

      setTimeout(() => {
        activeJob.value.progress = 90
        activeJob.value.currentStep = 'Merging assets & rendering MP4 via FFmpeg...'
        addLog('FFmpeg', 'Encoding 9:16 vertical video at 30fps using CPU libx264...')
      }, 7500)

      setTimeout(() => {
        activeJob.value.progress = 100
        activeJob.value.status = 'completed'
        activeJob.value.data = {
            video_url: '#',
            subtitle_url: '#'
        }
        previewVideoUrl.value = 'https://assets.mixkit.co/videos/preview/mixkit-holding-a-smartphone-playing-a-video-game-41852-large.mp4'
        addLog('System', 'SUCCESS! Video rendered & ready for upload.', 'success')
      }, 9500)

    } catch (error) {
      activeJob.value.status = 'failed'
      activeJob.value.currentStep = 'Error during pipeline execution'
      addLog('System', `Failed to generate video: ${error.message}`, 'error')
    }
  }

  return {
    apiStatus,
    formData,
    generationLogs,
    activeJob,
    previewVideoUrl,
    selectedVideoStyle,
    stylePresets,
    checkSystemHealth,
    triggerVideoGeneration
  }
}
