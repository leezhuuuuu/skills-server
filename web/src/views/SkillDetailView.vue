<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { getSkillDetail } from '@/api/client'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'

const route = useRoute()
const skill = ref(null)
const loading = ref(true)
const activeTab = ref('readme') // 'readme' or 'files'

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value;
      } catch (__) {}
    }
    return ''; // use external default escaping
  }
})

const renderedReadme = computed(() => {
  return skill.value?.readme ? md.render(skill.value.readme) : ''
})

onMounted(async () => {
  try {
    const res = await getSkillDetail(route.params.name)
    skill.value = res.data
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
})

const getDownloadUrl = () => {
  return `/api/v1/download/${skill.value.name}`
}
</script>

<template>
  <div v-if="loading" class="animate-pulse">
    <div class="h-8 bg-base-200 w-1/3 mb-4 rounded"></div>
    <div class="h-4 bg-base-200 w-1/2 mb-8 rounded"></div>
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-2 h-96 bg-base-200 rounded"></div>
      <div class="h-64 bg-base-200 rounded"></div>
    </div>
  </div>

  <div v-else-if="!skill" class="text-center py-12">
    <h2 class="text-2xl font-bold text-error">Skill not found</h2>
  </div>

  <div v-else>
    <!-- Header -->
    <div class="mb-8">
      <div class="text-sm breadcrumbs mb-2">
        <ul>
          <li><router-link to="/">Home</router-link></li>
          <li>{{ skill.name }}</li>
        </ul>
      </div>
      <div class="flex flex-col md:flex-row gap-4 justify-between items-start md:items-center">
        <div>
          <h1 class="text-3xl font-bold flex items-center gap-3">
            {{ skill.name }}
            <div class="badge badge-primary">{{ skill.version || 'v1.0' }}</div>
          </h1>
          <p class="text-base-content/60 mt-2 text-lg">{{ skill.description }}</p>
        </div>
        <div class="flex gap-2">
          <a :href="getDownloadUrl()" class="btn btn-primary" download>
            Download ZIP
          </a>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Main Content -->
      <div class="lg:col-span-2 space-y-6">
        <!-- Tabs -->
        <div class="tabs tabs-boxed bg-base-200">
          <a
            class="tab"
            :class="{ 'tab-active': activeTab === 'readme' }"
            @click="activeTab = 'readme'"
          >README / SKILL.md</a>
          <a
            class="tab"
            :class="{ 'tab-active': activeTab === 'files' }"
            @click="activeTab = 'files'"
          >File Structure</a>
        </div>

        <!-- Readme Tab -->
        <div v-show="activeTab === 'readme'" class="card bg-base-200 p-6">
          <div class="prose prose-sm md:prose-base max-w-none prose-invert" v-html="renderedReadme"></div>
        </div>

        <!-- Files Tab -->
        <div v-show="activeTab === 'files'" class="card bg-base-200 p-6">
          <pre class="font-mono text-sm overflow-x-auto">{{ skill.file_tree }}</pre>
        </div>
      </div>

      <!-- Sidebar -->
      <div class="space-y-6">
        <!-- Install Card -->
        <div class="card bg-base-300 shadow-xl">
          <div class="card-body p-5">
            <h3 class="card-title text-sm uppercase text-base-content/50">Installation</h3>
            <div class="mockup-code bg-black text-xs scale-95 -ml-4 w-[110%]">
              <pre data-prefix="$"><code>uv tool install skills-mcp</code></pre>
              <pre data-prefix="$" class="text-warning"><code>skills install {{ skill.name }}</code></pre>
            </div>
          </div>
        </div>

        <!-- Metadata Card -->
        <div class="card bg-base-200 border border-base-300">
          <div class="card-body p-5 space-y-3">
            <h3 class="card-title text-sm uppercase text-base-content/50">Metadata</h3>

            <div class="flex justify-between border-b border-base-content/10 pb-2">
              <span class="text-base-content/70">Author</span>
              <span class="font-medium">{{ skill.author || 'Unknown' }}</span>
            </div>

            <div class="flex justify-between border-b border-base-content/10 pb-2">
              <span class="text-base-content/70">Last Updated</span>
              <span class="font-medium">{{ new Date(skill.updated_at).toLocaleDateString() }}</span>
            </div>

            <div class="pt-2">
              <div class="text-base-content/70 mb-2">Tags</div>
              <div class="flex flex-wrap gap-1">
                <div v-for="tag in skill.tags" :key="tag" class="badge badge-outline">{{ tag }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Raw Link -->
        <a :href="`/skill/${skill.name}.md`" target="_blank" class="btn btn-ghost btn-block btn-sm">
          View Raw Context (for LLM)
        </a>
      </div>
    </div>
  </div>
</template>
