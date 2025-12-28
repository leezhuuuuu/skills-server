<template>
  <div class="flex flex-col min-h-screen">
    <!-- Navbar -->
    <div class="navbar bg-base-100 border-b border-base-300 px-4 sm:px-8">
      <div class="flex-1">
        <router-link to="/" class="btn btn-ghost text-xl normal-case gap-2">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-primary"><path d="m12 3-1.912 5.813a2 2 0 0 1-1.275 1.275L3 12l5.813 1.912a2 2 0 0 1 1.275 1.275L12 21l1.912-5.813a2 2 0 0 1 1.275-1.275L21 12l-5.813-1.912a2 2 0 0 1-1.275-1.275L12 3Z"/></svg>
          Skills Registry
        </router-link>
      </div>
      <div class="flex-none gap-2">
        <a href="https://github.com/leezhuuuuu/skills-server" target="_blank" class="btn btn-square btn-ghost">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 22v-4a4.8 4.8 0 0 0-1-3.5c3 0 6-2 6-5.5.08-1.25-.27-2.48-1-3.5.28-1.15.28-2.35 0-3.5 0 0-1 0-3 1.5-2.64-.5-5.36.5-8 0C6 2 5 2 5 2c-.3 1.15-.3 2.35 0 3.5A5.403 5.403 0 0 0 4 9c0 3.5 3 5.5 6 5.5-.39.49-.68 1.05-.85 1.65-.17.6-.22 1.23-.15 1.85v4"/><path d="M9 18c-4.51 2-5-2-7-2"/></svg>
        </a>
      </div>
    </div>

    <!-- Main Content -->
    <main class="flex-1 container mx-auto px-4 sm:px-8 py-8">
      <router-view></router-view>
    </main>

    <!-- Footer -->
    <footer class="footer items-center p-4 bg-neutral text-neutral-content">
      <aside class="items-center grid-flow-col">
        <p>Skills MCP Server</p>
      </aside>
    </footer>
  </div>
</template>
