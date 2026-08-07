<script setup lang="ts">
import { ChevronDown } from 'lucide'

// A collapsible panel built on native <details>. The element is keyboard
// operable and announced correctly without any script, so the disclosure still
// works if hydration never arrives — which matters here, because the content is
// delivery and returns policy a customer may need before they buy.
interface UiDisclosureProps {
  title: string
}

defineProps<UiDisclosureProps>()
</script>

<template>
  <details class="ui-disclosure">
    <summary class="ui-disclosure-summary">
      <span>{{ title }}</span>
      <StorefrontIcon
        class="ui-disclosure-chevron"
        :icon-node="ChevronDown"
        :size="16"
        :stroke-width="1.75"
      />
    </summary>

    <div class="ui-disclosure-panel">
      <slot />
    </div>
  </details>
</template>
