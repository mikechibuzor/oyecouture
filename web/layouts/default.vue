<template>
  <div class="storefront-shell">
    <StorefrontAnnouncementBar />
    <StorefrontNavBar />
    <main class="storefront-main">
      <slot />
    </main>
    <StorefrontFooter />
  </div>
</template>
