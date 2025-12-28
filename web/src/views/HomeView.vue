<script setup>
import { ref, onMounted, watch } from 'vue'
import { getSkills } from '@/api/client'
import { useRouter } from 'vue-router'
import { formatDistanceToNow } from 'date-fns'

const skills = ref([])
const loading = ref(false)
const searchQuery = ref('')
const router = useRouter()

const fetchSkills = async () => {
  loading.value = true
  try {
    const res = await getSkills(searchQuery.value)
    skills.value = res.data.skills || []
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

// Debounce search
let timeout
watch(searchQuery, () => {
  clearTimeout(timeout)
  timeout = setTimeout(fetchSkills, 300)
})

onMounted(fetchSkills)

const goToDetail = (name) => {
  router.push(`/skill/${name}`)
}
</script>

<template>
  <div class="flex flex-col items-center mb-12">
    <h1 class="text-4xl font-bold mb-4 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
      Discover Agent Skills
    </h1>
    <p class="text-base-content/60 mb-8 text-lg">
      Supercharge your AI Agents with deterministic capabilities.
    </p>

    <div class="form-control w-full max-w-lg">
      <input
        type="text"
        placeholder="Search skills..."
        class="input input-bordered w-full"
        v-model="searchQuery"
      />
    </div>
  </div>

  <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div v-for="i in 3" :key="i" class="card bg-base-200 h-48 animate-pulse"></div>
  </div>

  <div v-else-if="skills.length === 0" class="text-center py-12 text-base-content/50">
    No skills found matching your query.
  </div>

  <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
    <div
      v-for="skill in skills"
      :key="skill.name"
      class="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer border border-base-300 hover:border-primary/50"
      @click="goToDetail(skill.name)"
    >
      <div class="card-body">
        <div class="flex justify-between items-start">
          <h2 class="card-title text-primary">{{ skill.name }}</h2>
          <div class="badge badge-outline text-xs">{{ skill.version || 'v1.0' }}</div>
        </div>
        <p class="text-sm text-base-content/70 line-clamp-2 my-2">
          {{ skill.description }}
        </p>
        <div class="card-actions justify-between items-center mt-4">
          <div class="flex gap-1 flex-wrap">
            <span v-for="tag in skill.tags" :key="tag" class="badge badge-secondary badge-xs">
              {{ tag }}
            </span>
          </div>
          <span class="text-xs text-base-content/50">
            Updated {{ formatDistanceToNow(new Date(skill.updated_at)) }} ago
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
